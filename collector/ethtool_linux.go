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

// +build !noethtool

// The hard work of collecting data from the kernel via the ethtool interfaces is done by
// https://github.com/safchain/ethtool/
// by Sylvain Afchain. Used under the Apache license.

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"github.com/safchain/ethtool"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	receivedRegex    = regexp.MustCompile(`_rx_`)
	transmittedRegex = regexp.MustCompile(`_tx_`)
	ethtoolFixtures  = kingpin.Flag("collector.ethtool.fixtures", "test fixtures to use for ethtool collector end-to-end testing").Default("").String()
)

type EthtoolStats interface {
	Stats(string) (map[string]uint64, error)
}

type ethtoolStats struct {
}

func (e *ethtoolStats) Stats(intf string) (map[string]uint64, error) {
	return ethtool.Stats(intf)
}

type ethtoolCollector struct {
	fs      sysfs.FS
	entries map[string]*prometheus.Desc
	logger  log.Logger
	stats   EthtoolStats
}

type EthtoolFixture struct {
	fixturePath string
}

func (e *EthtoolFixture) Stats(intf string) (map[string]uint64, error) {
	res := make(map[string]uint64)

	fixtureFile, err := os.Open(filepath.Join(e.fixturePath, intf))
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
		// The fixture for this interface doesn't exist. That's OK because it replicates
		// an interface that doesn't support ethtool.
		return res, nil
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
		res[items[0]] = val
	}

	return res, err
}

func NewEthtoolTestCollector(logger log.Logger) (Collector, error) {
	collector, err := makeEthtoolCollector(logger)
	collector.stats = &EthtoolFixture{
		fixturePath: *ethtoolFixtures,
	}
	if err != nil {
		return nil, err
	}
	return collector, nil
}

// makeEthtoolCollector is the internal constructor for EthtoolCollector.
// This allows NewEthtoolTestCollector to override it's .stats interface
// for testing.
func makeEthtoolCollector(logger log.Logger) (*ethtoolCollector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	// Pre-populate some common ethtool metrics.
	return &ethtoolCollector{
		fs:    fs,
		stats: &ethtoolStats{},
		entries: map[string]*prometheus.Desc{
			"rx_bytes": prometheus.NewDesc(
				"node_ethtool_received_bytes_total",
				"Network interface bytes received",
				[]string{"device"}, nil,
			),
			"rx_dropped": prometheus.NewDesc(
				"node_ethtool_received_dropped_total",
				"Number of received frames dropped",
				[]string{"device"}, nil,
			),
			"rx_errors": prometheus.NewDesc(
				"node_ethtool_received_errors_total",
				"Number of received frames with errors",
				[]string{"device"}, nil,
			),
			"rx_packets": prometheus.NewDesc(
				"node_ethtool_received_packets_total",
				"Network interface packets received",
				[]string{"device"}, nil,
			),
			"tx_bytes": prometheus.NewDesc(
				"node_ethtool_transmitted_bytes_total",
				"Network interface bytes sent",
				[]string{"device"}, nil,
			),
			"tx_errors": prometheus.NewDesc(
				"node_ethtool_transmitted_errors_total",
				"Number of sent frames with errors",
				[]string{"device"}, nil,
			),
			"tx_packets": prometheus.NewDesc(
				"node_ethtool_transmitted_packets_total",
				"Network interface packets sent",
				[]string{"device"}, nil,
			),
		},
		logger: logger,
	}, nil
}

func init() {
	registerCollector("ethtool", defaultDisabled, NewEthtoolCollector)
}

// NewEthtoolCollector returns a new Collector exposing ethtool stats.
func NewEthtoolCollector(logger log.Logger) (Collector, error) {
	// Specifying --collector.ethtool.fixtures on the command line activates
	// the test fixtures. This is for `end-to-end-test.sh`
	if *ethtoolFixtures != "" {
		return NewEthtoolTestCollector(logger)
	}
	return makeEthtoolCollector(logger)
}

func (c *ethtoolCollector) Update(ch chan<- prometheus.Metric) error {
	netClass, err := c.fs.NetClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			level.Debug(c.logger).Log("msg", "Could not read netclass file", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get net class info: %w", err)
	}

	if len(netClass) == 0 {
		return fmt.Errorf("no network devices found")
	}

	for device := range netClass {
		var stats map[string]uint64
		var err error

		stats, err = c.stats.Stats(device)
		if err != nil {
			// Suppressing errors because it's hard to tell what interfaces support ethtool and which don't.
			continue
		}

		// Sort metric names so that the test fixtures will match up
		keys := make([]string, 0, len(stats))
		for k := range stats {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, metric := range keys {
			val := stats[metric]
			metricFQName := prometheus.BuildFQName(namespace, "ethtool", metric)
			metricFQName = receivedRegex.ReplaceAllString(metricFQName, "_received_")
			metricFQName = transmittedRegex.ReplaceAllString(metricFQName, "_transmitted_")

			// Check to see if this metric exists; if not then create it and store it in c.entries.
			entry, exists := c.entries[metric]
			if !exists {
				entry = prometheus.NewDesc(
					metricFQName,
					fmt.Sprintf("Network interface %s", metric),
					[]string{"device"}, nil,
				)
				c.entries[metric] = entry
			}
			ch <- prometheus.MustNewConstMetric(
				entry, prometheus.UntypedValue, float64(val), device)
		}
	}

	return nil
}
