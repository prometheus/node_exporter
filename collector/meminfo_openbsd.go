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

//go:build !nomeminfo && !amd64
// +build !nomeminfo,!amd64

package collector

import (
	"fmt"
)

/*
#include <sys/param.h>
#include <sys/types.h>
#include <sys/mount.h>
#include <sys/sysctl.h>

int
sysctl_uvmexp(struct uvmexp *uvmexp)
{
        static int uvmexp_mib[] = {CTL_VM, VM_UVMEXP};
        size_t sz = sizeof(struct uvmexp);

        if(sysctl(uvmexp_mib, 2, uvmexp, &sz, NULL, 0) < 0)
                return -1;

        return 0;
}

int
sysctl_bcstats(struct bcachestats *bcstats)
{
        static int bcstats_mib[] = {CTL_VFS, VFS_GENERIC, VFS_BCACHESTAT};
        size_t sz = sizeof(struct bcachestats);

        if(sysctl(bcstats_mib, 3, bcstats, &sz, NULL, 0) < 0)
                return -1;

        return 0;
}

*/
import "C"

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	var uvmexp C.struct_uvmexp
	var bcstats C.struct_bcachestats

	if _, err := C.sysctl_uvmexp(&uvmexp); err != nil {
		return nil, fmt.Errorf("sysctl CTL_VM VM_UVMEXP failed: %w", err)
	}

	if _, err := C.sysctl_bcstats(&bcstats); err != nil {
		return nil, fmt.Errorf("sysctl CTL_VFS VFS_GENERIC VFS_BCACHESTAT failed: %w", err)
	}

	ps := float64(uvmexp.pagesize)

	// see uvm(9)
	return map[string]float64{
		"active_bytes":                  ps * float64(uvmexp.active),
		"cache_bytes":                   ps * float64(bcstats.numbufpages),
		"free_bytes":                    ps * float64(uvmexp.free),
		"inactive_bytes":                ps * float64(uvmexp.inactive),
		"size_bytes":                    ps * float64(uvmexp.npages),
		"swap_size_bytes":               ps * float64(uvmexp.swpages),
		"swap_used_bytes":               ps * float64(uvmexp.swpginuse),
		"swapped_in_pages_bytes_total":  ps * float64(uvmexp.pgswapin),
		"swapped_out_pages_bytes_total": ps * float64(uvmexp.pgswapout),
		"wired_bytes":                   ps * float64(uvmexp.wired),
	}, nil
}
