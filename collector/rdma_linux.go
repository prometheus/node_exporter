// Copyright 2024 The Prometheus Authors
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

//go:build !nordma
// +build !nordma

// The hard work of collecting data from the kernel via the MLNX_OFED interfaces is done by
// https://github.com/Mellanox/rdmamap
// by Mellanox. Used under the Apache 2.0 license.

package collector

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Mellanox/rdmamap"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rdmaDeviceInclude   = kingpin.Flag("collector.rdma.device-include", "Regexp of rdma devices to include (mutually exclusive to device-exclude).").String()
	rdmaDeviceExclude   = kingpin.Flag("collector.rdma.device-exclude", "Regexp of rdma devices to exclude (mutually exclusive to device-include).").String()
	rdmaIncludedMetrics = kingpin.Flag("collector.rdma.metrics-include", "Regexp of rdma stats to include.").Default(".*").String()

	lookupTable = map[string]string{
		"port_rcv_data":                   "mlx5_port_rcv_data_total",
		"port_rcv_packets":                "mlx5_port_rcv_packets_total",
		"port_multicast_rcv_packets":      "mlx5_port_multicast_rcv_packets_total",
		"port_unicast_rcv_packets":        "mlx5_port_unicast_rcv_packets_total",
		"port_xmit_data":                  "mlx5_port_xmit_data_total",
		"port_xmit_packets":               "mlx5_port_xmit_packets_total",
		"port_rcv_switch_relay_errors":    "mlx5_port_rcv_switch_relay_errors_total",
		"port_rcv_errors":                 "mlx5_port_rcv_errors_total",
		"port_rcv_constraint_errors":      "mlx5_port_rcv_constraint_errors_total",
		"local_link_integrity_errors":     "mlx5_local_link_integrity_errors_total",
		"port_xmit_wait":                  "mlx5_port_xmit_wait_total",
		"port_multicast_xmit_packets":     "mlx5_port_multicast_xmit_packets_total",
		"port_unicast_xmit_packets":       "mlx5_port_unicast_xmit_packets_total",
		"port_xmit_discards":              "mlx5_port_xmit_discards_total",
		"port_xmit_constraint_errors":     "mlx5_port_xmit_constraint_errors_total",
		"port_rcv_remote_physical_errors": "mlx5_port_rcv_remote_physical_errors_total",
		"symbol_error":                    "mlx5_symbol_error_total",
		"VL15_dropped":                    "mlx5_vl15_dropped_total",
		"link_error_recovery":             "mlx5_link_error_recovery_total",
		"link_downed":                     "mlx5_link_downed_total",
		"duplicate_request":               "mlx5_duplicate_request_total",
		"implied_nak_seq_err":             "mlx5_implied_nak_seq_err_total",
		"lifespan":                        "mlx5_lifespan_ms",
		"local_ack_timeout_err":           "mlx5_local_ack_timeout_err_total",
		"np_cnp_sent":                     "mlx5_np_cnp_sent_total",
		"np_ecn_marked_roce_packets":      "mlx5_np_ecn_marked_roce_packets_total",
		"out_of_buffer":                   "mlx5_out_of_buffer_total",
		"out_of_sequence":                 "mlx5_out_of_sequence_total",
		"packet_seq_err":                  "mlx5_packet_seq_err_total",
		"req_cqe_error":                   "mlx5_req_cqe_error_total",
		"req_cqe_flush_error":             "mlx5_req_cqe_flush_error_total",
		"req_remote_access_errors":        "mlx5_req_remote_access_errors_total",
		"req_remote_invalid_request":      "mlx5_req_remote_invalid_request_total",
		"resp_cqe_error":                  "mlx5_resp_cqe_error_total",
		"resp_cqe_flush_error":            "mlx5_resp_cqe_flush_error_total",
		"resp_local_length_error":         "mlx5_resp_local_length_error_total",
		"resp_remote_access_errors":       "mlx5_resp_remote_access_errors_total",
		"rnr_nak_retry_err":               "mlx5_rnr_nak_retry_err_total",
		"rp_cnp_handled":                  "mlx5_rp_cnp_handled_total",
		"rp_cnp_ignored":                  "mlx5_rp_cnp_ignored_total",
		"rx_atomic_requests":              "mlx5_rx_atomic_requests_total",
		"rx_dct_connect":                  "mlx5_rx_dct_connect_total",
		"rx_read_requests":                "mlx5_rx_read_requests_total",
		"rx_write_requests":               "mlx5_rx_write_requests_total",
		"rx_icrc_encapsulated":            "mlx5_rx_icrc_encapsulated_total",
		"roce_adp_retrans":                "mlx5_roce_adp_retrans_total",
		"roce_adp_retrans_to":             "mlx5_roce_adp_retrans_to_total",
		"roce_slow_restart":               "mlx5_roce_slow_restart_total",
		"roce_slow_restart_cnps":          "mlx5_roce_slow_restart_cnps_total",
		"roce_slow_restart_trans":         "mlx5_roce_slow_restart_trans_total",
	}

	// https://enterprise-support.nvidia.com/s/article/understanding-mlx5-linux-counters-and-status-parameters
	portCounters = map[string]string{
		"mlx5_port_rcv_data_total":                   "Total number of data octets received on all VLs from the port (divided by 4, counting in double words)",
		"mlx5_port_rcv_packets_total":                "Total number of received packets (may include packets with errors)",
		"mlx5_port_multicast_rcv_packets_total":      "Total number of multicast packets received (including those with errors)",
		"mlx5_port_unicast_rcv_packets_total":        "Total number of unicast packets received (including those with errors)",
		"mlx5_port_xmit_data_total":                  "Total number of data octets transmitted on all VLs from the port (divided by 4, counting in double words)",
		"mlx5_port_xmit_packets_total":               "Total number of transmitted packets (may include packets with errors)",
		"mlx5_port_rcv_switch_relay_errors_total":    "Total number of packets discarded because they could not be forwarded by the switch relay",
		"mlx5_port_rcv_errors_total":                 "Total number of received packets with errors",
		"mlx5_port_rcv_constraint_errors_total":      "Total number of packets discarded due to constraints on the switch physical port",
		"mlx5_local_link_integrity_errors_total":     "Total number of times local physical errors exceeded the threshold and caused a local link integrity failure",
		"mlx5_port_xmit_wait_total":                  "Total number of ticks during which the port had data to transmit but no data was sent due to insufficient credits or lack of arbitration",
		"mlx5_port_multicast_xmit_packets_total":     "Total number of multicast packets transmitted (including those with errors)",
		"mlx5_port_unicast_xmit_packets_total":       "Total number of unicast packets transmitted (including those with errors)",
		"mlx5_port_xmit_discards_total":              "Total number of outbound packets discarded because the port is down or congested",
		"mlx5_port_xmit_constraint_errors_total":     "Total number of packets not transmitted due to constraints on the switch physical port",
		"mlx5_port_rcv_remote_physical_errors_total": "Total number of packets marked with the EBP delimiter received on the port",
		"mlx5_symbol_error_total":                    "Total number of minor link errors detected on one or more physical lanes",
		"mlx5_vl15_dropped_total":                    "Total number of incoming VL15 packets dropped due to resource limitations (e.g., lack of buffers)",
		"mlx5_link_error_recovery_total":             "Total number of successful link error recovery processes completed by the Port Training state machine",
		"mlx5_link_downed_total":                     "Total number of failed link error recovery processes that caused the link to be downed",
	}

	hwCounters = map[string]string{
		"mlx5_duplicate_request_total":          "Total number of received packets that were duplicates of previous requests",
		"mlx5_implied_nak_seq_err_total":        "Total number of times the requested ACK had a PSN larger than the expected PSN for an RDMA read or response",
		"mlx5_lifespan_ms":                      "Maximum period in milliseconds which defines the aging of counter reads",
		"mlx5_local_ack_timeout_err_total":      "Total number of times the QP's ACK timer expired for RC, XRC, or DCT QPs at the sender side (retry limit not exceeded)",
		"mlx5_np_cnp_sent_total":                "Total number of CNP packets sent by the Notification Point due to congestion in the RoCEv2 IP header (ECN bits)",
		"mlx5_np_ecn_marked_roce_packets_total": "Total number of RoCEv2 packets received marked with ECN (congestion experienced)",
		"mlx5_out_of_buffer_total":              "Total number of drops due to lack of WQE for the associated QPs",
		"mlx5_out_of_sequence_total":            "Total number of out-of-sequence packets received",
		"mlx5_packet_seq_err_total":             "Total number of received NAK sequence error packets (QP retry limit not exceeded)",
		"mlx5_req_cqe_error_total":              "Total number of times the requester detected CQEs completed with errors",
		"mlx5_req_cqe_flush_error_total":        "Total number of times the requester detected CQEs completed with flushed errors",
		"mlx5_req_remote_access_errors_total":   "Total number of times the requester detected remote access errors",
		"mlx5_req_remote_invalid_request_total": "Total number of times the requester detected remote invalid request errors",
		"mlx5_resp_cqe_error_total":             "Total number of times the responder detected CQEs completed with errors",
		"mlx5_resp_cqe_flush_error_total":       "Total number of times the responder detected CQEs completed with flushed errors",
		"mlx5_resp_local_length_error_total":    "Total number of times the responder detected local length errors",
		"mlx5_resp_remote_access_errors_total":  "Total number of times the responder detected remote access errors",
		"mlx5_rnr_nak_retry_err_total":          "Total number of received RNR NAK packets (QP retry limit not exceeded)",
		"mlx5_rp_cnp_handled_total":             "Total number of CNP packets handled by the Reaction Point HCA to throttle transmission rate",
		"mlx5_rp_cnp_ignored_total":             "Total number of CNP packets ignored by the Reaction Point HCA",
		"mlx5_rx_atomic_requests_total":         "Total number of received ATOMIC requests for associated QPs",
		"mlx5_rx_dct_connect_total":             "Total number of received connection requests for associated DCTs",
		"mlx5_rx_read_requests_total":           "Total number of received READ requests for associated QPs",
		"mlx5_rx_write_requests_total":          "Total number of received WRITE requests for associated QPs",
		"mlx5_rx_icrc_encapsulated_total":       "Total number of RoCE packets with ICRC errors",
		"mlx5_roce_adp_retrans_total":           "Total number of adaptive retransmissions for RoCE traffic",
		"mlx5_roce_adp_retrans_to_total":        "Total number of times RoCE traffic reached timeout due to adaptive retransmission",
		"mlx5_roce_slow_restart_total":          "Total number of times RoCE slow restart was used",
		"mlx5_roce_slow_restart_cnps_total":     "Total number of times RoCE slow restart generated CNP packets",
		"mlx5_roce_slow_restart_trans_total":    "Total number of times RoCE slow restart changed state to slow restart",
	}
)

type rdmaCollector struct {
	entries        map[string]*prometheus.Desc
	entriesMutex   sync.Mutex
	deviceFilter   deviceFilter
	infoDesc       *prometheus.Desc
	metricsPattern *regexp.Regexp
	logger         *slog.Logger
}

// makeRdmaCollector is the internal constructor for rdmaCollector.
func makeRdmaCollector(logger *slog.Logger) (*rdmaCollector, error) {
	if *rdmaDeviceInclude != "" {
		logger.Info("Parsed flag --collector.rdma.device-include", "flag", *rdmaDeviceInclude)
	}
	if *rdmaDeviceExclude != "" {
		logger.Info("Parsed flag --collector.rdma.device-exclude", "flag", *rdmaDeviceExclude)
	}
	if *rdmaIncludedMetrics != "" {
		logger.Info("Parsed flag --collector.rdma.metrics-include", "flag", *rdmaIncludedMetrics)
	}

	// Update paths to respect the mount points setup.
	for _, dir := range []*string{
		&rdmamap.RdmaClassDir,
		&rdmamap.RdmaIbUcmDir,
		&rdmamap.RdmaUmadDir,
		&rdmamap.RdmaUverbsDir,
		&rdmamap.PciDevDir,
		&rdmamap.AuxDevDir,
	} {
		*dir = strings.TrimPrefix(*dir, "/sys")
		*dir = sysFilePath(*dir)
	}
	for _, dir := range []*string{
		&rdmamap.RdmaUcmDevice,
		&rdmamap.RdmaDeviceDir,
	} {
		*dir = rootfsFilePath(*dir)
	}

	entries := make(map[string]*prometheus.Desc, len(portCounters)+len(hwCounters))
	for _, counters := range []map[string]string{portCounters, hwCounters} {
		for metric, help := range counters {
			entries[metric] = prometheus.NewDesc(
				buildRdmaFQName(metric),
				help,
				[]string{"device", "port", "interfaces"}, nil,
			)
		}
	}

	// Pre-populate some common rdma metrics.
	return &rdmaCollector{
		deviceFilter:   newDeviceFilter(*rdmaDeviceExclude, *rdmaDeviceInclude),
		metricsPattern: regexp.MustCompile(*rdmaIncludedMetrics),
		logger:         logger,
		entries:        entries,
		infoDesc: prometheus.NewDesc(
			buildRdmaFQName("info"),
			"A metric with a constant '1' value labeled by device, vendor_id, device_id, firmware_version, driver_version.",
			[]string{"device", "vendor_id", "device_id", "firmware_version", "driver_version"}, nil,
		),
	}, nil
}

func init() {
	registerCollector("rdma", defaultDisabled, NewRdmaCollector)
}

// Generate the fully-qualified metric name for the rdma metric.
func buildRdmaFQName(metric string) string {
	metricName := strings.TrimLeft(strings.ToLower(SanitizeMetricName(metric)), "_")
	return prometheus.BuildFQName(namespace, "rdma", metricName)
}

// NewRdmaCollector returns a new Collector exposing rdma stats.
func NewRdmaCollector(logger *slog.Logger) (Collector, error) {
	return makeRdmaCollector(logger)
}

func getNetworkInterfaces(rdmaDeviceName string) string {
	var ifs []string

	dir := filepath.Join(rdmamap.RdmaClassDir, rdmaDeviceName, "device", "net")
	fd, err := os.Open(dir)
	if err != nil {
		return ""
	}
	defer fd.Close()

	fileInfos, err := fd.Readdir(-1)
	if err != nil {
		return ""
	}

	for i := range fileInfos {
		if fileInfos[i].Name() == "." || fileInfos[i].Name() == ".." {
			continue
		}
		ifs = append(ifs, fileInfos[i].Name())
	}
	return strings.Join(ifs, ",")
}

func (c *rdmaCollector) Update(ch chan<- prometheus.Metric) error {
	rdmaDevices := rdmamap.GetRdmaDeviceList()
	if len(rdmaDevices) == 0 {
		return fmt.Errorf("no rdma devices found")
	}

	for _, device := range rdmaDevices {
		if c.deviceFilter.ignored(device) {
			continue
		}

		interfaces := getNetworkInterfaces(device)

		stats, err := rdmamap.GetRdmaSysfsAllPortsStats(device)
		if err != nil {
			c.logger.Error("rdma stats error", "err", err, "device", device)
			continue
		}

		updateFunc := func(key string, value float64, labelValues ...string) {
			metric, ok := lookupTable[key]
			if !ok {
				c.logger.Warn("rdma metric not found in lookup table", "key", key)
				return
			}
			if !c.metricsPattern.MatchString(metric) {
				c.logger.Debug("rdma metric excluded", "metric", metric)
				return
			}
			entry := c.entry(metric)
			if entry == nil {
				c.logger.Warn("rdma metric not found", "metric", metric)
				return
			}
			ch <- prometheus.MustNewConstMetric(entry, prometheus.GaugeValue,
				value, labelValues...)
		}

		for _, portstats := range stats.PortStats {
			for _, stat := range append(portstats.HwStats, portstats.Stats...) {
				updateFunc(stat.Name, float64(stat.Value), device, fmt.Sprintf("%d", portstats.Port), interfaces)
			}
		}

		vendorID := readStringFromFile(filepath.Join(rdmamap.RdmaClassDir, device, "device", "vendor"))
		deviceID := readStringFromFile(filepath.Join(rdmamap.RdmaClassDir, device, "device", "device"))
		firmwareVersion := readStringFromFile(filepath.Join(rdmamap.RdmaClassDir, "mlx5_0", "fw_ver"))
		driverVersion := readStringFromFile(sysFilePath("module/mlx5_core/version"))
		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0,
			device, vendorID, deviceID, firmwareVersion, driverVersion)
	}

	return nil
}

func (c *rdmaCollector) entry(key string) *prometheus.Desc {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	return c.entries[key]
}
