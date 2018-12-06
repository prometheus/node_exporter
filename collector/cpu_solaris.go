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

// +build solaris
// +build !nocpu

package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/siebenmann/go-kstat"
)

// #include <unistd.h>
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
	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	for cpu := 0; cpu < int(ncpus); cpu++ {
		ks, err := tok.Lookup("cpu", cpu, "sys")
		if err != nil {
			return err
		}

		idle_v, err := ks.GetNamed("cpu_ticks_idle")
		if err != nil {
			return err
		}

		kernel_v, err := ks.GetNamed("cpu_ticks_kernel")
		if err != nil {
			return err
		}

		user_v, err := ks.GetNamed("cpu_ticks_user")
		if err != nil {
			return err
		}

		wait_v, err := ks.GetNamed("cpu_ticks_wait")
		if err != nil {
			return err
		}

		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(idle_v.UintVal), lcpu, "idle")
		ch <- c.cpu.mustNewConstMetric(float64(kernel_v.UintVal), lcpu, "kernel")
		ch <- c.cpu.mustNewConstMetric(float64(user_v.UintVal), lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(wait_v.UintVal), lcpu, "wait")
	}
	return err
}
