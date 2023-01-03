// Copyright 2023 The Prometheus Authors
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

//go:build netbsd && !nomeminfo
// +build netbsd,!nomeminfo

package collector

import (
	"golang.org/x/sys/unix"
)

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	uvmexp, err := unix.SysctlUvmexp("vm.uvmexp2")
	if err != nil {
		return nil, err
	}

	ps := float64(uvmexp.Pagesize)

	// see uvm(9)
	return map[string]float64{
		"active_bytes":                  ps * float64(uvmexp.Active),
		"free_bytes":                    ps * float64(uvmexp.Free),
		"inactive_bytes":                ps * float64(uvmexp.Inactive),
		"size_bytes":                    ps * float64(uvmexp.Npages),
		"swap_size_bytes":               ps * float64(uvmexp.Swpages),
		"swap_used_bytes":               ps * float64(uvmexp.Swpginuse),
		"swapped_in_pages_bytes_total":  ps * float64(uvmexp.Pgswapin),
		"swapped_out_pages_bytes_total": ps * float64(uvmexp.Pgswapout),
		"wired_bytes":                   ps * float64(uvmexp.Wired),
	}, nil
}
