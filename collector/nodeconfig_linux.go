// Copyright 2025 The Prometheus Authors
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

//go:build linux && !nonodeconfig

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const (
	nodeconfigSubsystem = "nodeconfig"
	// PCI base class 0x02 = Network controller (see PCI SIG class codes).
	pciClassNetwork = 0x02
	// DMI/SMBIOS structure types (see DMTF DSP0134).
	dmiType16PhysicalMemoryArray = 16
	dmiType17MemoryDevice        = 17
	// Type 16 byte offset: Number of Memory Devices in this array.
	dmiType16NumDevicesOffset = 13
	// Type 17 byte offsets: Size (WORD, MB). 0 = no device, 0x7FFF = unknown.
	dmiType17SizeOffsetLo = 12
	dmiType17SizeOffsetHi = 13
	dmiSizeNotPopulated   = 0x7FFF
)

type nodeconfigCollector struct {
	fs                      sysfs.FS
	logger                  *slog.Logger
	pcieNICMinLinkWidthDesc *prometheus.Desc
	pcieSlotOkDesc          *prometheus.Desc
	coresDedicatedDesc      *prometheus.Desc
	memoryBanksFullDesc     *prometheus.Desc
}

func init() {
	registerCollector("nodeconfig", defaultDisabled, NewNodeconfigCollector)
}

// NewNodeconfigCollector returns a new Collector exposing node-level configuration
// facts useful for runbooks (e.g. DPDK troubleshooting: PCIe slot, memory banks, CPU isolation).
func NewNodeconfigCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &nodeconfigCollector{
		fs:     fs,
		logger: logger,
		pcieNICMinLinkWidthDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeconfigSubsystem, "pcie_nic_min_link_width"),
			"Minimum current PCIe link width (lanes) among PCI network controllers. Use in runbooks to infer PCIe slot correctness (e.g. expect >= 16 for x16 slots). -1 if no network PCIe devices or width unknown.",
			nil, nil,
		),
		pcieSlotOkDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeconfigSubsystem, "pcie_slot_ok"),
			"Whether PCIe slot/width is considered correct (1) or not (0). Derived from PCIe: 1 when minimum NIC link width >= 16, 0 otherwise. Absent if no network PCIe devices.",
			nil, nil,
		),
		coresDedicatedDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeconfigSubsystem, "cores_dedicated"),
			"Whether CPU cores are dedicated/isolated for workload (e.g. DPDK). 1 if at least one CPU is in /sys/devices/system/cpu/isolated, 0 otherwise.",
			nil, nil,
		),
		memoryBanksFullDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, nodeconfigSubsystem, "memory_banks_full"),
			"Whether memory channels/banks are fully populated (1) or not (0). Derived from DMI/SMBIOS: 1 when all memory device slots have a populated DIMM, 0 otherwise. Absent if DMI not available.",
			nil, nil,
		),
	}, nil
}

func (c *nodeconfigCollector) Update(ch chan<- prometheus.Metric) error {
	// PCIe: min link width among network-class PCI devices; pcie_slot_ok derived from it
	minWidth := c.pcieNICMinLinkWidth()
	if minWidth >= 0 {
		ch <- prometheus.MustNewConstMetric(c.pcieNICMinLinkWidthDesc, prometheus.GaugeValue, minWidth)
		pcieOk := 0.0
		if minWidth >= 16 {
			pcieOk = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.pcieSlotOkDesc, prometheus.GaugeValue, pcieOk)
	}

	// Cores dedicated: from sysfs isolated CPUs
	dedicated := 0.0
	isolated, err := c.fs.IsolatedCPUs()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		c.logger.Debug("nodeconfig: could not read isolated CPUs", "err", err)
	} else if len(isolated) > 0 {
		dedicated = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.coresDedicatedDesc, prometheus.GaugeValue, dedicated)

	// Memory banks full: from DMI/SMBIOS
	if full, ok := c.memoryBanksFullFromDMI(); ok {
		ch <- prometheus.MustNewConstMetric(c.memoryBanksFullDesc, prometheus.GaugeValue, full)
	}

	return nil
}

// pcieNICMinLinkWidth returns the minimum current link width among PCI network controllers,
// or -1 if none or unknown.
func (c *nodeconfigCollector) pcieNICMinLinkWidth() float64 {
	devices, err := c.fs.PciDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return -1
		}
		c.logger.Debug("nodeconfig: failed to list PCI devices", "err", err)
		return -1
	}

	var minWidth float64 = -1
	for _, device := range devices {
		baseClass := uint8((device.Class >> 16) & 0xff)
		if baseClass != pciClassNetwork {
			continue
		}
		if device.CurrentLinkWidth == nil {
			continue
		}
		w := *device.CurrentLinkWidth
		if w < 0 {
			continue
		}
		if minWidth < 0 || w < minWidth {
			minWidth = w
		}
	}
	return minWidth
}

// memoryBanksFullFromDMI reads DMI/SMBIOS from /sys/firmware/dmi/entries/ and returns
// (1.0, true) if all memory device slots are populated, (0.0, true) if not, (_, false) if unknown.
func (c *nodeconfigCollector) memoryBanksFullFromDMI() (float64, bool) {
	entriesPath := filepath.Join(*sysPath, "firmware", "dmi", "entries")
	entries, err := os.ReadDir(entriesPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("nodeconfig: could not read DMI entries", "path", entriesPath, "err", err)
		}
		return 0, false
	}

	var totalSlots int  // from Type 16 Number of Memory Devices
	var totalType17 int // count of Type 17 entries (one per slot)
	var populatedCount int

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		base := e.Name()
		typePath := filepath.Join(entriesPath, base, "type")
		dataPath := filepath.Join(entriesPath, base, "data")
		typeBuf, err := os.ReadFile(typePath)
		if err != nil {
			continue
		}
		var dmiType int
		if _, err := fmt.Sscanf(string(typeBuf), "%d", &dmiType); err != nil {
			continue
		}
		data, err := os.ReadFile(dataPath)
		if err != nil {
			continue
		}
		switch dmiType {
		case dmiType16PhysicalMemoryArray:
			if len(data) > dmiType16NumDevicesOffset {
				totalSlots += int(data[dmiType16NumDevicesOffset])
			}
		case dmiType17MemoryDevice:
			totalType17++
			if len(data) > dmiType17SizeOffsetHi {
				size := uint16(data[dmiType17SizeOffsetLo]) | uint16(data[dmiType17SizeOffsetHi])<<8
				if size > 0 && size != dmiSizeNotPopulated {
					populatedCount++
				}
			}
		}
	}

	// Use Type 16 total slots if present; else use Type 17 count as total (one entry per slot).
	if totalSlots == 0 {
		totalSlots = totalType17
	}
	if totalSlots == 0 {
		return 0, false
	}
	if populatedCount >= totalSlots {
		return 1.0, true
	}
	return 0.0, true
}
