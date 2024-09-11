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

//go:build !nobcache
// +build !nobcache

package collector

import (
	"fmt"
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/bcache"
)

var (
	priorityStats = kingpin.Flag("collector.bcache.priorityStats", "Expose expensive priority stats.").Bool()
)

func init() {
	registerCollector("bcache", defaultEnabled, NewBcacheCollector)
}

// A bcacheCollector is a Collector which gathers metrics from Linux bcache.
type bcacheCollector struct {
	fs     bcache.FS
	logger *slog.Logger
}

// NewBcacheCollector returns a newly allocated bcacheCollector.
// It exposes a number of Linux bcache statistics.
func NewBcacheCollector(logger *slog.Logger) (Collector, error) {
	fs, err := bcache.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &bcacheCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update reads and exposes bcache stats.
// It implements the Collector interface.
func (c *bcacheCollector) Update(ch chan<- prometheus.Metric) error {
	var stats []*bcache.Stats
	var err error
	if *priorityStats {
		stats, err = c.fs.Stats()
	} else {
		stats, err = c.fs.StatsWithoutPriority()
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve bcache stats: %w", err)
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

func bcachePeriodStatsToMetric(ps *bcache.PeriodStats, labelValue string) []bcacheMetric {
	label := []string{"backing_device"}

	metrics := []bcacheMetric{
		{
			name:            "bypassed_bytes_total",
			desc:            "Amount of IO (both reads and writes) that has bypassed the cache.",
			value:           float64(ps.Bypassed),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_hits_total",
			desc:            "Hits counted per individual IO as bcache sees them.",
			value:           float64(ps.CacheHits),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_misses_total",
			desc:            "Misses counted per individual IO as bcache sees them.",
			value:           float64(ps.CacheMisses),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_bypass_hits_total",
			desc:            "Hits for IO intended to skip the cache.",
			value:           float64(ps.CacheBypassHits),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_bypass_misses_total",
			desc:            "Misses for IO intended to skip the cache.",
			value:           float64(ps.CacheBypassMisses),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
		{
			name:            "cache_miss_collisions_total",
			desc:            "Instances where data insertion from cache miss raced with write (data already present).",
			value:           float64(ps.CacheMissCollisions),
			metricType:      prometheus.CounterValue,
			extraLabel:      label,
			extraLabelValue: labelValue,
		},
	}
	if ps.CacheReadaheads != 0 {
		bcacheReadaheadMetrics := []bcacheMetric{
			{
				name:            "cache_readaheads_total",
				desc:            "Count of times readahead occurred.",
				value:           float64(ps.CacheReadaheads),
				metricType:      prometheus.CounterValue,
				extraLabel:      label,
				extraLabelValue: labelValue,
			},
		}
		metrics = append(metrics, bcacheReadaheadMetrics...)
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
			name:       "btree_cache_size_bytes",
			desc:       "Amount of memory currently used by the btree cache.",
			value:      float64(s.Bcache.BtreeCacheSize),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "cache_available_percent",
			desc:       "Percentage of cache device without dirty data, usable for writeback (may contain clean cached data).",
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
			desc:       "Total nodes in the btree.",
			value:      float64(s.Bcache.Internal.BtreeNodes),
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "btree_read_average_duration_seconds",
			desc:       "Average btree read duration.",
			value:      float64(s.Bcache.Internal.BtreeReadAverageDurationNanoSeconds) * 1e-9,
			metricType: prometheus.GaugeValue,
		},
		{
			name:       "cache_read_races_total",
			desc:       "Counts instances where while data was being read from the cache, the bucket was reused and invalidated - i.e. where the pointer was stale after the read completed.",
			value:      float64(s.Bcache.Internal.CacheReadRaces),
			metricType: prometheus.CounterValue,
		},
	}

	for _, bdev := range s.Bdevs {
		// metrics in /sys/fs/bcache/<uuid>/<bdev>/
		metrics = []bcacheMetric{
			{
				name:            "dirty_data_bytes",
				desc:            "Amount of dirty data for this backing device in the cache.",
				value:           float64(bdev.DirtyData),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
			{
				name:            "dirty_target_bytes",
				desc:            "Current dirty data target threshold for this backing device in bytes.",
				value:           float64(bdev.WritebackRateDebug.Target),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
			{
				name:            "writeback_rate",
				desc:            "Current writeback rate for this backing device in bytes.",
				value:           float64(bdev.WritebackRateDebug.Rate),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
			{
				name:            "writeback_rate_proportional_term",
				desc:            "Current result of proportional controller, part of writeback rate",
				value:           float64(bdev.WritebackRateDebug.Proportional),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
			{
				name:            "writeback_rate_integral_term",
				desc:            "Current result of integral controller, part of writeback rate",
				value:           float64(bdev.WritebackRateDebug.Integral),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
			{
				name:            "writeback_change",
				desc:            "Last writeback rate change step for this backing device.",
				value:           float64(bdev.WritebackRateDebug.Change),
				metricType:      prometheus.GaugeValue,
				extraLabel:      []string{"backing_device"},
				extraLabelValue: bdev.Name,
			},
		}
		allMetrics = append(allMetrics, metrics...)

		// metrics in /sys/fs/bcache/<uuid>/<bdev>/stats_total
		metrics := bcachePeriodStatsToMetric(&bdev.Total, bdev.Name)
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
				extraLabel:      []string{"cache_device"},
				extraLabelValue: cache.Name,
			},
			{
				name:            "metadata_written_bytes_total",
				desc:            "Sum of all non data writes (btree writes and all other metadata).",
				value:           float64(cache.MetadataWritten),
				metricType:      prometheus.CounterValue,
				extraLabel:      []string{"cache_device"},
				extraLabelValue: cache.Name,
			},
			{
				name:            "written_bytes_total",
				desc:            "Sum of all data that has been written to the cache.",
				value:           float64(cache.Written),
				metricType:      prometheus.CounterValue,
				extraLabel:      []string{"cache_device"},
				extraLabelValue: cache.Name,
			},
		}
		if *priorityStats {
			// metrics in /sys/fs/bcache/<uuid>/<cache>/priority_stats
			priorityStatsMetrics := []bcacheMetric{
				{
					name:            "priority_stats_unused_percent",
					desc:            "The percentage of the cache that doesn't contain any data.",
					value:           float64(cache.Priority.UnusedPercent),
					metricType:      prometheus.GaugeValue,
					extraLabel:      []string{"cache_device"},
					extraLabelValue: cache.Name,
				},
				{
					name:            "priority_stats_metadata_percent",
					desc:            "Bcache's metadata overhead.",
					value:           float64(cache.Priority.MetadataPercent),
					metricType:      prometheus.GaugeValue,
					extraLabel:      []string{"cache_device"},
					extraLabelValue: cache.Name,
				},
			}
			metrics = append(metrics, priorityStatsMetrics...)
		}
		allMetrics = append(allMetrics, metrics...)
	}

	for _, m := range allMetrics {
		labels := append(devLabel, m.extraLabel...)

		desc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, m.name),
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
