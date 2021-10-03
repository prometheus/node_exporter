// Copyright 2018 The Prometheus Authors
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

//go:build solaris && !nozfs
// +build solaris,!nozfs

package collector

import (
	"strings"

	"github.com/go-kit/log"
	"github.com/illumos/go-kstat"
	"github.com/prometheus/client_golang/prometheus"
)

type zfsCollector struct {
	abdstatsLinearCount          *prometheus.Desc
	abdstatsLinearDataSize       *prometheus.Desc
	abdstatsScatterChunkWaste    *prometheus.Desc
	abdstatsScatterCount         *prometheus.Desc
	abdstatsScatterDataSize      *prometheus.Desc
	abdstatsStructSize           *prometheus.Desc
	arcstatsAnonSize             *prometheus.Desc
	arcstatsC                    *prometheus.Desc
	arcstatsCMax                 *prometheus.Desc
	arcstatsCMin                 *prometheus.Desc
	arcstatsDataSize             *prometheus.Desc
	arcstatsDemandDataHits       *prometheus.Desc
	arcstatsDemandDataMisses     *prometheus.Desc
	arcstatsDemandMetadataHits   *prometheus.Desc
	arcstatsDemandMetadataMisses *prometheus.Desc
	arcstatsHeaderSize           *prometheus.Desc
	arcstatsHits                 *prometheus.Desc
	arcstatsMisses               *prometheus.Desc
	arcstatsMFUGhostHits         *prometheus.Desc
	arcstatsMFUGhostSize         *prometheus.Desc
	arcstatsMFUSize              *prometheus.Desc
	arcstatsMRUGhostHits         *prometheus.Desc
	arcstatsMRUGhostSize         *prometheus.Desc
	arcstatsMRUSize              *prometheus.Desc
	arcstatsOtherSize            *prometheus.Desc
	arcstatsP                    *prometheus.Desc
	arcstatsSize                 *prometheus.Desc
	zfetchstatsHits              *prometheus.Desc
	zfetchstatsMisses            *prometheus.Desc
	logger                       log.Logger
}

const (
	zfsCollectorSubsystem = "zfs"
)

func init() {
	registerCollector("zfs", defaultEnabled, NewZfsCollector)
}

func NewZfsCollector(logger log.Logger) (Collector, error) {
	return &zfsCollector{
		abdstatsLinearCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_linear_count_total"),
			"ZFS ARC buffer data linear count", nil, nil,
		),
		abdstatsLinearDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_linear_data_bytes"),
			"ZFS ARC buffer data linear data size", nil, nil,
		),
		abdstatsScatterChunkWaste: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_chunk_waste_bytes"),
			"ZFS ARC buffer data scatter chunk waste", nil, nil,
		),
		abdstatsScatterCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_count_total"),
			"ZFS ARC buffer data scatter count", nil, nil,
		),
		abdstatsScatterDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_data_bytes"),
			"ZFS ARC buffer data scatter data size", nil, nil,
		),
		abdstatsStructSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_struct_bytes"),
			"ZFS ARC buffer data struct size", nil, nil,
		),
		arcstatsAnonSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_anon_bytes"),
			"ZFS ARC anon size", nil, nil,
		),
		arcstatsC: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_c_bytes"),
			"ZFS ARC target size", nil, nil,
		),
		arcstatsCMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_c_max_bytes"),
			"ZFS ARC maximum size", nil, nil,
		),
		arcstatsCMin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_c_min_bytes"),
			"ZFS ARC minimum size", nil, nil,
		),
		arcstatsDataSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_data_bytes"),
			"ZFS ARC data size", nil, nil,
		),
		arcstatsDemandDataHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_demand_data_hits_total"),
			"ZFS ARC demand data hits", nil, nil,
		),
		arcstatsDemandDataMisses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_demand_data_misses_total"),
			"ZFS ARC demand data misses", nil, nil,
		),
		arcstatsDemandMetadataHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_demand_metadata_hits_total"),
			"ZFS ARC demand metadata hits", nil, nil,
		),
		arcstatsDemandMetadataMisses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_demand_metadata_misses_total"),
			"ZFS ARC demand metadata misses", nil, nil,
		),
		arcstatsHeaderSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_hdr_bytes"),
			"ZFS ARC header size", nil, nil,
		),
		arcstatsHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_hits_total"),
			"ZFS ARC hits", nil, nil,
		),
		arcstatsMisses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_misses_total"),
			"ZFS ARC misses", nil, nil,
		),
		arcstatsMFUGhostHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_ghost_hits_total"),
			"ZFS ARC MFU ghost hits", nil, nil,
		),
		arcstatsMFUGhostSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_ghost_size"),
			"ZFS ARC MFU ghost size", nil, nil,
		),
		arcstatsMFUSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_bytes"),
			"ZFS ARC MFU size", nil, nil,
		),
		arcstatsMRUGhostHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_ghost_hits_total"),
			"ZFS ARC MRU ghost hits", nil, nil,
		),
		arcstatsMRUGhostSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_ghost_bytes"),
			"ZFS ARC MRU ghost size", nil, nil,
		),
		arcstatsMRUSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_bytes"),
			"ZFS ARC MRU size", nil, nil,
		),
		arcstatsOtherSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_other_bytes"),
			"ZFS ARC other size", nil, nil,
		),
		arcstatsP: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_p_bytes"),
			"ZFS ARC MRU target size", nil, nil,
		),
		arcstatsSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_size_bytes"),
			"ZFS ARC size", nil, nil,
		),
		zfetchstatsHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "zfetchstats_hits_total"),
			"ZFS cache fetch hits", nil, nil,
		),
		zfetchstatsMisses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "zfetchstats_misses_total"),
			"ZFS cache fetch misses", nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *zfsCollector) updateZfsAbdStats(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	ksZFSInfo, err := tok.Lookup("zfs", 0, "abdstats")
	if err != nil {
		return err
	}

	for k, v := range map[string]*prometheus.Desc{
		"linear_cnt":          c.abdstatsLinearCount,
		"linear_data_size":    c.abdstatsLinearDataSize,
		"scatter_chunk_waste": c.abdstatsScatterChunkWaste,
		"scatter_cnt":         c.abdstatsScatterCount,
		"scatter_data_size":   c.abdstatsScatterDataSize,
		"struct_size":         c.abdstatsStructSize,
	} {
		ksZFSInfoValue, err := ksZFSInfo.GetNamed(k)
		if err != nil {
			return err
		}

		if strings.HasSuffix(k, "_cnt") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}

		ch <- prometheus.MustNewConstMetric(
			v,
			metricType,
			float64(ksZFSInfoValue.UintVal),
		)
	}

	return nil
}

func (c *zfsCollector) updateZfsArcStats(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	ksZFSInfo, err := tok.Lookup("zfs", 0, "arcstats")
	if err != nil {
		return err
	}

	for k, v := range map[string]*prometheus.Desc{
		"anon_size":              c.arcstatsAnonSize,
		"c":                      c.arcstatsC,
		"c_max":                  c.arcstatsCMax,
		"c_min":                  c.arcstatsCMin,
		"data_size":              c.arcstatsDataSize,
		"demand_data_hits":       c.arcstatsDemandDataHits,
		"demand_data_misses":     c.arcstatsDemandDataMisses,
		"demand_metadata_hits":   c.arcstatsDemandMetadataHits,
		"demand_metadata_misses": c.arcstatsDemandMetadataMisses,
		"hdr_size":               c.arcstatsHeaderSize,
		"hits":                   c.arcstatsHits,
		"misses":                 c.arcstatsMisses,
		"mfu_ghost_hits":         c.arcstatsMFUGhostHits,
		"mfu_ghost_size":         c.arcstatsMFUGhostSize,
		"mfu_size":               c.arcstatsMFUSize,
		"mru_ghost_hits":         c.arcstatsMRUGhostHits,
		"mru_ghost_size":         c.arcstatsMRUGhostSize,
		"mru_size":               c.arcstatsMRUSize,
		"other_size":             c.arcstatsOtherSize,
		"p":                      c.arcstatsP,
		"size":                   c.arcstatsSize,
	} {
		ksZFSInfoValue, err := ksZFSInfo.GetNamed(k)
		if err != nil {
			return err
		}

		if strings.HasSuffix(k, "_hits") || strings.HasSuffix(k, "_misses") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}

		ch <- prometheus.MustNewConstMetric(
			v,
			metricType,
			float64(ksZFSInfoValue.UintVal),
		)
	}

	return nil
}

func (c *zfsCollector) updateZfsFetchStats(ch chan<- prometheus.Metric) error {
	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	ksZFSInfo, err := tok.Lookup("zfs", 0, "zfetchstats")
	if err != nil {
		return err
	}

	for k, v := range map[string]*prometheus.Desc{
		"hits":   c.zfetchstatsHits,
		"misses": c.zfetchstatsMisses,
	} {
		ksZFSInfoValue, err := ksZFSInfo.GetNamed(k)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			v,
			prometheus.CounterValue,
			float64(ksZFSInfoValue.UintVal),
		)
	}

	return nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateZfsAbdStats(ch); err != nil {
		return err
	}
	if err := c.updateZfsArcStats(ch); err != nil {
		return err
	}
	if err := c.updateZfsFetchStats(ch); err != nil {
		return err
	}
	return nil
}
