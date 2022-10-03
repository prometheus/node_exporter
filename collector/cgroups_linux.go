// Copyright 2022 The Prometheus Authors
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

//go:build !nostat
// +build !nostat

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const cgroupsCollectorSubsystem = "cgroups"

type cgroupSummaryCollector struct {
	fs      procfs.FS
	cgroups *prometheus.Desc
	enabled *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector(cgroupsCollectorSubsystem, defaultDisabled, NewCgroupSummaryCollector)
}

// NewCgroupSummaryCollector returns a new Collector exposing a summary of cgroups.
func NewCgroupSummaryCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &cgroupSummaryCollector{
		fs: fs,
		cgroups: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cgroupsCollectorSubsystem, "cgroups"),
			"Current cgroup number of the subsystem.",
			[]string{"subsys_name"}, nil,
		),
		enabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cgroupsCollectorSubsystem, "enabled"),
			"Current cgroup number of the subsystem.",
			[]string{"subsys_name"}, nil,
		),
		logger: logger,
	}, nil
}

// Update implements Collector and exposes cgroup statistics.
func (c *cgroupSummaryCollector) Update(ch chan<- prometheus.Metric) error {
	cgroupSummarys, err := c.fs.CgroupSummarys()
	if err != nil {
		return err
	}
	for _, cs := range cgroupSummarys {
		ch <- prometheus.MustNewConstMetric(c.cgroups, prometheus.GaugeValue, float64(cs.Cgroups), cs.SubsysName)
		ch <- prometheus.MustNewConstMetric(c.enabled, prometheus.GaugeValue, float64(cs.Enabled), cs.SubsysName)
	}
	return nil
}
