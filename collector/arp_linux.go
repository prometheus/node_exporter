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

package collector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jsimonetti/rtnetlink/v2/rtnl"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"

	"github.com/prometheus/node_exporter/config"
)

type arpCollector struct {
	fs           procfs.FS
	deviceFilter deviceFilter
	logger       *slog.Logger
	config       config.ARPConfig
}

func init() {
	registerCollector("arp", defaultEnabled, NewARPCollector)
}

var arpEntries = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "arp", "entries"),
	"ARP entries by device",
	[]string{"device"}, nil,
)

// NewARPCollector returns a new Collector exposing ARP stats.
func NewARPCollector(cfg config.Config, logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(cfg.Path.ProcFS)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &arpCollector{
		fs:           fs,
		deviceFilter: newDeviceFilter(cfg.ARP.DeviceExclude, cfg.ARP.DeviceInclude),
		logger:       logger,
		config:       cfg.ARP,
	}, nil
}

func getTotalArpEntries(deviceEntries []procfs.ARPEntry) map[string]uint32 {
	entries := make(map[string]uint32)

	for _, device := range deviceEntries {
		entries[device.Device]++
	}

	return entries
}

func getTotalArpEntriesRTNL() (map[string]uint32, error) {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Neighbors will also contain IPv6 neighbors, but since this is purely an ARP collector,
	// restrict to AF_INET.
	neighbors, err := conn.Neighbours(nil, unix.AF_INET)
	if err != nil {
		return nil, err
	}

	// Map of interface name to ARP neighbor count.
	entries := make(map[string]uint32)

	for _, n := range neighbors {
		// Skip entries which have state NUD_NOARP to conform to output of /proc/net/arp.
		if n.State&unix.NUD_NOARP == 0 {
			entries[n.Interface.Name]++
		}
	}

	return entries, nil
}

func (c *arpCollector) Update(_ context.Context, ch chan<- prometheus.Metric) error {
	var enumeratedEntry map[string]uint32

	if c.config.Netlink {
		var err error

		enumeratedEntry, err = getTotalArpEntriesRTNL()
		if err != nil {
			return fmt.Errorf("could not get ARP entries: %w", err)
		}
	} else {
		entries, err := c.fs.GatherARPEntries()
		if err != nil {
			return fmt.Errorf("could not get ARP entries: %w", err)
		}

		enumeratedEntry = getTotalArpEntries(entries)
	}

	for device, entryCount := range enumeratedEntry {
		if c.deviceFilter.ignored(device) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			arpEntries, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
