// Copyright 2019 The Prometheus Authors
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

// +build !nopressure

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

var (
	psiResources = []string{"cpu", "io", "memory"}
)

type pressureStatsCollector struct {
	cpu     *prometheus.Desc
	io      *prometheus.Desc
	ioFull  *prometheus.Desc
	mem     *prometheus.Desc
	memFull *prometheus.Desc

	fs procfs.FS
}

func init() {
	registerCollector("pressure", defaultEnabled, NewPressureStatsCollector)
}

// NewPressureStatsCollector returns a Collector exposing pressure stall information
func NewPressureStatsCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}

	return &pressureStatsCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "cpu_waiting_seconds_total"),
			"Total time in seconds that processes have waited for CPU time",
			nil, nil,
		),
		io: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "io_waiting_seconds_total"),
			"Total time in seconds that processes have waited due to IO congestion",
			nil, nil,
		),
		ioFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "io_stalled_seconds_total"),
			"Total time in seconds no process could make progress due to IO congestion",
			nil, nil,
		),
		mem: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "memory_waiting_seconds_total"),
			"Total time in seconds that processes have waited for memory",
			nil, nil,
		),
		memFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "memory_stalled_seconds_total"),
			"Total time in seconds no process could make progress due to memory congestion",
			nil, nil,
		),
		fs: fs,
	}, nil
}

// Update calls procfs.NewPSIStatsForResource for the different resources and updates the values
func (c *pressureStatsCollector) Update(ch chan<- prometheus.Metric) error {
	for _, res := range psiResources {
		log.Debugf("collecting statistics for resource: %s", res)
		vals, err := c.fs.PSIStatsForResource(res)
		if err != nil {
			log.Debug("pressure information is unavailable, you need a Linux kernel >= 4.20 and/or CONFIG_PSI enabled for your kernel")
			return nil
		}
		switch res {
		case "cpu":
			ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
		case "io":
			ch <- prometheus.MustNewConstMetric(c.io, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
			ch <- prometheus.MustNewConstMetric(c.ioFull, prometheus.CounterValue, float64(vals.Full.Total)/1000.0/1000.0)
		case "memory":
			ch <- prometheus.MustNewConstMetric(c.mem, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
			ch <- prometheus.MustNewConstMetric(c.memFull, prometheus.CounterValue, float64(vals.Full.Total)/1000.0/1000.0)
		default:
			log.Debugf("did not account for resource: %s", res)
		}
	}

	return nil
}
