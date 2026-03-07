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

//go:build !nonvme

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type nvmeCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("nvme", defaultEnabled, NewNVMeCollector)
}

var (
	nvmeInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "info"),
		"Non-numeric data from /sys/class/nvme/<device>, value is always 1.",
		[]string{"device", "firmware_revision", "model", "serial", "state", "cntlid"},
		nil,
	)
	nvmeNamespaceInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_info"),
		"Information about NVMe namespaces. Value is always 1",
		[]string{"device", "nsid", "ana_state"}, nil,
	)

	nvmeNamespaceCapacityBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_capacity_bytes"),
		"Capacity of the NVMe namespace in bytes. Computed as namespace_size * namespace_logical_block_size",
		[]string{"device", "nsid"}, nil,
	)

	nvmeNamespaceSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_size_bytes"),
		"Size of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/size",
		[]string{"device", "nsid"}, nil,
	)

	nvmeNamespaceUsedBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_used_bytes"),
		"Used space of the NVMe namespace in bytes. Available in /sys/class/nvme/<device>/<namespace>/nuse",
		[]string{"device", "nsid"}, nil,
	)

	nvmeNamespaceLogicalBlockSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nvme", "namespace_logical_block_size_bytes"),
		"Logical block size of the NVMe namespace in bytes. Usually 4Kb. Available in /sys/class/nvme/<device>/<namespace>/queue/logical_block_size",
		[]string{"device", "nsid"}, nil,
	)
)

// NewNVMeCollector returns a new Collector exposing NVMe stats.
func NewNVMeCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	return &nvmeCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

func (c *nvmeCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.NVMeClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("nvme statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining NVMe class info: %w", err)
	}

	for _, device := range devices {
		// Export device-level metrics
		ch <- prometheus.MustNewConstMetric(
			nvmeInfo,
			prometheus.GaugeValue,
			1.0,
			device.Name,
			device.FirmwareRevision,
			device.Model,
			device.Serial,
			device.State,
			device.ControllerID,
		)

		// Export namespace-level metrics
		for _, namespace := range device.Namespaces {
			// Namespace info metric
			ch <- prometheus.MustNewConstMetric(
				nvmeNamespaceInfo,
				prometheus.GaugeValue,
				1.0,
				device.Name,
				namespace.ID,
				namespace.ANAState,
			)

			// Namespace capacity in bytes
			ch <- prometheus.MustNewConstMetric(
				nvmeNamespaceCapacityBytes,
				prometheus.GaugeValue,
				float64(namespace.CapacityBytes),
				device.Name,
				namespace.ID,
			)

			// Namespace size in bytes
			ch <- prometheus.MustNewConstMetric(
				nvmeNamespaceSizeBytes,
				prometheus.GaugeValue,
				float64(namespace.SizeBytes),
				device.Name,
				namespace.ID,
			)

			// Namespace used space in bytes
			ch <- prometheus.MustNewConstMetric(
				nvmeNamespaceUsedBytes,
				prometheus.GaugeValue,
				float64(namespace.UsedBytes),
				device.Name,
				namespace.ID,
			)

			// Namespace logical block size in bytes
			ch <- prometheus.MustNewConstMetric(
				nvmeNamespaceLogicalBlockSizeBytes,
				prometheus.GaugeValue,
				float64(namespace.LogicalBlockSize),
				device.Name,
				namespace.ID,
			)
		}
	}

	return nil
}
