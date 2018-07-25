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
	"math"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/sys/unix"
)

type clockinfo struct {
	hz     int32 // clock frequency
	tick   int32 // micro-seconds per hz tick
	spare  int32
	stathz int32 // statistics clock frequency
	profhz int32 // profiling clock frequency
}

type cputime struct {
	user float64
	nice float64
	sys  float64
	intr float64
	idle float64
}

func getCPUTimes() ([]cputime, error) {
	const states = 5

	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return nil, err
	}
	clock := *(*clockinfo)(unsafe.Pointer(&clockb[0]))
	cpb, err := unix.SysctlRaw("kern.cp_times")
	if err != nil {
		return nil, err
	}

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

	cpus := make([]cputime, len(times)/states)
	for i := 0; i < len(times); i += states {
		cpu := &cpus[i/states]
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
		cpu: typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		temp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "temperature_celsius"),
			"CPU temperature",
			[]string{"cpu"}, nil,
		), prometheus.GaugeValue},
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

		temp, err := unix.SysctlUint32(fmt.Sprintf("dev.cpu.%d.temperature", cpu))
		if err != nil {
			if err == unix.ENOENT {
				// No temperature information for this CPU
				log.Debugf("no temperature information for CPU %d", cpu)
			} else {
				// Unexpected error
				ch <- c.temp.mustNewConstMetric(math.NaN(), lcpu)
				log.Errorf("failed to query CPU temperature for CPU %d: %s", cpu, err)
			}
			continue
		}
		ch <- c.temp.mustNewConstMetric(float64(temp-2732)/10, lcpu)
	}
	return err
}
