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
	"errors"
	"fmt"
	"net"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/jsimonetti/rtnetlink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"
)

var (
	arpDeviceInclude = kingpin.Flag("collector.arp.device-include", "Regexp of arp devices to include (mutually exclusive to device-exclude).").String()
	arpDeviceExclude = kingpin.Flag("collector.arp.device-exclude", "Regexp of arp devices to exclude (mutually exclusive to device-include).").String()
	arpNetlink       = kingpin.Flag("collector.arp.netlink", "Use netlink to gather stats instead of /proc/net/arp.").Default("true").Bool()
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

func getTotalArpEntriesRTNL() (map[string]uint32, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	neighbors, err := conn.Neigh.List()
	if err != nil {
		return nil, err
	}

	ifIndexEntries := make(map[uint32]uint32)

	for _, n := range neighbors {
		// Neighbors will also contain IPv6 neighbors, but since this is purely an ARP collector,
		// restrict to AF_INET. Also skip entries which have state NUD_NOARP to conform to output
		// of /proc/net/arp.
		if n.Family == unix.AF_INET && n.State&unix.NUD_NOARP == 0 {
			ifIndexEntries[n.Index]++
		}
	}

	enumEntries := make(map[string]uint32)

	// Convert interface indexes to names.
	for ifIndex, entryCount := range ifIndexEntries {
		iface, err := net.InterfaceByIndex(int(ifIndex))
		if err != nil {
			if errors.Unwrap(err).Error() == "no such network interface" {
				continue
			}
			return nil, err
		}

		enumEntries[iface.Name] = entryCount
	}

	return enumEntries, nil
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	var enumeratedEntry map[string]uint32

	if *arpNetlink {
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
			c.entries, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
