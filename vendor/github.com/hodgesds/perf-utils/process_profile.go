// +build linux

package perf

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// PERF_SAMPLE_IDENTIFIER is not defined in x/sys/unix.
	PERF_SAMPLE_IDENTIFIER = 1 << 16

	// PERF_IOC_FLAG_GROUP is not defined in x/sys/unix.
	PERF_IOC_FLAG_GROUP = 1 << 0
)

var (
	// ErrNoProfiler is returned when no profiler is available for profiling.
	ErrNoProfiler = fmt.Errorf("No profiler available")
)

// Profiler is a profiler.
type Profiler interface {
	Start() error
	Reset() error
	Stop() error
	Close() error
	Profile() (*ProfileValue, error)
}

// HardwareProfiler is a hardware profiler.
type HardwareProfiler interface {
	Start() error
	Reset() error
	Stop() error
	Close() error
	Profile() (*HardwareProfile, error)
}

// HardwareProfile is returned by a HardwareProfiler. Depending on kernel
// configuration some fields may return nil.
type HardwareProfile struct {
	CPUCycles             *uint64 `json:"cpu_cycles,omitempty"`
	Instructions          *uint64 `json:"instructions,omitempty"`
	CacheRefs             *uint64 `json:"cache_refs,omitempty"`
	CacheMisses           *uint64 `json:"cache_misses,omitempty"`
	BranchInstr           *uint64 `json:"branch_instr,omitempty"`
	BranchMisses          *uint64 `json:"branch_misses,omitempty"`
	BusCycles             *uint64 `json:"bus_cycles,omitempty"`
	StalledCyclesFrontend *uint64 `json:"stalled_cycles_frontend,omitempty"`
	StalledCyclesBackend  *uint64 `json:"stalled_cycles_backend,omitempty"`
	RefCPUCycles          *uint64 `json:"ref_cpu_cycles,omitempty"`
	TimeEnabled           *uint64 `json:"time_enabled,omitempty"`
	TimeRunning           *uint64 `json:"time_running,omitempty"`
}

// SoftwareProfiler is a software profiler.
type SoftwareProfiler interface {
	Start() error
	Reset() error
	Stop() error
	Close() error
	Profile() (*SoftwareProfile, error)
}

// SoftwareProfile is returned by a SoftwareProfiler.
type SoftwareProfile struct {
	CPUClock        *uint64 `json:"cpu_clock,omitempty"`
	TaskClock       *uint64 `json:"task_clock,omitempty"`
	PageFaults      *uint64 `json:"page_faults,omitempty"`
	ContextSwitches *uint64 `json:"context_switches,omitempty"`
	CPUMigrations   *uint64 `json:"cpu_migrations,omitempty"`
	MinorPageFaults *uint64 `json:"minor_page_faults,omitempty"`
	MajorPageFaults *uint64 `json:"major_page_faults,omitempty"`
	AlignmentFaults *uint64 `json:"alignment_faults,omitempty"`
	EmulationFaults *uint64 `json:"emulation_faults,omitempty"`
	TimeEnabled     *uint64 `json:"time_enabled,omitempty"`
	TimeRunning     *uint64 `json:"time_running,omitempty"`
}

// CacheProfiler is a cache profiler.
type CacheProfiler interface {
	Start() error
	Reset() error
	Stop() error
	Close() error
	Profile() (*CacheProfile, error)
}

// CacheProfile is returned by a CacheProfiler.
type CacheProfile struct {
	L1DataReadHit      *uint64 `json:"l1_data_read_hit,omitempty"`
	L1DataReadMiss     *uint64 `json:"l1_data_read_miss,omitempty"`
	L1DataWriteHit     *uint64 `json:"l1_data_write_hit,omitempty"`
	L1InstrReadMiss    *uint64 `json:"l1_instr_read_miss,omitempty"`
	LastLevelReadHit   *uint64 `json:"last_level_read_hit,omitempty"`
	LastLevelReadMiss  *uint64 `json:"last_level_read_miss,omitempty"`
	LastLevelWriteHit  *uint64 `json:"last_level_write_hit,omitempty"`
	LastLevelWriteMiss *uint64 `json:"last_level_write_miss,omitempty"`
	DataTLBReadHit     *uint64 `json:"data_tlb_read_hit,omitempty"`
	DataTLBReadMiss    *uint64 `json:"data_tlb_read_miss,omitempty"`
	DataTLBWriteHit    *uint64 `json:"data_tlb_write_hit,omitempty"`
	DataTLBWriteMiss   *uint64 `json:"data_tlb_write_miss,omitempty"`
	InstrTLBReadHit    *uint64 `json:"instr_tlb_read_hit,omitempty"`
	InstrTLBReadMiss   *uint64 `json:"instr_tlb_read_miss,omitempty"`
	BPUReadHit         *uint64 `json:"bpu_read_hit,omitempty"`
	BPUReadMiss        *uint64 `json:"bpu_read_miss,omitempty"`
	NodeReadHit        *uint64 `json:"node_read_hit,omitempty"`
	NodeReadMiss       *uint64 `json:"node_read_miss,omitempty"`
	NodeWriteHit       *uint64 `json:"node_write_hit,omitempty"`
	NodeWriteMiss      *uint64 `json:"node_write_miss,omitempty"`
	TimeEnabled        *uint64 `json:"time_enabled,omitempty"`
	TimeRunning        *uint64 `json:"time_running,omitempty"`
}

// ProfileValue is a value returned by a profiler.
type ProfileValue struct {
	Value       uint64
	TimeEnabled uint64
	TimeRunning uint64
}

// profiler is used to profile a process.
type profiler struct {
	fd int
}

// NewProfiler creates a new hardware profiler. It does not support grouping.
func NewProfiler(profilerType uint32, config uint64, pid, cpu int, opts ...int) (Profiler, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:        profilerType,
		Config:      config,
		Size:        uint32(unsafe.Sizeof(unix.PerfEventAttr{})),
		Bits:        unix.PerfBitDisabled | unix.PerfBitExcludeHv | unix.PerfBitInherit,
		Read_format: unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_TOTAL_TIME_ENABLED,
		Sample_type: PERF_SAMPLE_IDENTIFIER,
	}
	var eventOps int
	if len(opts) > 0 {
		eventOps = opts[0]
	}
	fd, err := unix.PerfEventOpen(
		eventAttr,
		pid,
		cpu,
		-1,
		eventOps,
	)
	if err != nil {
		return nil, err
	}

	return &profiler{
		fd: fd,
	}, nil
}

// NewCPUCycleProfiler returns a Profiler that profiles CPU cycles.
func NewCPUCycleProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_CPU_CYCLES,
		pid,
		cpu,
		opts...,
	)
}

// NewInstrProfiler returns a Profiler that profiles CPU instructions.
func NewInstrProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_INSTRUCTIONS,
		pid,
		cpu,
		opts...,
	)
}

// NewCacheRefProfiler returns a Profiler that profiles cache references.
func NewCacheRefProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_CACHE_REFERENCES,
		pid,
		cpu,
		opts...,
	)
}

// NewCacheMissesProfiler returns a Profiler that profiles cache misses.
func NewCacheMissesProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_CACHE_MISSES,
		pid,
		cpu,
		opts...,
	)
}

// NewBranchInstrProfiler returns a Profiler that profiles branch instructions.
func NewBranchInstrProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_BRANCH_INSTRUCTIONS,
		pid,
		cpu,
		opts...,
	)
}

// NewBranchMissesProfiler returns a Profiler that profiles branch misses.
func NewBranchMissesProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_BRANCH_MISSES,
		pid,
		cpu,
		opts...,
	)
}

// NewBusCyclesProfiler returns a Profiler that profiles bus cycles.
func NewBusCyclesProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_BUS_CYCLES,
		pid,
		cpu,
		opts...,
	)
}

// NewStalledCyclesFrontProfiler returns a Profiler that profiles stalled
// frontend cycles.
func NewStalledCyclesFrontProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_STALLED_CYCLES_FRONTEND,
		pid,
		cpu,
		opts...,
	)
}

// NewStalledCyclesBackProfiler returns a Profiler that profiles stalled
// backend cycles.
func NewStalledCyclesBackProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_STALLED_CYCLES_BACKEND,
		pid,
		cpu,
		opts...,
	)
}

// NewRefCPUCyclesProfiler returns a Profiler that profiles CPU cycles, it
// is not affected by frequency scaling.
func NewRefCPUCyclesProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HARDWARE,
		unix.PERF_COUNT_HW_REF_CPU_CYCLES,
		pid,
		cpu,
		opts...,
	)
}

// NewCPUClockProfiler returns a Profiler that profiles CPU clock speed.
func NewCPUClockProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_CPU_CLOCK,
		pid,
		cpu,
		opts...,
	)
}

// NewTaskClockProfiler returns a Profiler that profiles clock count of the
// running task.
func NewTaskClockProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_TASK_CLOCK,
		pid,
		cpu,
		opts...,
	)
}

// NewPageFaultProfiler returns a Profiler that profiles the number of page
// faults.
func NewPageFaultProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_PAGE_FAULTS,
		pid,
		cpu,
		opts...,
	)
}

// NewCtxSwitchesProfiler returns a Profiler that profiles the number of context
// switches.
func NewCtxSwitchesProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_CONTEXT_SWITCHES,
		pid,
		cpu,
		opts...,
	)
}

// NewCPUMigrationsProfiler returns a Profiler that profiles the number of times
// the process has migrated to a new CPU.
func NewCPUMigrationsProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_CPU_MIGRATIONS,
		pid,
		cpu,
		opts...,
	)
}

// NewMinorFaultsProfiler returns a Profiler that profiles the number of minor
// page faults.
func NewMinorFaultsProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_PAGE_FAULTS_MIN,
		pid,
		cpu,
		opts...,
	)
}

// NewMajorFaultsProfiler returns a Profiler that profiles the number of major
// page faults.
func NewMajorFaultsProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_PAGE_FAULTS_MAJ,
		pid,
		cpu,
		opts...,
	)
}

// NewAlignFaultsProfiler returns a Profiler that profiles the number of
// alignment faults.
func NewAlignFaultsProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_ALIGNMENT_FAULTS,
		pid,
		cpu,
		opts...,
	)
}

// NewEmulationFaultsProfiler returns a Profiler that profiles the number of
// alignment faults.
func NewEmulationFaultsProfiler(pid, cpu int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_SOFTWARE,
		unix.PERF_COUNT_SW_EMULATION_FAULTS,
		pid,
		cpu,
		opts...,
	)
}

// NewL1DataProfiler returns a Profiler that profiles L1 cache data.
func NewL1DataProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {

	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_L1D)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewL1InstrProfiler returns a Profiler that profiles L1 instruction data.
func NewL1InstrProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_L1I)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewLLCacheProfiler returns a Profiler that profiles last level cache.
func NewLLCacheProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_LL)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewDataTLBProfiler returns a Profiler that profiles the data TLB.
func NewDataTLBProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_DTLB)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewInstrTLBProfiler returns a Profiler that profiles the instruction TLB.
func NewInstrTLBProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_ITLB)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewBPUProfiler returns a Profiler that profiles the BPU (branch prediction unit).
func NewBPUProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_BPU)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// NewNodeCacheProfiler returns a Profiler that profiles the node cache accesses.
func NewNodeCacheProfiler(pid, cpu, op, result int, opts ...int) (Profiler, error) {
	return NewProfiler(
		unix.PERF_TYPE_HW_CACHE,
		uint64((unix.PERF_COUNT_HW_CACHE_NODE)|(op<<8)|(result<<16)),
		pid,
		cpu,
		opts...,
	)
}

// Reset is used to reset the counters of the profiler.
func (p *profiler) Reset() error {
	return unix.IoctlSetInt(p.fd, unix.PERF_EVENT_IOC_RESET, 0)
}

// Start is used to Start the profiler.
func (p *profiler) Start() error {
	return unix.IoctlSetInt(p.fd, unix.PERF_EVENT_IOC_ENABLE, 0)
}

// Stop is used to stop the profiler.
func (p *profiler) Stop() error {
	return unix.IoctlSetInt(p.fd, unix.PERF_EVENT_IOC_DISABLE, 0)
}

// Profile returns the current Profile.
func (p *profiler) Profile() (*ProfileValue, error) {
	// The underlying struct that gets read from the profiler looks like:
	/*
		     struct read_format {
			 u64 value;         // The value of the event
			 u64 time_enabled;  // if PERF_FORMAT_TOTAL_TIME_ENABLED
			 u64 time_running;  // if PERF_FORMAT_TOTAL_TIME_RUNNING
			 u64 id;            // if PERF_FORMAT_ID
		     };
	*/

	// read 24 bytes since PERF_FORMAT_TOTAL_TIME_ENABLED and
	// PERF_FORMAT_TOTAL_TIME_RUNNING are always set.
	// XXX: allow profile ids?
	buf := make([]byte, 24, 24)
	_, err := syscall.Read(p.fd, buf)
	if err != nil {
		return nil, err
	}

	return &ProfileValue{
		Value:       binary.LittleEndian.Uint64(buf[0:8]),
		TimeEnabled: binary.LittleEndian.Uint64(buf[8:16]),
		TimeRunning: binary.LittleEndian.Uint64(buf[16:24]),
	}, nil
}

// Close is used to close the perf context.
func (p *profiler) Close() error {
	return unix.Close(p.fd)
}
