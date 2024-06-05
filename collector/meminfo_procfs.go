// Copyright 2024 The Prometheus Authors
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

//go:build linux && !nomeminfo_procfs
// +build linux,!nomeminfo_procfs

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	memInfoProcfsSubsystem = "memory"
)

var (
	memoryBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memInfoProcfsSubsystem, "bytes"),
		"Value in bytes for the labeled field in /proc/meminfo.",
		[]string{"field"}, nil,
	)
)

type meminfoProcfsCollector struct {
	memoryBytesDesc *prometheus.Desc
	fs              procfs.FS
	logger          log.Logger
}

func init() {
	registerCollector("meminfo_procfs", defaultDisabled, NewMeminfoProcfsCollector)
}

// NewMeminfoProcfsCollector returns a new Collector exposing memory stats.
func NewMeminfoProcfsCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &meminfoProcfsCollector{
		fs:              fs,
		logger:          logger,
		memoryBytesDesc: memoryBytesDesc,
	}, nil
}

// Update calls (*meminfoProcfsCollector).getMemInfo to get the platform specific
// memory metrics.
func (c *meminfoProcfsCollector) Update(ch chan<- prometheus.Metric) error {
	memInfo, err := c.getMemInfo()
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %w", err)
	}

	level.Debug(c.logger).Log("msg", "Set node_mem", "memInfoProcfs", fmt.Sprintf("%v", memInfo))
	for k, v := range memInfo {
		ch <- prometheus.MustNewConstMetric(c.memoryBytesDesc, prometheus.GaugeValue, v, k)
	}
	return nil
}
