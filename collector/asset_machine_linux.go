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

//go:build linux && !noasset_machine

package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

type assetMachineCollector struct {
	info   *prometheus.Desc
	uptime *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("asset_machine", defaultEnabled, NewAssetMachineCollector)
}

// NewAssetMachineCollector returns a collector exposing machine hardware
// identity (vendor/product/serial/board/kernel/...) under siliconflow_asset_*.
func NewAssetMachineCollector(logger *slog.Logger) (Collector, error) {
	return &assetMachineCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "machine_info"),
			"A metric with a constant '1' value labeled by machine hardware identity "+
				"(vendor, product, serial, SMBIOS product UUID, board, kernel, OS, type, k8s_node).",
			[]string{
				assetUUIDLabel, "vendor", "product", "version", "serial", "machine_uuid",
				"hostname", "kernel", "kernel_arch", "os", "os_version", "type", "k8s_node",
				"board_vendor", "board_name", "board_version", "board_serial",
			},
			nil,
		),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "machine_uptime_seconds"),
			"System uptime in seconds.",
			[]string{assetUUIDLabel}, nil,
		),
		logger: logger,
	}, nil
}

func (c *assetMachineCollector) Update(ch chan<- prometheus.Metric) error {
	uuid, err := readAssetUUID()
	if err != nil {
		return err
	}
	m, err := cmdb.CollectMachine()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1,
		uuid,
		assetLabel(m.Vendor),
		assetLabel(m.Product),
		assetLabel(m.Version),
		assetLabel(m.Serial),
		assetLabel(m.UUID),
		assetLabel(m.Hostname),
		assetLabel(m.Kernel),
		assetLabel(m.KernelArch),
		assetLabel(m.OS),
		assetLabel(m.OSVersion),
		assetLabel(m.Type),
		assetBool(m.K8sNode),
		assetLabel(m.BoardVendor),
		assetLabel(m.BoardName),
		assetLabel(m.BoardVersion),
		assetLabel(m.BoardSerial),
	)
	ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.GaugeValue, float64(m.Uptime), uuid)
	return nil
}
