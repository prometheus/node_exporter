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

//go:build !nomeminfo
// +build !nomeminfo

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/procfs"
)

type meminfoCollector struct {
	fs     procfs.FS
	logger log.Logger
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewMeminfoCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &meminfoCollector{
		logger: logger,
		fs:     fs,
	}, nil
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	meminfo, err := c.fs.Meminfo()
	if err != nil {
		return nil, fmt.Errorf("Failed to get memory info: %s", err)
	}

	return map[string]float64{
		"Active_bytes":            uint64PtrToFloat(meminfo.ActiveBytes),
		"Active_anon_bytes":       uint64PtrToFloat(meminfo.ActiveAnonBytes),
		"Active_file_bytes":       uint64PtrToFloat(meminfo.ActiveFileBytes),
		"AnonHugePages_bytes":     uint64PtrToFloat(meminfo.AnonHugePagesBytes),
		"AnonPages_bytes":         uint64PtrToFloat(meminfo.AnonPagesBytes),
		"Bounce_bytes":            uint64PtrToFloat(meminfo.BounceBytes),
		"Buffers_bytes":           uint64PtrToFloat(meminfo.BuffersBytes),
		"Cached_bytes":            uint64PtrToFloat(meminfo.CachedBytes),
		"CmaFree_bytes":           uint64PtrToFloat(meminfo.CmaFreeBytes),
		"CmaTotal_bytes":          uint64PtrToFloat(meminfo.CmaTotalBytes),
		"CommitLimit_bytes":       uint64PtrToFloat(meminfo.CommitLimitBytes),
		"Committed_AS_bytes":      uint64PtrToFloat(meminfo.CommittedASBytes),
		"DirectMap1G_bytes":       uint64PtrToFloat(meminfo.DirectMap1GBytes),
		"DirectMap2M_bytes":       uint64PtrToFloat(meminfo.DirectMap2MBytes),
		"DirectMap4k_bytes":       uint64PtrToFloat(meminfo.DirectMap4kBytes),
		"Dirty_bytes":             uint64PtrToFloat(meminfo.DirtyBytes),
		"HardwareCorrupted_bytes": uint64PtrToFloat(meminfo.HardwareCorruptedBytes),
		"Hugepagesize_bytes":      uint64PtrToFloat(meminfo.HugepagesizeBytes),
		"Inactive_bytes":          uint64PtrToFloat(meminfo.InactiveBytes),
		"Inactive_anon_bytes":     uint64PtrToFloat(meminfo.InactiveAnonBytes),
		"Inactive_file_bytes":     uint64PtrToFloat(meminfo.InactiveFileBytes),
		"KernelStack_bytes":       uint64PtrToFloat(meminfo.KernelStackBytes),
		"Mapped_bytes":            uint64PtrToFloat(meminfo.MappedBytes),
		"MemAvailable_bytes":      uint64PtrToFloat(meminfo.MemAvailableBytes),
		"MemFree_bytes":           uint64PtrToFloat(meminfo.MemFreeBytes),
		"MemTotal_bytes":          uint64PtrToFloat(meminfo.MemTotalBytes),
		"Mlocked_bytes":           uint64PtrToFloat(meminfo.MlockedBytes),
		"NFS_Unstable_bytes":      uint64PtrToFloat(meminfo.NFSUnstableBytes),
		"PageTables_bytes":        uint64PtrToFloat(meminfo.PageTablesBytes),
		"Percpu_bytes":            uint64PtrToFloat(meminfo.PercpuBytes),
		"SReclaimable_bytes":      uint64PtrToFloat(meminfo.SReclaimableBytes),
		"SUnreclaim_bytes":        uint64PtrToFloat(meminfo.SUnreclaimBytes),
		"Shmem_bytes":             uint64PtrToFloat(meminfo.ShmemBytes),
		"ShmemHugePages_bytes":    uint64PtrToFloat(meminfo.ShmemHugePagesBytes),
		"ShmemPmdMapped_bytes":    uint64PtrToFloat(meminfo.ShmemPmdMappedBytes),
		"Slab_bytes":              uint64PtrToFloat(meminfo.SlabBytes),
		"SwapCached_bytes":        uint64PtrToFloat(meminfo.SwapCachedBytes),
		"SwapFree_bytes":          uint64PtrToFloat(meminfo.SwapFreeBytes),
		"SwapTotal_bytes":         uint64PtrToFloat(meminfo.SwapTotalBytes),
		"Unevictable_bytes":       uint64PtrToFloat(meminfo.UnevictableBytes),
		"VmallocChunk_bytes":      uint64PtrToFloat(meminfo.VmallocChunkBytes),
		"VmallocTotal_bytes":      uint64PtrToFloat(meminfo.VmallocTotalBytes),
		"VmallocUsed_bytes":       uint64PtrToFloat(meminfo.VmallocUsedBytes),
		"Writeback_bytes":         uint64PtrToFloat(meminfo.WritebackBytes),
		"WritebackTmp_bytes":      uint64PtrToFloat(meminfo.WritebackTmpBytes),
		// These fields are always in bytes and do not have `Bytes`
		// suffixed counterparts in the procfs.Meminfo struct, nor do
		// they have `_bytes` suffix on the metric names.
		"HugePages_Free":  uint64PtrToFloat(meminfo.HugePagesFree),
		"HugePages_Rsvd":  uint64PtrToFloat(meminfo.HugePagesRsvd),
		"HugePages_Surp":  uint64PtrToFloat(meminfo.HugePagesSurp),
		"HugePages_Total": uint64PtrToFloat(meminfo.HugePagesTotal),
	}, nil
}
