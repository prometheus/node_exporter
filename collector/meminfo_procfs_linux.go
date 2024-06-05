// Copyright 2024 The Prometheus Authors
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

//go:build !nomeminfo_procfs
// +build !nomeminfo_procfs

package collector

import (
	"fmt"
)

func (c *meminfoProcfsCollector) getMemInfo() (map[string]float64, error) {
	meminfo, err := c.fs.Meminfo()
	if err != nil {
		return nil, fmt.Errorf("Failed to get memory info: %s", err)
	}

	return map[string]float64{
		"MemTotal":          uint64PtrToFloat(meminfo.MemTotalBytes),
		"MemFree":           uint64PtrToFloat(meminfo.MemFreeBytes),
		"MemAvailable":      uint64PtrToFloat(meminfo.MemAvailableBytes),
		"Buffers":           uint64PtrToFloat(meminfo.BuffersBytes),
		"Cached":            uint64PtrToFloat(meminfo.CachedBytes),
		"SwapCached":        uint64PtrToFloat(meminfo.SwapCachedBytes),
		"Active":            uint64PtrToFloat(meminfo.ActiveBytes),
		"Inactive":          uint64PtrToFloat(meminfo.InactiveBytes),
		"ActiveAnon":        uint64PtrToFloat(meminfo.ActiveAnonBytes),
		"InactiveAnon":      uint64PtrToFloat(meminfo.InactiveAnonBytes),
		"ActiveFile":        uint64PtrToFloat(meminfo.ActiveFileBytes),
		"InactiveFile":      uint64PtrToFloat(meminfo.InactiveFileBytes),
		"Unevictable":       uint64PtrToFloat(meminfo.UnevictableBytes),
		"Mlocked":           uint64PtrToFloat(meminfo.MlockedBytes),
		"SwapTotal":         uint64PtrToFloat(meminfo.SwapTotalBytes),
		"SwapFree":          uint64PtrToFloat(meminfo.SwapFreeBytes),
		"Dirty":             uint64PtrToFloat(meminfo.DirtyBytes),
		"Writeback":         uint64PtrToFloat(meminfo.WritebackBytes),
		"AnonPages":         uint64PtrToFloat(meminfo.AnonPagesBytes),
		"Mapped":            uint64PtrToFloat(meminfo.MappedBytes),
		"Shmem":             uint64PtrToFloat(meminfo.ShmemBytes),
		"Slab":              uint64PtrToFloat(meminfo.SlabBytes),
		"SReclaimable":      uint64PtrToFloat(meminfo.SReclaimableBytes),
		"SUnreclaim":        uint64PtrToFloat(meminfo.SUnreclaimBytes),
		"KernelStack":       uint64PtrToFloat(meminfo.KernelStackBytes),
		"PageTables":        uint64PtrToFloat(meminfo.PageTablesBytes),
		"NFSUnstable":       uint64PtrToFloat(meminfo.NFSUnstableBytes),
		"Bounce":            uint64PtrToFloat(meminfo.BounceBytes),
		"WritebackTmp":      uint64PtrToFloat(meminfo.WritebackTmpBytes),
		"CommitLimit":       uint64PtrToFloat(meminfo.CommitLimitBytes),
		"CommittedAS":       uint64PtrToFloat(meminfo.CommittedASBytes),
		"VmallocTotal":      uint64PtrToFloat(meminfo.VmallocTotalBytes),
		"VmallocUsed":       uint64PtrToFloat(meminfo.VmallocUsedBytes),
		"VmallocChunk":      uint64PtrToFloat(meminfo.VmallocChunkBytes),
		"Percpu":            uint64PtrToFloat(meminfo.PercpuBytes),
		"HardwareCorrupted": uint64PtrToFloat(meminfo.HardwareCorruptedBytes),
		"AnonHugePages":     uint64PtrToFloat(meminfo.AnonHugePagesBytes),
		"ShmemHugePages":    uint64PtrToFloat(meminfo.ShmemHugePagesBytes),
		"ShmemPmdMapped":    uint64PtrToFloat(meminfo.ShmemPmdMappedBytes),
		"CmaTotal":          uint64PtrToFloat(meminfo.CmaTotalBytes),
		"CmaFree":           uint64PtrToFloat(meminfo.CmaFreeBytes),
		"Hugepagesize":      uint64PtrToFloat(meminfo.HugepagesizeBytes),
		"DirectMap4k":       uint64PtrToFloat(meminfo.DirectMap4kBytes),
		"DirectMap2M":       uint64PtrToFloat(meminfo.DirectMap2MBytes),
		"DirectMap1G":       uint64PtrToFloat(meminfo.DirectMap1GBytes),
	}, nil
}
