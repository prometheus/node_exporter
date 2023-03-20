// Copyright 2017 The Prometheus Authors
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

//go:build !noarp
// +build !noarp

package collector

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	arpDeviceInclude = kingpin.Flag("collector.arp.device-include", "Regexp of arp devices to include (mutually exclusive to device-exclude).").String()
	arpDeviceExclude = kingpin.Flag("collector.arp.device-exclude", "Regexp of arp devices to exclude (mutually exclusive to device-include).").String()
)

type arpCollector struct {
	fs           procfs.FS
	deviceFilter deviceFilter
	entries      *prometheus.Desc
	logger       log.Logger
}

func init() {
	registerCollector("arp", defaultEnabled, NewARPCollector)
}

// NewARPCollector returns a new Collector exposing ARP stats.
func NewARPCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &arpCollector{
		fs:           fs,
		deviceFilter: newDeviceFilter(*arpDeviceExclude, *arpDeviceInclude),
		entries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "arp", "entries"),
			"ARP entries by device",
			[]string{"device"}, nil,
		),
		logger: logger,
	}, nil
}

func getTotalArpEntries(deviceEntries []procfs.ARPEntry) map[string]uint32 {
	entries := make(map[string]uint32)

	for _, device := range deviceEntries {
		entries[device.Device]++
	}

	return entries
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := c.fs.GatherARPEntries()
	if err != nil {
		return fmt.Errorf("could not get ARP entries: %w", err)
	}

	enumeratedEntry := getTotalArpEntries(entries)

	for device, entryCount := range enumeratedEntry {
		if c.deviceFilter.ignored(device) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.entries, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
