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

package collector

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/hodgesds/perf-utils"
	"github.com/prometheus/client_golang/prometheus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	perfSubsystem = "perf"
)

var (
	perfCPUsFlag = kingpin.Flag("collector.perf.cpus", "List of CPUs from which perf metrics should be collected").Default("").String()
)

func init() {
	registerCollector(perfSubsystem, defaultDisabled, NewPerfCollector)
}

// perfCollector is a Collector that uses the perf subsystem to collect
// metrics. It uses perf_event_open an ioctls for profiling. Due to the fact
// that the perf subsystem is highly dependent on kernel configuration and
// settings not all profiler values may be exposed on the target system at any
// given time.
type perfCollector struct {
	hwProfilerCPUMap    map[*perf.HardwareProfiler]int
	swProfilerCPUMap    map[*perf.SoftwareProfiler]int
	cacheProfilerCPUMap map[*perf.CacheProfiler]int
	perfHwProfilers     map[int]*perf.HardwareProfiler
	perfSwProfilers     map[int]*perf.SoftwareProfiler
	perfCacheProfilers  map[int]*perf.CacheProfiler
	desc                map[string]*prometheus.Desc
	logger              log.Logger
}

// perfCPUFlagToCPUs returns a set of CPUs for the perf collectors to monitor.
func perfCPUFlagToCPUs(cpuFlag string) ([]int, error) {
	var err error
	cpus := []int{}
	for _, subset := range strings.Split(cpuFlag, ",") {
		// First parse a single CPU.
		if !strings.Contains(subset, "-") {
			cpu, err := strconv.Atoi(subset)
			if err != nil {
				return nil, err
			}
			cpus = append(cpus, cpu)
			continue
		}

		stride := 1
		// Handle strides, ie 1-10:5 should yield 1,5,10
		strideSet := strings.Split(subset, ":")
		if len(strideSet) == 2 {
			stride, err = strconv.Atoi(strideSet[1])
			if err != nil {
				return nil, err
			}
		}

		rangeSet := strings.Split(strideSet[0], "-")
		if len(rangeSet) != 2 {
			return nil, fmt.Errorf("invalid flag value %q", cpuFlag)
		}
		start, err := strconv.Atoi(rangeSet[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(rangeSet[1])
		if err != nil {
			return nil, err
		}
		for i := start; i <= end; i += stride {
			cpus = append(cpus, i)
		}
	}

	return cpus, nil
}

// NewPerfCollector returns a new perf based collector, it creates a profiler
// per CPU.
func NewPerfCollector(logger log.Logger) (Collector, error) {
	collector := &perfCollector{
		perfHwProfilers:     map[int]*perf.HardwareProfiler{},
		perfSwProfilers:     map[int]*perf.SoftwareProfiler{},
		perfCacheProfilers:  map[int]*perf.CacheProfiler{},
		hwProfilerCPUMap:    map[*perf.HardwareProfiler]int{},
		swProfilerCPUMap:    map[*perf.SoftwareProfiler]int{},
		cacheProfilerCPUMap: map[*perf.CacheProfiler]int{},
		logger:              logger,
	}

	if perfCPUsFlag != nil && *perfCPUsFlag != "" {
		cpus, err := perfCPUFlagToCPUs(*perfCPUsFlag)
		if err != nil {
			return nil, err
		}
		for _, cpu := range cpus {
			// Use -1 to profile all processes on the CPU, see:
			// man perf_event_open
			hwProf := perf.NewHardwareProfiler(-1, cpu)
			if err := hwProf.Start(); err != nil {
				return nil, err
			}
			collector.perfHwProfilers[cpu] = &hwProf

			swProf := perf.NewSoftwareProfiler(-1, cpu)
			if err := swProf.Start(); err != nil {
				return nil, err
			}
			collector.perfSwProfilers[cpu] = &swProf

			cacheProf := perf.NewCacheProfiler(-1, cpu)
			if err := cacheProf.Start(); err != nil {
				return nil, err
			}
			collector.perfCacheProfilers[cpu] = &cacheProf
		}
	} else {
		for i := 0; i < runtime.NumCPU(); i++ {
			hwProf := perf.NewHardwareProfiler(-1, i)
			if err := hwProf.Start(); err != nil {
				return nil, err
			}
			collector.perfHwProfilers[i] = &hwProf

			swProf := perf.NewSoftwareProfiler(-1, i)
			if err := swProf.Start(); err != nil {
				return nil, err
			}
			collector.perfSwProfilers[i] = &swProf

			cacheProf := perf.NewCacheProfiler(-1, i)
			if err := cacheProf.Start(); err != nil {
				return nil, err
			}
			collector.perfCacheProfilers[i] = &cacheProf
		}
	}

	collector.desc = map[string]*prometheus.Desc{
		"cpucycles_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cpucycles_total",
			),
			"Number of CPU cycles (frequency scaled)",
			[]string{"cpu"},
			nil,
		),
		"instructions_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"instructions_total",
			),
			"Number of CPU instructions",
			[]string{"cpu"},
			nil,
		),
		"branch_instructions_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"branch_instructions_total",
			),
			"Number of CPU branch instructions",
			[]string{"cpu"},
			nil,
		),
		"branch_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"branch_misses_total",
			),
			"Number of CPU branch misses",
			[]string{"cpu"},
			nil,
		),
		"cache_refs_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_refs_total",
			),
			"Number of cache references (non frequency scaled)",
			[]string{"cpu"},
			nil,
		),
		"cache_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_misses_total",
			),
			"Number of cache misses",
			[]string{"cpu"},
			nil,
		),
		"ref_cpucycles_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"ref_cpucycles_total",
			),
			"Number of CPU cycles",
			[]string{"cpu"},
			nil,
		),
		"page_faults_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"page_faults_total",
			),
			"Number of page faults",
			[]string{"cpu"},
			nil,
		),
		"context_switches_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"context_switches_total",
			),
			"Number of context switches",
			[]string{"cpu"},
			nil,
		),
		"cpu_migrations_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cpu_migrations_total",
			),
			"Number of CPU process migrations",
			[]string{"cpu"},
			nil,
		),
		"minor_faults_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"minor_faults_total",
			),
			"Number of minor page faults",
			[]string{"cpu"},
			nil,
		),
		"major_faults_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"major_faults_total",
			),
			"Number of major page faults",
			[]string{"cpu"},
			nil,
		),
		"cache_l1d_read_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_l1d_read_hits_total",
			),
			"Number L1 data cache read hits",
			[]string{"cpu"},
			nil,
		),
		"cache_l1d_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_l1d_read_misses_total",
			),
			"Number L1 data cache read misses",
			[]string{"cpu"},
			nil,
		),
		"cache_l1d_write_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_l1d_write_hits_total",
			),
			"Number L1 data cache write hits",
			[]string{"cpu"},
			nil,
		),
		"cache_l1_instr_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_l1_instr_read_misses_total",
			),
			"Number instruction L1 instruction read misses",
			[]string{"cpu"},
			nil,
		),
		"cache_tlb_instr_read_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_instr_read_hits_total",
			),
			"Number instruction TLB read hits",
			[]string{"cpu"},
			nil,
		),
		"cache_tlb_instr_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_instr_read_misses_total",
			),
			"Number instruction TLB read misses",
			[]string{"cpu"},
			nil,
		),
		"cache_ll_read_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_ll_read_hits_total",
			),
			"Number last level read hits",
			[]string{"cpu"},
			nil,
		),
		"cache_ll_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_ll_read_misses_total",
			),
			"Number last level read misses",
			[]string{"cpu"},
			nil,
		),
		"cache_ll_write_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_ll_write_hits_total",
			),
			"Number last level write hits",
			[]string{"cpu"},
			nil,
		),
		"cache_ll_write_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_ll_write_misses_total",
			),
			"Number last level write misses",
			[]string{"cpu"},
			nil,
		),
		"cache_bpu_read_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_bpu_read_hits_total",
			),
			"Number BPU read hits",
			[]string{"cpu"},
			nil,
		),
		"cache_bpu_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_bpu_read_misses_total",
			),
			"Number BPU read misses",
			[]string{"cpu"},
			nil,
		),
	}

	return collector, nil
}

// Update implements the Collector interface and will collect metrics per CPU.
func (c *perfCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateHardwareStats(ch); err != nil {
		return err
	}

	if err := c.updateSoftwareStats(ch); err != nil {
		return err
	}

	if err := c.updateCacheStats(ch); err != nil {
		return err
	}

	return nil
}

func (c *perfCollector) updateHardwareStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfHwProfilers {
		cpuid := c.hwProfilerCPUMap[profiler]
		cpuStr := fmt.Sprintf("%d", cpuid)
		hwProfile, err := (*profiler).Profile()
		if err != nil {
			return err
		}
		if hwProfile == nil {
			continue
		}

		if hwProfile.CPUCycles != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cpucycles_total"],
				prometheus.CounterValue, float64(*hwProfile.CPUCycles),
				cpuStr,
			)
		}

		if hwProfile.Instructions != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["instructions_total"],
				prometheus.CounterValue, float64(*hwProfile.Instructions),
				cpuStr,
			)
		}

		if hwProfile.BranchInstr != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["branch_instructions_total"],
				prometheus.CounterValue, float64(*hwProfile.BranchInstr),
				cpuStr,
			)
		}

		if hwProfile.BranchMisses != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["branch_misses_total"],
				prometheus.CounterValue, float64(*hwProfile.BranchMisses),
				cpuStr,
			)
		}

		if hwProfile.CacheRefs != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_refs_total"],
				prometheus.CounterValue, float64(*hwProfile.CacheRefs),
				cpuStr,
			)
		}

		if hwProfile.CacheMisses != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_misses_total"],
				prometheus.CounterValue, float64(*hwProfile.CacheMisses),
				cpuStr,
			)
		}

		if hwProfile.RefCPUCycles != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["ref_cpucycles_total"],
				prometheus.CounterValue, float64(*hwProfile.RefCPUCycles),
				cpuStr,
			)
		}
	}

	return nil
}

func (c *perfCollector) updateSoftwareStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfSwProfilers {
		cpuid := c.swProfilerCPUMap[profiler]
		cpuStr := fmt.Sprintf("%d", cpuid)
		swProfile, err := (*profiler).Profile()
		if err != nil {
			return err
		}
		if swProfile == nil {
			continue
		}

		if swProfile.PageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["page_faults_total"],
				prometheus.CounterValue, float64(*swProfile.PageFaults),
				cpuStr,
			)
		}

		if swProfile.ContextSwitches != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["context_switches_total"],
				prometheus.CounterValue, float64(*swProfile.ContextSwitches),
				cpuStr,
			)
		}

		if swProfile.CPUMigrations != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cpu_migrations_total"],
				prometheus.CounterValue, float64(*swProfile.CPUMigrations),
				cpuStr,
			)
		}

		if swProfile.MinorPageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["minor_faults_total"],
				prometheus.CounterValue, float64(*swProfile.MinorPageFaults),
				cpuStr,
			)
		}

		if swProfile.MajorPageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["major_faults_total"],
				prometheus.CounterValue, float64(*swProfile.MajorPageFaults),
				cpuStr,
			)
		}
	}

	return nil
}

func (c *perfCollector) updateCacheStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfCacheProfilers {
		cpuid := c.cacheProfilerCPUMap[profiler]
		cpuStr := fmt.Sprintf("%d", cpuid)
		cacheProfile, err := (*profiler).Profile()
		if err != nil {
			return err
		}
		if cacheProfile == nil {
			continue
		}

		if cacheProfile.L1DataReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataReadHit),
				cpuStr,
			)
		}

		if cacheProfile.L1DataReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataReadMiss),
				cpuStr,
			)
		}

		if cacheProfile.L1DataWriteHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_write_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataWriteHit),
				cpuStr,
			)
		}

		if cacheProfile.L1InstrReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1_instr_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1InstrReadMiss),
				cpuStr,
			)
		}

		if cacheProfile.InstrTLBReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_instr_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.InstrTLBReadHit),
				cpuStr,
			)
		}

		if cacheProfile.InstrTLBReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_instr_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.InstrTLBReadMiss),
				cpuStr,
			)
		}

		if cacheProfile.LastLevelReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelReadHit),
				cpuStr,
			)
		}

		if cacheProfile.LastLevelReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelReadMiss),
				cpuStr,
			)
		}

		if cacheProfile.LastLevelWriteHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_write_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelWriteHit),
				cpuStr,
			)
		}

		if cacheProfile.LastLevelWriteMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_write_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelWriteMiss),
				cpuStr,
			)
		}

		if cacheProfile.BPUReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_bpu_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.BPUReadHit),
				cpuStr,
			)
		}

		if cacheProfile.BPUReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_bpu_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.BPUReadMiss),
				cpuStr,
			)
		}
	}

	return nil
}
