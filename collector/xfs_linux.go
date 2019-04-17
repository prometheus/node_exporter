// Copyright 2017 The Prometheus Authors
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

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/xfs"
)

// An xfsCollector is a Collector which gathers metrics from XFS filesystems.
type xfsCollector struct {
	fs xfs.FS
}

func init() {
	registerCollector("xfs", defaultEnabled, NewXFSCollector)
}

// NewXFSCollector returns a new Collector exposing XFS statistics.
func NewXFSCollector() (Collector, error) {
	fs, err := xfs.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %v", err)
	}

	return &xfsCollector{
		fs: fs,
	}, nil
}

// Update implements Collector.
func (c *xfsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.SysStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve XFS stats: %v", err)
	}

	for _, s := range stats {
		c.updateXFSStats(ch, s)
	}

	return nil
}

// updateXFSStats collects statistics for a single XFS filesystem.
func (c *xfsCollector) updateXFSStats(ch chan<- prometheus.Metric, s *xfs.Stats) {
	const (
		subsystem = "xfs"
	)

	var (
		labels = []string{"device"}
	)

	// Metric names and descriptions are sourced from:
	// http://xfs.org/index.php/Runtime_Stats.
	//
	// Each metric has a name that roughly follows the pattern of
	// "node_xfs_category_value_total", using the categories and value names
	// found on the XFS wiki.
	//
	// Note that statistics for more than one internal B-tree are measured,
	// and as such, each one must be differentiated by name.
	metrics := []struct {
		name  string
		desc  string
		value float64
	}{
		{
			name:  "extent_allocation_extents_allocated_total",
			desc:  "Number of extents allocated for a filesystem.",
			value: float64(s.ExtentAllocation.ExtentsAllocated),
		},
		{
			name:  "extent_allocation_blocks_allocated_total",
			desc:  "Number of blocks allocated for a filesystem.",
			value: float64(s.ExtentAllocation.BlocksAllocated),
		},
		{
			name:  "extent_allocation_extents_freed_total",
			desc:  "Number of extents freed for a filesystem.",
			value: float64(s.ExtentAllocation.ExtentsFreed),
		},
		{
			name:  "extent_allocation_blocks_freed_total",
			desc:  "Number of blocks freed for a filesystem.",
			value: float64(s.ExtentAllocation.BlocksFreed),
		},
		{
			name:  "allocation_btree_lookups_total",
			desc:  "Number of allocation B-tree lookups for a filesystem.",
			value: float64(s.AllocationBTree.Lookups),
		},
		{
			name:  "allocation_btree_compares_total",
			desc:  "Number of allocation B-tree compares for a filesystem.",
			value: float64(s.AllocationBTree.Compares),
		},
		{
			name:  "allocation_btree_records_inserted_total",
			desc:  "Number of allocation B-tree records inserted for a filesystem.",
			value: float64(s.AllocationBTree.RecordsInserted),
		},
		{
			name:  "allocation_btree_records_deleted_total",
			desc:  "Number of allocation B-tree records deleted for a filesystem.",
			value: float64(s.AllocationBTree.RecordsDeleted),
		},
		{
			name:  "block_mapping_reads_total",
			desc:  "Number of block map for read operations for a filesystem.",
			value: float64(s.BlockMapping.Reads),
		},
		{
			name:  "block_mapping_writes_total",
			desc:  "Number of block map for write operations for a filesystem.",
			value: float64(s.BlockMapping.Writes),
		},
		{
			name:  "block_mapping_unmaps_total",
			desc:  "Number of block unmaps (deletes) for a filesystem.",
			value: float64(s.BlockMapping.Unmaps),
		},
		{
			name:  "block_mapping_extent_list_insertions_total",
			desc:  "Number of extent list insertions for a filesystem.",
			value: float64(s.BlockMapping.ExtentListInsertions),
		},
		{
			name:  "block_mapping_extent_list_deletions_total",
			desc:  "Number of extent list deletions for a filesystem.",
			value: float64(s.BlockMapping.ExtentListDeletions),
		},
		{
			name:  "block_mapping_extent_list_lookups_total",
			desc:  "Number of extent list lookups for a filesystem.",
			value: float64(s.BlockMapping.ExtentListLookups),
		},
		{
			name:  "block_mapping_extent_list_compares_total",
			desc:  "Number of extent list compares for a filesystem.",
			value: float64(s.BlockMapping.ExtentListCompares),
		},
		{
			name:  "block_map_btree_lookups_total",
			desc:  "Number of block map B-tree lookups for a filesystem.",
			value: float64(s.BlockMapBTree.Lookups),
		},
		{
			name:  "block_map_btree_compares_total",
			desc:  "Number of block map B-tree compares for a filesystem.",
			value: float64(s.BlockMapBTree.Compares),
		},
		{
			name:  "block_map_btree_records_inserted_total",
			desc:  "Number of block map B-tree records inserted for a filesystem.",
			value: float64(s.BlockMapBTree.RecordsInserted),
		},
		{
			name:  "block_map_btree_records_deleted_total",
			desc:  "Number of block map B-tree records deleted for a filesystem.",
			value: float64(s.BlockMapBTree.RecordsDeleted),
		},
	}

	for _, m := range metrics {
		desc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, m.name),
			m.desc,
			labels,
			nil,
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.CounterValue,
			m.value,
			s.Name,
		)
	}
}
