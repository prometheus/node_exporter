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

	tmp32, err := unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return nil, fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}
	size := uint64(tmp32)
	fromPage := func(v uint64) uint64 {
		return v * size
	}

	for _, ctl := range []bsdSysctl{
		{name: "active_bytes", mib: "vm.stats.vm.v_active_count", conversion: fromPage},
		{name: "inactive_bytes", mib: "vm.stats.vm.v_inactive_count", conversion: fromPage},
		{name: "wired_bytes", mib: "vm.stats.vm.v_wire_count", conversion: fromPage},
		{name: "cache_bytes", mib: "vm.stats.vm.v_cache_count", conversion: fromPage},
		{name: "buffer_bytes", mib: "vfs.bufspace"},
		{name: "free_bytes", mib: "vm.stats.vm.v_free_count", conversion: fromPage},
		{name: "size_bytes", mib: "vm.stats.vm.v_page_count", conversion: fromPage},
		{name: "swap_in_bytes_total", mib: "vm.stats.vm.v_swappgsin", conversion: fromPage},
		{name: "swap_out_bytes_total", mib: "vm.stats.vm.v_swappgsout", conversion: fromPage},
		{name: "swap_size_bytes", mib: "vm.swap_total", dataType: bsdSysctlTypeUint64},
	} {
		v, err := ctl.Value()
		if err != nil {
			return nil, err
		}

		info[ctl.name] = v
	}

	return info, nil
}
