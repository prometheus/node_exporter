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

//go:build !nocpu
// +build !nocpu

package collector

import (
	"errors"
	"strconv"
	"unsafe"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

/*
#cgo LDFLAGS:
#include <sys/sysctl.h>
#include <kinfo.h>
#include <stdlib.h>
#include <stdio.h>

int
getCPUTimes(uint64_t **cputime, size_t *cpu_times_len) {
	size_t len;

	// Get number of cpu cores.
	int mib[2];
	int ncpu;
	mib[0] = CTL_HW;
	mib[1] = HW_NCPU;
	len = sizeof(ncpu);
	if (sysctl(mib, 2, &ncpu, &len, NULL, 0)) {
		return -1;
	}

	// Get the cpu times.
	struct kinfo_cputime cp_t[ncpu];
	bzero(cp_t, sizeof(struct kinfo_cputime)*ncpu);
	len = sizeof(cp_t[0])*ncpu;
	if (sysctlbyname("kern.cputime", &cp_t, &len, NULL, 0)) {
		return -1;
	}

	*cpu_times_len = ncpu*CPUSTATES;

	uint64_t user, nice, sys, intr, idle;
	user = nice = sys = intr = idle = 0;
	*cputime = (uint64_t *) malloc(sizeof(uint64_t)*(*cpu_times_len));
	for (int i = 0; i < ncpu; ++i) {
		int offset = CPUSTATES * i;
		(*cputime)[offset] = cp_t[i].cp_user;
		(*cputime)[offset+1] = cp_t[i].cp_nice;
		(*cputime)[offset+2] = cp_t[i].cp_sys;
		(*cputime)[offset+3] = cp_t[i].cp_intr;
		(*cputime)[offset+4] = cp_t[i].cp_idle;
	}

	return 0;

}
*/
import "C"

const maxCPUTimesLen = C.MAXCPU * C.CPUSTATES

type statCollector struct {
	cpu    *prometheus.Desc
	logger log.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewStatCollector)
}

// NewStatCollector returns a new Collector exposing CPU stats.
func NewStatCollector(logger log.Logger) (Collector, error) {
	return &statCollector{
		cpu:    nodeCPUSecondsDesc,
		logger: logger,
	}, nil
}

func getDragonFlyCPUTimes() ([]float64, error) {
	// We want time spent per-CPU per CPUSTATE.
	// CPUSTATES (number of CPUSTATES) is defined as 5U.
	// States: CP_USER | CP_NICE | CP_SYS | CP_IDLE | CP_INTR
	//
	// Each value is in microseconds
	//
	// Look into sys/kern/kern_clock.c for details.

	var (
		cpuTimesC      *C.uint64_t
		cpuTimesLength C.size_t
	)

	if C.getCPUTimes(&cpuTimesC, &cpuTimesLength) == -1 {
		return nil, errors.New("could not retrieve CPU times")
	}
	defer C.free(unsafe.Pointer(cpuTimesC))

	cput := (*[maxCPUTimesLen]C.uint64_t)(unsafe.Pointer(cpuTimesC))[:cpuTimesLength:cpuTimesLength]

	cpuTimes := make([]float64, cpuTimesLength)
	for i, value := range cput {
		cpuTimes[i] = float64(value) / float64(1000000)
	}
	return cpuTimes, nil
}

// Expose CPU stats using sysctl.
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	var fieldsCount = 5
	cpuTimes, err := getDragonFlyCPUTimes()
	if err != nil {
		return err
	}

	// Export order: user nice sys intr idle
	cpuFields := []string{"user", "nice", "sys", "interrupt", "idle"}
	for i, value := range cpuTimes {
		cpux := strconv.Itoa(i / fieldsCount)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, value, cpux, cpuFields[i%fieldsCount])
	}

	return nil
}
