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
	"log/slog"

	"github.com/prometheus/procfs"
)

type meminfoCollector struct {
	fs     procfs.FS
	logger *slog.Logger
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewMeminfoCollector(logger *slog.Logger) (Collector, error) {
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
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	metrics := make(map[string]float64)

	if meminfo.ActiveBytes != nil {
		metrics["Active_bytes"] = float64(*meminfo.ActiveBytes)
	}
	if meminfo.ActiveAnonBytes != nil {
		metrics["Active_anon_bytes"] = float64(*meminfo.ActiveAnonBytes)
	}
	if meminfo.ActiveFileBytes != nil {
		metrics["Active_file_bytes"] = float64(*meminfo.ActiveFileBytes)
	}
	if meminfo.AnonHugePagesBytes != nil {
		metrics["AnonHugePages_bytes"] = float64(*meminfo.AnonHugePagesBytes)
	}
	if meminfo.AnonPagesBytes != nil {
		metrics["AnonPages_bytes"] = float64(*meminfo.AnonPagesBytes)
	}
	if meminfo.BounceBytes != nil {
		metrics["Bounce_bytes"] = float64(*meminfo.BounceBytes)
	}
	if meminfo.BuffersBytes != nil {
		metrics["Buffers_bytes"] = float64(*meminfo.BuffersBytes)
	}
	if meminfo.CachedBytes != nil {
		metrics["Cached_bytes"] = float64(*meminfo.CachedBytes)
	}
	if meminfo.CmaFreeBytes != nil {
		metrics["CmaFree_bytes"] = float64(*meminfo.CmaFreeBytes)
	}
	if meminfo.CmaTotalBytes != nil {
		metrics["CmaTotal_bytes"] = float64(*meminfo.CmaTotalBytes)
	}
	if meminfo.CommitLimitBytes != nil {
		metrics["CommitLimit_bytes"] = float64(*meminfo.CommitLimitBytes)
	}
	if meminfo.CommittedASBytes != nil {
		metrics["Committed_AS_bytes"] = float64(*meminfo.CommittedASBytes)
	}
	if meminfo.DirectMap1GBytes != nil {
		metrics["DirectMap1G_bytes"] = float64(*meminfo.DirectMap1GBytes)
	}
	if meminfo.DirectMap2MBytes != nil {
		metrics["DirectMap2M_bytes"] = float64(*meminfo.DirectMap2MBytes)
	}
	if meminfo.DirectMap4kBytes != nil {
		metrics["DirectMap4k_bytes"] = float64(*meminfo.DirectMap4kBytes)
	}
	if meminfo.DirtyBytes != nil {
		metrics["Dirty_bytes"] = float64(*meminfo.DirtyBytes)
	}
	if meminfo.HardwareCorruptedBytes != nil {
		metrics["HardwareCorrupted_bytes"] = float64(*meminfo.HardwareCorruptedBytes)
	}
	if meminfo.HugepagesizeBytes != nil {
		metrics["Hugepagesize_bytes"] = float64(*meminfo.HugepagesizeBytes)
	}
	if meminfo.InactiveBytes != nil {
		metrics["Inactive_bytes"] = float64(*meminfo.InactiveBytes)
	}
	if meminfo.InactiveAnonBytes != nil {
		metrics["Inactive_anon_bytes"] = float64(*meminfo.InactiveAnonBytes)
	}
	if meminfo.InactiveFileBytes != nil {
		metrics["Inactive_file_bytes"] = float64(*meminfo.InactiveFileBytes)
	}
	if meminfo.KernelStackBytes != nil {
		metrics["KernelStack_bytes"] = float64(*meminfo.KernelStackBytes)
	}
	if meminfo.MappedBytes != nil {
		metrics["Mapped_bytes"] = float64(*meminfo.MappedBytes)
	}
	if meminfo.MemAvailableBytes != nil {
		metrics["MemAvailable_bytes"] = float64(*meminfo.MemAvailableBytes)
	}
	if meminfo.MemFreeBytes != nil {
		metrics["MemFree_bytes"] = float64(*meminfo.MemFreeBytes)
	}
	if meminfo.MemTotalBytes != nil {
		metrics["MemTotal_bytes"] = float64(*meminfo.MemTotalBytes)
	}
	if meminfo.MlockedBytes != nil {
		metrics["Mlocked_bytes"] = float64(*meminfo.MlockedBytes)
	}
	if meminfo.NFSUnstableBytes != nil {
		metrics["NFS_Unstable_bytes"] = float64(*meminfo.NFSUnstableBytes)
	}
	if meminfo.PageTablesBytes != nil {
		metrics["PageTables_bytes"] = float64(*meminfo.PageTablesBytes)
	}
	if meminfo.PercpuBytes != nil {
		metrics["Percpu_bytes"] = float64(*meminfo.PercpuBytes)
	}
	if meminfo.SReclaimableBytes != nil {
		metrics["SReclaimable_bytes"] = float64(*meminfo.SReclaimableBytes)
	}
	if meminfo.SUnreclaimBytes != nil {
		metrics["SUnreclaim_bytes"] = float64(*meminfo.SUnreclaimBytes)
	}
	if meminfo.ShmemBytes != nil {
		metrics["Shmem_bytes"] = float64(*meminfo.ShmemBytes)
	}
	if meminfo.ShmemHugePagesBytes != nil {
		metrics["ShmemHugePages_bytes"] = float64(*meminfo.ShmemHugePagesBytes)
	}
	if meminfo.ShmemPmdMappedBytes != nil {
		metrics["ShmemPmdMapped_bytes"] = float64(*meminfo.ShmemPmdMappedBytes)
	}
	if meminfo.SlabBytes != nil {
		metrics["Slab_bytes"] = float64(*meminfo.SlabBytes)
	}
	if meminfo.SwapCachedBytes != nil {
		metrics["SwapCached_bytes"] = float64(*meminfo.SwapCachedBytes)
	}
	if meminfo.SwapFreeBytes != nil {
		metrics["SwapFree_bytes"] = float64(*meminfo.SwapFreeBytes)
	}
	if meminfo.SwapTotalBytes != nil {
		metrics["SwapTotal_bytes"] = float64(*meminfo.SwapTotalBytes)
	}
	if meminfo.UnevictableBytes != nil {
		metrics["Unevictable_bytes"] = float64(*meminfo.UnevictableBytes)
	}
	if meminfo.VmallocChunkBytes != nil {
		metrics["VmallocChunk_bytes"] = float64(*meminfo.VmallocChunkBytes)
	}
	if meminfo.VmallocTotalBytes != nil {
		metrics["VmallocTotal_bytes"] = float64(*meminfo.VmallocTotalBytes)
	}
	if meminfo.VmallocUsedBytes != nil {
		metrics["VmallocUsed_bytes"] = float64(*meminfo.VmallocUsedBytes)
	}
	if meminfo.WritebackBytes != nil {
		metrics["Writeback_bytes"] = float64(*meminfo.WritebackBytes)
	}
	if meminfo.WritebackTmpBytes != nil {
		metrics["WritebackTmp_bytes"] = float64(*meminfo.WritebackTmpBytes)
	}
	if meminfo.ZswapBytes != nil {
		metrics["Zswap_bytes"] = float64(*meminfo.ZswapBytes)
	}
	if meminfo.ZswappedBytes != nil {
		metrics["Zswapped_bytes"] = float64(*meminfo.ZswappedBytes)
	}

	// These fields are always in bytes and do not have `Bytes`
	// suffixed counterparts in the procfs.Meminfo struct, nor do
	// they have `_bytes` suffix on the metric names.
	if meminfo.HugePagesFree != nil {
		metrics["HugePages_Free"] = float64(*meminfo.HugePagesFree)
	}
	if meminfo.HugePagesRsvd != nil {
		metrics["HugePages_Rsvd"] = float64(*meminfo.HugePagesRsvd)
	}
	if meminfo.HugePagesSurp != nil {
		metrics["HugePages_Surp"] = float64(*meminfo.HugePagesSurp)
	}
	if meminfo.HugePagesTotal != nil {
		metrics["HugePages_Total"] = float64(*meminfo.HugePagesTotal)
	}

	return metrics, nil
}
