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

// +build freebsd darwin,amd64 dragonfly
// +build !nomeminfo

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	memInfoSubsystem = "memory"
)

type meminfoCollector struct{}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Memory stats.
func NewMeminfoCollector() (Collector, error) {
	return &meminfoCollector{}, nil
}

func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	pages := make(map[string]uint32)

	size, err := unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}
	pages["active"], _ = unix.SysctlUint32("vm.stats.vm.v_active_count")
	pages["inactive"], _ = unix.SysctlUint32("vm.stats.vm.v_inactive_count")
	pages["wire"], _ = unix.SysctlUint32("vm.stats.vm.v_wire_count")
	pages["cache"], _ = unix.SysctlUint32("vm.stats.vm.v_cache_count")
	pages["free"], _ = unix.SysctlUint32("vm.stats.vm.v_free_count")
	pages["swappgsin"], _ = unix.SysctlUint32("vm.stats.vm.v_swappgsin")
	pages["swappgsout"], _ = unix.SysctlUint32("vm.stats.vm.v_swappgsout")
	pages["total"], _ = unix.SysctlUint32("vm.stats.vm.v_page_count")

	for k, v := range pages {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, memInfoSubsystem, k),
				k+" from sysctl()",
				nil, nil,
			),
			// Convert metrics to kB (same as Linux meminfo).
			prometheus.UntypedValue, float64(v)*float64(size),
		)
	}
	return err
}
