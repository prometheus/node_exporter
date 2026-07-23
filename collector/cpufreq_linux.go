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

//go:build !nocpu

package collector

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type cpuFreqCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("cpufreq", defaultEnabled, NewCPUFreqCollector)
}

// NewCPUFreqCollector returns a new Collector exposing kernel/system statistics.
func NewCPUFreqCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &cpuFreqCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuFreqCollector) Update(ch chan<- prometheus.Metric) error {
	cpuFreqs, err := c.fs.SystemCpufreq()
	if err != nil {
		return err
	}

	// The current frequency metric is also provided by the cpuinfo-based collector.
	// If that collector is enabled and set to provide frequency information,
	// we skip this metric here (otherwise the two collectors would race for it)
	cpuCollectorEnabled, ok := collectorState["cpu"]
	skipCurrentFreq := false
	if !ok || cpuCollectorEnabled == nil {
		c.logger.Debug("cpu key missing or nil value in collectorState map")
	} else {
		skipCurrentFreq = *cpuCollectorEnabled && *enableCPUInfo
	}

	// sysfs cpufreq values are kHz, thus multiply by 1000 to export base units (hz).
	// See https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt
	for _, stats := range cpuFreqs {
		if stats.CpuinfoCurrentFrequency != nil && !skipCurrentFreq {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqHertzDesc,
				prometheus.GaugeValue,
				float64(*stats.CpuinfoCurrentFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.CpuinfoAverageFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqAvgDesc,
				prometheus.GaugeValue,
				float64(*stats.CpuinfoAverageFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.CpuinfoMinimumFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqMinDesc,
				prometheus.GaugeValue,
				float64(*stats.CpuinfoMinimumFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.CpuinfoMaximumFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqMaxDesc,
				prometheus.GaugeValue,
				float64(*stats.CpuinfoMaximumFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.ScalingCurrentFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqScalingFreqDesc,
				prometheus.GaugeValue,
				float64(*stats.ScalingCurrentFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.ScalingMinimumFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqScalingFreqMinDesc,
				prometheus.GaugeValue,
				float64(*stats.ScalingMinimumFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.ScalingMaximumFrequency != nil {
			ch <- prometheus.MustNewConstMetric(
				cpuFreqScalingFreqMaxDesc,
				prometheus.GaugeValue,
				float64(*stats.ScalingMaximumFrequency)*1000.0,
				stats.Name,
			)
		}
		if stats.Governor != "" {
			availableGovernors := strings.SplitSeq(stats.AvailableGovernors, " ")
			for g := range availableGovernors {
				state := 0
				if g == stats.Governor {
					state = 1
				}
				ch <- prometheus.MustNewConstMetric(
					cpuFreqScalingGovernorDesc,
					prometheus.GaugeValue,
					float64(state),
					stats.Name,
					g,
				)
			}
		}
	}
	return nil
}
