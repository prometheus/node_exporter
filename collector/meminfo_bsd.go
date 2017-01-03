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
// +build !nomeminfo

package collector

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	info := make(map[string]float64)

	size, err := unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return nil, fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}

	for key, v := range map[string]string{
		"active":     "vm.stats.vm.v_active_count",
		"inactive":   "vm.stats.vm.v_inactive_count",
		"wire":       "vm.stats.vm.v_wire_count",
		"cache":      "vm.stats.vm.v_cache_count",
		"free":       "vm.stats.vm.v_free_count",
		"swappgsin":  "vm.stats.vm.v_swappgsin",
		"swappgsout": "vm.stats.vm.v_swappgsout",
		"total":      "vm.stats.vm.v_page_count",
	} {
		value, err := unix.SysctlUint32(v)
		if err != nil {
			return nil, err
		}
		// Convert metrics to kB (same as Linux meminfo).
		info[key] = float64(value) * float64(size)
	}
	return info, nil
}
