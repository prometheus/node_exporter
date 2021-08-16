// Copyright 2021 The Prometheus Authors
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

package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/safchain/ethtool"
	"golang.org/x/sys/unix"
)

type EthtoolFixture struct {
	fixturePath string
}

func (e *EthtoolFixture) DriverInfo(intf string) (ethtool.DrvInfo, error) {
	res := ethtool.DrvInfo{}

	fixtureFile, err := os.Open(filepath.Join(e.fixturePath, intf, "driver"))
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
		// The fixture for this interface doesn't exist. Translate that to unix.EOPNOTSUPP
		// to replicate an interface that doesn't support ethtool driver info
		return res, unix.EOPNOTSUPP
	}
	if err != nil {
		return res, err
	}
	defer fixtureFile.Close()

	scanner := bufio.NewScanner(fixtureFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.Trim(line, " ")
		items := strings.Split(line, ": ")
		switch items[0] {
		case "driver":
			res.Driver = items[1]
		case "version":
			res.Version = items[1]
		case "firmware-version":
			res.FwVersion = items[1]
		case "bus-info":
			res.BusInfo = items[1]
		case "expansion-rom-version":
			res.EromVersion = items[1]
		}
	}

	return res, err
}

func (e *EthtoolFixture) Stats(intf string) (map[string]uint64, error) {
	res := make(map[string]uint64)

	fixtureFile, err := os.Open(filepath.Join(e.fixturePath, intf, "statistics"))
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
		// The fixture for this interface doesn't exist. Translate that to unix.EOPNOTSUPP
		// to replicate an interface that doesn't support ethtool stats
		return res, unix.EOPNOTSUPP
	}
	if err != nil {
		return res, err
	}
	defer fixtureFile.Close()

	scanner := bufio.NewScanner(fixtureFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "NIC statistics:") {
			continue
		}
		line = strings.Trim(line, " ")
		items := strings.Split(line, ": ")
		val, err := strconv.ParseUint(items[1], 10, 64)
		if err != nil {
			return res, err
		}
		if items[0] == "ERROR" {
			return res, unix.Errno(val)
		}
		res[items[0]] = val
	}

	return res, err
}

func NewEthtoolTestCollector(logger log.Logger) (Collector, error) {
	collector, err := makeEthtoolCollector(logger)
	collector.ethtool = &EthtoolFixture{
		fixturePath: "fixtures/ethtool/",
	}
	if err != nil {
		return nil, err
	}
	return collector, nil
}

func TestSanitizeMetricName(t *testing.T) {
	testcases := map[string]string{
		"":                             "",
		"rx_errors":                    "rx_errors",
		"Queue[0] AllocFails":          "Queue_0_AllocFails",
		"Tx LPI entry count":           "Tx_LPI_entry_count",
		"port.VF_admin_queue_requests": "port_VF_admin_queue_requests",
		"[3]: tx_bytes":                "_3_tx_bytes",
	}

	for metricName, expected := range testcases {
		got := SanitizeMetricName(metricName)
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		}
	}
}

func TestBuildEthtoolFQName(t *testing.T) {
	testcases := map[string]string{
		"rx_errors":                    "node_ethtool_received_errors",
		"Queue[0] AllocFails":          "node_ethtool_queue_0_allocfails",
		"Tx LPI entry count":           "node_ethtool_transmitted_lpi_entry_count",
		"port.VF_admin_queue_requests": "node_ethtool_port_vf_admin_queue_requests",
		"[3]: tx_bytes":                "node_ethtool_3_transmitted_bytes",
	}

	for metric, expected := range testcases {
		got := buildEthtoolFQName(metric)
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		}
	}
}

func TestEthtoolCollector(t *testing.T) {
	testcases := []string{
		prometheus.NewDesc("node_ethtool_info",
			"A metric with a constant '1' value labeled by bus_info, device, driver, expansion_rom_version, firmware_version, version.",
			[]string{"bus_info", "device", "driver", "expansion_rom_version", "firmware_version", "version"}, nil).String(),
		prometheus.NewDesc("node_ethtool_align_errors", "Network interface align_errors", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_broadcast", "Network interface rx_broadcast", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_errors_total", "Number of received frames with errors", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_missed", "Network interface rx_missed", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_multicast", "Network interface rx_multicast", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_packets_total", "Network interface packets received", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_received_unicast", "Network interface rx_unicast", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_transmitted_aborted", "Network interface tx_aborted", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_transmitted_errors_total", "Number of sent frames with errors", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_transmitted_multi_collisions", "Network interface tx_multi_collisions", []string{"device"}, nil).String(),
		prometheus.NewDesc("node_ethtool_transmitted_packets_total", "Network interface packets sent", []string{"device"}, nil).String(),
	}

	*ethtoolIgnoredDevices = "^$"
	*sysPath = "fixtures/sys"

	collector, err := NewEthtoolTestCollector(log.NewNopLogger())
	if err != nil {
		panic(err)
	}

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for _, expected := range testcases {
		metric := (<-sink)
		if metric == nil {
			t.Fatalf("Expected '%s' but got nothing (nil).", expected)
		}

		got := metric.Desc().String()
		metric.Desc()
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		} else {
			t.Logf("Successfully got '%s'", got)
		}
	}
}
