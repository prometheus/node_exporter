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
	gpuUtil        *prometheus.Desc // percentage of time during kernels are executing on the GPU.
	gpuMemUtil     *prometheus.Desc // percentage of time during memory is being read or written.
	gpuMemUsage    *prometheus.Desc // percentage of used memory size
	gpuTemperature *prometheus.Desc // GPU temperature in Celsius degrees
	gpuClockMhz    *prometheus.Desc // GPU graphics clock in Mhz
	gpuMemClockMhz *prometheus.Desc // GPU memory clock in Mhz
	gpuThrottle    *prometheus.Desc // throttle reason
	gpuPerfState   *prometheus.Desc // performance state    C.uint 0: max / 15: min
}

func init() {
	registerCollector("gpu", defaultEnabled, NewGPUCollector)
}

func NewGPUCollector() (Collector, error) {
	return &gpuCollector{
		gpuUtil: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_util"),
			"Percentage of time during kernels are executing on the GPU.",
			[]string{"gpu"}, nil,
		),
		gpuMemUtil: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_util"),
			"Percentage of time during memory is being read or written.",
			[]string{"gpu"}, nil,
		),
		gpuMemUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_usage"),
			"Percentage of used memory size.",
			[]string{"gpu"}, nil,
		),
		gpuTemperature: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_temperature"),
			"GPU temperature in Celsius degrees.",
			[]string{"gpu"}, nil,
		),
		gpuClockMhz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_clock_mhz"),
			"GPU graphics clock in Mhz.",
			[]string{"gpu"}, nil,
		),
		gpuMemClockMhz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_mem_clock_mhz"),
			"GPU memory clock in Mhz.",
			[]string{"gpu"}, nil,
		),
		gpuThrottle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_throttle"),
			"Throttle reason.",
			[]string{"gpu"}, nil,
		),
		gpuPerfState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "gpu_perf_state"),
			"Performance state. C.uint 0: max / 15: min.",
			[]string{"gpu"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes gpu related metrics with nvml library
func (c *gpuCollector) Update(ch chan<- prometheus.Metric) error {
	stats := nvml.GetGPUStats()

	for _, v := range stats {
		gpuID := fmt.Sprintf("gpu%d", v.ID)
		ch <- prometheus.MustNewConstMetric(c.gpuUtil, prometheus.CounterValue, float64(v.UtilGPU), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemUtil, prometheus.CounterValue, float64(v.UtilMem), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemUsage, prometheus.CounterValue, float64(v.MemUsage), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuTemperature, prometheus.CounterValue, float64(v.Temperature), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuClockMhz, prometheus.CounterValue, float64(v.ClockGraphics), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuMemClockMhz, prometheus.CounterValue, float64(v.ClockMem), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuThrottle, prometheus.CounterValue, float64(v.Throttle), gpuID)
		ch <- prometheus.MustNewConstMetric(c.gpuPerfState, prometheus.CounterValue, float64(v.PerfState), gpuID)
	}

	return nil
}
