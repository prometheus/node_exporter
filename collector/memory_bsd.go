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
// +build !nomeminfo

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	memorySubsystem = "memory"
)

type memoryCollector struct {
	pageSize uint64
	sysctls  []bsdSysctl
	kvm      kvm
}

func init() {
	registerCollector("meminfo", defaultEnabled, NewMemoryCollector)
}

// NewMemoryCollector returns a new Collector exposing memory stats.
func NewMemoryCollector() (Collector, error) {
	tmp32, err := unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return nil, fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}
	size := float64(tmp32)

	mibSwapTotal := "vm.swap_total"
	/* swap_total is FreeBSD specific. Fall back to Dfly specific mib if not present. */
	_, err = unix.SysctlUint64(mibSwapTotal)
	if err != nil {
		mibSwapTotal = "vm.swap_size"
	}

	fromPage := func(v float64) float64 {
		return v * size
	}

	return &memoryCollector{
		pageSize: uint64(tmp32),
		sysctls: []bsdSysctl{
			// Descriptions via: https://wiki.freebsd.org/Memory
			{
				name:        "active_bytes",
				description: "Recently used by userland",
				mib:         "vm.stats.vm.v_active_count",
				conversion:  fromPage,
			},
			{
				name:        "inactive_bytes",
				description: "Not recently used by userland",
				mib:         "vm.stats.vm.v_inactive_count",
				conversion:  fromPage,
			},
			{
				name:        "wired_bytes",
				description: "Locked in memory by kernel, mlock, etc",
				mib:         "vm.stats.vm.v_wire_count",
				conversion:  fromPage,
			},
			{
				name:        "cache_bytes",
				description: "Almost free, backed by swap or files, available for re-allocation",
				mib:         "vm.stats.vm.v_cache_count",
				conversion:  fromPage,
			},
			{
				name:        "buffer_bytes",
				description: "Disk IO Cache entries for non ZFS filesystems, only usable by kernel",
				mib:         "vfs.bufspace",
				dataType:    bsdSysctlTypeCLong,
			},
			{
				name:        "free_bytes",
				description: "Unallocated, available for allocation",
				mib:         "vm.stats.vm.v_free_count",
				conversion:  fromPage,
			},
			{
				name:        "size_bytes",
				description: "Total physical memory size",
				mib:         "vm.stats.vm.v_page_count",
				conversion:  fromPage,
			},
			{
				name:        "swap_size_bytes",
				description: "Total swap memory size",
				mib:         mibSwapTotal,
				dataType:    bsdSysctlTypeUint64,
			},
			// Descriptions via: top(1)
			{
				name:        "swap_in_bytes_total",
				description: "Bytes paged in from swap devices",
				mib:         "vm.stats.vm.v_swappgsin",
				valueType:   prometheus.CounterValue,
				conversion:  fromPage,
			},
			{
				name:        "swap_out_bytes_total",
				description: "Bytes paged out to swap devices",
				mib:         "vm.stats.vm.v_swappgsout",
				valueType:   prometheus.CounterValue,
				conversion:  fromPage,
			},
		},
	}, nil
}

// Update checks relevant sysctls for current memory usage, and kvm for swap
// usage.
func (c *memoryCollector) Update(ch chan<- prometheus.Metric) error {
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			return fmt.Errorf("couldn't get memory: %s", err)
		}

		// Most are gauges.
		if m.valueType == 0 {
			m.valueType = prometheus.GaugeValue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, memorySubsystem, m.name),
				m.description,
				nil, nil,
			), m.valueType, v)
	}

	swapUsed, err := c.kvm.SwapUsedPages()
	if err != nil {
		return fmt.Errorf("couldn't get kvm: %s", err)
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, memorySubsystem, "swap_used_bytes"),
			"Currently allocated swap",
			nil, nil,
		), prometheus.GaugeValue, float64(swapUsed*c.pageSize))

	return nil
}
