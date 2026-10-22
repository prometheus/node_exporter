// Copyright 2026 The Prometheus Authors
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

//go:build !nondisc

package collector

import (
	"fmt"
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jsimonetti/rtnetlink/v2/rtnl"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

var (
	ndiscDeviceInclude = kingpin.Flag("collector.ndisc.device-include", "Regexp of ndisc devices to include (mutually exclusive to device-exclude).").String()
	ndiscDeviceExclude = kingpin.Flag("collector.ndisc.device-exclude", "Regexp of ndisc devices to exclude (mutually exclusive to device-include).").String()
)

type ndiscCollector struct {
	deviceFilter deviceFilter
	logger       *slog.Logger
}

func init() {
	registerCollector("ndisc", defaultEnabled, NewNdiscCollector)
}

var (
	ndiscEntries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "ndisc", "entries"),
		"NDISC entries by device",
		[]string{"device"}, nil,
	)
)

// NewNdiscCollector returns a new Collector exposing NDISC stats.
func NewNdiscCollector(logger *slog.Logger) (Collector, error) {
	return &ndiscCollector{
		deviceFilter: newDeviceFilter(*ndiscDeviceExclude, *ndiscDeviceInclude),
		logger:       logger,
	}, nil
}

func getTotalNdiscEntries(neighbors []*rtnl.Neigh) map[string]uint32 {
	entries := make(map[string]uint32)

	for _, n := range neighbors {
		if n.State&unix.NUD_NOARP == 0 && n.Interface != nil {
			entries[n.Interface.Name]++
		}
	}

	return entries
}

func getTotalNdiscEntriesRTNL() (map[string]uint32, error) {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	neighbors, err := conn.Neighbours(nil, unix.AF_INET6)
	if err != nil {
		return nil, err
	}

	return getTotalNdiscEntries(neighbors), nil
}

func (c *ndiscCollector) Update(ch chan<- prometheus.Metric) error {
	enumeratedEntries, err := getTotalNdiscEntriesRTNL()
	if err != nil {
		return fmt.Errorf("could not get NDISC entries: %w", err)
	}

	for device, entryCount := range enumeratedEntries {
		if c.deviceFilter.ignored(device) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			ndiscEntries, prometheus.GaugeValue, float64(entryCount), device,
		)
	}

	return nil
}
