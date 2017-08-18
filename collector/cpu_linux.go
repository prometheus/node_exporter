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

// +build !nocpu

package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

const (
	cpuCollectorNamespace = "cpu"
)

type cpuCollector struct {
	cpu                *prometheus.Desc
	cpuFreq            *prometheus.Desc
	cpuFreqMin         *prometheus.Desc
	cpuFreqMax         *prometheus.Desc
	cpuCoreThrottle    *prometheus.Desc
	cpuPackageThrottle *prometheus.Desc
}

func init() {
	Factories["cpu"] = NewCPUCollector
}

// NewCPUCollector returns a new Collector exposing kernel/system statistics.
func NewCPUCollector() (Collector, error) {
	return &cpuCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", cpuCollectorNamespace),
			"Seconds the cpus spent in each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		cpuFreq: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpuCollectorNamespace, "frequency_hertz"),
			"Current cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuFreqMin: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpuCollectorNamespace, "frequency_min_hertz"),
			"Minimum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuFreqMax: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpuCollectorNamespace, "frequency_max_hertz"),
			"Maximum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuCoreThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpuCollectorNamespace, "core_throttles_total"),
			"Number of times this cpu core has been throttled.",
			[]string{"cpu"}, nil,
		),
		cpuPackageThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, cpuCollectorNamespace, "package_throttles_total"),
			"Number of times this cpu package has been throttled.",
			[]string{"cpu"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateStat(ch); err != nil {
		return err
	}
	if err := c.updateCPUfreq(ch); err != nil {
		return err
	}
	return nil
}

// updateCPUfreq reads /sys/bus/cpu/devices/cpu* and expose cpu frequency statistics.
func (c *cpuCollector) updateCPUfreq(ch chan<- prometheus.Metric) error {
	cpus, err := filepath.Glob(sysFilePath("bus/cpu/devices/cpu[0-9]*"))
	if err != nil {
		return err
	}

	var value uint64

	for _, cpu := range cpus {
		_, cpuname := filepath.Split(cpu)

		if _, err := os.Stat(filepath.Join(cpu, "cpufreq")); os.IsNotExist(err) {
			log.Debugf("CPU %q is missing cpufreq", cpu)
		} else {
			// sysfs cpufreq values are kHz, thus multiply by 1000 to export base units (hz).
			// See https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt
			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_cur_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreq, prometheus.GaugeValue, float64(value)*1000.0, cpuname)

			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_min_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreqMin, prometheus.GaugeValue, float64(value)*1000.0, cpuname)

			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq/scaling_max_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreqMax, prometheus.GaugeValue, float64(value)*1000.0, cpuname)
		}

		if _, err := os.Stat(filepath.Join(cpu, "thermal_throttle")); os.IsNotExist(err) {
			log.Debugf("CPU %q is missing thermal_throttle", cpu)
		} else {
			if value, err = readUintFromFile(filepath.Join(cpu, "thermal_throttle/core_throttle_count")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuCoreThrottle, prometheus.CounterValue, float64(value), cpuname)

			if value, err = readUintFromFile(filepath.Join(cpu, "thermal_throttle/package_throttle_count")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuPackageThrottle, prometheus.CounterValue, float64(value), cpuname)
		}
	}

	return nil
}

// updateStat reads /proc/stat through procfs and exports cpu related metrics.
func (c *cpuCollector) updateStat(ch chan<- prometheus.Metric) error {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %v", err)
	}
	stats, err := fs.NewStat()
	if err != nil {
		return err
	}

	for cpuID, cpuStat := range stats.CPU {
		cpuName := fmt.Sprintf("cpu%d", cpuID)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.User, cpuName, "user")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Nice, cpuName, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.System, cpuName, "system")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Idle, cpuName, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Iowait, cpuName, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.IRQ, cpuName, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.SoftIRQ, cpuName, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Steal, cpuName, "steal")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Guest, cpuName, "guest")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.GuestNice, cpuName, "guest_nice")
	}

	return nil
}
