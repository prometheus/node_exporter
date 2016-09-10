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
	"errors"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS:
#include <stdio.h>
#include <kinfo.h>
#include <sys/sysctl.h>

int
getCPUTimes(int *ncpu, double **cpu_times, size_t *cp_times_length)
{
	int cpu, mib[2];
	struct kinfo_cputime cp_t[ncpu];
	uint64_t user, nice, sys, intr, idle;

	size_t len;

	mib[0] = CTL_HW;
	mib[1] = HW_NCPU;
	len = sizeof(*ncpu);
	if (-1 == sysctl(mib, 2, &ncpu, &len, NULL, 0))
		return -1;

	bzero(cp_t, sizeof(struct kinfo_cputime)*(*ncpu));

	len = sizeof(cp_t[0])*(*ncpu);
	if (sysctlbyname("kern.cputime", &cp_t, &len, NULL, 0))
		return -1;

	// Retrieve clockrate
	struct clockinfo clockrate;
	int mib_kern_clockrate[] = {CTL_KERN, KERN_CLOCKRATE};
	size_t mib_kern_clockrate_len = 2;
	size_t clockrate_size = sizeof(clockrate);

	if (sysctl(mib_kern_clockrate, mib_kern_clockrate_len, &clockrate, &clockrate_size, NULL, 0) == -1)
		return -1;

	long cpufreq = clockrate.stathz > 0 ? clockrate.stathz : clockrate.hz;
	cpu_times = (double *) malloc(sizeof(double)*(*cp_times_length));
	for (int i = 0; i < (*cp_times_length); i++) {
		(*cpu_times)[i] = ((double) cp_times[i]) / cpufreq;
	}

	return 0;
}

void freeCPUTimes(double *cpu_times) {
	free(cpu_times);
}
*/

import "C"

type statCollector struct {
	cpu *prometheus.CounterVec
}

func init() {
	Factories["cpu"] = NewStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// CPU stats.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "cpu_seconds_total",
				Help:      "Seconds the CPU spent in each mode.",
			},
			[]string{"cpu", "mode"},
		),
	}, nil
}

// Expose CPU stats using sysctl.
func (c *statCollector) Update(ch chan<- prometheus.Metric) (err error) {
	// Adapted from
	// https://www.dragonflybsd.org/mailarchive/users/2010-04/msg00056.html

	var ncpu C.int
	var cpuTimesC *C.double
	var cpuTimesLength C.size_t
	if C.getCPUTimes(&ncpu, &cpuTimesC, &cpuTimesLength) == -1 {
		return errors.New("could not retrieve CPU times")
	}
	defer C.freeCPUTimes(cpuTimesC)

	// Convert C.double array to Go array (https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices).
	cpuTimes := (*[maxCPUTimesLen]C.double)(unsafe.Pointer(cpuTimesC))[:cpuTimesLength:cpuTimesLength]

	for cpu := 0; cpu < int(ncpu); cpu++ {
		base_idx := C.CPUSTATES * cpu
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(cpu), "mode": "user"}).Set(float64(cpuTimes[base_idx+C.CP_USER]))
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(cpu), "mode": "nice"}).Set(float64(cpuTimes[base_idx+C.CP_NICE]))
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(cpu), "mode": "system"}).Set(float64(cpuTimes[base_idx+C.CP_SYS]))
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(cpu), "mode": "interrupt"}).Set(float64(cpuTimes[base_idx+C.CP_INTR]))
		c.cpu.With(prometheus.Labels{"cpu": strconv.Itoa(cpu), "mode": "idle"}).Set(float64(cpuTimes[base_idx+C.CP_IDLE]))
	}

	c.cpu.Collect(ch)
	return err
}
