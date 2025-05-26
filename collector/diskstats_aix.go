// Copyright 2024 The Prometheus Authors
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

	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

const diskstatsDefaultIgnoredDevices = ""

type diskstatsCollector struct {
	rbytes typedDesc
	wbytes typedDesc
	time   typedDesc
	bsize  typedDesc
	qdepth typedDesc

	rserv typedDesc
	wserv typedDesc

	xfers typedDesc
	xrate typedDesc

	deviceFilter deviceFilter
	logger       *slog.Logger

	tickPerSecond float64
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger *slog.Logger) (Collector, error) {
	ticks, err := tickPerSecond()
	if err != nil {
		return nil, err
	}
	deviceFilter, err := newDiskstatsDeviceFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device filter flags: %w", err)
	}

	return &diskstatsCollector{
		rbytes: typedDesc{readBytesDesc, prometheus.CounterValue},
		wbytes: typedDesc{writtenBytesDesc, prometheus.CounterValue},
		time:   typedDesc{ioTimeSecondsDesc, prometheus.CounterValue},

		bsize: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "block_size_bytes"),
				"Size of the block device in bytes.",
				diskLabelNames, nil,
			),
			prometheus.GaugeValue,
		},
		qdepth: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "queue_depth"),
				"Number of requests in the queue.",
				diskLabelNames, nil,
			),
			prometheus.GaugeValue,
		},
		rserv: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"),
				"The total time spent servicing read requests.",
				diskLabelNames, nil,
			),
			prometheus.CounterValue,
		},
		wserv: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"),
				"The total time spent servicing write requests.",
				diskLabelNames, nil,
			),
			prometheus.CounterValue,
		},
		xfers: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "transfers_total"),
				"The total number of transfers to/from disk.",
				diskLabelNames, nil,
			),
			prometheus.CounterValue,
		},
		xrate: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, diskSubsystem, "transfers_to_disk_total"),
				"The total number of transfers from disk.",
				diskLabelNames, nil,
			),
			prometheus.CounterValue,
		},
		deviceFilter: deviceFilter,
		logger:       logger,

		tickPerSecond: ticks,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.DiskStat()
	if err != nil {
		return err
	}

	for _, stat := range stats {
		if c.deviceFilter.ignored(stat.Name) {
			continue
		}
		ch <- c.rbytes.mustNewConstMetric(float64(stat.Rblks*512), stat.Name)
		ch <- c.wbytes.mustNewConstMetric(float64(stat.Wblks*512), stat.Name)
		ch <- c.time.mustNewConstMetric(float64(stat.Time)/float64(c.tickPerSecond), stat.Name)

		ch <- c.bsize.mustNewConstMetric(float64(stat.BSize), stat.Name)
		ch <- c.qdepth.mustNewConstMetric(float64(stat.QDepth), stat.Name)
		ch <- c.rserv.mustNewConstMetric(float64(stat.Rserv)/1e9, stat.Name)
		ch <- c.wserv.mustNewConstMetric(float64(stat.Wserv)/1e9, stat.Name)
		ch <- c.xfers.mustNewConstMetric(float64(stat.Xfers), stat.Name)
		ch <- c.xrate.mustNewConstMetric(float64(stat.XRate), stat.Name)
	}
	return nil
}
