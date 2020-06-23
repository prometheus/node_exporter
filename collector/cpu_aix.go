// Copyright 2020 The Prometheus Authors
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

/*
#cgo LDFLAGS: -lperfstat
#include <stdlib.h>
#include <stdio.h>
#include <libperfstat.h>
#include <string.h>
#include <time.h>

u_longlong_t **ref;

int getCPUTicks(uint64_t **cputicks, size_t *cpu_ticks_len) {
	int i, ncpus, cputotal;
	perfstat_id_t firstcpu;
	perfstat_cpu_t *statp;

	cputotal = perfstat_cpu(NULL, NULL, sizeof(perfstat_cpu_t), 0);
	if (cputotal <= 0){
		return -1;
	}

	statp = calloc(cputotal, sizeof(perfstat_cpu_t));
	if(statp==NULL){
		return -1;
	}
	strcpy(firstcpu.name, FIRST_CPU);
	ncpus = perfstat_cpu(&firstcpu, statp, sizeof(perfstat_cpu_t), cputotal);

	*cpu_ticks_len = ncpus*4;
	*cputicks = (uint64_t *) malloc(sizeof(uint64_t)*(*cpu_ticks_len));
	for (i = 0; i < ncpus; i++) {
		int offset = 4 * i;
		(*cputicks)[offset] = statp[i].user;
		(*cputicks)[offset+1] = statp[i].sys;
		(*cputicks)[offset+2] = statp[i].wait;
		(*cputicks)[offset+3] = statp[i].idle;
	}
	return 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"unsafe"
)

const ClocksPerSec = float64(C.CLK_TCK)
const maxCPUTimesLen = 1024 * 4

type statCollector struct {
	cpu *prometheus.Desc
}

func init() {
	registerCollector("cpu", true, NewCPUCollector)
}

func NewCPUCollector(logger log.Logger) (Collector, error) {
	return &statCollector{
		cpu: nodeCPUSecondsDesc,
	}, nil
}

func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	var fieldsCount = 4
	cpuFields := []string{"user", "sys", "wait", "idle"}

	var (
		cpuTimesC      *C.uint64_t
		cpuTimesLength C.size_t
	)

	if C.getCPUTicks(&cpuTimesC, &cpuTimesLength) == -1 {
		return errors.New("could not retrieve CPU times")
	}
	defer C.free(unsafe.Pointer(cpuTimesC))
	cput := (*[maxCPUTimesLen]C.u_longlong_t)(unsafe.Pointer(cpuTimesC))[:cpuTimesLength:cpuTimesLength]

	cpuTicks := make([]float64, cpuTimesLength)
	for i, value := range cput {
		cpuTicks[i] = float64(value) / ClocksPerSec
	}

	for i, value := range cpuTicks {
		cpux := fmt.Sprintf("CPU %d", i/fieldsCount)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, value, cpux, cpuFields[i%fieldsCount])
	}

	return nil
}
