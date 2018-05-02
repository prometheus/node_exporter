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

// +build linux,gpu

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	nvml "github.com/prometheus/node_exporter/collector/nvml"
)

const (
	gpuCollectorSubsystem = "gpu"
)

type gpuCollector struct {
	gpuTimePercent    *prometheus.Desc // percentage of time during kernels are executing on the GPU.
	gpuMemTimePercent *prometheus.Desc // percentage of time during memory is being read or written.
	gpuMemUsage       *prometheus.Desc // percentage of used memory size
	gpuTemperature    *prometheus.Desc // GPU temperature in Celsius degrees
	gpuClockHz        *prometheus.Desc // GPU graphics clock in Hz
	gpuMemClockHz     *prometheus.Desc // GPU memory clock in Hz
	// GPU clock throttle reason.
	// The descriptions of the values can be seen in NvmlClocksThrottleReasons section in NVML API Reference.
	gpuThrottleReason *prometheus.Desc
	// GPU performance state (C.uint). 0 to 15. 0 for max performance, 15 for min performance. 32 for unknown.
	// The descriptions of the values can be seen in nvmlPstates_t in Device Enums section in NVML API Reference.
	gpuPerfState *prometheus.Desc
}

func init() {
	registerCollector("gpu", defaultEnabled, NewGPUCollector)
}

func NewGPUCollector() (Collector, error) {
	return &gpuCollector{
		gpuTimePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_processor_time_percent"),
			"Percentage of time during kernels are executing on the GPU.",
			[]string{"gpu"}, nil,
		),
		gpuMemTimePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_time_percent"),
			"Percentage of time during memory is being read or written.",
			[]string{"gpu"}, nil,
		),
		gpuMemUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_usage_percent"),
			"Percentage of used memory size.",
			[]string{"gpu"}, nil,
		),
		gpuTemperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_temperature"),
			"GPU temperature in Celsius degrees.",
			[]string{"gpu"}, nil,
		),
		gpuClockHz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_processor_clock_hz"),
			"GPU graphics clock in Hz.",
			[]string{"gpu"}, nil,
		),
		gpuMemClockHz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_clock_hz"),
			"GPU memory clock in Hz.",
			[]string{"gpu"}, nil,
		),
		gpuThrottleReason: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_throttle_reason"),
			"GPU clock throttle reason. (See NVML API Reference for details.)",
			[]string{"gpu"}, nil,
		),
		gpuPerfState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_perf_state"),
			"Performance state. 0, 1, 2,..15 or 32. 0 for max and 5 for min performance. 32 for unknown. (See NVML API Reference for details.)",
			[]string{"gpu"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes gpu related metrics with nvml library
func (c *gpuCollector) Update(ch chan<- prometheus.Metric) error {
	stats := nvml.GetGPUStats()

	for _, v := range stats {
		gpuID := fmt.Sprintf("gpu%d", v.ID)
		ch <- prometheus.MustNewConstMetric(c.gpuTimePercent, prometheus.CounterValue, float64(v.UtilGPU), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemTimePercent, prometheus.CounterValue, float64(v.UtilMem), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemUsage, prometheus.CounterValue, float64(v.MemUsage), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuTemperature, prometheus.CounterValue, float64(v.Temperature), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuClockHz, prometheus.CounterValue, float64(v.ClockGraphics*1e6), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemClockHz, prometheus.CounterValue, float64(v.ClockMem*1e6), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuThrottleReason, prometheus.CounterValue, float64(v.Throttle), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuPerfState, prometheus.CounterValue, float64(v.PerfState), gpuID)
	}

	return nil
}
