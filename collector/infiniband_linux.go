// Copyright 2017-2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !noinfiniband

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

// efaVendorID is the PCI vendor ID for AWS Elastic Fabric Adapter.
// EFA devices register under /sys/class/infiniband but do NOT populate the
// IB-spec port_xmit_data / port_rcv_data counters. Bytes/packets live in
// hw_counters/{tx,rx}_{bytes,pkts} as raw values (no IB 4-octet scaling).
const efaVendorID = "0x1d0f"

type infinibandCollector struct {
	fs          sysfs.FS
	sysPath     string
	metricDescs map[string]*prometheus.Desc
	logger      *slog.Logger
	subsystem   string
}

func init() {
	registerCollector("infiniband", defaultEnabled, NewInfiniBandCollector)
}

// NewInfiniBandCollector returns a new Collector exposing InfiniBand stats.
func NewInfiniBandCollector(logger *slog.Logger) (Collector, error) {
	var i infinibandCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.sysPath = *sysPath
	i.logger = logger

	// Detailed description for all metrics.
	descriptions := map[string]string{
		"legacy_multicast_packets_received_total":    "Number of multicast packets received",
		"legacy_multicast_packets_transmitted_total": "Number of multicast packets transmitted",
		"legacy_data_received_bytes_total":           "Number of data octets received on all links",
		"legacy_packets_received_total":              "Number of data packets received on all links",
		"legacy_unicast_packets_received_total":      "Number of unicast packets received",
		"legacy_unicast_packets_transmitted_total":   "Number of unicast packets transmitted",
		"legacy_data_transmitted_bytes_total":        "Number of data octets transmitted on all links",
		"legacy_packets_transmitted_total":           "Number of data packets received on all links",
		"excessive_buffer_overrun_errors_total":      "Number of times that OverrunErrors consecutive flow control update periods occurred, each having at least one overrun error.",
		"link_downed_total":                          "Number of times the link failed to recover from an error state and went down",
		"link_error_recovery_total":                  "Number of times the link successfully recovered from an error state",
		"local_link_integrity_errors_total":          "Number of times that the count of local physical errors exceeded the threshold specified by LocalPhyErrors.",
		"multicast_packets_received_total":           "Number of multicast packets received (including errors)",
		"multicast_packets_transmitted_total":        "Number of multicast packets transmitted (including errors)",
		"physical_state_id":                          "Physical state of the InfiniBand port (0: no change, 1: sleep, 2: polling, 3: disable, 4: shift, 5: link up, 6: link error recover, 7: phytest)",
		"port_constraint_errors_received_total":      "Number of packets received on the switch physical port that are discarded",
		"port_constraint_errors_transmitted_total":   "Number of packets not transmitted from the switch physical port",
		"port_data_received_bytes_total":             "Number of data octets received on all links",
		"port_data_transmitted_bytes_total":          "Number of data octets transmitted on all links",
		"port_discards_received_total":               "Number of inbound packets discarded by the port because the port is down or congested",
		"port_discards_transmitted_total":            "Number of outbound packets discarded by the port because the port is down or congested",
		"port_errors_received_total":                 "Number of packets containing an error that were received on this port",
		"port_packets_received_total":                "Number of packets received on all VLs by this port (including errors)",
		"port_packets_transmitted_total":             "Number of packets transmitted on all VLs from this port (including errors)",
		"port_transmit_wait_total":                   "Number of ticks during which the port had data to transmit but no data was sent during the entire tick",
		"rate_bytes_per_second":                      "Maximum signal transfer rate",
		"state_id":                                   "State of the InfiniBand port (0: no change, 1: down, 2: init, 3: armed, 4: active, 5: act defer)",
		"unicast_packets_received_total":             "Number of unicast packets received (including errors)",
		"unicast_packets_transmitted_total":          "Number of unicast packets transmitted (including errors)",
		"port_receive_remote_physical_errors_total":  "Number of packets marked with the EBP (End of Bad Packet) delimiter received on the port.",
		"port_receive_switch_relay_errors_total":     "Number of packets that could not be forwarded by the switch.",
		"symbol_error_total":                         "Number of minor link errors detected on one or more physical lanes.",
		"vl15_dropped_total":                         "Number of incoming VL15 packets dropped due to resource limitations.",

		// https://enterprise-support.nvidia.com/s/article/understanding-mlx5-linux-counters-and-status-parameters
		"duplicate_requests_packets_total":          "The number of received packets. A duplicate request is a request that had been previously executed.",
		"implied_nak_seq_errors_total":              "The number of time the requested decided an ACK. with a PSN larger than the expected PSN for an RDMA read or response.",
		"lifespan_seconds":                          "The maximum period in ms which defines the aging of the counter reads. Two consecutive reads within this period might return the same values.",
		"local_ack_timeout_errors_total":            "The number of times QP's ack timer expired for RC, XRC, DCT QPs at the sender side. The QP retry limit was not exceed, therefore it is still recoverable error.",
		"np_cnp_packets_sent_total":                 "The number of CNP packets sent by the Notification Point when it noticed congestion experienced in the RoCEv2 IP header (ECN bits). The counters was added in MLNX_OFED 4.1",
		"np_ecn_marked_roce_packets_received_total": "The number of RoCEv2 packets received by the notification point which were marked for experiencing the congestion (ECN bits where '11' on the ingress RoCE traffic) . The counters was added in MLNX_OFED 4.1",
		"out_of_buffer_drops_total":                 "The number of drops occurred due to lack of WQE for the associated QPs.",
		"out_of_sequence_packets_received_total":    "The number of out of sequence packets received.",
		"packet_sequence_errors_total":              "The number of received NAK sequence error packets. The QP retry limit was not exceeded.",
		"req_cqes_errors_total":                     "The number of times requester detected CQEs completed with errors. The counters was added in MLNX_OFED 4.1",
		"req_cqes_flush_errors_total":               "The number of times requester detected CQEs completed with flushed errors. The counters was added in MLNX_OFED 4.1",
		"req_remote_access_errors_total":            "The number of times requester detected remote access errors. The counters was added in MLNX_OFED 4.1",
		"req_remote_invalid_request_errors_total":   "The number of times requester detected remote invalid request errors. The counters was added in MLNX_OFED 4.1",
		"resp_cqes_errors_total":                    "The number of times responder detected CQEs completed with errors. The counters was added in MLNX_OFED 4.1",
		"resp_cqes_flush_errors_total":              "The number of times responder detected CQEs completed with flushed errors. The counters was added in MLNX_OFED 4.1",
		"resp_local_length_errors_total":            "The number of times responder detected local length errors. The counters was added in MLNX_OFED 4.1",
		"resp_remote_access_errors_total":           "The number of times responder detected remote access errors. The counters was added in MLNX_OFED 4.1",
		"rnr_nak_retry_packets_received_total":      "The number of received RNR NAK packets. The QP retry limit was not exceeded.",
		"roce_adp_retransmits_total":                "The number of adaptive retransmissions for RoCE traffic. The counter was added in MLNX_OFED rev 5.0-1.0.0.0 and kernel v5.6.0",
		"roce_adp_retransmits_timeout_total":        "The number of times RoCE traffic reached timeout due to adaptive retransmission. The counter was added in MLNX_OFED rev 5.0-1.0.0.0 and kernel v5.6.0",
		"roce_slow_restart_used_total":              "The number of times RoCE slow restart was used. The counter was added in MLNX_OFED rev 5.0-1.0.0.0 and kernel v5.6.0",
		"roce_slow_restart_cnps_total":              "The number of times RoCE slow restart generated CNP packets. The counter was added in MLNX_OFED rev 5.0-1.0.0.0 and kernel v5.6.0",
		"roce_slow_restart_total":                   "The number of times RoCE slow restart changed state to slow restart. The counter was added in MLNX_OFED rev 5.0-1.0.0.0 and kernel v5.6.0",
		"rp_cnp_packets_handled_total":              "The number of CNP packets handled by the Reaction Point HCA to throttle the transmission rate. The counters was added in MLNX_OFED 4.1",
		"rp_cnp_ignored_packets_received_total":     "The number of CNP packets received and ignored by the Reaction Point HCA. This counter should not raise if RoCE Congestion Control was enabled in the network. If this counter raise, verify that ECN was enabled on the adapter.",
		"rx_atomic_requests_total":                  "The number of received ATOMIC request for the associated QPs.",
		"rx_dct_connect_requests_total":             "The number of received connection requests for the associated DCTs.",
		"rx_read_requests_total":                    "The number of received READ requests for the associated QPs.",
		"rx_write_requests_total":                   "The number of received WRITE requests for the associated QPs.",
		"rx_icrc_encapsulated_errors_total":         "The number of RoCE packets with ICRC errors. This counter was added in MLNX_OFED 4.4 and kernel 4.19",

		// EFA-specific hw_counters (vendor 0x1d0f). EFA NICs do not follow the
		// IB spec for port_xmit_data / port_rcv_data, so the IB code path leaves
		// port_data_*_bytes_total empty. The EFA branch in Update() fills those
		// from hw_counters/{tx,rx}_bytes and additionally emits the diagnostic
		// counters listed here under the efa_ prefix to avoid clashing with the
		// Mellanox-specific hw_counters above.
		"efa_rx_drops_total":                    "EFA: packets dropped on receive (hw_counters/rx_drops).",
		"efa_retrans_packets_total":             "EFA: retransmitted packets (hw_counters/retrans_pkts).",
		"efa_retrans_bytes_total":               "EFA: retransmitted bytes (hw_counters/retrans_bytes).",
		"efa_retrans_timeout_events_total":      "EFA: retransmit timeout events (hw_counters/retrans_timeout_events).",
		"efa_unresponsive_remote_events_total":  "EFA: unresponsive remote events (hw_counters/unresponsive_remote_events).",
		"efa_impaired_remote_conn_events_total": "EFA: impaired remote connection events (hw_counters/impaired_remote_conn_events).",
		"efa_rdma_read_bytes_total":             "EFA: RDMA read bytes (hw_counters/rdma_read_bytes).",
		"efa_rdma_write_bytes_total":            "EFA: RDMA write bytes (hw_counters/rdma_write_bytes).",
	}

	i.metricDescs = make(map[string]*prometheus.Desc)
	i.subsystem = "infiniband"

	for metricName, description := range descriptions {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, i.subsystem, metricName),
			description,
			[]string{"device", "port"},
			nil,
		)
	}

	return &i, nil
}

func (c *infinibandCollector) pushMetric(ch chan<- prometheus.Metric, name string, value uint64, deviceName string, port string, valueType prometheus.ValueType) {
	ch <- prometheus.MustNewConstMetric(c.metricDescs[name], valueType, float64(value), deviceName, port)
}

func (c *infinibandCollector) pushCounter(ch chan<- prometheus.Metric, name string, value *uint64, deviceName string, port string) {
	if value != nil {
		c.pushMetric(ch, name, *value, deviceName, port, prometheus.CounterValue)
	}
}

// isEFADevice reports whether the InfiniBand-class device is an AWS EFA NIC
// by checking its PCI vendor ID (0x1d0f). EFA NICs register under
// /sys/class/infiniband but do not follow the IB spec for byte/packet
// counters, so they need a different read path (hw_counters/).
func (c *infinibandCollector) isEFADevice(deviceName string) bool {
	path := filepath.Join(c.sysPath, "class", "infiniband", deviceName, "device", "vendor")
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == efaVendorID
}

// readEFAHWCounter reads a single uint64 counter from
// /sys/class/infiniband/<device>/ports/<port>/hw_counters/<counter>.
// Returns nil if the file is missing or unparseable, so pushCounter can skip
// emitting absent series.
func (c *infinibandCollector) readEFAHWCounter(deviceName string, port uint, counter string) *uint64 {
	path := filepath.Join(c.sysPath, "class", "infiniband", deviceName,
		"ports", strconv.FormatUint(uint64(port), 10), "hw_counters", counter)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	v, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		c.logger.Debug("failed to parse EFA hw_counter",
			"path", path, "err", err)
		return nil
	}
	return &v
}

// pushEFACounter is a convenience wrapper that reads a single hw_counter and
// pushes it as a Prometheus counter if present.
func (c *infinibandCollector) pushEFACounter(ch chan<- prometheus.Metric, metricName, counterFile, deviceName string, port uint, portStr string) {
	c.pushCounter(ch, metricName, c.readEFAHWCounter(deviceName, port, counterFile), deviceName, portStr)
}

func (c *infinibandCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.InfiniBandClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("infiniband statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining InfiniBand class info: %w", err)
	}

	for _, device := range devices {
		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/infiniband/<device>, value is always 1.",
			[]string{"device", "board_id", "firmware_version", "hca_type"},
			nil,
		)
		infoValue := 1.0
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, device.Name, device.BoardID, device.FirmwareVersion, device.HCAType)

		// EFA NICs share /sys/class/infiniband layout with IB but use
		// hw_counters/ for byte/packet stats (raw values, no IB ×4 scaling).
		// Detect once per device to avoid stat'ing /sys repeatedly per port.
		isEFA := c.isEFADevice(device.Name)

		for _, port := range device.Ports {
			portStr := strconv.FormatUint(uint64(port.Port), 10)

			c.pushMetric(ch, "state_id", uint64(port.StateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "physical_state_id", uint64(port.PhysStateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "rate_bytes_per_second", port.Rate, port.Name, portStr, prometheus.GaugeValue)

			if isEFA {
				// EFA path: port.Counters (from procfs/sysfs IB-spec parser)
				// is empty/zero, so we read hw_counters/ directly and emit
				// under the existing port_data_* / port_packets_* metric
				// names so existing IB dashboards transparently see EFA data.
				c.pushEFACounter(ch, "port_data_transmitted_bytes_total", "tx_bytes", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "port_data_received_bytes_total", "rx_bytes", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "port_packets_transmitted_total", "tx_pkts", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "port_packets_received_total", "rx_pkts", port.Name, port.Port, portStr)

				// EFA-only diagnostic counters — emitted under efa_* names to
				// avoid colliding with IB-spec semantics. Useful for tracking
				// fabric retransmissions and unresponsive peers.
				c.pushEFACounter(ch, "efa_rx_drops_total", "rx_drops", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_retrans_packets_total", "retrans_pkts", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_retrans_bytes_total", "retrans_bytes", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_retrans_timeout_events_total", "retrans_timeout_events", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_unresponsive_remote_events_total", "unresponsive_remote_events", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_impaired_remote_conn_events_total", "impaired_remote_conn_events", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_rdma_read_bytes_total", "rdma_read_bytes", port.Name, port.Port, portStr)
				c.pushEFACounter(ch, "efa_rdma_write_bytes_total", "rdma_write_bytes", port.Name, port.Port, portStr)
				continue
			}

			c.pushCounter(ch, "legacy_multicast_packets_received_total", port.Counters.LegacyPortMulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_multicast_packets_transmitted_total", port.Counters.LegacyPortMulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_received_bytes_total", port.Counters.LegacyPortRcvData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_received_total", port.Counters.LegacyPortRcvPackets64, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_received_total", port.Counters.LegacyPortUnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_transmitted_total", port.Counters.LegacyPortUnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_transmitted_bytes_total", port.Counters.LegacyPortXmitData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_transmitted_total", port.Counters.LegacyPortXmitPackets64, port.Name, portStr)
			c.pushCounter(ch, "excessive_buffer_overrun_errors_total", port.Counters.ExcessiveBufferOverrunErrors, port.Name, portStr)
			c.pushCounter(ch, "link_downed_total", port.Counters.LinkDowned, port.Name, portStr)
			c.pushCounter(ch, "link_error_recovery_total", port.Counters.LinkErrorRecovery, port.Name, portStr)
			c.pushCounter(ch, "local_link_integrity_errors_total", port.Counters.LocalLinkIntegrityErrors, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_received_total", port.Counters.MulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_transmitted_total", port.Counters.MulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_received_total", port.Counters.PortRcvConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_transmitted_total", port.Counters.PortXmitConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_data_received_bytes_total", port.Counters.PortRcvData, port.Name, portStr)
			c.pushCounter(ch, "port_data_transmitted_bytes_total", port.Counters.PortXmitData, port.Name, portStr)
			c.pushCounter(ch, "port_discards_received_total", port.Counters.PortRcvDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_discards_transmitted_total", port.Counters.PortXmitDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_errors_received_total", port.Counters.PortRcvErrors, port.Name, portStr)
			c.pushCounter(ch, "port_packets_received_total", port.Counters.PortRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "port_packets_transmitted_total", port.Counters.PortXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_transmit_wait_total", port.Counters.PortXmitWait, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_received_total", port.Counters.UnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_transmitted_total", port.Counters.UnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_receive_remote_physical_errors_total", port.Counters.PortRcvRemotePhysicalErrors, port.Name, portStr)
			c.pushCounter(ch, "port_receive_switch_relay_errors_total", port.Counters.PortRcvSwitchRelayErrors, port.Name, portStr)
			c.pushCounter(ch, "symbol_error_total", port.Counters.SymbolError, port.Name, portStr)
			c.pushCounter(ch, "vl15_dropped_total", port.Counters.VL15Dropped, port.Name, portStr)

			// port.HwCounters
			if port.HwCounters.Lifespan != nil {
				c.pushMetric(ch, "lifespan_seconds", *(port.HwCounters.Lifespan)/1000, port.Name, portStr, prometheus.GaugeValue)
			}

			c.pushCounter(ch, "duplicate_requests_packets_total", port.HwCounters.DuplicateRequest, port.Name, portStr)
			c.pushCounter(ch, "implied_nak_seq_errors_total", port.HwCounters.ImpliedNakSeqErr, port.Name, portStr)
			c.pushCounter(ch, "local_ack_timeout_errors_total", port.HwCounters.LocalAckTimeoutErr, port.Name, portStr)
			c.pushCounter(ch, "np_cnp_packets_sent_total", port.HwCounters.NpCnpSent, port.Name, portStr)
			c.pushCounter(ch, "np_ecn_marked_roce_packets_received_total", port.HwCounters.NpEcnMarkedRocePackets, port.Name, portStr)
			c.pushCounter(ch, "out_of_buffer_drops_total", port.HwCounters.OutOfBuffer, port.Name, portStr)
			c.pushCounter(ch, "out_of_sequence_packets_received_total", port.HwCounters.OutOfSequence, port.Name, portStr)
			c.pushCounter(ch, "packet_sequence_errors_total", port.HwCounters.PacketSeqErr, port.Name, portStr)
			c.pushCounter(ch, "req_cqes_errors_total", port.HwCounters.ReqCqeError, port.Name, portStr)
			c.pushCounter(ch, "req_cqes_flush_errors_total", port.HwCounters.ReqCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "req_remote_access_errors_total", port.HwCounters.ReqRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "req_remote_invalid_request_errors_total", port.HwCounters.ReqRemoteInvalidRequest, port.Name, portStr)
			c.pushCounter(ch, "resp_cqes_errors_total", port.HwCounters.RespCqeError, port.Name, portStr)
			c.pushCounter(ch, "resp_cqes_flush_errors_total", port.HwCounters.RespCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "resp_local_length_errors_total", port.HwCounters.RespLocalLengthError, port.Name, portStr)
			c.pushCounter(ch, "resp_remote_access_errors_total", port.HwCounters.RespRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "rnr_nak_retry_packets_received_total", port.HwCounters.RnrNakRetryErr, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retransmits_total", port.HwCounters.RoceAdpRetrans, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retransmits_timeout_total", port.HwCounters.RoceAdpRetransTo, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_used_total", port.HwCounters.RoceSlowRestart, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_cnps_total", port.HwCounters.RoceSlowRestartCnps, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_total", port.HwCounters.RoceSlowRestartTrans, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_packets_handled_total", port.HwCounters.RpCnpHandled, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_ignored_packets_received_total", port.HwCounters.RpCnpIgnored, port.Name, portStr)
			c.pushCounter(ch, "rx_atomic_requests_total", port.HwCounters.RxAtomicRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_dct_connect_requests_total", port.HwCounters.RxDctConnect, port.Name, portStr)
			c.pushCounter(ch, "rx_read_requests_total", port.HwCounters.RxReadRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_write_requests_total", port.HwCounters.RxWriteRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_icrc_encapsulated_errors_total", port.HwCounters.RxIcrcEncapsulated, port.Name, portStr)
		}
	}

	return nil
}
