// +build linux

package perf

import (
	"encoding/binary"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	// EventAttrSize is the size of a PerfEventAttr
	EventAttrSize = uint32(unsafe.Sizeof(unix.PerfEventAttr{}))
)

// profileFn is a helper function to profile a function.
func profileFn(eventAttr *unix.PerfEventAttr, f func() error) (*ProfileValue, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	fd, err := unix.PerfEventOpen(
		eventAttr,
		unix.Gettid(),
		-1,
		-1,
		0,
	)
	if err != nil {
		return nil, err
	}
	if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_RESET, 0); err != nil {
		return nil, err
	}
	if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_ENABLE, 0); err != nil {
		return nil, err
	}
	if err := f(); err != nil {
		return nil, err
	}
	if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_DISABLE, 0); err != nil {
		return nil, err
	}
	buf := make([]byte, 24)
	if _, err := syscall.Read(fd, buf); err != nil {
		return nil, err
	}
	return &ProfileValue{
		Value:       binary.LittleEndian.Uint64(buf[0:8]),
		TimeEnabled: binary.LittleEndian.Uint64(buf[8:16]),
		TimeRunning: binary.LittleEndian.Uint64(buf[16:24]),
	}, unix.Close(fd)
}

// CPUInstructions is used to profile a function and return the number of CPU instructions.
// Note that it will call runtime.LockOSThread to ensure accurate profilng.
func CPUInstructions(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_INSTRUCTIONS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPUInstructionsEventAttr returns a unix.PerfEventAttr configured for CPUInstructions.
func CPUInstructionsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_INSTRUCTIONS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// CPUCycles is used to profile a function and return the number of CPU cycles.
// Note that it will call runtime.LockOSThread to ensure accurate profilng.
func CPUCycles(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CPU_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPUCyclesEventAttr returns a unix.PerfEventAttr configured for CPUCycles.
func CPUCyclesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CPU_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// CacheRef is used to profile a function and return the number of cache
// references. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func CacheRef(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CACHE_REFERENCES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CacheRefEventAttr returns a unix.PerfEventAttr configured for CacheRef.
func CacheRefEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CACHE_REFERENCES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// CacheMiss is used to profile a function and return the number of cache
// misses. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func CacheMiss(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CACHE_MISSES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CacheMissEventAttr returns a unix.PerfEventAttr configured for CacheMisses.
func CacheMissEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_CACHE_MISSES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// BusCycles is used to profile a function and return the number of bus
// cycles. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func BusCycles(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_BUS_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// BusCyclesEventAttr returns a unix.PerfEventAttr configured for BusCycles.
func BusCyclesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_BUS_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// StalledFrontendCycles is used to profile a function and return the number of
// stalled frontend cycles. Note that it will call runtime.LockOSThread to
// ensure accurate profilng.
func StalledFrontendCycles(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_STALLED_CYCLES_FRONTEND,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// StalledFrontendCyclesEventAttr returns a unix.PerfEventAttr configured for StalledFrontendCycles.
func StalledFrontendCyclesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_STALLED_CYCLES_FRONTEND,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// StalledBackendCycles is used to profile a function and return the number of
// stalled backend cycles. Note that it will call runtime.LockOSThread to
// ensure accurate profilng.
func StalledBackendCycles(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_STALLED_CYCLES_BACKEND,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// StalledBackendCyclesEventAttr returns a unix.PerfEventAttr configured for StalledBackendCycles.
func StalledBackendCyclesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_STALLED_CYCLES_BACKEND,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// CPURefCycles is used to profile a function and return the number of CPU
// references cycles which are not affected by frequency scaling. Note that it
// will call runtime.LockOSThread to ensure accurate profilng.
func CPURefCycles(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_REF_CPU_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPURefCyclesEventAttr returns a unix.PerfEventAttr configured for CPURefCycles.
func CPURefCyclesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HARDWARE,
		Config:      unix.PERF_COUNT_HW_REF_CPU_CYCLES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// CPUClock is used to profile a function and return the CPU clock timer. Note
// that it will call runtime.LockOSThread to ensure accurate profilng.
func CPUClock(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CPU_CLOCK,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPUClockEventAttr returns a unix.PerfEventAttr configured for CPUClock.
func CPUClockEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CPU_CLOCK,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// CPUTaskClock is used to profile a function and return the CPU clock timer
// for the running task. Note that it will call runtime.LockOSThread to ensure
// accurate profilng.
func CPUTaskClock(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_TASK_CLOCK,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPUTaskClockEventAttr returns a unix.PerfEventAttr configured for CPUTaskClock.
func CPUTaskClockEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_TASK_CLOCK,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// PageFaults is used to profile a function and return the number of page
// faults. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func PageFaults(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// PageFaultsEventAttr returns a unix.PerfEventAttr configured for PageFaults.
func PageFaultsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// ContextSwitches is used to profile a function and return the number of
// context switches. Note that it will call runtime.LockOSThread to ensure
// accurate profilng.
func ContextSwitches(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CONTEXT_SWITCHES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// ContextSwitchesEventAttr returns a unix.PerfEventAttr configured for ContextSwitches.
func ContextSwitchesEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CONTEXT_SWITCHES,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// CPUMigrations is used to profile a function and return the number of times
// the thread has been migrated to a new CPU. Note that it will call
// runtime.LockOSThread to ensure accurate profilng.
func CPUMigrations(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CPU_MIGRATIONS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// CPUMigrationsEventAttr returns a unix.PerfEventAttr configured for CPUMigrations.
func CPUMigrationsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_CPU_MIGRATIONS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// MinorPageFaults is used to profile a function and return the number of minor
// page faults. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func MinorPageFaults(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS_MIN,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// MinorPageFaultsEventAttr returns a unix.PerfEventAttr configured for MinorPageFaults.
func MinorPageFaultsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS_MIN,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// MajorPageFaults is used to profile a function and return the number of major
// page faults. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func MajorPageFaults(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS_MAJ,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// MajorPageFaultsEventAttr returns a unix.PerfEventAttr configured for MajorPageFaults.
func MajorPageFaultsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_PAGE_FAULTS_MAJ,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// AlignmentFaults is used to profile a function and return the number of alignment
// faults. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func AlignmentFaults(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_ALIGNMENT_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// AlignmentFaultsEventAttr returns a unix.PerfEventAttr configured for AlignmentFaults.
func AlignmentFaultsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_ALIGNMENT_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// EmulationFaults is used to profile a function and return the number of emulation
// faults. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func EmulationFaults(f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_EMULATION_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// EmulationFaultsEventAttr returns a unix.PerfEventAttr configured for EmulationFaults.
func EmulationFaultsEventAttr() unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_SOFTWARE,
		Config:      unix.PERF_COUNT_SW_EMULATION_FAULTS,
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// L1Data is used to profile a function and the L1 data cache faults. Use
// PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func L1Data(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_L1D) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// L1DataEventAttr returns a unix.PerfEventAttr configured for L1Data.
func L1DataEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_L1D) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// L1Instructions is used to profile a function for the instruction level L1
// cache. Use PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func L1Instructions(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_L1I) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// L1InstructionsEventAttr returns a unix.PerfEventAttr configured for L1Instructions.
func L1InstructionsEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_L1I) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// LLCache is used to profile a function and return the number of emulation
// PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func LLCache(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_LL) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// LLCacheEventAttr returns a unix.PerfEventAttr configured for LLCache.
func LLCacheEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_LL) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// DataTLB is used to profile the data TLB. Use PERF_COUNT_HW_CACHE_OP_READ,
// PERF_COUNT_HW_CACHE_OP_WRITE, or PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt
// and PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for
// the result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func DataTLB(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_DTLB) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// DataTLBEventAttr returns a unix.PerfEventAttr configured for DataTLB.
func DataTLBEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_DTLB) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// InstructionTLB is used to profile the instruction TLB. Use
// PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func InstructionTLB(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_ITLB) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// InstructionTLBEventAttr returns a unix.PerfEventAttr configured for InstructionTLB.
func InstructionTLBEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_ITLB) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}

}

// BPU is used to profile a function for the Branch Predictor Unit.
// Use PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func BPU(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_BPU) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// BPUEventAttr returns a unix.PerfEventAttr configured for BPU events.
func BPUEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_BPU) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}

// NodeCache is used to profile a function for NUMA operations. Use Use
// PERF_COUNT_HW_CACHE_OP_READ, PERF_COUNT_HW_CACHE_OP_WRITE, or
// PERF_COUNT_HW_CACHE_OP_PREFETCH for the opt and
// PERF_COUNT_HW_CACHE_RESULT_ACCESS or PERF_COUNT_HW_CACHE_RESULT_MISS for the
// result. Note that it will call runtime.LockOSThread to ensure accurate
// profilng.
func NodeCache(op, result int, f func() error) (*ProfileValue, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_NODE) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
	return profileFn(eventAttr, f)
}

// NodeCacheEventAttr returns a unix.PerfEventAttr configured for NUMA cache operations.
func NodeCacheEventAttr(op, result int) unix.PerfEventAttr {
	return unix.PerfEventAttr{
		Type:        unix.PERF_TYPE_HW_CACHE,
		Config:      uint64((unix.PERF_COUNT_HW_CACHE_NODE) | (op << 8) | (result << 16)),
		Size:        EventAttrSize,
		Bits:        unix.PerfBitExcludeKernel | unix.PerfBitExcludeHv,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
	}
}
