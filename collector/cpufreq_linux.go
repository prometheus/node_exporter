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

// +build !nocpufreq

package collector

import (
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	cpufreqSubsystem = "cpufreq"
)

type cpufreqCollector struct {
	curFreq         *prometheus.Desc
	minFreq         *prometheus.Desc
	maxFreq         *prometheus.Desc
	coreThrottle    *prometheus.Desc
	packageThrottle *prometheus.Desc
}

func init() {
	Factories["cpufreq"] = NewCpufreqCollector
}

// NewEdacCollector returns a new Collector exposing edac stats.
func NewCpufreqCollector() (Collector, error) {
	return &cpufreqCollector{
		curFreq: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpufreqSubsystem, "current_frequency_hertz"),
			"Current cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		minFreq: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpufreqSubsystem, "minimum_frequency_hertz"),
			"Minimum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		maxFreq: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpufreqSubsystem, "maximum_frequency_hertz"),
			"Minimum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		coreThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpufreqSubsystem, "core_throttles_total"),
			"Number of times this cpu core has been throttled.",
			[]string{"cpu"}, nil,
		),
		packageThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpufreqSubsystem, "package_throttles_total"),
			"Number of times this cpu package has been throttled.",
			[]string{"cpu"}, nil,
		),
	}, nil
}

func (c *cpufreqCollector) Update(ch chan<- prometheus.Metric) error {
	cpus, err := filepath.Glob(sysFilePath("bus/cpu/devices/cpu[0-9]*"))
	if err != nil {
		return err
	}

	for _, cpu := range cpus {
		_, cpuname := filepath.Split(cpu)

		if value, err := readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_cur_freq")); err != nil {
			return err
		} else {
			ch <- prometheus.MustNewConstMetric(c.curFreq, prometheus.GaugeValue, float64(value), cpuname)
		}

		if value, err := readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_min_freq")); err != nil {
			return err
		} else {
			ch <- prometheus.MustNewConstMetric(c.minFreq, prometheus.GaugeValue, float64(value), cpuname)
		}

		if value, err := readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_max_freq")); err != nil {
			return err
		} else {
			ch <- prometheus.MustNewConstMetric(c.maxFreq, prometheus.GaugeValue, float64(value), cpuname)
		}

		if value, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle/core_throttle_count")); err != nil {
			return err
		} else {
			ch <- prometheus.MustNewConstMetric(c.coreThrottle, prometheus.CounterValue, float64(value), cpuname)
		}

		if value, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle/package_throttle_count")); err != nil {
			return err
		} else {
			ch <- prometheus.MustNewConstMetric(c.packageThrottle, prometheus.CounterValue, float64(value), cpuname)
		}

	}

	return nil
}
