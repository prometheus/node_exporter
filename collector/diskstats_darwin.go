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

	"github.com/lufia/iostat"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	diskSubsystem = "disk"
)

type typedDescFunc struct {
	typedDesc
	value func(stat *iostat.DriveStats) float64
}

type diskstatsCollector struct {
	descs []typedDescFunc
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector() (Collector, error) {
	var diskLabelNames = []string{"device"}

	return &diskstatsCollector{
		descs: []typedDescFunc{
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"),
						"The total number of reads completed successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.NumRead)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "read_sectors_total"),
						"The total number of sectors read successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.NumRead) / float64(stat.BlockSize)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"),
						"The total number of seconds spent by all reads.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return stat.TotalReadTime.Seconds()
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"),
						"The total number of writes completed successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.NumWrite)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "written_sectors_total"),
						"The total number of sectors written successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.NumWrite) / float64(stat.BlockSize)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"),
						"This is the total number of seconds spent by all writes.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return stat.TotalWriteTime.Seconds()
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"),
						"The total number of bytes read successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.BytesRead)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"),
						"The total number of bytes written successfully.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.BytesWritten)
				},
			},
		},
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	diskStats, err := iostat.ReadDriveStats()
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
