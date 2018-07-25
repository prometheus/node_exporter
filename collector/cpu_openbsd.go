// Copyright 2018 The Prometheus Authors
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

// +build !nocpu

package collector

import (
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

/*
#include <sys/param.h>
#include <sys/sched.h>
*/
import "C"

type cpuCollector struct {
	cpu typedDesc
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func NewCpuCollector() (Collector, error) {
	return &cpuCollector{
		cpu: typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) (err error) {
	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return err
	}
	clock := *(*C.struct_clockinfo)(unsafe.Pointer(&clockb[0]))
	hz := float64(clock.stathz)

	ncpus, err := unix.SysctlUint32("hw.ncpu")
	if err != nil {
		return err
	}

	var cp_time [][C.CPUSTATES]C.int64_t
	for i := 0; i < int(ncpus); i++ {
		cp_timeb, err := unix.SysctlRaw("kern.cp_time2", i)
		if err != nil {
			return err
		}
		cp_time = append(cp_time, *(*[C.CPUSTATES]C.int64_t)(unsafe.Pointer(&cp_timeb[0])))
	}

	for cpu, time := range cp_time {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_USER])/hz, lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_NICE])/hz, lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_SYS])/hz, lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_INTR])/hz, lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_IDLE])/hz, lcpu, "idle")
	}
	return err
}
