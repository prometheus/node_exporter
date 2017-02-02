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

type sysctlType uint8

const (
	t_uint32 sysctlType = iota
	t_uint64
)

type meminfoSysctl struct {
	name       string
	dataType   sysctlType
	conversion func(uint64) uint64
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	var tmp32 uint32
	var tmp64 uint64
	var err error

	info := make(map[string]float64)

	tmp32, err = unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return nil, fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}
	size := uint64(tmp32)
	fromPage := func(v uint64) uint64 {
		return v * size
	}

	for key, v := range map[string]meminfoSysctl{
		"active_bytes":         {"vm.stats.vm.v_active_count", t_uint32, fromPage},
		"inactive_bytes":       {"vm.stats.vm.v_inactive_count", t_uint32, fromPage},
		"wire_bytes":           {"vm.stats.vm.v_wire_count", t_uint32, fromPage},
		"cache_bytes":          {"vm.stats.vm.v_cache_count", t_uint32, fromPage},
		"buffer_bytes":         {"vfs.bufspace", t_uint32, nil},
		"free_bytes":           {"vm.stats.vm.v_free_count", t_uint32, fromPage},
		"size_bytes":           {"vm.stats.vm.v_page_count", t_uint32, fromPage},
		"swap_in_bytes_total":  {"vm.stats.vm.v_swappgsin", t_uint32, fromPage},
		"swap_out_bytes_total": {"vm.stats.vm.v_swappgsout", t_uint32, fromPage},
		"swap_size_bytes":      {"vm.swap_total", t_uint64, nil},
	} {
		switch v.dataType {
		case t_uint32:
			tmp32, err = unix.SysctlUint32(v.name)
			tmp64 = uint64(tmp32)
		case t_uint64:
			tmp64, err = unix.SysctlUint64(v.name)
		}
		if err != nil {
			return nil, err
		}

		if v.conversion != nil {
			// convert to bytes
			info[key] = float64(v.conversion(tmp64))
			continue
		}

		info[key] = float64(tmp64)
	}

	return info, nil
}
