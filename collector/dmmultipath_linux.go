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
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/blockdevice"
)

// isPathActive returns true for device states that indicate a healthy,
// usable path. This covers SCSI ("running") and NVMe ("live") devices.
func isPathActive(state string) bool {
	return state == "running" || state == "live"
}

func init() {
	registerCollector("dmmultipath", defaultEnabled, NewDMMultipathCollector)
}

var (
	dmmultipathDeviceLabels = []string{"device", "sysfs_name"}

	dmmultipathDeviceInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_info"),
		"Non-numeric information about a DM-multipath device.",
		[]string{"device", "sysfs_name", "uuid"}, nil,
	)
	dmmultipathDeviceActive = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_active"),
		"Whether the multipath device-mapper device is active (1) or suspended (0).",
		dmmultipathDeviceLabels, nil,
	)
	dmmultipathDeviceSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_size_bytes"),
		"Size of the multipath device in bytes, read from /sys/block/<dm>/size.",
		dmmultipathDeviceLabels, nil,
	)
	dmmultipathDevicePaths = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_paths"),
		"Number of paths for a multipath device.",
		dmmultipathDeviceLabels, nil,
	)
	dmmultipathDevicePathsActive = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_paths_active"),
		"Number of paths in active state (SCSI running or NVMe live) for a multipath device.",
		dmmultipathDeviceLabels, nil,
	)
	dmmultipathDevicePathsFailed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "device_paths_failed"),
		"Number of paths not in active state for a multipath device.",
		dmmultipathDeviceLabels, nil,
	)
	dmmultipathPathState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dmmultipath", "path_state"),
		"Reports the underlying device state for a multipath path, as read from /sys/block/<dev>/device/state.",
		[]string{"device", "path", "state"}, nil,
	)
)

type dmMultipathCollector struct {
	fs     blockdevice.FS
	logger *slog.Logger
}

// NewDMMultipathCollector returns a new Collector exposing Device Mapper
// multipath device metrics from /sys/block/dm-*.
func NewDMMultipathCollector(logger *slog.Logger) (Collector, error) {
	fs, err := blockdevice.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &dmMultipathCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

func (c *dmMultipathCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.DMMultipathDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			c.logger.Debug("Could not read DM-multipath devices", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to scan DM-multipath devices: %w", err)
	}

	for _, dev := range devices {
		ch <- prometheus.MustNewConstMetric(dmmultipathDeviceInfo, prometheus.GaugeValue, 1,
			dev.Name, dev.SysfsName, dev.UUID)

		active := 0.0
		if !dev.Suspended {
			active = 1.0
		}
		ch <- prometheus.MustNewConstMetric(dmmultipathDeviceActive, prometheus.GaugeValue, active, dev.Name, dev.SysfsName)
		ch <- prometheus.MustNewConstMetric(dmmultipathDeviceSizeBytes, prometheus.GaugeValue, float64(dev.SizeBytes), dev.Name, dev.SysfsName)

		var activePaths, failedPaths float64
		for _, p := range dev.Paths {
			if isPathActive(p.State) {
				activePaths++
			} else {
				failedPaths++
			}

			ch <- prometheus.MustNewConstMetric(dmmultipathPathState, prometheus.GaugeValue, 1,
				dev.Name, p.Device, p.State)
		}

		ch <- prometheus.MustNewConstMetric(dmmultipathDevicePaths, prometheus.GaugeValue, float64(len(dev.Paths)), dev.Name, dev.SysfsName)
		ch <- prometheus.MustNewConstMetric(dmmultipathDevicePathsActive, prometheus.GaugeValue, activePaths, dev.Name, dev.SysfsName)
		ch <- prometheus.MustNewConstMetric(dmmultipathDevicePathsFailed, prometheus.GaugeValue, failedPaths, dev.Name, dev.SysfsName)
	}

	return nil
}
