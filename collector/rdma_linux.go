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

	rdmaHwCounters = map[string]string{
		"roce_slow_restart_cnps":     "RDMA RoCE slow restart CNPS",
		"rp_cnp_ignored":             "RDMA RP CNP ignored",
		"roce_adp_retrans_to":        "RDMA RoCE adaptive retransmission timeout",
		"rx_icrc_encapsulated":       "RDMA RX ICRC encapsulated",
		"resp_local_length_error":    "RDMA response local length error",
		"np_ecn_marked_roce_packets": "RDMA NP ECN marked RoCE packets",
		"roce_slow_restart_trans":    "RDMA RoCE slow restart transactions",
		"req_remote_invalid_request": "RDMA request remote invalid request",
		"local_ack_timeout_err":      "RDMA local ACK timeout error",
		"lifespan":                   "RDMA lifespan",
		"req_cqe_error":              "RDMA request CQE error",
		"rnr_nak_retry_err":          "RDMA RNR NAK retry error",
		"np_cnp_sent":                "RDMA NP CNP sent",
		"rx_dct_connect":             "RDMA RX DCT connect",
		"rp_cnp_handled":             "RDMA RP CNP handled",
		"implied_nak_seq_err":        "RDMA implied NAK sequence error",
		"roce_slow_restart":          "RDMA RoCE slow restart",
		"req_cqe_flush_error":        "RDMA request CQE flush error",
		"packet_seq_err":             "RDMA packet sequence error",
		"duplicate_request":          "RDMA duplicate request",
		"roce_adp_retrans":           "RDMA RoCE adaptive retransmission",
		"out_of_buffer":              "RDMA out of buffer",
		"resp_cqe_error":             "RDMA response CQE error",
		"resp_cqe_flush_error":       "RDMA response CQE flush error",
		"out_of_sequence":            "RDMA out of sequence",
		"rx_read_requests":           "RDMA RX read requests",
		"rx_atomic_requests":         "RDMA RX atomic requests",
		"req_remote_access_errors":   "RDMA request remote access errors",
		"rx_write_requests":          "RDMA RX write requests",
		"resp_remote_access_errors":  "RDMA response remote access errors",
	}
	rdmaCounters = map[string]string{
		"unicast_rcv_packets":             "RDMA unicast received packets",
		"port_xmit_data":                  "RDMA port transmit data",
		"port_xmit_constraint_errors":     "RDMA port transmit constraint errors",
		"VL15_dropped":                    "RDMA VL15 dropped",
		"port_rcv_errors":                 "RDMA port receive errors",
		"port_xmit_wait":                  "RDMA port transmit wait",
		"link_error_recovery":             "RDMA link error recovery",
		"multicast_rcv_packets":           "RDMA multicast received packets",
		"multicast_xmit_packets":          "RDMA multicast transmitted packets",
		"port_rcv_remote_physical_errors": "RDMA port receive remote physical errors",
		"port_rcv_packets":                "RDMA port receive packets",
		"unicast_xmit_packets":            "RDMA unicast transmitted packets",
		"excessive_buffer_overrun_errors": "RDMA excessive buffer overrun errors",
		"port_rcv_data":                   "RDMA port receive data",
		"port_rcv_constraint_errors":      "RDMA port receive constraint errors",
		"link_downed":                     "RDMA link downed",
		"local_link_integrity_errors":     "RDMA local link integrity errors",
		"port_xmit_discards":              "RDMA port transmit discards",
		"port_rcv_switch_relay_errors":    "RDMA port receive switch relay errors",
		"port_xmit_packets":               "RDMA port transmit packets",
		"symbol_error":                    "RDMA symbol error",
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

	entries := make(map[string]*prometheus.Desc, len(rdmaHwCounters)+len(rdmaCounters))
	for metric, help := range rdmaHwCounters {
		entries[metric] = prometheus.NewDesc(
			buildRdmaFQName(fmt.Sprintf("hw_%s", metric)),
			help,
			[]string{"device", "port", "interfaces"}, nil,
		)
	}
	for metric, help := range rdmaCounters {
		entries[metric] = prometheus.NewDesc(
			buildRdmaFQName(metric),
			help,
			[]string{"device", "port", "interfaces"}, nil,
		)
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

		updateFunc := func(name string, value float64, labelValues ...string) {
			if !c.metricsPattern.MatchString(name) {
				return
			}
			ch <- prometheus.MustNewConstMetric(c.entry(name), prometheus.GaugeValue,
				value, labelValues...)
		}

		for _, portstats := range stats.PortStats {
			for _, stat := range portstats.HwStats {
				updateFunc(stat.Name, float64(stat.Value), device, fmt.Sprintf("%d", portstats.Port), interfaces)
			}
			for _, stat := range portstats.Stats {
				updateFunc(stat.Name, float64(stat.Value), device, fmt.Sprintf("%d", portstats.Port), interfaces)
			}
		}

		vendorID := readStringFromFile(filepath.Join(rdmamap.RdmaClassDir, device, "device", "vendor"))
		deviceID := readStringFromFile(filepath.Join(rdmamap.RdmaClassDir, device, "device", "device"))
		firmwareVersion := readStringFromFile("/sys/class/infiniband/mlx5_0/fw_ver")
		driverVersion := readStringFromFile("/sys/module/mlx5_core/version")
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
