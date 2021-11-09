// Copyright 2016 The Prometheus Authors
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

// Based on gopsutil/cpu/cpu_darwin_cgo.go @ ae251eb which is licensed under
// BSD. See https://github.com/shirou/gopsutil/blob/master/LICENSE for details.

//go:build !nocpu
// +build !nocpu

package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS:
#include <stdlib.h>
#include <limits.h>
#include <sys/sysctl.h>
#include <sys/mount.h>
#include <mach/mach_init.h>
#include <mach/mach_host.h>
#include <mach/host_info.h>
#include <TargetConditionals.h>
#if TARGET_OS_MAC
#include <libproc.h>
#endif
#include <mach/processor_info.h>
#include <mach/vm_map.h>
*/
import "C"

// ClocksPerSec default value. from time.h
const ClocksPerSec = float64(C.CLK_TCK)

type statCollector struct {
	cpu    *prometheus.Desc
	logger log.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

// NewCPUCollector returns a new Collector exposing CPU stats.
func NewCPUCollector(logger log.Logger) (Collector, error) {
	return &statCollector{
		cpu:    nodeCPUSecondsDesc,
		logger: logger,
	}, nil
}

func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	var (
		count   C.mach_msg_type_number_t
		cpuload *C.processor_cpu_load_info_data_t
		ncpu    C.natural_t
	)

	status := C.host_processor_info(C.host_t(C.mach_host_self()),
		C.PROCESSOR_CPU_LOAD_INFO,
		&ncpu,
		(*C.processor_info_array_t)(unsafe.Pointer(&cpuload)),
		&count)

	if status != C.KERN_SUCCESS {
		return fmt.Errorf("host_processor_info error=%d", status)
	}

	// jump through some cgo casting hoops and ensure we properly free
	// the memory that cpuload points to
	target := C.vm_map_t(C.mach_task_self_)
	address := C.vm_address_t(uintptr(unsafe.Pointer(cpuload)))
	defer C.vm_deallocate(target, address, C.vm_size_t(ncpu))

	// the body of struct processor_cpu_load_info
	// aka processor_cpu_load_info_data_t
	var cpuTicks [C.CPU_STATE_MAX]uint32

	// copy the cpuload array to a []byte buffer
	// where we can binary.Read the data
	size := int(ncpu) * binary.Size(cpuTicks)
	buf := (*[1 << 30]byte)(unsafe.Pointer(cpuload))[:size:size]

	bbuf := bytes.NewBuffer(buf)

	for i := 0; i < int(ncpu); i++ {
		err := binary.Read(bbuf, binary.LittleEndian, &cpuTicks)
		if err != nil {
			return err
		}
		for k, v := range map[string]int{
			"user":   C.CPU_STATE_USER,
			"system": C.CPU_STATE_SYSTEM,
			"nice":   C.CPU_STATE_NICE,
			"idle":   C.CPU_STATE_IDLE,
		} {
			ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(cpuTicks[v])/ClocksPerSec, strconv.Itoa(i), k)
		}
	}
	return nil
}
