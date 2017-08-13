// Copyright 2017 The Prometheus Authors
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
	registerCollector("exec", defaultEnabled, NewExecCollector)
}

// NewExecCollector returns a new Collector exposing system execution statistics.
func NewExecCollector() (Collector, error) {
	// From sys/vm/vm_meter.c:
	// All are of type CTLTYPE_UINT.
	//
	// vm.stats.sys.v_swtch: Context switches
	// vm.stats.sys.v_trap: Traps
	// vm.stats.sys.v_syscall: System calls
	// vm.stats.sys.v_intr: Device interrupts
	// vm.stats.sys.v_soft: Software interrupts
	// vm.stats.vm.v_forks: Number of fork() calls
	//
	// From sys/kern/kern_tc.c:
	// kern.boottime is an S,timeval

	return &execCollector{
		sysctls: []bsdSysctl{
			{
				name:        "context_switches_total",
				description: "Context switches since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.sys.v_swtch",
			},
			{
				name:        "traps_total",
				description: "Traps since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.sys.v_trap",
			},
			{
				name:        "system_calls_total",
				description: "System calls since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.sys.v_syscall",
			},
			{
				name:        "device_interrupts_total",
				description: "Device interrupts since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.sys.v_intr",
			},
			{
				name:        "software_interrupts_total",
				description: "Software interrupts since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.sys.v_soft",
			},
			{
				name:        "forks_total",
				description: "Number of fork() calls since system boot.  Resets at architeture unsigned integer.",
				mib:         "vm.stats.vm.v_forks",
			},
			{
				name:        "boot_timestamp_seconds",
				description: "Unix time of last boot, including microseconds.",
				mib:         "kern.boottime",
				dataType:    bsdSysctlTypeStructTimeval,
			},
		},
	}, nil
}

// Update pushes exec statistics onto ch
func (c *execCollector) Update(ch chan<- prometheus.Metric) error {
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "exec", m.name),
				m.description,
				nil, nil,
			), prometheus.CounterValue, v)
	}

	return nil
}
