// Copyright The Prometheus Authors
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

//go:build !nodmmultipath

package collector

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type dmMultipathDevice struct {
	Name      string
	SysfsName string
	UUID      string
	Suspended bool
	SizeBytes uint64
	Paths     []dmMultipathPath
}

type dmMultipathPath struct {
	Device string
	State  string
}

var dmPathStates = []string{
	"running", "offline", "blocked", "transport-offline", "unknown",
}

func normalizeDMPathState(raw string) string {
	switch raw {
	case "running", "offline", "blocked", "transport-offline":
		return raw
	case "created", "quiesce":
		return raw
	default:
		return "unknown"
	}
}

type dmMultipathCollector struct {
	logger      *slog.Logger
	scanDevices func() ([]dmMultipathDevice, error)

	deviceInfo        *prometheus.Desc
	deviceActive      *prometheus.Desc
	deviceSizeBytes   *prometheus.Desc
	devicePathsTotal  *prometheus.Desc
	devicePathsActive *prometheus.Desc
	devicePathsFailed *prometheus.Desc
	pathState         *prometheus.Desc
}

func init() {
	registerCollector("dmmultipath", defaultDisabled, NewDMMultipathCollector)
}

// NewDMMultipathCollector returns a new Collector exposing Device Mapper
// multipath device metrics from /sys/block/dm-*.
func NewDMMultipathCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "dmmultipath"

	deviceLabels := []string{"device"}

	c := &dmMultipathCollector{
		logger: logger,
		deviceInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_info"),
			"Non-numeric information about a DM-multipath device.",
			[]string{"device", "sysfs_name", "uuid"}, nil,
		),
		deviceActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_active"),
			"Whether the multipath device-mapper device is active (1) or suspended (0).",
			deviceLabels, nil,
		),
		deviceSizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_size_bytes"),
			"Size of the multipath device in bytes, read from /sys/block/<dm>/size.",
			deviceLabels, nil,
		),
		devicePathsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_total"),
			"Total number of paths for a multipath device.",
			deviceLabels, nil,
		),
		devicePathsActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_active"),
			"Number of paths in running state for a multipath device.",
			deviceLabels, nil,
		),
		devicePathsFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_failed"),
			"Number of paths in non-running state for a multipath device.",
			deviceLabels, nil,
		),
		pathState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "path_state"),
			"Current SCSI device state for a multipath path (1 for the current state, 0 for all others).",
			[]string{"device", "path", "state"}, nil,
		),
	}

	c.scanDevices = func() ([]dmMultipathDevice, error) {
		return scanDMMultipathDevices(*sysPath)
	}

	return c, nil
}

func (c *dmMultipathCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.scanDevices()
	if err != nil {
		return fmt.Errorf("failed to scan DM-multipath devices: %w", err)
	}

	for _, dev := range devices {
		ch <- prometheus.MustNewConstMetric(c.deviceInfo, prometheus.GaugeValue, 1,
			dev.Name, dev.SysfsName, dev.UUID)

		active := 0.0
		if !dev.Suspended {
			active = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.deviceActive, prometheus.GaugeValue, active, dev.Name)
		ch <- prometheus.MustNewConstMetric(c.deviceSizeBytes, prometheus.GaugeValue, float64(dev.SizeBytes), dev.Name)

		var activePaths, failedPaths float64
		for _, p := range dev.Paths {
			state := normalizeDMPathState(p.State)
			if state == "running" {
				activePaths++
			} else {
				failedPaths++
			}

			for _, s := range dmPathStates {
				val := 0.0
				if s == state {
					val = 1.0
				}
				ch <- prometheus.MustNewConstMetric(c.pathState, prometheus.GaugeValue, val,
					dev.Name, p.Device, s)
			}
		}

		ch <- prometheus.MustNewConstMetric(c.devicePathsTotal, prometheus.GaugeValue, float64(len(dev.Paths)), dev.Name)
		ch <- prometheus.MustNewConstMetric(c.devicePathsActive, prometheus.GaugeValue, activePaths, dev.Name)
		ch <- prometheus.MustNewConstMetric(c.devicePathsFailed, prometheus.GaugeValue, failedPaths, dev.Name)
	}

	return nil
}

// scanDMMultipathDevices discovers DM-multipath devices by scanning
// /sys/block/dm-* and filtering on dm/uuid prefix "mpath-".
func scanDMMultipathDevices(sysfsBase string) ([]dmMultipathDevice, error) {
	blockDir := filepath.Join(sysfsBase, "block")

	entries, err := os.ReadDir(blockDir)
	if err != nil {
		return nil, err
	}

	var devices []dmMultipathDevice
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "dm-") {
			continue
		}

		dmDir := filepath.Join(blockDir, entry.Name())
		uuid := readBlockAttr(filepath.Join(dmDir, "dm", "uuid"))
		if !strings.HasPrefix(uuid, "mpath-") {
			continue
		}

		name := readBlockAttr(filepath.Join(dmDir, "dm", "name"))
		if name == "" {
			name = entry.Name()
		}

		suspended := readBlockAttr(filepath.Join(dmDir, "dm", "suspended")) == "1"

		var sizeBytes uint64
		if sectors, err := strconv.ParseUint(readBlockAttr(filepath.Join(dmDir, "size")), 10, 64); err == nil {
			sizeBytes = sectors * uint64(unixSectorSize)
		}

		paths := scanDMPaths(sysfsBase, filepath.Join(dmDir, "slaves"))

		devices = append(devices, dmMultipathDevice{
			Name:      name,
			SysfsName: entry.Name(),
			UUID:      uuid,
			Suspended: suspended,
			SizeBytes: sizeBytes,
			Paths:     paths,
		})
	}

	return devices, nil
}

func scanDMPaths(sysfsBase, slavesDir string) []dmMultipathPath {
	entries, err := os.ReadDir(slavesDir)
	if err != nil {
		return nil
	}

	var paths []dmMultipathPath
	for _, entry := range entries {
		state := readBlockAttr(filepath.Join(sysfsBase, "block", entry.Name(), "device", "state"))
		if state == "" {
			state = "unknown"
		}
		paths = append(paths, dmMultipathPath{
			Device: entry.Name(),
			State:  state,
		})
	}

	return paths
}

// readBlockAttr reads a single sysfs attribute file, returning its
// trimmed content or an empty string on error.
func readBlockAttr(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
