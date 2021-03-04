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

// +build !nosoftnet

package collector

import (
	"fmt"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type softnetCollector struct {
	fs           procfs.FS
	processed    *prometheus.Desc
	dropped      *prometheus.Desc
	timeSqueezed *prometheus.Desc
	logger       log.Logger
}

const (
	softnetSubsystem = "softnet"
)

func init() {
	registerCollector("softnet", defaultEnabled, NewSoftnetCollector)
}

// NewSoftnetCollector returns a new Collector exposing softnet metrics.
func NewSoftnetCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &softnetCollector{
		fs: fs,
		processed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, softnetSubsystem, "processed_total"),
			"Number of processed packets",
			[]string{"cpu"}, nil,
		),
		dropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, softnetSubsystem, "dropped_total"),
			"Number of dropped packets",
			[]string{"cpu"}, nil,
		),
		timeSqueezed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, softnetSubsystem, "times_squeezed_total"),
			"Number of times processing packets ran out of quota",
			[]string{"cpu"}, nil,
		),
		logger: logger,
	}, nil
}

// Update gets parsed softnet statistics using procfs.
func (c *softnetCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.NetSoftnetStat()
	if err != nil {
		return fmt.Errorf("could not get softnet statistics: %w", err)
	}

	for cpuNumber, cpuStats := range stats {
		cpu := strconv.Itoa(cpuNumber)

		ch <- prometheus.MustNewConstMetric(
			c.processed,
			prometheus.CounterValue,
			float64(cpuStats.Processed),
			cpu,
		)
		ch <- prometheus.MustNewConstMetric(
			c.dropped,
			prometheus.CounterValue,
			float64(cpuStats.Dropped),
			cpu,
		)
		ch <- prometheus.MustNewConstMetric(
			c.timeSqueezed,
			prometheus.CounterValue,
			float64(cpuStats.TimeSqueezed),
			cpu,
		)
	}

	return nil
}
