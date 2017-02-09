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

// +build freebsd dragonfly
// +build !noexec

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type execCollector struct {
	sysctls []bsdSysctl
}

func init() {
	Factories["exec"] = NewExecCollector
}

// Returns a new Collector exposing system execution statistics
func NewExecCollector() (Collector, error) {
	// from sys/vm/vm_meter.c:
	// vm.stats.sys.v_swtch: Context switches
	// vm.stats.sys.v_trap: Traps
	// vm.stats.sys.v_syscall: System calls
	// vm.stats.sys.v_intr: Device interrupts
	// vm.stats.sys.v_soft: Software interrupts
	// vm.stats.vm.v_forks: Number of fork() calls

	return &execCollector{
		sysctls: []bsdSysctl{
			{
				name:        "context_switches_total",
				description: "Context switches",
				mib:         "vm.stats.sys.v_swtch",
			},
			{
				name:        "traps_total",
				description: "Traps",
				mib:         "vm.stats.sys.v_trap",
			},
			{
				name:        "system_calls_total",
				description: "System calls",
				mib:         "vm.stats.sys.v_syscall",
			},
			{
				name:        "device_interrupts_total",
				description: "Device interrupts",
				mib:         "vm.stats.sys.v_intr",
			},
			{
				name:        "software_interrupts_total",
				description: "Software interrupts",
				mib:         "vm.stats.sys.v_soft",
			},
			{
				name:        "forks_total",
				description: "Number of fork() calls",
				mib:         "vm.stats.vm.v_forks",
			},
		},
	}, nil
}

// Expose kernel and system execistics.
func (c *execCollector) Update(ch chan<- prometheus.Metric) (err error) {
	for i := range c.sysctls {
		vt := c.sysctls[i].valueType
		if vt == 0 {
			// Make good use of the zero value.
			vt = prometheus.CounterValue
		}

		v, err := c.sysctls[i].GetValue()
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(c.sysctls[i].GetDesc("exec"), vt, v)
	}

	return nil
}
