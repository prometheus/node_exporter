// Copyright 2020 The Prometheus Authors
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
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/shirou/gopsutil/cpu"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cpuCollector struct {
	cpu           *prometheus.Desc
	cpuInfo       *prometheus.Desc
	cpuTimes      *prometheus.Desc
	cpuGuest      *prometheus.Desc
	perCPUPercent *prometheus.Desc
	logger        log.Logger
	cpuStats      []procfs.CPUStat
	cpuStatsMutex sync.Mutex
}

var (
	enableCPUInfo = kingpin.Flag("collector.cpu.info", "Enables metric cpu_info").Bool()
)

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

// NewCPUCollector returns a new Collector exposing kernel/system statistics.
func NewCPUCollector(logger log.Logger) (Collector, error) {
	return &cpuCollector{
		cpu: nodeCPUSecondsDesc,
		cpuInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "info"),
			"CPU information",
			[]string{"package", "core", "cpu", "vendor", "family", "model", "model_name", "microcode", "stepping", "cachesize"}, nil,
		),
		cpuGuest: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "guest_seconds_total"),
			"Seconds the cpus spent in guests (VMs) for each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		cpuTimes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "per_seconds_total"),
			"The amounts of time the CPU has spent performing different kinds of work.",
			[]string{"cpu", "mode"}, nil,
		),
		perCPUPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "percentage"),
			"Percent calculates the percentage of cpu used either per CPU",
			[]string{"cpu"}, nil,
		),
		logger: logger,
	}, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	if *enableCPUInfo {
		if err := c.updateInfo(ch); err != nil {
			return err
		}
	}
	if err := c.updateTimesStat(ch); err != nil {
		return err
	}
	if err := c.updatePercent(ch); err != nil {
		return err
	}
	return nil
}

// updateInfo cpuinfo
func (c *cpuCollector) updateInfo(ch chan<- prometheus.Metric) error {
	info, err := cpu.Info()
	if err != nil {
		return err
	}
	for _, cpu := range info {
		ch <- prometheus.MustNewConstMetric(c.cpuInfo,
			prometheus.GaugeValue,
			1,
			cpu.CoreID,
			strconv.Itoa(int(cpu.Cores)),
			cpu.VendorID,
			cpu.Family,
			cpu.Model,
			cpu.ModelName,
			cpu.Microcode,
			strconv.Itoa(int(cpu.Stepping)),
			strconv.Itoa(int(cpu.CacheSize)))
	}
	return nil
}

// updatePercent do percent calculates the percentage of cpu used either per CPU or combined.
// If an interval of 0 is given it will compare the current cpu times against the last call.
// Returns one value per cpu, or a single value if percpu is set to false.
func (c *cpuCollector) updatePercent(ch chan<- prometheus.Metric) error {
	percents, err := cpu.Percent(time.Duration(1)*time.Second, true)
	if err != nil {
		return err
	}

	var allCPU float64
	for cpuID, percent := range percents {
		ch <- prometheus.MustNewConstMetric(c.perCPUPercent, prometheus.CounterValue, percent, fmt.Sprintf("cpu%v", cpuID))
		allCPU += percent
	}
	ch <- prometheus.MustNewConstMetric(c.perCPUPercent, prometheus.CounterValue, allCPU/float64(len(percents)), "cpu")

	return nil
}

// updateStat exports cpu related metrics.
func (c *cpuCollector) updateTimesStat(ch chan<- prometheus.Metric) error {
	c.cpuStatsMutex.Lock()
	defer c.cpuStatsMutex.Unlock()

	cpuStats, err := cpu.Times(true)
	if err != nil {
		return err
	}
	for cpuID, cpuStat := range cpuStats {
		cpuNum := strconv.Itoa(cpuID)
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.User, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Nice, cpuNum, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.System, cpuNum, "system")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Idle, cpuNum, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Iowait, cpuNum, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Irq, cpuNum, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Softirq, cpuNum, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpuTimes, prometheus.CounterValue, cpuStat.Steal, cpuNum, "steal")
	}

	totalCPUStat, err := cpu.Times(false)
	if err != nil {
		return err
	}
	for _, total := range totalCPUStat {
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.User, total.CPU, "user")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Nice, total.CPU, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.System, total.CPU, "system")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Idle, total.CPU, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Iowait, total.CPU, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Irq, total.CPU, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Softirq, total.CPU, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, total.Steal, total.CPU, "steal")

		// Guest CPU is also accounted for in cpuStat.User and cpuStat.Nice, expose these as separate metrics.
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, total.Guest, total.CPU, "user")
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, total.GuestNice, total.CPU, "nice")
	}
	return nil
}
