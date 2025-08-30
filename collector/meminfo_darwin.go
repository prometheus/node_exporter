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

//go:build !nomeminfo
// +build !nomeminfo

package collector

// #include <mach/mach_host.h>
// #include <sys/sysctl.h>
// typedef struct xsw_usage xsw_usage_t;
import "C"

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"unsafe"

	"golang.org/x/sys/unix"
)

type meminfoCollector struct {
	logger *slog.Logger
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewMeminfoCollector(logger *slog.Logger) (Collector, error) {
	return &meminfoCollector{
		logger: logger,
	}, nil
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	host := C.mach_host_self()
	infoCount := C.mach_msg_type_number_t(C.HOST_VM_INFO64_COUNT)
	vmstat := C.vm_statistics64_data_t{}
	ret := C.host_statistics64(
		C.host_t(host),
		C.HOST_VM_INFO64,
		C.host_info_t(unsafe.Pointer(&vmstat)),
		&infoCount,
	)
	if ret != C.KERN_SUCCESS {
		return nil, fmt.Errorf("couldn't get memory statistics, host_statistics returned %d", ret)
	}
	totalb, err := unix.Sysctl("hw.memsize")
	if err != nil {
		return nil, err
	}

	swapraw, err := unix.SysctlRaw("vm.swapusage")
	if err != nil {
		return nil, err
	}
	swap := (*C.xsw_usage_t)(unsafe.Pointer(&swapraw[0]))

	// Syscall removes terminating NUL which we need to cast to uint64
	total := binary.LittleEndian.Uint64([]byte(totalb + "\x00"))

	var pageSize C.vm_size_t
	C.host_page_size(C.host_t(host), &pageSize)

	ps := float64(pageSize)
	return map[string]float64{
		"active_bytes":            ps * float64(vmstat.active_count),
		"compressed_bytes":        ps * float64(vmstat.compressor_page_count),
		"inactive_bytes":          ps * float64(vmstat.inactive_count),
		"wired_bytes":             ps * float64(vmstat.wire_count),
		"free_bytes":              ps * float64(vmstat.free_count),
		"swapped_in_bytes_total":  ps * float64(vmstat.pageins),
		"swapped_out_bytes_total": ps * float64(vmstat.pageouts),
		"internal_bytes":          ps * float64(vmstat.internal_page_count),
		"purgeable_bytes":         ps * float64(vmstat.purgeable_count),
		"total_bytes":             float64(total),
		"swap_used_bytes":         float64(swap.xsu_used),
		"swap_total_bytes":        float64(swap.xsu_total),
	}, nil
}
