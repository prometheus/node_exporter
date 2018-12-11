// Copyright 2016 The Prometheus Authors
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

// +build solaris

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/siebenmann/go-kstat"
)

type zfsCollector struct {
	abdstatsLinearCount	*prometheus.Desc
	abdstatsLinearDataSize	*prometheus.Desc
	abdstatsScatterChunkWaste	*prometheus.Desc
	abdstatsScatterCount	*prometheus.Desc
	abdstatsScatterDataSize	*prometheus.Desc
	abdstatsStructSize	*prometheus.Desc
	arcstatsAnonSize	*prometheus.Desc
	arcstatsHeaderSize	*prometheus.Desc
	arcstatsHits		*prometheus.Desc
	arcstatsMisses		*prometheus.Desc
	arcstatsMFUGhostHits	*prometheus.Desc
	arcstatsMFUGhostSize	*prometheus.Desc
	arcstatsMFUSize		*prometheus.Desc
	arcstatsMRUGhostHits	*prometheus.Desc
	arcstatsMRUGhostSize	*prometheus.Desc
	arcstatsMRUSize		*prometheus.Desc
	arcstatsOtherSize	*prometheus.Desc
	arcstatsSize		*prometheus.Desc
	zfetchstatsHits		*prometheus.Desc
	zfetchstatsMisses	*prometheus.Desc
}

const (
        zfsCollectorSubsystem = "zfs"
)

func init() {
        registerCollector("zfs", defaultEnabled, NewZfsCollector)
}

func NewZfsCollector() (Collector, error) {
	return &zfsCollector{
                abdstatsLinearCount: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_linear_count"),
			"ZFS ARC buffer data linear count", nil, nil,
		),
                abdstatsLinearDataSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_linear_data_size"),
			"ZFS ARC buffer data linear data size", nil, nil,
		),
                abdstatsScatterChunkWaste: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_chunk_waste"),
			"ZFS ARC buffer data scatter chunk waste", nil, nil,
		),
                abdstatsScatterCount: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_count"),
			"ZFS ARC buffer data scatter count", nil, nil,
		),
                abdstatsScatterDataSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_scatter_data_size"),
			"ZFS ARC buffer data scatter data size", nil, nil,
		),
                abdstatsStructSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "abdstats_struct_size"),
			"ZFS ARC buffer data struct size", nil, nil,
		),
                arcstatsAnonSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_anon_size"),
			"ZFS ARC anon size", nil, nil,
		),
                arcstatsHeaderSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_hdr_size"),
			"ZFS ARC header size", nil, nil,
		),
                arcstatsHits: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_hits"),
			"ZFS ARC hits", nil, nil,
		),
                arcstatsMisses: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_misses"),
			"ZFS ARC misses", nil, nil,
		),
                arcstatsMFUGhostHits: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_ghost_hits"),
			"ZFS ARC MFU ghost hits", nil, nil,
		),
                arcstatsMFUGhostSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_ghost_size"),
			"ZFS ARC MFU ghost size", nil, nil,
		),
                arcstatsMFUSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mfu_size"),
			"ZFS ARC MFU size", nil, nil,
		),
                arcstatsMRUGhostHits: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_ghost_hits"),
			"ZFS ARC MRU ghost hits", nil, nil,
		),
                arcstatsMRUGhostSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_ghost_size"),
			"ZFS ARC MRU ghost size", nil, nil,
		),
                arcstatsMRUSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_mru_size"),
			"ZFS ARC MRU size", nil, nil,
		),
                arcstatsOtherSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_other_size"),
			"ZFS ARC other size", nil, nil,
		),
                arcstatsSize: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "arcstats_size"),
			"ZFS ARC size", nil, nil,
		),
                zfetchstatsHits: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "zfetchstats_hits"),
			"ZFS cache fetch hits", nil, nil,
		),
                zfetchstatsMisses: prometheus.NewDesc(
                        prometheus.BuildFQName(namespace, zfsCollectorSubsystem, "zfetchstats_misses"),
			"ZFS cache fetch misses", nil, nil,
		),
	}, nil
}

func (c *zfsCollector) updateZfsAbdStats(ch chan<- prometheus.Metric) (err error) {
	tok, err := kstat.Open()
        if err != nil {
                return err
        }

        defer tok.Close()

	ks_zfs_info, err := tok.Lookup("zfs", 0, "abdstats")
	if err != nil {
		return err
	}

	zfs_abdstats_linear_count_v, err := ks_zfs_info.GetNamed("linear_cnt")
	if err != nil {
		return err
	}

	zfs_abdstats_linear_data_size_v, err := ks_zfs_info.GetNamed("linear_data_size")
	if err != nil {
		return err
	}

	zfs_abdstats_scatter_chunk_waste_v, err := ks_zfs_info.GetNamed("scatter_chunk_waste")
	if err != nil {
		return err
	}

	zfs_abdstats_scatter_count_v, err := ks_zfs_info.GetNamed("scatter_cnt")
	if err != nil {
		return err
	}

	zfs_abdstats_scatter_data_size_v, err := ks_zfs_info.GetNamed("scatter_data_size")
	if err != nil {
		return err
	}

	zfs_abdstats_struct_size_v, err := ks_zfs_info.GetNamed("struct_size")
	if err != nil {
		return err
	}


	ch <- prometheus.MustNewConstMetric(
		c.abdstatsLinearCount,
		prometheus.GaugeValue,
		float64(zfs_abdstats_linear_count_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.abdstatsLinearDataSize,
		prometheus.GaugeValue,
		float64(zfs_abdstats_linear_data_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.abdstatsScatterChunkWaste,
		prometheus.GaugeValue,
		float64(zfs_abdstats_scatter_chunk_waste_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.abdstatsScatterCount,
		prometheus.GaugeValue,
		float64(zfs_abdstats_scatter_count_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.abdstatsScatterDataSize,
		prometheus.GaugeValue,
		float64(zfs_abdstats_scatter_data_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.abdstatsStructSize,
		prometheus.GaugeValue,
		float64(zfs_abdstats_struct_size_v.UintVal),
	)

	return err
}

func (c *zfsCollector) updateZfsArcStats(ch chan<- prometheus.Metric) (err error) {
	tok, err := kstat.Open()
        if err != nil {
                return err
        }

        defer tok.Close()

	ks_zfs_info, err := tok.Lookup("zfs", 0, "arcstats")
	if err != nil {
		return err
	}

	zfs_arcstats_anon_size_v, err := ks_zfs_info.GetNamed("anon_size")
	if err != nil {
		return err
	}

	zfs_arcstats_hdr_size_v, err := ks_zfs_info.GetNamed("hdr_size")
	if err != nil {
		return err
	}

	zfs_arcstats_hits_v, err := ks_zfs_info.GetNamed("hits")
	if err != nil {
		return err
	}

	zfs_arcstats_misses_v, err := ks_zfs_info.GetNamed("misses")
	if err != nil {
		return err
	}

	zfs_arcstats_mfu_ghost_hits_v, err := ks_zfs_info.GetNamed("mfu_ghost_hits")
	if err != nil {
		return err
	}

	zfs_arcstats_mfu_ghost_size_v, err := ks_zfs_info.GetNamed("mfu_ghost_size")
	if err != nil {
		return err
	}

	zfs_arcstats_mfu_size_v, err := ks_zfs_info.GetNamed("mfu_size")
	if err != nil {
		return err
	}

	zfs_arcstats_mru_ghost_hits_v, err := ks_zfs_info.GetNamed("mru_ghost_hits")
	if err != nil {
		return err
	}

	zfs_arcstats_mru_ghost_size_v, err := ks_zfs_info.GetNamed("mru_ghost_size")
	if err != nil {
		return err
	}

	zfs_arcstats_mru_size_v, err := ks_zfs_info.GetNamed("mru_size")
	if err != nil {
		return err
	}

	zfs_arcstats_size_v, err := ks_zfs_info.GetNamed("size")
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.arcstatsAnonSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_anon_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsHeaderSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_hdr_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsHits,
		prometheus.GaugeValue,
		float64(zfs_arcstats_hits_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMisses,
		prometheus.GaugeValue,
		float64(zfs_arcstats_misses_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMFUGhostHits,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mfu_ghost_hits_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMFUGhostSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mfu_ghost_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMFUSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mfu_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMRUGhostHits,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mru_ghost_hits_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMRUGhostSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mru_ghost_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsMRUSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_mru_size_v.UintVal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.arcstatsSize,
		prometheus.GaugeValue,
		float64(zfs_arcstats_size_v.UintVal),
	)

	return err
}

func (c *zfsCollector) updateZfsFetchStats(ch chan<- prometheus.Metric) (err error) {
	tok, err := kstat.Open()
        if err != nil {
                return err
        }

        defer tok.Close()

	ks_zfs_info, err := tok.Lookup("zfs", 0, "zfetchstats")

	zfs_fetchstats_hits_v, err := ks_zfs_info.GetNamed("hits")
	if err != nil {
		return err
	}

	zfs_fetchstats_misses_v, err := ks_zfs_info.GetNamed("misses")
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.zfetchstatsHits,
		prometheus.GaugeValue,
		float64(zfs_fetchstats_hits_v.UintVal),
	)

	ch <- prometheus.MustNewConstMetric(
		c.zfetchstatsMisses,
		prometheus.GaugeValue,
		float64(zfs_fetchstats_misses_v.UintVal),
	)

	return err
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {
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
