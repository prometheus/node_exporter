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

// +build !nobcache

package collector

import (
	"fmt"

	// https://godoc.org/github.com/prometheus/client_golang/prometheus
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/bcache"
	"github.com/prometheus/procfs/sysfs"
)

func init() {
	Factories["bcache"] = NewBcacheCollector
}

// A bcacheCollector is a Collector which gathers metrics from Linux bcache.
type bcacheCollector struct {
	fs sysfs.FS
}

// NewBcacheCollector returns a newly allocated bcacheCollector.
// It exposes a number of Linux bcache statistics.
func NewBcacheCollector() (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %v", err)
	}

	return &bcacheCollector{
		fs: fs,
	}, nil
}

// Update reads and exposes bcache stats.
// It implements the Collector interface.
func (c *bcacheCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.BcacheStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve bcache stats: %v", err)
	}

	for _, s := range stats {
		c.updateBcacheStats(ch, s)
	}
	return nil
}

type bcacheMetric struct {
	name            string
	desc            string
	value           float64
	metricType      prometheus.ValueType
	extraLabel      []string
	extraLabelValue string
}

func bcachePeriodStatsToMetric(ps *bcache.PeriodStats, period string, labelValue string) []bcacheMetric {
	label := []string{"bdev_no"}
	var mT prometheus.ValueType

	if period == "total" {
		mT = prometheus.CounterValue
	} else {
		mT = prometheus.GaugeValue
	}

	metrics := []bcacheMetric{
		{
			name:            "bypassed_" + period,
			desc:            "Amount of IO (both reads and writes) that has bypassed the cache.",
			value:           float64(ps.Bypassed),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_hits_" + period,
			desc:            "Hits counted per individual IO as bcache sees them.",
			value:           float64(ps.CacheHits),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_misses_" + period,
			desc:            "Misses counted per individual IO as bcache sees them.",
			value:           float64(ps.CacheMisses),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_bypass_hits_" + period,
			desc:            "Hits for IO intended to skip the cache.",
			value:           float64(ps.CacheBypassHits),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_bypass_misses_" + period,
			desc:            "Misses for IO intended to skip the cache.",
			value:           float64(ps.CacheBypassMisses),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_miss_collisions_" + period,
			desc:            "Instances where data insertion from cache miss raced with write (data already present).",
			value:           float64(ps.CacheMissCollisions),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_readaheads_" + period,
			desc:            "Count of times readahead occurred.",
			value:           float64(ps.CacheReadaheads),
			metricType:      mT,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
	}
	return metrics
}

// UpdateBcacheStats collects statistics for one bcache ID.
func (c *bcacheCollector) updateBcacheStats(ch chan<- prometheus.Metric, s *bcache.Stats) {

	const (
		subsystem = "bcache"
	)

	var (
		devLabel   = []string{"uuid"}
		allMetrics []bcacheMetric
		metrics    []bcacheMetric
	)

	allMetrics = []bcacheMetric{
		// metrics in /sys/fs/bcache/<uuid>/
		{
			name:       "average_key_size_sectors",
			desc:       "Average data per key in the btree (sectors).",
			value:      float64(s.Bcache.AverageKeySize),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "btree_cache_size",
			desc:       "Amount of memory currently used by the btree cache.",
			value:      float64(s.Bcache.BtreeCacheSize),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "cache_available_percent",
			desc:       "Percentage of cache device without dirty data, useable for writeback (may contain clean cached data).",
			value:      float64(s.Bcache.CacheAvailablePercent),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "congested",
			desc:       "Congestion.",
			value:      float64(s.Bcache.Congested),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "root_usage_percent",
			desc:       "Percentage of the root btree node in use (tree depth increases if too high).",
			value:      float64(s.Bcache.RootUsagePercent),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "tree_depth",
			desc:       "Depth of the btree.",
			value:      float64(s.Bcache.TreeDepth),
			metricType: prometheus.GaugeValue,
		},
		// metrics in /sys/fs/bcache/<uuid>/internal/
		{
			name:       "active_journal_entries",
			desc:       "Number of journal entries that are newer than the index.",
			value:      float64(s.Bcache.Internal.ActiveJournalEntries),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "btree_nodes",
			desc:       "Depth of the btree.",
			value:      float64(s.Bcache.Internal.BtreeNodes),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "btree_read_average_duration_nanoseconds",
			desc:       "Average btree read duration.",
			value:      float64(s.Bcache.Internal.BtreeReadAverageDurationNanoSeconds),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "cache_read_races",
			desc:       "Counts instances where while data was being read from the cache, the bucket was reused and invalidated - i.e. where the pointer was stale after the read completed.",
			value:      float64(s.Bcache.Internal.CacheReadRaces),
			metricType: prometheus.CounterValue,
		},
	}

	metrics = bcachePeriodStatsToMetric(&s.Bcache.Total, "total", "all")
	allMetrics = append(allMetrics, metrics...)
	metrics = bcachePeriodStatsToMetric(&s.Bcache.FiveMin, "5_minutes", "all")
	allMetrics = append(allMetrics, metrics...)

	for _, bdev := range s.Bdevs {
		// metrics in /sys/fs/bcache/<uuid>/<bdev>/
		metrics = []bcacheMetric{
			{
				name:            "dirty_data",
				desc:            "Amount of dirty data for this backing device in the cache.",
				value:           float64(bdev.DirtyData),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"bdev_no"},
				extraLabelValue: bdev.Name,
			},
		}
		allMetrics = append(allMetrics, metrics...)

		// metrics in /sys/fs/bcache/<uuid>/<bdev>/stats_total
		metrics := bcachePeriodStatsToMetric(&bdev.Total, "total", bdev.Name)
		allMetrics = append(allMetrics, metrics...)

		// metrics in /sys/fs/bcache/<uuid>/<bdev>/stats_five_minute
		metrics = bcachePeriodStatsToMetric(&bdev.FiveMin, "5_minutes", bdev.Name)
		allMetrics = append(allMetrics, metrics...)
	}

	for _, cache := range s.Caches {
		metrics = []bcacheMetric{
			// metrics in /sys/fs/bcache/<uuid>/<cache>/
			{
				name:            "io_errors",
				desc:            "Number of errors that have occurred, decayed by io_error_halflife.",
				value:           float64(cache.IOErrors),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"cache_no"},
				extraLabelValue: cache.Name,
			},
			{
				name:            "metadata_written",
				desc:            "Sum of all non data writes (btree writes and all other metadata).",
				value:           float64(cache.MetadataWritten),
				metricType:      prometheus.CounterValue,
				extraLabel:      []string{"cache_no"},
				extraLabelValue: cache.Name,
			},
			{
				name:            "written",
				desc:            "Sum of all data that has been written to the cache.",
				value:           float64(cache.Written),
				metricType:      prometheus.CounterValue,
				extraLabel:      []string{"cache_no"},
				extraLabelValue: cache.Name,
			},
			// metrics in /sys/fs/bcache/<uuid>/<cache>/priority_stats
			{
				name:            "priority_stats_unused_percent",
				desc:            "The percentage of the cache that doesn't contain any data.",
				value:           float64(cache.Priority.UnusedPercent),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"cache_no"},
				extraLabelValue: cache.Name,
			},
			{
				name:            "priority_stats_metadata_percent",
				desc:            "Bcache's metadata overhead.",
				value:           float64(cache.Priority.MetadataPercent),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"cache_no"},
				extraLabelValue: cache.Name,
			},
		}
		allMetrics = append(allMetrics, metrics...)
	}

	for _, m := range allMetrics {
		labels := append(devLabel, m.extraLabel...)

		desc := prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, m.name),
			m.desc,
			labels,
			nil,
		)

		labelValues := []string{s.Name}
		if m.extraLabelValue != "" {
			labelValues = append(labelValues, m.extraLabelValue)
		}

		ch <- prometheus.MustNewConstMetric(
			desc,
			m.metricType,
			m.value,
			labelValues...,
		)
	}
}
