// +build linux,gpu

package nvml

/*
#cgo CFLAGS: -std=c99 -I .
#cgo LDFLAGS: -L${SRCDIR}/ -lnvidia-ml
#include "nvidia-smi.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Nvidia_GPU_Stats struct {
	ClockMem      uint // GPU memory clock in Mhz
	ClockGraphics uint // GPU graphics clock in Mhz
	Throttle      uint // throttle reason
	PerfState     uint // performance state    C.uint 0: max / 15: min
	Temperature   uint // GPU temperature in Celsius degrees
	UtilGPU       uint // percentage of time during kernels are executing on the GPU.
	UtilMem       uint // percentage of time during memory is being read or written.
	MemUsage      uint // percentage of used memory size
	ID            uint // device ID
}

var fn_get_stats = get_stats_real

// get_stats_real calls real function to get GPU stats
func get_stats_real(cnt *C.int, cs **C.GPU_Stats) C.int {
	return C.get_gpu_stats(cnt, cs)
}

// GetGPUStats() returns stats of each GPU on host
func GetGPUStats() []Nvidia_GPU_Stats {
	var stats []Nvidia_GPU_Stats
	var cs *C.GPU_Stats
	var cnt C.int
	if fn_get_stats(&cnt, &cs) == -1 {
		return nil
	}
	defer C.free(unsafe.Pointer(cs))

	// convert C array to slice
	slice := (*[1 << 30]C.GPU_Stats)(unsafe.Pointer(cs))[:cnt:cnt]

	stats = make([]Nvidia_GPU_Stats, cnt)
	for i, val := range slice {
		stats[i].ClockMem = uint(val.clock_mem)
		stats[i].ClockGraphics = uint(val.clock_graphics)
		if val.throttle >= 0xFFFF {
			fmt.Printf("nvml: unknown throttle reason 0x%x\n", val.throttle)
			stats[i].Throttle = 0xFFFF
		} else {
			stats[i].Throttle = uint(val.throttle)
		}
		stats[i].PerfState = uint(val.perf_state)
		stats[i].Temperature = uint(val.temperature)
		stats[i].UtilGPU = uint(val.util_gpu)
		stats[i].MemUsage = uint(val.mem_usage)
		stats[i].ID = uint(val.id)
	}

	return stats
}
