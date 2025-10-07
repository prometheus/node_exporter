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
// +build !noinfiniband

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type infinibandCollector struct {
	fs          sysfs.FS
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
		"rx_write_requests":                          "Number of write Requests from hwcounters.",
		"rx_read_requests":                           "Number of read Requests from hwcounters.",
		"rx_icrc_encapsulated":                       "Number of RxIcrcEncapsulated packets from hwcounters",
		"rx_dct_connect":                             "Number of DCT connect requests received from hwcounters",
		"rx_atomic_requests":                         "Number of atomic requests received from hwcounters",
		"rp_cnp_ignored":                             "Number of CNP packets ignored by rate limiter from hwcounters",
		"rp_cnp_handled":                             "Number of CNP packets handled by rate limiter from hwcounters",
		"roce_slow_restart":                          "Number of RoCE slow restart events from hwcounters",
		"roce_slow_restart_cnps":                     "Number of RoCE slow restart CNP packets from hwcounters",
		"roce_slow_restart_trans":                    "Number of RoCE slow restart transmissions from hwcounters",
		"roce_adp_retrans_to":                        "Number of RoCE adaptive retransmission timeouts from hwcounters",
		"roce_adp_retrans":                           "Number of RoCE adaptive retransmissions from hwcounters",
		"duplicate_request":                          "Number of duplicate requests from hwcounters",
		"implied_nak_seq_err":                        "Number of implied NAK sequence errors from hwcounters",
		"lifespan":                                   "Packet lifespan counter from hwcounters",
		"local_ack_timeout_err":                      "Number of local ACK timeout errors from hwcounters",
		"np_cnp_sent":                                "Number of CNP packets sent by notification point from hwcounters",
		"np_ecn_marked_roce_packets":                 "Number of ECN marked RoCE packets from hwcounters",
		"out_of_buffer":                              "Number of out of buffer events from hwcounters",
		"out_of_sequence":                            "Number of out of sequence packets from hwcounters",
		"packet_seq_err":                             "Number of packet sequence errors from hwcounters",
		"req_cqe_error":                              "Number of request completion queue errors from hwcounters",
		"req_cqe_flush_error":                        "Number of request completion queue flush errors from hwcounters",
		"req_remote_access_errors":                   "Number of request remote access errors from hwcounters",
		"req_remote_invalid_request":                 "Number of request remote invalid requests from hwcounters",
		"resp_cqe_error":                             "Number of response completion queue errors from hwcounters",
		"resp_cqe_flush_error":                       "Number of response completion queue flush errors from hwcounters",
		"resp_local_length_error":                    "Number of response local length errors from hwcounters",
		"resp_remote_access_errors":                  "Number of response remote access errors from hwcounters",
		"rnr_nak_retry_err":                          "Number of RNR NAK retry errors from hwcounters",
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

		for _, port := range device.Ports {
			portStr := strconv.FormatUint(uint64(port.Port), 10)

			c.pushMetric(ch, "state_id", uint64(port.StateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "physical_state_id", uint64(port.PhysStateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "rate_bytes_per_second", port.Rate, port.Name, portStr, prometheus.GaugeValue)

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
			c.pushCounter(ch, "duplicate_request", port.HwCounters.DuplicateRequest, port.Name, portStr)                 
			c.pushCounter(ch, "implied_nak_seq_err", port.HwCounters.ImpliedNakSeqErr, port.Name, portStr)
			c.pushCounter(ch, "lifespan", port.HwCounters.Lifespan, port.Name, portStr)
			c.pushCounter(ch, "local_ack_timeout_err", port.HwCounters.LocalAckTimeoutErr, port.Name, portStr)
			c.pushCounter(ch, "np_cnp_sent", port.HwCounters.NpCnpSent, port.Name, portStr)
			c.pushCounter(ch, "np_ecn_marked_roce_packets", port.HwCounters.NpEcnMarkedRocePackets, port.Name, portStr)
			c.pushCounter(ch, "out_of_buffer", port.HwCounters.OutOfBuffer, port.Name, portStr)
			c.pushCounter(ch, "out_of_sequence", port.HwCounters.OutOfSequence, port.Name, portStr)
			c.pushCounter(ch, "packet_seq_err", port.HwCounters.PacketSeqErr, port.Name, portStr)
			c.pushCounter(ch, "req_cqe_error", port.HwCounters.ReqCqeError, port.Name, portStr)
			c.pushCounter(ch, "req_cqe_flush_error", port.HwCounters.ReqCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "req_remote_access_errors", port.HwCounters.ReqRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "req_remote_invalid_request", port.HwCounters.ReqRemoteInvalidRequest, port.Name, portStr)
			c.pushCounter(ch, "resp_cqe_error", port.HwCounters.RespCqeError, port.Name, portStr)
			c.pushCounter(ch, "resp_cqe_flush_error", port.HwCounters.RespCqeFlushError, port.Name, portStr)
			c.pushCounter(ch, "resp_local_length_error", port.HwCounters.RespLocalLengthError, port.Name, portStr)
			c.pushCounter(ch, "resp_remote_access_errors", port.HwCounters.RespRemoteAccessErrors, port.Name, portStr)
			c.pushCounter(ch, "rnr_nak_retry_err", port.HwCounters.RnrNakRetryErr, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retrans", port.HwCounters.RoceAdpRetrans, port.Name, portStr)
			c.pushCounter(ch, "roce_adp_retrans_to", port.HwCounters.RoceAdpRetransTo, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart", port.HwCounters.RoceSlowRestart, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_cnps", port.HwCounters.RoceSlowRestartCnps, port.Name, portStr)
			c.pushCounter(ch, "roce_slow_restart_trans", port.HwCounters.RoceSlowRestartTrans, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_handled", port.HwCounters.RpCnpHandled, port.Name, portStr)
			c.pushCounter(ch, "rp_cnp_ignored", port.HwCounters.RpCnpIgnored, port.Name, portStr)
			c.pushCounter(ch, "rx_atomic_requests", port.HwCounters.RxAtomicRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_dct_connect", port.HwCounters.RxDctConnect, port.Name, portStr)
			c.pushCounter(ch, "rx_icrc_encapsulated", port.HwCounters.RxIcrcEncapsulated, port.Name, portStr)
			c.pushCounter(ch, "rx_write_requests", port.HwCounters.RxWriteRequests, port.Name, portStr)
			c.pushCounter(ch, "rx_read_requests", port.HwCounters.RxReadRequests, port.Name, portStr)
		}
	}

	return nil
}
