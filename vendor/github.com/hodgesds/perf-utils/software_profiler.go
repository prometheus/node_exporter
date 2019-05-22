// +build linux

package perf

import (
	"go.uber.org/multierr"
	"golang.org/x/sys/unix"
)

type softwareProfiler struct {
	// map of perf counter type to file descriptor
	profilers map[int]Profiler
}

// NewSoftwareProfiler returns a new software profiler.
func NewSoftwareProfiler(pid, cpu int, opts ...int) SoftwareProfiler {
	profilers := map[int]Profiler{}

	cpuClockProfiler, err := NewCPUClockProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_CPU_CLOCK] = cpuClockProfiler
	}

	taskClockProfiler, err := NewTaskClockProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_TASK_CLOCK] = taskClockProfiler
	}

	pageFaultProfiler, err := NewPageFaultProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_PAGE_FAULTS] = pageFaultProfiler
	}

	ctxSwitchesProfiler, err := NewCtxSwitchesProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_CONTEXT_SWITCHES] = ctxSwitchesProfiler
	}

	cpuMigrationsProfiler, err := NewCPUMigrationsProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_CPU_MIGRATIONS] = cpuMigrationsProfiler
	}

	minorFaultProfiler, err := NewMinorFaultsProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_PAGE_FAULTS_MIN] = minorFaultProfiler
	}

	majorFaultProfiler, err := NewMajorFaultsProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_PAGE_FAULTS_MAJ] = majorFaultProfiler
	}

	alignFaultsFrontProfiler, err := NewAlignFaultsProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_ALIGNMENT_FAULTS] = alignFaultsFrontProfiler
	}

	emuFaultProfiler, err := NewEmulationFaultsProfiler(pid, cpu, opts...)
	if err == nil {
		profilers[unix.PERF_COUNT_SW_EMULATION_FAULTS] = emuFaultProfiler
	}

	return &softwareProfiler{
		profilers: profilers,
	}
}

// Start is used to start the SoftwareProfiler.
func (p *softwareProfiler) Start() error {
	if len(p.profilers) == 0 {
		return ErrNoProfiler
	}
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Start())
	}
	return err
}

// Reset is used to reset the SoftwareProfiler.
func (p *softwareProfiler) Reset() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Reset())
	}
	return err
}

// Stop is used to reset the SoftwareProfiler.
func (p *softwareProfiler) Stop() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Stop())
	}
	return err
}

// Close is used to reset the SoftwareProfiler.
func (p *softwareProfiler) Close() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Close())
	}
	return err
}

// Profile is used to read the SoftwareProfiler SoftwareProfile it returns an
// error only if all profiles fail.
func (p *softwareProfiler) Profile() (*SoftwareProfile, error) {
	var err error
	swProfile := &SoftwareProfile{}
	for profilerType, profiler := range p.profilers {
		profileVal, err2 := profiler.Profile()
		err = multierr.Append(err, err2)
		if err2 == nil {
			if swProfile.TimeEnabled == nil {
				swProfile.TimeEnabled = &profileVal.TimeEnabled
			}
			if swProfile.TimeRunning == nil {
				swProfile.TimeRunning = &profileVal.TimeRunning
			}
			switch profilerType {
			case unix.PERF_COUNT_SW_CPU_CLOCK:
				swProfile.CPUClock = &profileVal.Value
			case unix.PERF_COUNT_SW_TASK_CLOCK:
				swProfile.TaskClock = &profileVal.Value
			case unix.PERF_COUNT_SW_PAGE_FAULTS:
				swProfile.PageFaults = &profileVal.Value
			case unix.PERF_COUNT_SW_CONTEXT_SWITCHES:
				swProfile.ContextSwitches = &profileVal.Value
			case unix.PERF_COUNT_SW_CPU_MIGRATIONS:
				swProfile.CPUMigrations = &profileVal.Value
			case unix.PERF_COUNT_SW_PAGE_FAULTS_MIN:
				swProfile.MinorPageFaults = &profileVal.Value
			case unix.PERF_COUNT_SW_PAGE_FAULTS_MAJ:
				swProfile.MajorPageFaults = &profileVal.Value
			case unix.PERF_COUNT_SW_ALIGNMENT_FAULTS:
				swProfile.AlignmentFaults = &profileVal.Value
			case unix.PERF_COUNT_SW_EMULATION_FAULTS:
				swProfile.EmulationFaults = &profileVal.Value
			default:
			}
		}
	}
	if len(multierr.Errors(err)) == len(p.profilers) {
		return nil, err
	}

	return swProfile, nil
}
