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

//go:build linux && !nogpu

package collector

import (
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/collector/asset/cmdb"
)

const gpuCollectorSubsystem = "gpu"

type gpuCollector struct {
	health           *prometheus.Desc
	memTotal         *prometheus.Desc
	memUsed          *prometheus.Desc
	memFree          *prometheus.Desc
	utilizationRatio *prometheus.Desc
	temperature      *prometheus.Desc
	power            *prometheus.Desc
	logger           *slog.Logger
}

func init() {
	registerCollector("gpu", defaultEnabled, NewGPUCollector)
}

func NewGPUCollector(logger *slog.Logger) (Collector, error) {
	return &gpuCollector{
		health: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "health"),
			"GPU health status as a constant 1 labeled by health state.",
			[]string{"device", "health"}, nil,
		),
		memTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "memory_total_bytes"),
			"Total GPU memory in bytes.",
			[]string{"device"}, nil,
		),
		memUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "memory_used_bytes"),
			"Used GPU memory in bytes.",
			[]string{"device"}, nil,
		),
		memFree: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "memory_free_bytes"),
			"Free GPU memory in bytes.",
			[]string{"device"}, nil,
		),
		utilizationRatio: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "utilization_ratio"),
			"GPU utilization as a 0-1 ratio.",
			[]string{"device"}, nil,
		),
		temperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "temperature_celsius"),
			"GPU temperature in degrees Celsius.",
			[]string{"device"}, nil,
		),
		power: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "power_watts"),
			"GPU power draw in watts.",
			[]string{"device"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *gpuCollector) Update(ch chan<- prometheus.Metric) error {
	g, err := cmdb.CollectGPU()
	if err != nil {
		return err
	}

	for _, dev := range g.Devices {
		if !dev.RuntimeMetrics {
			continue
		}
		idx := strconv.Itoa(dev.Index)
		ch <- prometheus.MustNewConstMetric(c.health, prometheus.GaugeValue, 1, idx, assetLabel(dev.Health))
		ch <- prometheus.MustNewConstMetric(c.memTotal, prometheus.GaugeValue, float64(dev.MemoryTotalMB)*1024*1024, idx)
		ch <- prometheus.MustNewConstMetric(c.memUsed, prometheus.GaugeValue, float64(dev.MemoryUsedMB)*1024*1024, idx)
		ch <- prometheus.MustNewConstMetric(c.memFree, prometheus.GaugeValue, float64(dev.MemoryFreeMB)*1024*1024, idx)
		ch <- prometheus.MustNewConstMetric(c.utilizationRatio, prometheus.GaugeValue, dev.Utilization/100, idx)
		ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, dev.Temperature, idx)
		ch <- prometheus.MustNewConstMetric(c.power, prometheus.GaugeValue, dev.PowerW, idx)
	}
	return nil
}
