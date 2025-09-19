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

//go:build !noext4
// +build !noext4

package collector

import (
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/ext4"
)

// An ext4Collector is a Collector which gathers metrics from ext4 filesystems.
type ext4Collector struct {
	fs     ext4.FS
	logger *slog.Logger
}

func init() {
	registerCollector("ext4", defaultEnabled, NewExt4Collector)
}

// NewExt4Collector returns a new Collector exposing ext4 statistics.
func NewExt4Collector(logger *slog.Logger) (Collector, error) {
	fs, err := ext4.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &ext4Collector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update implements Collector.
func (c *ext4Collector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.ProcStat()
	if err != nil {
		return fmt.Errorf("failed to retrieve ext4 stats: %w", err)
	}

	for _, s := range stats {
		c.updateExt4Stats(ch, s)
	}

	return nil
}

// updateExt4Stats collects statistics for a single ext4 filesystem.
func (c *ext4Collector) updateExt4Stats(ch chan<- prometheus.Metric, s *ext4.Stats) {
	const (
		subsystem = "ext4"
	)
	var (
		labels = []string{"device"}
	)

	metrics := []struct {
		name  string
		desc  string
		value float64
	}{
		{
			name:  "errors",
			desc:  "Number of ext4 filesystem errors.",
			value: float64(s.Errors),
		},
		{
			name:  "warnings",
			desc:  "Number of ext4 filesystem warnings.",
			value: float64(s.Warnings),
		},
		{
			name:  "messages",
			desc:  "Number of ext4 filesystem log messages.",
			value: float64(s.Messages),
		},
	}

	for _, m := range metrics {
		desc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, m.name),
			m.desc,
			labels,
			nil,
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.CounterValue,
			m.value,
			s.Name,
		)
	}
}
