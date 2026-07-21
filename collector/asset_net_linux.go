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

//go:build linux && !noasset_net

package collector

import (
	"log/slog"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

type assetNetCollector struct {
	info      *prometheus.Desc
	mtuBytes  *prometheus.Desc
	speedMbps *prometheus.Desc
	logger    *slog.Logger
}

func init() {
	registerCollector("asset_net", defaultEnabled, NewAssetNetCollector)
}

// NewAssetNetCollector returns a collector exposing physical NIC and bond
// identity and link attributes under siliconflow_asset_*. Only physical NICs and
// bond interfaces are reported (container/K8s virtual interfaces excluded by the
// vendored cmdb collector).
func NewAssetNetCollector(logger *slog.Logger) (Collector, error) {
	return &assetNetCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "net_info"),
			"A metric with a constant '1' value labeled by NIC identity (mac, up, physical, bond master, slaves, vendor, driver).",
			[]string{
				assetUUIDLabel, "name", "mac", "up", "physical", "master",
				"slaves", "vendor", "driver",
			},
			nil,
		),
		mtuBytes: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "net_mtu_bytes"),
			"NIC maximum transmission unit in bytes.",
			[]string{assetUUIDLabel, "name"}, nil,
		),
		speedMbps: prometheus.NewDesc(
			prometheus.BuildFQName(assetNamespace, "", "net_speed_mbps"),
			"NIC link speed in megabits per second (0 if unavailable).",
			[]string{assetUUIDLabel, "name"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *assetNetCollector) Update(ch chan<- prometheus.Metric) error {
	uuid, err := readAssetUUID()
	if err != nil {
		return err
	}
	n, err := cmdb.CollectNet()
	if err != nil {
		return err
	}

	for _, dev := range n.Devices {
		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1,
			uuid,
			assetLabel(dev.Name),
			assetLabel(dev.Mac),
			assetBool(dev.Up),
			assetBool(dev.Physical),
			assetLabel(dev.Master),
			strings.Join(dev.Slaves, ","),
			assetLabel(dev.Vendor),
			assetLabel(dev.Driver),
		)
		ch <- prometheus.MustNewConstMetric(c.mtuBytes, prometheus.GaugeValue, float64(dev.MTU), uuid, dev.Name)
		ch <- prometheus.MustNewConstMetric(c.speedMbps, prometheus.GaugeValue, float64(dev.SpeedMbps), uuid, dev.Name)
	}
	return nil
}
