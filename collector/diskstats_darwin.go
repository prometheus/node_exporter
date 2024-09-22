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

//go:build !nodiskstats
// +build !nodiskstats

package collector

import (
	"fmt"
	"log/slog"

	"github.com/lufia/iostat"
	"github.com/prometheus/client_golang/prometheus"
)

const diskstatsDefaultIgnoredDevices = ""

type typedDescFunc struct {
	typedDesc
	value func(stat *iostat.DriveStats) float64
}

type diskstatsCollector struct {
	descs []typedDescFunc

	deviceFilter deviceFilter
	logger       *slog.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger *slog.Logger) (Collector, error) {
	var diskLabelNames = []string{"device"}

	deviceFilter, err := newDiskstatsDeviceFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device filter flags: %w", err)
	}

	return &diskstatsCollector{
		descs: []typedDescFunc{
			{
				typedDesc: typedDesc{
					desc:      readsCompletedDesc,
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
					desc:      readTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return stat.TotalReadTime.Seconds()
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writesCompletedDesc,
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
					desc:      writeTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return stat.TotalWriteTime.Seconds()
				},
			},
			{
				typedDesc: typedDesc{
					desc:      readBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.BytesRead)
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writtenBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.BytesWritten)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "read_errors_total"),
						"The total number of read errors.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.ReadErrors)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "write_errors_total"),
						"The total number of write errors.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.WriteErrors)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "read_retries_total"),
						"The total number of read retries.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.ReadRetries)
				},
			},
			{
				typedDesc: typedDesc{
					desc: prometheus.NewDesc(
						prometheus.BuildFQName(namespace, diskSubsystem, "write_retries_total"),
						"The total number of write retries.",
						diskLabelNames,
						nil,
					),
					valueType: prometheus.CounterValue,
				},
				value: func(stat *iostat.DriveStats) float64 {
					return float64(stat.WriteRetries)
				},
			},
		},

		deviceFilter: deviceFilter,
		logger:       logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	diskStats, err := iostat.ReadDriveStats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %w", err)
	}

	for _, stats := range diskStats {
		if c.deviceFilter.ignored(stats.Name) {
			continue
		}
		for _, desc := range c.descs {
			v := desc.value(stats)
			ch <- desc.mustNewConstMetric(v, stats.Name)
		}
	}
	return nil
}
