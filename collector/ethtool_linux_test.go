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
	"golang.org/x/sys/unix"
)

type EthtoolFixture struct {
	fixturePath string
}

func (e *EthtoolFixture) Stats(intf string) (map[string]uint64, error) {
	res := make(map[string]uint64)

	fixtureFile, err := os.Open(filepath.Join(e.fixturePath, intf))
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
	collector.stats = &EthtoolFixture{
		fixturePath: "fixtures/ethtool/",
	}
	if err != nil {
		return nil, err
	}
	return collector, nil
}

func TestEthtoolCollector(t *testing.T) {
	testcases := []string{
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
