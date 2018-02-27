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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

var (
	digitRegexp = regexp.MustCompile("[0-9]+")
)

type cpuCollector struct {
	cpu                *prometheus.Desc
	cpuGuest           *prometheus.Desc
	cpuFreq            *prometheus.Desc
	cpuFreqMin         *prometheus.Desc
	cpuFreqMax         *prometheus.Desc
	cpuCoreThrottle    *prometheus.Desc
	cpuPackageThrottle *prometheus.Desc
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

// NewCPUCollector returns a new Collector exposing kernel/system statistics.
func NewCPUCollector() (Collector, error) {
	return &cpuCollector{
		cpu: nodeCpuSecondsDesc,
		cpuGuest: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "guest_seconds_total"),
			"Seconds the cpus spent in guests (VMs) for each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		cpuFreq: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz"),
			"Current cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuFreqMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_min_hertz"),
			"Minimum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuFreqMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_max_hertz"),
			"Maximum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuCoreThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "core_throttles_total"),
			"Number of times this cpu core has been throttled.",
			[]string{"core"}, nil,
		),
		cpuPackageThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "package_throttles_total"),
			"Number of times this cpu package has been throttled.",
			[]string{"node"}, nil,
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

	cpu_core_throttles := make(map[int]uint64)

	// cpu loop
	for _, cpu := range cpus {
		_, cpuName := filepath.Split(cpu)
		cpuNum := strings.TrimPrefix(cpuName, "cpu")

		if _, err := os.Stat(filepath.Join(cpu, "cpufreq")); os.IsNotExist(err) {
			log.Debugf("CPU %v is missing cpufreq", cpu)
		} else {
			// sysfs cpufreq values are kHz, thus multiply by 1000 to export base units (hz).
			// See https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt
			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq", "scaling_cur_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreq, prometheus.GaugeValue, float64(value)*1000.0, cpuNum)

			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq", "scaling_min_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreqMin, prometheus.GaugeValue, float64(value)*1000.0, cpuNum)

			if value, err = readUintFromFile(filepath.Join(cpu, "cpufreq", "scaling_max_freq")); err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.cpuFreqMax, prometheus.GaugeValue, float64(value)*1000.0, cpuNum)
		}

		if _, err := os.Stat(filepath.Join(cpu, "thermal_throttle")); os.IsNotExist(err) {
			log.Debugf("CPU %v is missing thermal_throttle", cpu)
			continue
		}

		if value, err := readUintFromFile(filepath.Join(cpu, "topology/core_id")); err != nil {
			log.Debugf("CPU %v is misssing topology/core_id", cpu)
		} else {
			core_id := int(value)
			if value, err = readUintFromFile(filepath.Join(cpu, "thermal_throttle", "core_throttle_count")); err != nil {
				return err
			}
			cpu_core_throttles[core_id] = value
		}
	}

	// core throttles
	for core_id, value := range cpu_core_throttles {
		ch <- prometheus.MustNewConstMetric(c.cpuCoreThrottle, prometheus.CounterValue, float64(value), strconv.Itoa(core_id))
	}

	nodes, err := filepath.Glob(sysFilePath("bus/node/devices/node[0-9]*"))
	if err != nil {
		return err
	}

	// package / NUMA node loop
	for _, node := range nodes {
		if _, err := os.Stat(filepath.Join(node, "cpulist")); os.IsNotExist(err) {
			log.Debugf("NUMA node %v is missing cpulist", node)
			continue
		}
		cpulist, err := ioutil.ReadFile(filepath.Join(node, "cpulist"))
		if err != nil {
			log.Debugf("could not read cpulist of NUMA node %v", node)
			return err
		}
		// cpulist example of one package/node with HT: "0-11,24-35"
		line := strings.Split(string(cpulist), "\n")[0]
		if line == "" {
			// Skip processor-less (memory-only) NUMA nodes.
			// E.g. RAM expansion with Intel Optane Drive(s) using
			// Intel Memory Drive Technology (IMDT).
			log.Debugf("skipping processor-less (memory-only) NUMA node %v", node)
			continue
		}
		firstCPU := strings.FieldsFunc(line, func(r rune) bool {
			return r == '-' || r == ','
		})[0]
		if _, err := os.Stat(filepath.Join(node, "cpu"+firstCPU, "thermal_throttle", "package_throttle_count")); os.IsNotExist(err) {
			log.Debugf("Node %v CPU %v is missing package_throttle", node, firstCPU)
			continue
		}
		if value, err = readUintFromFile(filepath.Join(node, "cpu"+firstCPU, "thermal_throttle", "package_throttle_count")); err != nil {
			return err
		}
		nodeno := digitRegexp.FindAllString(node, 1)[0]
		ch <- prometheus.MustNewConstMetric(c.cpuPackageThrottle, prometheus.CounterValue, float64(value), nodeno)
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
		cpuNum := fmt.Sprintf("%d", cpuID)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.User, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Nice, cpuNum, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.System, cpuNum, "system")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Idle, cpuNum, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Iowait, cpuNum, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.IRQ, cpuNum, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.SoftIRQ, cpuNum, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Steal, cpuNum, "steal")

		// Guest CPU is also accounted for in cpuStat.User and cpuStat.Nice, expose these as separate metrics.
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.Guest, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.GuestNice, cpuNum, "nice")
	}

	return nil
}
