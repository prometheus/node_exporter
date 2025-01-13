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
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jsimonetti/rtnetlink/v2/rtnl"
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
	states       *prometheus.Desc
	logger       *slog.Logger
}

var neighborStatesMap = map[uint16]string{
	unix.NUD_INCOMPLETE: "incomplete",
	unix.NUD_REACHABLE:  "reachable",
	unix.NUD_STALE:      "stale",
	unix.NUD_DELAY:      "delay",
	unix.NUD_PROBE:      "probe",
	unix.NUD_FAILED:     "failed",
	unix.NUD_PERMANENT:  "permanent",
}

type neighborState struct {
	ip    string
	state string
}

func init() {
	registerCollector("arp", defaultEnabled, NewARPCollector)
}

// NewARPCollector returns a new Collector exposing ARP stats.
func NewARPCollector(logger *slog.Logger) (Collector, error) {
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
		states: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "arp", "states"),
			"ARP states by device",
			[]string{"device", "state"}, nil,
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

func getArpEntriesRTNL() (map[string]uint32, map[string]map[string]int, error) {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	// Neighbors will also contain IPv6 neighbors, but since this is purely an ARP collector,
	// restrict to AF_INET.
	neighbors, err := conn.Neighbours(nil, unix.AF_INET)
	if err != nil {
		return nil, nil, err
	}

	// Map of interface name to ARP neighbor count.
	entries := make(map[string]uint32)
	// Map of map[InterfaceName]map[StateName]int
	states := make(map[string]map[string]int)

	for _, n := range neighbors {
		// Skip entries which have state NUD_NOARP to conform to output of /proc/net/arp.
		if n.State&unix.NUD_NOARP != unix.NUD_NOARP {
			continue
		}

		entries[n.Interface.Name]++

		_, ok := states[n.Interface.Name]
		if !ok {
			states[n.Interface.Name] = make(map[string]int)
		}

		states[n.Interface.Name][neighborStatesMap[n.State]]++
	}

	return entries, states, nil
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	var (
		enumeratedEntry map[string]uint32
		enumStates      map[string]map[string]int
	)

	if *arpNetlink {
		var err error

		enumeratedEntry, enumStates, err = getArpEntriesRTNL()
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

		if *arpNetlink {
			states := enumStates[device]
			for state, count := range states {
				ch <- prometheus.MustNewConstMetric(
					c.states, prometheus.GaugeValue, float64(count), device, state)
			}
		}
	}

	return nil
}
