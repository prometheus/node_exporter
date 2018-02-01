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

// +build !nocpu

package collector

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

// https://man.openbsd.org/hz.9
type clockinfo struct {
	hz      int32
	tick    int32
	tickadj int32
	stathz  int32
	profhz  int32
}

// src/sys/sys/sched.h
type cputime struct {
	user float64
	nice float64
	sys  float64
	intr float64
	idle float64
}

func getCPUTimes() ([]cputime, error) {
	// defined in sys/sys/sched.h
	const CPUSTATES = 5
	var (
		cpb []byte
		err error
	)

	// OpenBSD doesn't have a handy sysctl like kern.cp_times for both single and multiple cores.
	sysctlNcpu, err := unix.SysctlUint32("hw.ncpu")
	if err != nil {
		return nil, err
	}
	ncpu := int(sysctlNcpu)

	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return nil, err
	}
	clock := *(*clockinfo)(unsafe.Pointer(&clockb[0]))

	// compute the kern.cp_times equivalent
	if ncpu > 1 {
		for i := 0; i < ncpu; i++ {
			cpTime2, err := unix.SysctlRaw("kern.cp_time2", i)
			if err != nil {
				return nil, fmt.Errorf("sysctl(kern.cp_time2) failed: %s", err)
			}
			cpb = append(cpb, cpTime2...)
		}
	} else {
		cpb, err = unix.SysctlRaw("kern.cp_time")
		if err != nil {
			return nil, fmt.Errorf("sysctl(kern.cp_time) failed: %s", err)
		}
	}

	// same bloc/logic as for FreeBSD
	var cpufreq float64
	if clock.stathz > 0 {
		cpufreq = float64(clock.stathz)
	} else {
		cpufreq = float64(clock.hz)
	}

	var times []float64
	for len(cpb) >= int(unsafe.Sizeof(int(0))) {
		t := *(*int)(unsafe.Pointer(&cpb[0]))
		times = append(times, float64(t)/cpufreq)
		cpb = cpb[unsafe.Sizeof(int(0)):]
	}

	cpus := make([]cputime, len(times)/CPUSTATES)
	for i := 0; i < len(times); i += CPUSTATES {
		cpu := &cpus[i/CPUSTATES]
		cpu.user = times[i]
		cpu.nice = times[i+1]
		cpu.sys = times[i+2]
		cpu.intr = times[i+3]
		cpu.idle = times[i+4]
	}
	return cpus, nil
}

type statCollector struct {
	cpu  typedDesc
	temp typedDesc
}

func init() {
	registerCollector("cpu", defaultEnabled, NewStatCollector)
}

// NewStatCollector returns a new Collector exposing CPU stats.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cpu", "seconds_total"),
			"Seconds the CPU spent in each mode.",
			[]string{"cpu", "mode"}, nil,
		), prometheus.CounterValue},
	}, nil
}

// Expose CPU stats using sysctl.
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	// We want time spent per-cpu per CPUSTATE.
	// CPUSTATES (number of CPUSTATES) is defined as 5U.
	// Order: CP_USER | CP_NICE | CP_SYS | CP_IDLE | CP_INTR
	// sysctl kern.cp_times provides hw.ncpu * CPUSTATES long integers:
	//   hw.ncpu * (space-separated list of the above variables)
	//
	// Each value is a counter incremented at frequency
	//   kern.clockrate.(stathz | hz)
	//
	// Look into sys/kern/kern_clock.c for details.

	cpuTimes, err := getCPUTimes()
	if err != nil {
		return err
	}
	for cpu, t := range cpuTimes {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(t.user), lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(t.nice), lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(t.sys), lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(t.intr), lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(t.idle), lcpu, "idle")
	}
	return err
}
