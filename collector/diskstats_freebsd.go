// Copyright 2017 The Prometheus Authors
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

// +build !nodiskstats

package collector

import (
	"fmt"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/shirou/gopsutil/disk"  
	"github.com/prometheus/client_golang/prometheus"
)

type typedDescFunc struct {
	typedDesc
	value func(stat disk.IOCountersStat) float64
}

type diskstatsCollector struct {
	descs  []typedDescFunc
	logger log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	var diskLabelNames = []string{"device"}

	return &diskstatsCollector{
		descs: []typedDescFunc{
			{
				typedDesc: typedDesc{
					desc:      readsCompletedDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.ReadCount)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      readTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.ReadTime)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writesCompletedDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.WriteCount)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writeTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.WriteTime)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      readBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.ReadBytes)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writtenBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.WriteBytes)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writtenBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.WriteBytes)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "io_now"),
						"The number of I/Os currently in progress.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.GaugeValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.IopsInProgress)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"),
						"The weighted # of seconds spent doing I/Os.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.GaugeValue,
				},
				value: func(stat disk.IOCountersStat) float64 {
					return float64(stat.WeightedIO)
				},
			},
		},
		logger: logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	diskStats, err := disk.IOCountersWithContext(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %s", err)
	}

	for _, stats := range diskStats { 
		for _, desc := range c.descs {
			v := desc.value(stats)
			ch <- desc.mustNewConstMetric(v, stats.Name)
		}
	}
	return nil
}
