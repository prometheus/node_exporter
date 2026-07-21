// Copyright 2025 The Prometheus Authors
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

//go:build linux && !noasset_disk

package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

type assetDiskCollector struct {
	info   *prometheus.Desc
	sizeGB *prometheus.Desc
	usedGB *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("asset_disk", defaultEnabled, NewAssetDiskCollector)
}

// NewAssetDiskCollector returns a collector exposing block device identity and
// capacity/usage under siliconflow_asset_*.
func NewAssetDiskCollector(logger *slog.Logger) (Collector, error) {
	return &assetDiskCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "disk_info"),
			"A metric with a constant '1' value labeled by block device identity (name, model, vendor, serial, mountpoint, fs type).",
			[]string{
				assetUUIDLabel, "name", "type", "model", "vendor", "serial",
				"mountpoint", "fs_type",
			},
			nil,
		),
		sizeGB: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "disk_size_gb"),
			"Block device capacity in gigabytes (1 GB = 10^9 bytes).",
			[]string{assetUUIDLabel, "name"}, nil,
		),
		usedGB: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "disk_used_gb"),
			"Used gigabytes on the block device's mounted filesystem (0 if unmounted). 1 GB = 10^9 bytes.",
			[]string{assetUUIDLabel, "name"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *assetDiskCollector) Update(ch chan<- prometheus.Metric) error {
	uuid, err := readAssetUUID()
	if err != nil {
		return err
	}
	d, err := cmdb.CollectDisk()
	if err != nil {
		return err
	}

	for _, dev := range d.Devices {
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1,
			uuid,
			assetLabel(dev.Name),
			assetLabel(dev.Type),
			assetLabel(dev.Model),
			assetLabel(dev.Vendor),
			assetLabel(dev.Serial),
			assetLabel(dev.Mountpoint),
			assetLabel(dev.FsType),
		)
		ch <- prometheus.MustNewConstMetric(c.sizeGB, prometheus.GaugeValue, float64(dev.SizeBytes)/1e9, uuid, dev.Name)
		ch <- prometheus.MustNewConstMetric(c.usedGB, prometheus.GaugeValue, float64(dev.UsedBytes)/1e9, uuid, dev.Name)
	}
	return nil
}
