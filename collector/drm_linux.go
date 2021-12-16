// Copyright 2021 The Prometheus Authors
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

//go:build !nogpu
// +build !nogpu

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const (
	drmCollectorSubsystem = "drm"
)

type drmCollector struct {
	fs                    sysfs.FS
	logger                log.Logger
	CardInfo              *prometheus.Desc
	GPUBusyPercent        *prometheus.Desc
	MemoryGTTSize         *prometheus.Desc
	MemoryGTTUsed         *prometheus.Desc
	MemoryVisibleVRAMSize *prometheus.Desc
	MemoryVisibleVRAMUsed *prometheus.Desc
	MemoryVRAMSize        *prometheus.Desc
	MemoryVRAMUsed        *prometheus.Desc
}

func init() {
	registerCollector("drm", defaultDisabled, NewDrmCollector)
}

// NewDrmCollector returns a new Collector exposing /sys/class/drm/card?/device stats.
func NewDrmCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &drmCollector{
		fs:     fs,
		logger: logger,
		CardInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "card_info"),
			"Card information",
			[]string{"card", "memory_vendor", "power_performance_level", "unique_id", "vendor"}, nil,
		),
		GPUBusyPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "gpu_busy_percent"),
			"How busy the GPU is as a percentage.",
			[]string{"card"}, nil,
		),
		MemoryGTTSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_gtt_size_bytes"),
			"The size of the graphics translation table (GTT) block in bytes.",
			[]string{"card"}, nil,
		),
		MemoryGTTUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_gtt_used_bytes"),
			"The used amount of the graphics translation table (GTT) block in bytes.",
			[]string{"card"}, nil,
		),
		MemoryVisibleVRAMSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_vis_vram_size_bytes"),
			"The size of visible VRAM in bytes.",
			[]string{"card"}, nil,
		),
		MemoryVisibleVRAMUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_vis_vram_used_bytes"),
			"The used amount of visible VRAM in bytes.",
			[]string{"card"}, nil,
		),
		MemoryVRAMSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_vram_size_bytes"),
			"The size of VRAM in bytes.",
			[]string{"card"}, nil,
		),
		MemoryVRAMUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, drmCollectorSubsystem, "memory_vram_used_bytes"),
			"The used amount of VRAM in bytes.",
			[]string{"card"}, nil,
		),
	}, nil
}

func (c *drmCollector) Update(ch chan<- prometheus.Metric) error {
	return c.updateAMDCards(ch)
}

func (c *drmCollector) updateAMDCards(ch chan<- prometheus.Metric) error {
	vendor := "amd"
	stats, err := c.fs.ClassDRMCardAMDGPUStats()
	if err != nil {
		return err
	}

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.CardInfo, prometheus.GaugeValue, 1,
			s.Name, s.MemoryVRAMVendor, s.PowerDPMForcePerformanceLevel, s.UniqueID, vendor)

		ch <- prometheus.MustNewConstMetric(
			c.GPUBusyPercent, prometheus.GaugeValue, float64(s.GPUBusyPercent), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryGTTSize, prometheus.GaugeValue, float64(s.MemoryGTTSize), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryGTTUsed, prometheus.GaugeValue, float64(s.MemoryGTTUsed), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryVRAMSize, prometheus.GaugeValue, float64(s.MemoryVRAMSize), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryVRAMUsed, prometheus.GaugeValue, float64(s.MemoryVRAMUsed), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryVisibleVRAMSize, prometheus.GaugeValue, float64(s.MemoryVisibleVRAMSize), s.Name)

		ch <- prometheus.MustNewConstMetric(
			c.MemoryVisibleVRAMUsed, prometheus.GaugeValue, float64(s.MemoryVisibleVRAMUsed), s.Name)
	}

	return nil
}
