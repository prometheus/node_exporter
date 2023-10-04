// Copyright 2019 The Prometheus Authors
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

//go:build !nozfs
// +build !nozfs

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/sys/unix"
)

type zfsCollector struct {
	sysctls []bsdSysctl
	logger  log.Logger
}

const (
	zfsCollectorSubsystem = "zfs"
)

func init() {
	registerCollector("zfs", defaultEnabled, NewZfsCollector)
}

func NewZfsCollector(logger log.Logger) (Collector, error) {
	return &zfsCollector{
		sysctls: []bsdSysctl{
			{
				name:        "abdstats_linear_count_total",
				description: "ZFS ARC buffer data linear count",
				mib:         "kstat.zfs.misc.abdstats.linear_cnt",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "abdstats_linear_data_bytes",
				description: "ZFS ARC buffer data linear data size",
				mib:         "kstat.zfs.misc.abdstats.linear_data_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "abdstats_scatter_chunk_waste_bytes",
				description: "ZFS ARC buffer data scatter chunk waste",
				mib:         "kstat.zfs.misc.abdstats.scatter_chunk_waste",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "abdstats_scatter_count_total",
				description: "ZFS ARC buffer data scatter count",
				mib:         "kstat.zfs.misc.abdstats.scatter_cnt",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "abdstats_scatter_data_bytes",
				description: "ZFS ARC buffer data scatter data size",
				mib:         "kstat.zfs.misc.abdstats.scatter_data_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "abdstats_struct_bytes",
				description: "ZFS ARC buffer data struct size",
				mib:         "kstat.zfs.misc.abdstats.struct_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_anon_bytes",
				description: "ZFS ARC anon size",
				mib:         "kstat.zfs.misc.arcstats.anon_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_c_bytes",
				description: "ZFS ARC target size",
				mib:         "kstat.zfs.misc.arcstats.c",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_c_max_bytes",
				description: "ZFS ARC maximum size",
				mib:         "kstat.zfs.misc.arcstats.c_max",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_c_min_bytes",
				description: "ZFS ARC minimum size",
				mib:         "kstat.zfs.misc.arcstats.c_min",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_data_bytes",
				description: "ZFS ARC data size",
				mib:         "kstat.zfs.misc.arcstats.data_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_demand_data_hits_total",
				description: "ZFS ARC demand data hits",
				mib:         "kstat.zfs.misc.arcstats.demand_data_hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_demand_data_misses_total",
				description: "ZFS ARC demand data misses",
				mib:         "kstat.zfs.misc.arcstats.demand_data_misses",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_demand_metadata_hits_total",
				description: "ZFS ARC demand metadata hits",
				mib:         "kstat.zfs.misc.arcstats.demand_metadata_hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_demand_metadata_misses_total",
				description: "ZFS ARC demand metadata misses",
				mib:         "kstat.zfs.misc.arcstats.demand_metadata_misses",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_hdr_bytes",
				description: "ZFS ARC header size",
				mib:         "kstat.zfs.misc.arcstats.hdr_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_hits_total",
				description: "ZFS ARC hits",
				mib:         "kstat.zfs.misc.arcstats.hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_misses_total",
				description: "ZFS ARC misses",
				mib:         "kstat.zfs.misc.arcstats.misses",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mfu_ghost_hits_total",
				description: "ZFS ARC MFU ghost hits",
				mib:         "kstat.zfs.misc.arcstats.mfu_ghost_hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mfu_ghost_size",
				description: "ZFS ARC MFU ghost size",
				mib:         "kstat.zfs.misc.arcstats.mfu_ghost_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mfu_bytes",
				description: "ZFS ARC MFU size",
				mib:         "kstat.zfs.misc.arcstats.mfu_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mru_ghost_hits_total",
				description: "ZFS ARC MRU ghost hits",
				mib:         "kstat.zfs.misc.arcstats.mru_ghost_hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mru_ghost_bytes",
				description: "ZFS ARC MRU ghost size",
				mib:         "kstat.zfs.misc.arcstats.mru_ghost_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_mru_bytes",
				description: "ZFS ARC MRU size",
				mib:         "kstat.zfs.misc.arcstats.mru_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_other_bytes",
				description: "ZFS ARC other size",
				mib:         "kstat.zfs.misc.arcstats.other_size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			// when FreeBSD 14.0+, `meta/pm/pd` install of `p`.
			{
				name:        "arcstats_p_bytes",
				description: "ZFS ARC MRU target size",
				mib:         "kstat.zfs.misc.arcstats.p",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "arcstats_meta_bytes",
				description: "ZFS ARC metadata target frac ",
				mib:         "kstat.zfs.misc.arcstats.meta",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "arcstats_pd_bytes",
				description: "ZFS ARC data MRU target frac",
				mib:         "kstat.zfs.misc.arcstats.pd",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "arcstats_pm_bytes",
				description: "ZFS ARC meta MRU target frac",
				mib:         "kstat.zfs.misc.arcstats.pm",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "arcstats_size_bytes",
				description: "ZFS ARC size",
				mib:         "kstat.zfs.misc.arcstats.size",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.GaugeValue,
				labels:      nil,
			},
			{
				name:        "zfetchstats_hits_total",
				description: "ZFS cache fetch hits",
				mib:         "kstat.zfs.misc.zfetchstats.hits",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
			{
				name:        "zfetchstats_misses_total",
				description: "ZFS cache fetch misses",
				mib:         "kstat.zfs.misc.zfetchstats.misses",
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      nil,
			},
		},
		logger: logger,
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) error {
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			// debug logging
			level.Debug(c.logger).Log("name", m.name, "couldn't get sysctl:", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, zfsCollectorSubsystem, m.name),
				m.description,
				nil, nil,
			), m.valueType, v)
	}

	return nil
}

func (c *zfsCollector) parseFreeBSDPoolObjsetStats() error {

	sysCtlMetrics := []string{
		"nunlinked", "nunlinks", "nread", "reads", "nwritten", "writes",
	}
	zfsPoolMibPrefix := "kstat.zfs.pool.dataset"
	zfsDatasetNames := []string{}

	zfsDatasets, err := unix.Sysctl(zfsPoolMibPrefix)
	if err != nil {
		return fmt.Errorf("couldn't get sysctl: %w", err)
	}

	for dataset, _ := range zfsDatasets {
		if strings.HasSuffix(dataset, ".dataset_name") {
			zfsDatasetNames = append(zfsDatasetNames, strings.SplitAfter(dataset, ".")[3])
		}
	}

	for zpoolDataset := range zfsDatasetsNames {
		zfsDatasetLabels := map[string]string{
			"dataset": zpoolDataset,
			"zpool":   strings.SplitAfter(zpoolDataset, "/")[0],
		}
		for metric := range sysCtlMetrics {
			c.sysctls = append(c.sysctls, bsdSysctl{
				name:        fmt.SprintF("node_zfs_zpool_dataset_%s", metric),
				description: fmt.SprintF("node_zfs_zpool_dataset_%s", metric),
				mib:         fmt.Sprintf("%s.%s.%s", zfsPoolMibPrefix, poolObj, metric),
				dataType:    bsdSysctlTypeUint64,
				valueType:   prometheus.CounterValue,
				labels:      zfsDatasetLabels,
			})
		}
	}

	return nil
}
