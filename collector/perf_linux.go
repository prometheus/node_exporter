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

//go:build !noperf
// +build !noperf

package collector

import (
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hodgesds/perf-utils"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	perfSubsystem = "perf"
)

var (
	perfCPUsFlag       = kingpin.Flag("collector.perf.cpus", "List of CPUs from which perf metrics should be collected").Default("").String()
	perfTracepointFlag = kingpin.Flag("collector.perf.tracepoint", "perf tracepoint that should be collected").Strings()
	perfNoHwProfiler   = kingpin.Flag("collector.perf.disable-hardware-profilers", "disable perf hardware profilers").Default("false").Bool()
	perfHwProfilerFlag = kingpin.Flag("collector.perf.hardware-profilers", "perf hardware profilers that should be collected").Strings()
	perfNoSwProfiler   = kingpin.Flag("collector.perf.disable-software-profilers", "disable perf software profilers").Default("false").Bool()
	perfSwProfilerFlag = kingpin.Flag("collector.perf.software-profilers", "perf software profilers that should be collected").Strings()
	perfNoCaProfiler   = kingpin.Flag("collector.perf.disable-cache-profilers", "disable perf cache profilers").Default("false").Bool()
	perfCaProfilerFlag = kingpin.Flag("collector.perf.cache-profilers", "perf cache profilers that should be collected").Strings()
)

func init() {
	registerCollector(perfSubsystem, defaultDisabled, NewPerfCollector)
}

var (
	perfHardwareProfilerMap = map[string]perf.HardwareProfilerType{
		"CpuCycles":             perf.CpuCyclesProfiler,
		"CpuInstr":              perf.CpuInstrProfiler,
		"CacheRef":              perf.CacheRefProfiler,
		"CacheMisses":           perf.CacheMissesProfiler,
		"BranchInstr":           perf.BranchInstrProfiler,
		"BranchMisses":          perf.BranchMissesProfiler,
		"StalledCyclesBackend":  perf.StalledCyclesBackendProfiler,
		"StalledCyclesFrontend": perf.StalledCyclesFrontendProfiler,
		"RefCpuCycles":          perf.RefCpuCyclesProfiler,
		// "BusCycles":             perf.BusCyclesProfiler,
	}
	perfSoftwareProfilerMap = map[string]perf.SoftwareProfilerType{
		"PageFault":     perf.PageFaultProfiler,
		"ContextSwitch": perf.ContextSwitchProfiler,
		"CpuMigration":  perf.CpuMigrationProfiler,
		"MinorFault":    perf.MinorFaultProfiler,
		"MajorFault":    perf.MajorFaultProfiler,
		// "CpuClock":      perf.CpuClockProfiler,
		// "TaskClock":     perf.TaskClockProfiler,
		// "AlignFault":    perf.AlignFaultProfiler,
		// "EmuFault":      perf.EmuFaultProfiler,
	}
	perfCacheProfilerMap = map[string]perf.CacheProfilerType{
		"L1DataReadHit":    perf.L1DataReadHitProfiler,
		"L1DataReadMiss":   perf.L1DataReadMissProfiler,
		"L1DataWriteHit":   perf.L1DataWriteHitProfiler,
		"L1InstrReadMiss":  perf.L1InstrReadMissProfiler,
		"LLReadHit":        perf.LLReadHitProfiler,
		"LLReadMiss":       perf.LLReadMissProfiler,
		"LLWriteHit":       perf.LLWriteHitProfiler,
		"LLWriteMiss":      perf.LLWriteMissProfiler,
		"InstrTLBReadHit":  perf.InstrTLBReadHitProfiler,
		"InstrTLBReadMiss": perf.InstrTLBReadMissProfiler,
		"BPUReadHit":       perf.BPUReadHitProfiler,
		"BPUReadMiss":      perf.BPUReadMissProfiler,
		// "L1InstrReadHit":     perf.L1InstrReadHitProfiler,
		"DataTLBReadHit":   perf.DataTLBReadHitProfiler,
		"DataTLBReadMiss":  perf.DataTLBReadMissProfiler,
		"DataTLBWriteHit":  perf.DataTLBWriteHitProfiler,
		"DataTLBWriteMiss": perf.DataTLBWriteMissProfiler,
		// "NodeCacheReadHit":   perf.NodeCacheReadHitProfiler,
		// "NodeCacheReadMiss":  perf.NodeCacheReadMissProfiler,
		// "NodeCacheWriteHit":  perf.NodeCacheWriteHitProfiler,
		// "NodeCacheWriteMiss": perf.NodeCacheWriteMissProfiler,
	}
)

// perfTracepointFlagToTracepoints returns the set of configured tracepoints.
func perfTracepointFlagToTracepoints(tracepointsFlag []string) ([]*perfTracepoint, error) {
	tracepoints := make([]*perfTracepoint, len(tracepointsFlag))

	for i, tracepoint := range tracepointsFlag {
		split := strings.Split(tracepoint, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid tracepoint config %v", tracepoint)
		}
		tracepoints[i] = &perfTracepoint{
			subsystem: split[0],
			event:     split[1],
		}
	}
	return tracepoints, nil
}

// perfCPUFlagToCPUs returns a set of CPUs for the perf collectors to monitor.
func perfCPUFlagToCPUs(cpuFlag string) ([]int, error) {
	var err error
	cpus := []int{}
	for subset := range strings.SplitSeq(cpuFlag, ",") {
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

// perfTracepoint is a struct for holding tracepoint information.
type perfTracepoint struct {
	subsystem string
	event     string
}

// label returns the tracepoint name in the format of subsystem_tracepoint.
func (t *perfTracepoint) label() string {
	return t.subsystem + "_" + t.event
}

// tracepoint returns the tracepoint name in the format of subsystem:tracepoint.
func (t *perfTracepoint) tracepoint() string {
	return t.subsystem + ":" + t.event
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
	logger              *slog.Logger
	tracepointCollector *perfTracepointCollector
}

type perfTracepointCollector struct {
	// desc is the mapping of subsystem to tracepoint *prometheus.Desc.
	descs map[string]map[string]*prometheus.Desc
	// collection order is the sorted configured collection order of the profiler.
	collectionOrder []string

	logger    *slog.Logger
	profilers map[int]perf.GroupProfiler
}

// update is used collect all tracepoints across all tracepoint profilers.
func (c *perfTracepointCollector) update(ch chan<- prometheus.Metric) error {
	for cpu := range c.profilers {
		if err := c.updateCPU(cpu, ch); err != nil {
			return err
		}
	}
	return nil
}

// updateCPU is used to update metrics per CPU profiler.
func (c *perfTracepointCollector) updateCPU(cpu int, ch chan<- prometheus.Metric) error {
	profiler := c.profilers[cpu]
	p := &perf.GroupProfileValue{}
	if err := profiler.Profile(p); err != nil {
		c.logger.Error("Failed to collect tracepoint profile", "err", err)
		return err
	}

	cpuid := strconv.Itoa(cpu)

	for i, value := range p.Values {
		// Get the Desc from the ordered group value.
		descKey := c.collectionOrder[i]
		descKeySlice := strings.Split(descKey, ":")
		ch <- prometheus.MustNewConstMetric(
			c.descs[descKeySlice[0]][descKeySlice[1]],
			prometheus.CounterValue,
			float64(value),
			cpuid,
		)
	}
	return nil
}

// newPerfTracepointCollector returns a configured perfTracepointCollector.
func newPerfTracepointCollector(
	logger *slog.Logger,
	tracepointsFlag []string,
	cpus []int,
) (*perfTracepointCollector, error) {
	tracepoints, err := perfTracepointFlagToTracepoints(tracepointsFlag)
	if err != nil {
		return nil, err
	}

	collectionOrder := make([]string, len(tracepoints))
	descs := map[string]map[string]*prometheus.Desc{}
	eventAttrs := make([]unix.PerfEventAttr, len(tracepoints))

	for i, tracepoint := range tracepoints {
		eventAttr, err := perf.TracepointEventAttr(tracepoint.subsystem, tracepoint.event)
		if err != nil {
			return nil, err
		}
		eventAttrs[i] = *eventAttr
		collectionOrder[i] = tracepoint.tracepoint()
		if _, ok := descs[tracepoint.subsystem]; !ok {
			descs[tracepoint.subsystem] = map[string]*prometheus.Desc{}
		}
		descs[tracepoint.subsystem][tracepoint.event] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				tracepoint.label(),
			),
			"Perf tracepoint "+tracepoint.tracepoint(),
			[]string{"cpu"},
			nil,
		)
	}

	profilers := make(map[int]perf.GroupProfiler, len(cpus))
	for _, cpu := range cpus {
		profiler, err := perf.NewGroupProfiler(-1, cpu, 0, eventAttrs...)
		if err != nil {
			return nil, err
		}
		profilers[cpu] = profiler
	}

	c := &perfTracepointCollector{
		descs:           descs,
		collectionOrder: collectionOrder,
		profilers:       profilers,
		logger:          logger,
	}

	for _, profiler := range c.profilers {
		if err := profiler.Start(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// NewPerfCollector returns a new perf based collector, it creates a profiler
// per CPU.
func NewPerfCollector(logger *slog.Logger) (Collector, error) {
	collector := &perfCollector{
		perfHwProfilers:     map[int]*perf.HardwareProfiler{},
		perfSwProfilers:     map[int]*perf.SoftwareProfiler{},
		perfCacheProfilers:  map[int]*perf.CacheProfiler{},
		hwProfilerCPUMap:    map[*perf.HardwareProfiler]int{},
		swProfilerCPUMap:    map[*perf.SoftwareProfiler]int{},
		cacheProfilerCPUMap: map[*perf.CacheProfiler]int{},
		logger:              logger,
	}

	var (
		cpus []int
		err  error
	)
	if perfCPUsFlag != nil && *perfCPUsFlag != "" {
		cpus, err = perfCPUFlagToCPUs(*perfCPUsFlag)
		if err != nil {
			return nil, err
		}
	} else {
		cpus = make([]int, runtime.NumCPU())
		for i := range cpus {
			cpus[i] = i
		}
	}

	// First configure any tracepoints.
	if *perfTracepointFlag != nil && len(*perfTracepointFlag) > 0 {
		tracepointCollector, err := newPerfTracepointCollector(logger, *perfTracepointFlag, cpus)
		if err != nil {
			return nil, err
		}
		collector.tracepointCollector = tracepointCollector
	}

	// Configure perf profilers
	hardwareProfilers := perf.AllHardwareProfilers
	if *perfHwProfilerFlag != nil && len(*perfHwProfilerFlag) > 0 {
		// hardwareProfilers = 0
		for _, hf := range *perfHwProfilerFlag {
			if v, ok := perfHardwareProfilerMap[hf]; ok {
				hardwareProfilers |= v
			}
		}
	}
	softwareProfilers := perf.AllSoftwareProfilers
	if *perfSwProfilerFlag != nil && len(*perfSwProfilerFlag) > 0 {
		// softwareProfilers = 0
		for _, sf := range *perfSwProfilerFlag {
			if v, ok := perfSoftwareProfilerMap[sf]; ok {
				softwareProfilers |= v
			}
		}
	}
	cacheProfilers := perf.L1DataReadHitProfiler | perf.L1DataReadMissProfiler | perf.L1DataWriteHitProfiler | perf.L1InstrReadMissProfiler | perf.InstrTLBReadHitProfiler | perf.InstrTLBReadMissProfiler | perf.DataTLBReadHitProfiler | perf.DataTLBReadMissProfiler | perf.DataTLBWriteHitProfiler | perf.DataTLBWriteMissProfiler | perf.LLReadHitProfiler | perf.LLReadMissProfiler | perf.LLWriteHitProfiler | perf.LLWriteMissProfiler | perf.BPUReadHitProfiler | perf.BPUReadMissProfiler
	if *perfCaProfilerFlag != nil && len(*perfCaProfilerFlag) > 0 {
		cacheProfilers = 0
		for _, cf := range *perfCaProfilerFlag {
			if v, ok := perfCacheProfilerMap[cf]; ok {
				cacheProfilers |= v
			}
		}
	}

	// Configure all profilers for the specified CPUs.
	for _, cpu := range cpus {
		// Use -1 to profile all processes on the CPU, see:
		// man perf_event_open
		if !*perfNoHwProfiler {
			hwProf, err := perf.NewHardwareProfiler(
				-1,
				cpu,
				hardwareProfilers,
			)
			if err != nil && !hwProf.HasProfilers() {
				return nil, err
			}
			if err := hwProf.Start(); err != nil {
				return nil, err
			}
			collector.perfHwProfilers[cpu] = &hwProf
			collector.hwProfilerCPUMap[&hwProf] = cpu
		}

		if !*perfNoSwProfiler {
			swProf, err := perf.NewSoftwareProfiler(-1, cpu, softwareProfilers)
			if err != nil && !swProf.HasProfilers() {
				return nil, err
			}
			if err := swProf.Start(); err != nil {
				return nil, err
			}
			collector.perfSwProfilers[cpu] = &swProf
			collector.swProfilerCPUMap[&swProf] = cpu
		}

		if !*perfNoCaProfiler {
			cacheProf, err := perf.NewCacheProfiler(
				-1,
				cpu,
				cacheProfilers,
			)
			if err != nil && !cacheProf.HasProfilers() {
				return nil, err
			}
			if err := cacheProf.Start(); err != nil {
				return nil, err
			}
			collector.perfCacheProfilers[cpu] = &cacheProf
			collector.cacheProfilerCPUMap[&cacheProf] = cpu
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
		"stalled_cycles_backend_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"stalled_cycles_backend_total",
			),
			"Number of stalled backend CPU cycles",
			[]string{"cpu"},
			nil,
		),
		"stalled_cycles_frontend_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"stalled_cycles_frontend_total",
			),
			"Number of stalled frontend CPU cycles",
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
		"cache_tlb_data_read_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_data_read_hits_total",
			),
			"Number of data TLB read hits",
			[]string{"cpu"},
			nil,
		),
		"cache_tlb_data_read_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_data_read_misses_total",
			),
			"Number of data TLB read misses",
			[]string{"cpu"},
			nil,
		),
		"cache_tlb_data_write_hits_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_data_write_hits_total",
			),
			"Number of data TLB write hits",
			[]string{"cpu"},
			nil,
		),
		"cache_tlb_data_write_misses_total": prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				perfSubsystem,
				"cache_tlb_data_write_misses_total",
			),
			"Number of data TLB write misses",
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
	if c.tracepointCollector != nil {
		return c.tracepointCollector.update(ch)
	}

	return nil
}

func (c *perfCollector) updateHardwareStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfHwProfilers {
		hwProfile := &perf.HardwareProfile{}
		if err := (*profiler).Profile(hwProfile); err != nil {
			return err
		}

		cpuid := strconv.Itoa(c.hwProfilerCPUMap[profiler])

		if hwProfile.CPUCycles != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cpucycles_total"],
				prometheus.CounterValue, float64(*hwProfile.CPUCycles),
				cpuid,
			)
		}

		if hwProfile.Instructions != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["instructions_total"],
				prometheus.CounterValue, float64(*hwProfile.Instructions),
				cpuid,
			)
		}

		if hwProfile.BranchInstr != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["branch_instructions_total"],
				prometheus.CounterValue, float64(*hwProfile.BranchInstr),
				cpuid,
			)
		}

		if hwProfile.BranchMisses != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["branch_misses_total"],
				prometheus.CounterValue, float64(*hwProfile.BranchMisses),
				cpuid,
			)
		}

		if hwProfile.CacheRefs != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_refs_total"],
				prometheus.CounterValue, float64(*hwProfile.CacheRefs),
				cpuid,
			)
		}

		if hwProfile.CacheMisses != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_misses_total"],
				prometheus.CounterValue, float64(*hwProfile.CacheMisses),
				cpuid,
			)
		}

		if hwProfile.RefCPUCycles != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["ref_cpucycles_total"],
				prometheus.CounterValue, float64(*hwProfile.RefCPUCycles),
				cpuid,
			)
		}

		if hwProfile.StalledCyclesBackend != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["stalled_cycles_backend_total"],
				prometheus.CounterValue, float64(*hwProfile.StalledCyclesBackend),
				cpuid,
			)
		}

		if hwProfile.StalledCyclesFrontend != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["stalled_cycles_frontend_total"],
				prometheus.CounterValue, float64(*hwProfile.StalledCyclesFrontend),
				cpuid,
			)
		}
	}

	return nil
}

func (c *perfCollector) updateSoftwareStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfSwProfilers {
		swProfile := &perf.SoftwareProfile{}
		if err := (*profiler).Profile(swProfile); err != nil {
			return err
		}

		cpuid := strconv.Itoa(c.swProfilerCPUMap[profiler])

		if swProfile.PageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["page_faults_total"],
				prometheus.CounterValue, float64(*swProfile.PageFaults),
				cpuid,
			)
		}

		if swProfile.ContextSwitches != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["context_switches_total"],
				prometheus.CounterValue, float64(*swProfile.ContextSwitches),
				cpuid,
			)
		}

		if swProfile.CPUMigrations != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cpu_migrations_total"],
				prometheus.CounterValue, float64(*swProfile.CPUMigrations),
				cpuid,
			)
		}

		if swProfile.MinorPageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["minor_faults_total"],
				prometheus.CounterValue, float64(*swProfile.MinorPageFaults),
				cpuid,
			)
		}

		if swProfile.MajorPageFaults != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["major_faults_total"],
				prometheus.CounterValue, float64(*swProfile.MajorPageFaults),
				cpuid,
			)
		}
	}

	return nil
}

func (c *perfCollector) updateCacheStats(ch chan<- prometheus.Metric) error {
	for _, profiler := range c.perfCacheProfilers {
		cacheProfile := &perf.CacheProfile{}
		if err := (*profiler).Profile(cacheProfile); err != nil {
			return err
		}

		cpuid := strconv.Itoa(c.cacheProfilerCPUMap[profiler])

		if cacheProfile.L1DataReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataReadHit),
				cpuid,
			)
		}

		if cacheProfile.L1DataReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataReadMiss),
				cpuid,
			)
		}

		if cacheProfile.L1DataWriteHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1d_write_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1DataWriteHit),
				cpuid,
			)
		}

		if cacheProfile.L1InstrReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_l1_instr_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.L1InstrReadMiss),
				cpuid,
			)
		}

		if cacheProfile.InstrTLBReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_instr_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.InstrTLBReadHit),
				cpuid,
			)
		}

		if cacheProfile.InstrTLBReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_instr_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.InstrTLBReadMiss),
				cpuid,
			)
		}

		if cacheProfile.DataTLBReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_data_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.DataTLBReadHit),
				cpuid,
			)
		}

		if cacheProfile.DataTLBReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_data_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.DataTLBReadMiss),
				cpuid,
			)
		}

		if cacheProfile.DataTLBWriteHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_data_write_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.DataTLBWriteHit),
				cpuid,
			)
		}

		if cacheProfile.DataTLBWriteMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_tlb_data_write_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.DataTLBWriteMiss),
				cpuid,
			)
		}

		if cacheProfile.LastLevelReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelReadHit),
				cpuid,
			)
		}

		if cacheProfile.LastLevelReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelReadMiss),
				cpuid,
			)
		}

		if cacheProfile.LastLevelWriteHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_write_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelWriteHit),
				cpuid,
			)
		}

		if cacheProfile.LastLevelWriteMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_ll_write_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.LastLevelWriteMiss),
				cpuid,
			)
		}

		if cacheProfile.BPUReadHit != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_bpu_read_hits_total"],
				prometheus.CounterValue, float64(*cacheProfile.BPUReadHit),
				cpuid,
			)
		}

		if cacheProfile.BPUReadMiss != nil {
			ch <- prometheus.MustNewConstMetric(
				c.desc["cache_bpu_read_misses_total"],
				prometheus.CounterValue, float64(*cacheProfile.BPUReadMiss),
				cpuid,
			)
		}
	}

	return nil
}
