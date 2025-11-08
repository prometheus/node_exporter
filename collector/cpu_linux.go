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

//go:build !nocpu
// +build !nocpu

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"golang.org/x/exp/maps"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/sysfs"
)

type cpuCollector struct {
	procfs             procfs.FS
	sysfs              sysfs.FS
	cpu                *prometheus.Desc
	cpuInfo            *prometheus.Desc
	cpuFrequencyHz     *prometheus.Desc
	cpuFlagsInfo       *prometheus.Desc
	cpuBugsInfo        *prometheus.Desc
	cpuGuest           *prometheus.Desc
	cpuCoreThrottle    *prometheus.Desc
	cpuPackageThrottle *prometheus.Desc
	cpuIsolated        *prometheus.Desc
	logger             *slog.Logger
	cpuOnline          *prometheus.Desc
	cpuStats           map[int64]procfs.CPUStat
	cpuStatsMutex      sync.Mutex
	isolatedCpus       []uint16

	cpuFlagsIncludeRegexp *regexp.Regexp
	cpuBugsIncludeRegexp  *regexp.Regexp
}

// Idle jump back limit in seconds.
const jumpBackSeconds = 3.0

var (
	enableCPUGuest       = kingpin.Flag("collector.cpu.guest", "Enables metric node_cpu_guest_seconds_total").Default("true").Bool()
	enableCPUInfo        = kingpin.Flag("collector.cpu.info", "Enables metric cpu_info").Bool()
	flagsInclude         = kingpin.Flag("collector.cpu.info.flags-include", "Filter the `flags` field in cpuInfo with a value that must be a regular expression").String()
	bugsInclude          = kingpin.Flag("collector.cpu.info.bugs-include", "Filter the `bugs` field in cpuInfo with a value that must be a regular expression").String()
	jumpBackDebugMessage = fmt.Sprintf("CPU Idle counter jumped backwards more than %f seconds, possible hotplug event, resetting CPU stats", jumpBackSeconds)
)

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

// NewCPUCollector returns a new Collector exposing kernel/system statistics.
func NewCPUCollector(logger *slog.Logger) (Collector, error) {
	pfs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	sfs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	isolcpus, err := sfs.IsolatedCPUs()
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("unable to get isolated cpus: %w", err)
		}
		logger.Debug("couldn't open isolated file", "error", err)
	}

	c := &cpuCollector{
		procfs: pfs,
		sysfs:  sfs,
		cpu:    nodeCPUSecondsDesc,
		cpuInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "info"),
			"CPU information from /proc/cpuinfo.",
			[]string{"package", "core", "cpu", "vendor", "family", "model", "model_name", "microcode", "stepping", "cachesize"}, nil,
		),
		cpuFrequencyHz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz"),
			"CPU frequency in hertz from /proc/cpuinfo.",
			[]string{"package", "core", "cpu"}, nil,
		),
		cpuFlagsInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "flag_info"),
			"The `flags` field of CPU information from /proc/cpuinfo taken from the first core.",
			[]string{"flag"}, nil,
		),
		cpuBugsInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "bug_info"),
			"The `bugs` field of CPU information from /proc/cpuinfo taken from the first core.",
			[]string{"bug"}, nil,
		),
		cpuGuest: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "guest_seconds_total"),
			"Seconds the CPUs spent in guests (VMs) for each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		cpuCoreThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "core_throttles_total"),
			"Number of times this CPU core has been throttled.",
			[]string{"package", "core"}, nil,
		),
		cpuPackageThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "package_throttles_total"),
			"Number of times this CPU package has been throttled.",
			[]string{"package"}, nil,
		),
		cpuIsolated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "isolated"),
			"Whether each core is isolated, information from /sys/devices/system/cpu/isolated.",
			[]string{"cpu"}, nil,
		),
		cpuOnline: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "online"),
			"CPUs that are online and being scheduled.",
			[]string{"cpu"}, nil,
		),
		logger:       logger,
		isolatedCpus: isolcpus,
		cpuStats:     make(map[int64]procfs.CPUStat),
	}
	err = c.compileIncludeFlags(flagsInclude, bugsInclude)
	if err != nil {
		return nil, fmt.Errorf("fail to compile --collector.cpu.info.flags-include and --collector.cpu.info.bugs-include, the values of them must be regular expressions: %w", err)
	}
	return c, nil
}

func (c *cpuCollector) compileIncludeFlags(flagsIncludeFlag, bugsIncludeFlag *string) error {
	if (*flagsIncludeFlag != "" || *bugsIncludeFlag != "") && !*enableCPUInfo {
		*enableCPUInfo = true
		c.logger.Info("--collector.cpu.info has been set to `true` because you set the following flags, like --collector.cpu.info.flags-include and --collector.cpu.info.bugs-include")
	}

	var err error
	if *flagsIncludeFlag != "" {
		c.cpuFlagsIncludeRegexp, err = regexp.Compile(*flagsIncludeFlag)
		if err != nil {
			return err
		}
	}
	if *bugsIncludeFlag != "" {
		c.cpuBugsIncludeRegexp, err = regexp.Compile(*bugsIncludeFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	if *enableCPUInfo {
		if err := c.updateInfo(ch); err != nil {
			return err
		}
	}
	if err := c.updateStat(ch); err != nil {
		return err
	}
	if c.isolatedCpus != nil {
		c.updateIsolated(ch)
	}
	err := c.updateThermalThrottle(ch)
	if err != nil {
		return err
	}
	err = c.updateOnline(ch)
	if err != nil {
		return err
	}

	return nil
}

// updateInfo reads /proc/cpuinfo
func (c *cpuCollector) updateInfo(ch chan<- prometheus.Metric) error {
	info, err := c.procfs.CPUInfo()
	if err != nil {
		return err
	}
	for _, cpu := range info {
		ch <- prometheus.MustNewConstMetric(c.cpuInfo,
			prometheus.GaugeValue,
			1,
			cpu.PhysicalID,
			cpu.CoreID,
			strconv.Itoa(int(cpu.Processor)),
			cpu.VendorID,
			cpu.CPUFamily,
			cpu.Model,
			cpu.ModelName,
			cpu.Microcode,
			cpu.Stepping,
			cpu.CacheSize)
	}

	cpuFreqEnabled, ok := collectorState["cpufreq"]
	if !ok || cpuFreqEnabled == nil {
		c.logger.Debug("cpufreq key missing or nil value in collectorState map")
	} else if *cpuFreqEnabled {
		for _, cpu := range info {
			ch <- prometheus.MustNewConstMetric(c.cpuFrequencyHz,
				prometheus.GaugeValue,
				cpu.CPUMHz*1e6,
				cpu.PhysicalID,
				cpu.CoreID,
				strconv.Itoa(int(cpu.Processor)))
		}
	}

	if len(info) != 0 {
		cpu := info[0]
		if err := updateFieldInfo(cpu.Flags, c.cpuFlagsIncludeRegexp, c.cpuFlagsInfo, ch); err != nil {
			return err
		}
		if err := updateFieldInfo(cpu.Bugs, c.cpuBugsIncludeRegexp, c.cpuBugsInfo, ch); err != nil {
			return err
		}
	}

	return nil
}

func updateFieldInfo(valueList []string, filter *regexp.Regexp, desc *prometheus.Desc, ch chan<- prometheus.Metric) error {
	if filter == nil {
		return nil
	}

	for _, val := range valueList {
		if !filter.MatchString(val) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(desc,
			prometheus.GaugeValue,
			1,
			val,
		)
	}
	return nil
}

// updateThermalThrottle reads /sys/devices/system/cpu/cpu* and expose thermal throttle statistics.
func (c *cpuCollector) updateThermalThrottle(ch chan<- prometheus.Metric) error {
	cpus, err := filepath.Glob(sysFilePath("devices/system/cpu/cpu[0-9]*"))
	if err != nil {
		return err
	}

	packageThrottles := make(map[uint64]uint64)
	packageCoreThrottles := make(map[uint64]map[uint64]uint64)

	// cpu loop
	for _, cpu := range cpus {
		// See
		// https://www.kernel.org/doc/Documentation/x86/topology.txt
		// https://www.kernel.org/doc/Documentation/cputopology.txt
		// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-devices-system-cpu
		var err error
		var physicalPackageID, coreID uint64

		// topology/physical_package_id
		if physicalPackageID, err = readUintFromFile(filepath.Join(cpu, "topology", "physical_package_id")); err != nil {
			c.logger.Debug("CPU is missing physical_package_id", "cpu", cpu)
			continue
		}
		// topology/core_id
		if coreID, err = readUintFromFile(filepath.Join(cpu, "topology", "core_id")); err != nil {
			c.logger.Debug("CPU is missing core_id", "cpu", cpu)
			continue
		}

		// metric node_cpu_core_throttles_total
		//
		// We process this metric before the package throttles as there
		// are CPU+kernel combinations that only present core throttles
		// but no package throttles.
		// Seen e.g. on an Intel Xeon E5472 system with RHEL 6.9 kernel.
		if _, present := packageCoreThrottles[physicalPackageID]; !present {
			packageCoreThrottles[physicalPackageID] = make(map[uint64]uint64)
		}
		if _, present := packageCoreThrottles[physicalPackageID][coreID]; !present {
			// Read thermal_throttle/core_throttle_count only once
			if coreThrottleCount, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle", "core_throttle_count")); err == nil {
				packageCoreThrottles[physicalPackageID][coreID] = coreThrottleCount
			} else {
				c.logger.Debug("CPU is missing core_throttle_count", "cpu", cpu)
			}
		}

		// metric node_cpu_package_throttles_total
		if _, present := packageThrottles[physicalPackageID]; !present {
			// Read thermal_throttle/package_throttle_count only once
			if packageThrottleCount, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle", "package_throttle_count")); err == nil {
				packageThrottles[physicalPackageID] = packageThrottleCount
			} else {
				c.logger.Debug("CPU is missing package_throttle_count", "cpu", cpu)
			}
		}
	}

	for physicalPackageID, packageThrottleCount := range packageThrottles {
		ch <- prometheus.MustNewConstMetric(c.cpuPackageThrottle,
			prometheus.CounterValue,
			float64(packageThrottleCount),
			strconv.FormatUint(physicalPackageID, 10))
	}

	for physicalPackageID, coreMap := range packageCoreThrottles {
		for coreID, coreThrottleCount := range coreMap {
			ch <- prometheus.MustNewConstMetric(c.cpuCoreThrottle,
				prometheus.CounterValue,
				float64(coreThrottleCount),
				strconv.FormatUint(physicalPackageID, 10),
				strconv.FormatUint(coreID, 10))
		}
	}
	return nil
}

// updateIsolated reads /sys/devices/system/cpu/isolated through sysfs and exports isolation level metrics.
func (c *cpuCollector) updateIsolated(ch chan<- prometheus.Metric) {
	for _, cpu := range c.isolatedCpus {
		cpuNum := strconv.Itoa(int(cpu))
		ch <- prometheus.MustNewConstMetric(c.cpuIsolated, prometheus.GaugeValue, 1.0, cpuNum)
	}
}

// updateOnline reads /sys/devices/system/cpu/cpu*/online through sysfs and exports online status metrics.
func (c *cpuCollector) updateOnline(ch chan<- prometheus.Metric) error {
	cpus, err := c.sysfs.CPUs()
	if err != nil {
		return err
	}
	// No-op if the system does not support CPU online stats.
	cpu0 := cpus[0]
	if _, err := cpu0.Online(); err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	}
	for _, cpu := range cpus {
		setOnline := float64(0)
		if online, _ := cpu.Online(); online {
			setOnline = 1
		}
		ch <- prometheus.MustNewConstMetric(c.cpuOnline, prometheus.GaugeValue, setOnline, cpu.Number())
	}

	return nil
}

// updateStat reads /proc/stat through procfs and exports CPU-related metrics.
func (c *cpuCollector) updateStat(ch chan<- prometheus.Metric) error {
	stats, err := c.procfs.Stat()
	if err != nil {
		return err
	}

	c.updateCPUStats(stats.CPU)

	// Acquire a lock to read the stats.
	c.cpuStatsMutex.Lock()
	defer c.cpuStatsMutex.Unlock()
	for cpuID, cpuStat := range c.cpuStats {
		cpuNum := strconv.Itoa(int(cpuID))
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.User, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Nice, cpuNum, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.System, cpuNum, "system")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Idle, cpuNum, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Iowait, cpuNum, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.IRQ, cpuNum, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.SoftIRQ, cpuNum, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Steal, cpuNum, "steal")

		if *enableCPUGuest {
			// Guest CPU is also accounted for in cpuStat.User and cpuStat.Nice, expose these as separate metrics.
			ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.Guest, cpuNum, "user")
			ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.GuestNice, cpuNum, "nice")
		}
	}

	return nil
}

// updateCPUStats updates the internal cache of CPU stats.
func (c *cpuCollector) updateCPUStats(newStats map[int64]procfs.CPUStat) {

	// Acquire a lock to update the stats.
	c.cpuStatsMutex.Lock()
	defer c.cpuStatsMutex.Unlock()

	// Reset the cache if the list of CPUs has changed.
	for i, n := range newStats {
		cpuStats := c.cpuStats[i]

		// If idle jumps backwards by more than X seconds, assume we had a hotplug event and reset the stats for this CPU.
		if (cpuStats.Idle - n.Idle) >= jumpBackSeconds {
			c.logger.Debug(jumpBackDebugMessage, "cpu", i, "old_value", cpuStats.Idle, "new_value", n.Idle)
			cpuStats = procfs.CPUStat{}
		}

		if n.Idle >= cpuStats.Idle {
			cpuStats.Idle = n.Idle
		} else {
			c.logger.Debug("CPU Idle counter jumped backwards", "cpu", i, "old_value", cpuStats.Idle, "new_value", n.Idle)
		}

		if n.User >= cpuStats.User {
			cpuStats.User = n.User
		} else {
			c.logger.Debug("CPU User counter jumped backwards", "cpu", i, "old_value", cpuStats.User, "new_value", n.User)
		}

		if n.Nice >= cpuStats.Nice {
			cpuStats.Nice = n.Nice
		} else {
			c.logger.Debug("CPU Nice counter jumped backwards", "cpu", i, "old_value", cpuStats.Nice, "new_value", n.Nice)
		}

		if n.System >= cpuStats.System {
			cpuStats.System = n.System
		} else {
			c.logger.Debug("CPU System counter jumped backwards", "cpu", i, "old_value", cpuStats.System, "new_value", n.System)
		}

		if n.Iowait >= cpuStats.Iowait {
			cpuStats.Iowait = n.Iowait
		} else {
			c.logger.Debug("CPU Iowait counter jumped backwards", "cpu", i, "old_value", cpuStats.Iowait, "new_value", n.Iowait)
		}

		if n.IRQ >= cpuStats.IRQ {
			cpuStats.IRQ = n.IRQ
		} else {
			c.logger.Debug("CPU IRQ counter jumped backwards", "cpu", i, "old_value", cpuStats.IRQ, "new_value", n.IRQ)
		}

		if n.SoftIRQ >= cpuStats.SoftIRQ {
			cpuStats.SoftIRQ = n.SoftIRQ
		} else {
			c.logger.Debug("CPU SoftIRQ counter jumped backwards", "cpu", i, "old_value", cpuStats.SoftIRQ, "new_value", n.SoftIRQ)
		}

		if n.Steal >= cpuStats.Steal {
			cpuStats.Steal = n.Steal
		} else {
			c.logger.Debug("CPU Steal counter jumped backwards", "cpu", i, "old_value", cpuStats.Steal, "new_value", n.Steal)
		}

		if n.Guest >= cpuStats.Guest {
			cpuStats.Guest = n.Guest
		} else {
			c.logger.Debug("CPU Guest counter jumped backwards", "cpu", i, "old_value", cpuStats.Guest, "new_value", n.Guest)
		}

		if n.GuestNice >= cpuStats.GuestNice {
			cpuStats.GuestNice = n.GuestNice
		} else {
			c.logger.Debug("CPU GuestNice counter jumped backwards", "cpu", i, "old_value", cpuStats.GuestNice, "new_value", n.GuestNice)
		}

		c.cpuStats[i] = cpuStats
	}

	// Remove offline CPUs.
	if len(newStats) != len(c.cpuStats) {
		onlineCPUIds := maps.Keys(newStats)
		maps.DeleteFunc(c.cpuStats, func(key int64, item procfs.CPUStat) bool {
			return !slices.Contains(onlineCPUIds, key)
		})
	}
}
