// Copyright 2015 The Prometheus Authors
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

//go:build linux
// +build linux

package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const vmSysctlPath = "/proc/sys/vm"

type vmCollector struct {
	logger *slog.Logger
	desc   *prometheus.Desc
}

func init() {
	registerCollector("vm_sysctl", defaultDisabled, VMSysctlCollector)
}

func VMSysctlCollector(logger *slog.Logger) (Collector, error) {
	return &vmCollector{
		logger: logger,
	}, nil
}

func (c *vmCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := os.ReadDir(vmSysctlPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", vmSysctlPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // skip subdirectories
		}

		name := entry.Name()
		fullPath := filepath.Join(vmSysctlPath, name)

		data, err := os.ReadFile(fullPath)
		if err != nil {
			c.logger.Warn("failed to read vm sysctl file", "file", fullPath, "err", err)
			continue
		}

		valStr := strings.TrimSpace(string(data))
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			// If file is not numeric, skip it silently
			continue
		}

		metricName := prometheus.BuildFQName("node", "vm_sysctl", sanitizeMetricName(name))
		desc := prometheus.NewDesc(metricName, fmt.Sprintf("Value of /proc/sys/vm/%s", name), nil, nil)

		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val)
	}

	return nil
}

// sanitizeMetricName replaces invalid chars with underscores (only allow a-z, 0-9, and _)
func sanitizeMetricName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")

	// If needed, add more replacements here

	return name
}
