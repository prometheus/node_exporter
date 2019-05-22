// +build linux

package perf

import (
	"go.uber.org/multierr"
	"golang.org/x/sys/unix"
)

const (
	// L1DataReadHit is a constant...
	L1DataReadHit = (unix.PERF_COUNT_HW_CACHE_L1D) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// L1DataReadMiss is a constant...
	L1DataReadMiss = (unix.PERF_COUNT_HW_CACHE_L1D) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)
	// L1DataWriteHit is a constant...
	L1DataWriteHit = (unix.PERF_COUNT_HW_CACHE_L1D) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// L1InstrReadMiss is a constant...
	L1InstrReadMiss = (unix.PERF_COUNT_HW_CACHE_L1I) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)

	// LLReadHit is a constant...
	LLReadHit = (unix.PERF_COUNT_HW_CACHE_LL) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// LLReadMiss is a constant...
	LLReadMiss = (unix.PERF_COUNT_HW_CACHE_LL) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)
	// LLWriteHit is a constant...
	LLWriteHit = (unix.PERF_COUNT_HW_CACHE_LL) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// LLWriteMiss is a constant...
	LLWriteMiss = (unix.PERF_COUNT_HW_CACHE_LL) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)

	// DataTLBReadHit is a constant...
	DataTLBReadHit = (unix.PERF_COUNT_HW_CACHE_DTLB) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// DataTLBReadMiss is a constant...
	DataTLBReadMiss = (unix.PERF_COUNT_HW_CACHE_DTLB) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)
	// DataTLBWriteHit is a constant...
	DataTLBWriteHit = (unix.PERF_COUNT_HW_CACHE_DTLB) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// DataTLBWriteMiss is a constant...
	DataTLBWriteMiss = (unix.PERF_COUNT_HW_CACHE_DTLB) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)

	// InstrTLBReadHit is a constant...
	InstrTLBReadHit = (unix.PERF_COUNT_HW_CACHE_ITLB) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// InstrTLBReadMiss is a constant...
	InstrTLBReadMiss = (unix.PERF_COUNT_HW_CACHE_ITLB) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)

	// BPUReadHit is a constant...
	BPUReadHit = (unix.PERF_COUNT_HW_CACHE_BPU) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// BPUReadMiss is a constant...
	BPUReadMiss = (unix.PERF_COUNT_HW_CACHE_BPU) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)

	// NodeCacheReadHit is a constant...
	NodeCacheReadHit = (unix.PERF_COUNT_HW_CACHE_NODE) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// NodeCacheReadMiss is a constant...
	NodeCacheReadMiss = (unix.PERF_COUNT_HW_CACHE_NODE) | (unix.PERF_COUNT_HW_CACHE_OP_READ << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)
	// NodeCacheWriteHit is a constant...
	NodeCacheWriteHit = (unix.PERF_COUNT_HW_CACHE_NODE) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)
	// NodeCacheWriteMiss is a constant...
	NodeCacheWriteMiss = (unix.PERF_COUNT_HW_CACHE_NODE) | (unix.PERF_COUNT_HW_CACHE_OP_WRITE << 8) | (unix.PERF_COUNT_HW_CACHE_RESULT_MISS << 16)
)

type cacheProfiler struct {
	// map of perf counter type to file descriptor
	profilers map[int]Profiler
}

// NewCacheProfiler returns a new cache profiler.
func NewCacheProfiler(pid, cpu int, opts ...int) CacheProfiler {
	profilers := map[int]Profiler{}

	// L1 data
	op := unix.PERF_COUNT_HW_CACHE_OP_READ
	result := unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	l1dataReadHit, err := NewL1DataProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[L1DataReadHit] = l1dataReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	l1dataReadMiss, err := NewL1DataProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[L1DataReadMiss] = l1dataReadMiss
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	l1dataWriteHit, err := NewL1DataProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[L1DataWriteHit] = l1dataWriteHit
	}

	// L1 instruction
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	l1InstrReadMiss, err := NewL1InstrProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[L1InstrReadMiss] = l1InstrReadMiss
	}

	// Last Level
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	llReadHit, err := NewLLCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[LLReadHit] = llReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	llReadMiss, err := NewLLCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[LLReadMiss] = llReadMiss
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	llWriteHit, err := NewLLCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[LLWriteHit] = llWriteHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	llWriteMiss, err := NewLLCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[LLWriteMiss] = llWriteMiss
	}

	// dTLB
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	dTLBReadHit, err := NewDataTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[DataTLBReadHit] = dTLBReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	dTLBReadMiss, err := NewDataTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[DataTLBReadMiss] = dTLBReadMiss
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	dTLBWriteHit, err := NewDataTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[DataTLBWriteHit] = dTLBWriteHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	dTLBWriteMiss, err := NewDataTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[DataTLBWriteMiss] = dTLBWriteMiss
	}

	// iTLB
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	iTLBReadHit, err := NewInstrTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[InstrTLBReadHit] = iTLBReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	iTLBReadMiss, err := NewInstrTLBProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[InstrTLBReadMiss] = iTLBReadMiss
	}

	// BPU
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	bpuReadHit, err := NewBPUProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[BPUReadHit] = bpuReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	bpuReadMiss, err := NewBPUProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[BPUReadMiss] = bpuReadMiss
	}

	// Node
	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	nodeReadHit, err := NewNodeCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[NodeCacheReadHit] = nodeReadHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_READ
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	nodeReadMiss, err := NewNodeCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[NodeCacheReadMiss] = nodeReadMiss
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_ACCESS
	nodeWriteHit, err := NewNodeCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[NodeCacheWriteHit] = nodeWriteHit
	}

	op = unix.PERF_COUNT_HW_CACHE_OP_WRITE
	result = unix.PERF_COUNT_HW_CACHE_RESULT_MISS
	nodeWriteMiss, err := NewNodeCacheProfiler(pid, cpu, op, result, opts...)
	if err == nil {
		profilers[NodeCacheWriteMiss] = nodeWriteMiss
	}

	return &cacheProfiler{
		profilers: profilers,
	}
}

// Start is used to start the CacheProfiler, it will return an error if no
// profilers are configured.
func (p *cacheProfiler) Start() error {
	if len(p.profilers) == 0 {
		return ErrNoProfiler
	}
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Start())
	}
	return err
}

// Reset is used to reset the CacheProfiler.
func (p *cacheProfiler) Reset() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Reset())
	}
	return err
}

// Stop is used to reset the CacheProfiler.
func (p *cacheProfiler) Stop() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Stop())
	}
	return err
}

// Close is used to reset the CacheProfiler.
func (p *cacheProfiler) Close() error {
	var err error
	for _, profiler := range p.profilers {
		err = multierr.Append(err, profiler.Close())
	}
	return err
}

// Profile is used to read the CacheProfiler CacheProfile it returns an
// error only if all profiles fail.
func (p *cacheProfiler) Profile() (*CacheProfile, error) {
	var err error
	cacheProfile := &CacheProfile{}
	for profilerType, profiler := range p.profilers {
		profileVal, err2 := profiler.Profile()
		err = multierr.Append(err, err2)
		if err2 == nil {
			if cacheProfile.TimeEnabled == nil {
				cacheProfile.TimeEnabled = &profileVal.TimeEnabled
			}
			if cacheProfile.TimeRunning == nil {
				cacheProfile.TimeRunning = &profileVal.TimeRunning
			}
			switch {
			// L1 data
			case (profilerType ^ L1DataReadHit) == 0:
				cacheProfile.L1DataReadHit = &profileVal.Value
			case (profilerType ^ L1DataReadMiss) == 0:
				cacheProfile.L1DataReadMiss = &profileVal.Value
			case (profilerType ^ L1DataWriteHit) == 0:
				cacheProfile.L1DataWriteHit = &profileVal.Value

			// L1 instruction
			case (profilerType ^ L1InstrReadMiss) == 0:
				cacheProfile.L1InstrReadMiss = &profileVal.Value

			// Last Level
			case (profilerType ^ LLReadHit) == 0:
				cacheProfile.LastLevelReadHit = &profileVal.Value
			case (profilerType ^ LLReadMiss) == 0:
				cacheProfile.LastLevelReadMiss = &profileVal.Value
			case (profilerType ^ LLWriteHit) == 0:
				cacheProfile.LastLevelWriteHit = &profileVal.Value
			case (profilerType ^ LLWriteMiss) == 0:
				cacheProfile.LastLevelWriteMiss = &profileVal.Value

			// dTLB
			case (profilerType ^ DataTLBReadHit) == 0:
				cacheProfile.DataTLBReadHit = &profileVal.Value
			case (profilerType ^ DataTLBReadMiss) == 0:
				cacheProfile.DataTLBReadMiss = &profileVal.Value
			case (profilerType ^ DataTLBWriteHit) == 0:
				cacheProfile.DataTLBWriteHit = &profileVal.Value
			case (profilerType ^ DataTLBWriteMiss) == 0:
				cacheProfile.DataTLBWriteMiss = &profileVal.Value

			// iTLB
			case (profilerType ^ InstrTLBReadHit) == 0:
				cacheProfile.InstrTLBReadHit = &profileVal.Value
			case (profilerType ^ InstrTLBReadMiss) == 0:
				cacheProfile.InstrTLBReadMiss = &profileVal.Value

			// BPU
			case (profilerType ^ BPUReadHit) == 0:
				cacheProfile.BPUReadHit = &profileVal.Value
			case (profilerType ^ BPUReadMiss) == 0:
				cacheProfile.BPUReadMiss = &profileVal.Value

			// node
			case (profilerType ^ NodeCacheReadHit) == 0:
				cacheProfile.NodeReadHit = &profileVal.Value
			case (profilerType ^ NodeCacheReadMiss) == 0:
				cacheProfile.NodeReadMiss = &profileVal.Value
			case (profilerType ^ NodeCacheWriteHit) == 0:
				cacheProfile.NodeWriteHit = &profileVal.Value
			case (profilerType ^ NodeCacheWriteMiss) == 0:
				cacheProfile.NodeWriteMiss = &profileVal.Value
			}
		}
	}
	if len(multierr.Errors(err)) == len(p.profilers) {
		return nil, err
	}

	return cacheProfile, nil
}
