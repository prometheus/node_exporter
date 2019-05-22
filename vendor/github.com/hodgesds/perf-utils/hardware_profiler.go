// +build linux

package perf

import (
	"go.uber.org/multierr"
	"golang.org/x/sys/unix"
)

type hardwareProfiler struct {
	// map of perf counter type to file descriptor
	profilers map[int]Profiler
}

// NewHardwareProfiler returns a new hardware profiler.
func NewHardwareProfiler(pid, cpu int, opts ...int) HardwareProfiler {
	profilers := map[int]Profiler{}

	cpuCycleProfiler, err := NewCPUCycleProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_CPU_CYCLES] = cpuCycleProfiler
	}

	instrProfiler, err := NewInstrProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_INSTRUCTIONS] = instrProfiler
	}

	cacheRefProfiler, err := NewCacheRefProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_CACHE_REFERENCES] = cacheRefProfiler
	}

	cacheMissesProfiler, err := NewCacheMissesProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_CACHE_MISSES] = cacheMissesProfiler
	}

	branchInstrProfiler, err := NewBranchInstrProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_BRANCH_INSTRUCTIONS] = branchInstrProfiler
	}

	branchMissesProfiler, err := NewBranchMissesProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_BRANCH_MISSES] = branchMissesProfiler
	}

	busCyclesProfiler, err := NewBusCyclesProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_BUS_CYCLES] = busCyclesProfiler
	}

	stalledCyclesFrontProfiler, err := NewStalledCyclesFrontProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_STALLED_CYCLES_FRONTEND] = stalledCyclesFrontProfiler
	}

	stalledCyclesBackProfiler, err := NewStalledCyclesBackProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_STALLED_CYCLES_BACKEND] = stalledCyclesBackProfiler
	}

	refCPUCyclesProfiler, err := NewRefCPUCyclesProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_HW_REF_CPU_CYCLES] = refCPUCyclesProfiler
	}

	return &hardwareProfiler{
		profilers: profilers,
	}
}

// Start is used to start the HardwareProfiler.
func (p *hardwareProfiler) Start() error {
	if len(p.profilers) == 0 {
		return ErrNoProfiler
	}
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Start())
	}
	return err
}

// Reset is used to reset the HardwareProfiler.
func (p *hardwareProfiler) Reset() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Reset())
	}
	return err
}

// Stop is used to reset the HardwareProfiler.
func (p *hardwareProfiler) Stop() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Stop())
	}
	return err
}

// Close is used to reset the HardwareProfiler.
func (p *hardwareProfiler) Close() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Close())
	}
	return err
}

// Profile is used to read the HardwareProfiler HardwareProfile it returns an
// error only if all profiles fail.
func (p *hardwareProfiler) Profile() (*HardwareProfile, error) {
	var err error
	hwProfile := &HardwareProfile{}
	for profilerType, profiler := range p.profilers {
		profileVal, err2 := profiler.Profile()
		err = multierr.Append(err, err2)
		if err2 == nil {
			if hwProfile.TimeEnabled == nil {
				hwProfile.TimeEnabled = &profileVal.TimeEnabled
			}
			if hwProfile.TimeRunning == nil {
				hwProfile.TimeRunning = &profileVal.TimeRunning
			}
			switch profilerType {
			case unix.PERF_COUNT_HW_CPU_CYCLES:
				hwProfile.CPUCycles = &profileVal.Value
			case unix.PERF_COUNT_HW_INSTRUCTIONS:
				hwProfile.Instructions = &profileVal.Value
			case unix.PERF_COUNT_HW_CACHE_REFERENCES:
				hwProfile.CacheRefs = &profileVal.Value
			case unix.PERF_COUNT_HW_CACHE_MISSES:
				hwProfile.CacheMisses = &profileVal.Value
			case unix.PERF_COUNT_HW_BRANCH_INSTRUCTIONS:
				hwProfile.BranchInstr = &profileVal.Value
			case unix.PERF_COUNT_HW_BRANCH_MISSES:
				hwProfile.BranchMisses = &profileVal.Value
			case unix.PERF_COUNT_HW_BUS_CYCLES:
				hwProfile.BusCycles = &profileVal.Value
			case unix.PERF_COUNT_HW_STALLED_CYCLES_FRONTEND:
				hwProfile.StalledCyclesFrontend = &profileVal.Value
			case unix.PERF_COUNT_HW_STALLED_CYCLES_BACKEND:
				hwProfile.StalledCyclesBackend = &profileVal.Value
			case unix.PERF_COUNT_HW_REF_CPU_CYCLES:
				hwProfile.RefCPUCycles = &profileVal.Value
			}
		}
	}
	if len(multierr.Errors(err)) == len(p.profilers) {
		return nil, err
	}

	return hwProfile, nil
}
