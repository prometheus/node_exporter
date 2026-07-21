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

//go:build linux && !noasset_memory

package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

type assetMemoryCollector struct {
	moduleInfo *prometheus.Desc
	totalMB    *prometheus.Desc
	cache      assetCache[*cmdb.Memory]
	logger     *slog.Logger
}

func init() {
	registerCollector("asset_memory", defaultEnabled, NewAssetMemoryCollector)
}

// NewAssetMemoryCollector returns a collector exposing memory size and per-DIMM
// module identity under siliconflow_asset_*. DIMM details are only collected on
// physical machines (synthetic on VMs), so it first determines the machine type
// via CollectMachine and forwards it to CollectMemory.
func NewAssetMemoryCollector(logger *slog.Logger) (Collector, error) {
	return &assetMemoryCollector{
		moduleInfo: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "memory_module_info"),
			"A metric with a constant '1' value labeled by per-DIMM identity (locator, size, type, speed, manufacturer, serial, part number).",
			[]string{
				assetUUIDLabel, "locator", "bank_locator", "size", "type", "speed",
				"manufacturer", "serial", "part_number",
			},
			nil,
		),
		totalMB: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "memory_total_mb"),
			"Total physical memory in megabytes.",
			[]string{assetUUIDLabel}, nil,
		),
		logger: logger,
	}, nil
}

func (c *assetMemoryCollector) Update(ch chan<- prometheus.Metric) error {
	uuid, err := readAssetUUID()
	if err != nil {
		return err
	}
	// CollectMemory takes the machine type so it can skip dmidecode on virtual
	// machines (whose SMBIOS data is synthetic). Determine it inside the cache
	// fetch so the resolved Memory (including the dmidecode result) is cached.
	mem, err := c.cache.get(*assetCacheTTL, func() (*cmdb.Memory, error) {
		machineType := ""
		if m, e := cmdb.CollectMachine(); e == nil {
			machineType = m.Type
		}
		return cmdb.CollectMemory(machineType)
	})
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.totalMB, prometheus.GaugeValue, float64(mem.TotalBytes)/1024/1024, uuid)

	for _, mod := range mem.Modules {
		ch <- prometheus.MustNewConstMetric(c.moduleInfo, prometheus.GaugeValue, 1,
			uuid,
			assetLabel(mod.Locator),
			assetLabel(mod.BankLocator),
			assetLabel(mod.Size),
			assetLabel(mod.Type),
			assetLabel(mod.Speed),
			assetLabel(mod.Manufacturer),
			assetLabel(mod.Serial),
			assetLabel(mod.PartNumber),
		)
	}
	return nil
}
