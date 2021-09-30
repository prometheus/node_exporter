// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License")
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

//go:build !nomeminfo
// +build !nomeminfo

package collector

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

const (
	CTL_VFS        = 10
	VFS_GENERIC    = 0
	VFS_BCACHESTAT = 3
)

type bcachestats struct {
	Numbufs        int64
	Numbufpages    int64
	Numdirtypages  int64
	Numcleanpages  int64
	Pendingwrites  int64
	Pendingreads   int64
	Numwrites      int64
	Numreads       int64
	Cachehits      int64
	Busymapped     int64
	Dmapages       int64
	Highpages      int64
	Delwribufs     int64
	Kvaslots       int64
	Kvaslots_avail int64
	Highflips      int64
	Highflops      int64
	Dmaflips       int64
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	uvmexpb, err := unix.SysctlRaw("vm.uvmexp")
	if err != nil {
		return nil, err
	}

	mib := [3]_C_int{CTL_VFS, VFS_GENERIC, VFS_BCACHESTAT}
	bcstatsb, err := sysctl(mib[:])
	if err != nil {
		return nil, err
	}

	uvmexp := *(*unix.Uvmexp)(unsafe.Pointer(&uvmexpb[0]))
	ps := float64(uvmexp.Pagesize)

	bcstats := *(*bcachestats)(unsafe.Pointer(&bcstatsb[0]))

	// see uvm(9)
	return map[string]float64{
		"active_bytes":                  ps * float64(uvmexp.Active),
		"cache_bytes":                   ps * float64(bcstats.Numbufpages),
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
