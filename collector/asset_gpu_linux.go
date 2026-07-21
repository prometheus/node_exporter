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

//go:build linux && !noasset_gpu

package collector

import (
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

type assetGPUCollector struct {
	info   *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("asset_gpu", defaultEnabled, NewAssetGPUCollector)
}

// NewAssetGPUCollector returns a collector exposing GPU/NPU identity under
// siliconflow_asset_*. Runtime metrics (memory/utilization/temperature/power)
// are exposed by the separate "gpu" collector under node_gpu_*. NVIDIA, Huawei
// NPU and a catch-all lspci fallback are handled by the vendored cmdb collector.
func NewAssetGPUCollector(logger *slog.Logger) (Collector, error) {
	return &assetGPUCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "gpu_info"),
			"A metric with a constant '1' value labeled by per-device GPU/NPU identity (vendor, name, serial, UUID, health, driver/firmware version).",
			[]string{
				assetUUIDLabel, "index", "vendor", "name", "serial", "gpu_uuid",
				"driver_version", "firmware_version",
			},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *assetGPUCollector) Update(ch chan<- prometheus.Metric) error {
	uuid, err := readAssetUUID()
	if err != nil {
		return err
	}
	g, err := cmdb.CollectGPU()
	if err != nil {
		return err
	}

	for _, dev := range g.Devices {
		idx := strconv.Itoa(dev.Index)
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1,
			uuid, idx,
			assetLabel(dev.Vendor),
			assetLabel(dev.Name),
			assetLabel(dev.Serial),
			assetLabel(dev.UUID),
			assetLabel(dev.DriverVersion),
			assetLabel(dev.FirmwareVersion),
		)
	}
	return nil
}
