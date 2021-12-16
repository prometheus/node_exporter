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

//go:build solaris && !nocpu
// +build solaris,!nocpu

package collector

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/illumos/go-kstat"
	"github.com/prometheus/client_golang/prometheus"
)

// #include <unistd.h>
import "C"

type cpuCollector struct {
	cpu    typedDesc
	logger log.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func NewCpuCollector(logger log.Logger) (Collector, error) {
	return &cpuCollector{
		cpu:    typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	for cpu := 0; cpu < int(ncpus); cpu++ {
		ksCPU, err := tok.Lookup("cpu", cpu, "sys")
		if err != nil {
			return err
		}

		for k, v := range map[string]string{
			"idle":   "cpu_ticks_idle",
			"kernel": "cpu_ticks_kernel",
			"user":   "cpu_ticks_user",
			"wait":   "cpu_ticks_wait",
		} {
			kstatValue, err := ksCPU.GetNamed(v)
			if err != nil {
				return err
			}

			ch <- c.cpu.mustNewConstMetric(float64(kstatValue.UintVal), strconv.Itoa(cpu), k)
		}
	}
	return nil
}
