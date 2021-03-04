// Copyright 2019 The Prometheus Authors
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

// +build !nobtrfs

package collector

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/btrfs"
)

// A btrfsCollector is a Collector which gathers metrics from Btrfs filesystems.
type btrfsCollector struct {
	fs     btrfs.FS
	logger log.Logger
}

func init() {
	registerCollector("btrfs", defaultEnabled, NewBtrfsCollector)
}

// NewBtrfsCollector returns a new Collector exposing Btrfs statistics.
func NewBtrfsCollector(logger log.Logger) (Collector, error) {
	fs, err := btrfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &btrfsCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update retrieves and exports Btrfs statistics.
// It implements Collector.
func (c *btrfsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.Stats()
	if err != nil {
		return fmt.Errorf("failed to retrieve Btrfs stats: %w", err)
	}

	for _, s := range stats {
		c.updateBtrfsStats(ch, s)
	}

	return nil
}

// btrfsMetric represents a single Btrfs metric that is converted into a Prometheus Metric.
type btrfsMetric struct {
	name            string
	desc            string
	value           float64
	extraLabel      []string
	extraLabelValue []string
}

// updateBtrfsStats collects statistics for one bcache ID.
func (c *btrfsCollector) updateBtrfsStats(ch chan<- prometheus.Metric, s *btrfs.Stats) {
	const subsystem = "btrfs"

	// Basic information about the filesystem.
	devLabels := []string{"uuid"}

	// Retrieve the metrics.
	metrics := c.getMetrics(s)

	// Convert all gathered metrics to Prometheus Metrics and add to channel.
	for _, m := range metrics {
		labels := append(devLabels, m.extraLabel...)

		desc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, m.name),
			m.desc,
			labels,
			nil,
		)

		labelValues := []string{s.UUID}
		if len(m.extraLabelValue) > 0 {
			labelValues = append(labelValues, m.extraLabelValue...)
		}

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			m.value,
			labelValues...,
		)
	}
}

// getMetrics returns metrics for the given Btrfs statistics.
func (c *btrfsCollector) getMetrics(s *btrfs.Stats) []btrfsMetric {
	metrics := []btrfsMetric{
		{
			name:            "info",
			desc:            "Filesystem information",
			value:           1,
			extraLabel:      []string{"label"},
			extraLabelValue: []string{s.Label},
		},
		{
			name:  "global_rsv_size_bytes",
			desc:  "Size of global reserve.",
			value: float64(s.Allocation.GlobalRsvSize),
		},
	}

	// Information about devices.
	for n, dev := range s.Devices {
		metrics = append(metrics, btrfsMetric{
			name:            "device_size_bytes",
			desc:            "Size of a device that is part of the filesystem.",
			value:           float64(dev.Size),
			extraLabel:      []string{"device"},
			extraLabelValue: []string{n},
		})
	}

	// Information about data, metadata and system data.
	metrics = append(metrics, c.getAllocationStats("data", s.Allocation.Data)...)
	metrics = append(metrics, c.getAllocationStats("metadata", s.Allocation.Metadata)...)
	metrics = append(metrics, c.getAllocationStats("system", s.Allocation.System)...)

	return metrics
}

// getAllocationStats returns allocation metrics for the given Btrfs Allocation statistics.
func (c *btrfsCollector) getAllocationStats(a string, s *btrfs.AllocationStats) []btrfsMetric {
	metrics := []btrfsMetric{
		{
			name:            "reserved_bytes",
			desc:            "Amount of space reserved for a data type",
			value:           float64(s.ReservedBytes),
			extraLabel:      []string{"block_group_type"},
			extraLabelValue: []string{a},
		},
	}

	// Add all layout statistics.
	for layout, stats := range s.Layouts {
		metrics = append(metrics, c.getLayoutStats(a, layout, stats)...)
	}

	return metrics
}

// getLayoutStats returns metrics for a data layout.
func (c *btrfsCollector) getLayoutStats(a, l string, s *btrfs.LayoutUsage) []btrfsMetric {
	return []btrfsMetric{
		{
			name:            "used_bytes",
			desc:            "Amount of used space by a layout/data type",
			value:           float64(s.UsedBytes),
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
		{
			name:            "size_bytes",
			desc:            "Amount of space allocated for a layout/data type",
			value:           float64(s.TotalBytes),
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
		{
			name:            "allocation_ratio",
			desc:            "Data allocation ratio for a layout/data type",
			value:           s.Ratio,
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
	}
}
