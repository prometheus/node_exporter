// Copyright 2015 The Prometheus Authors
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
	"regexp"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/blockdevice"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	secondsPerTick = 1.0 / 1000.0
)

var (
	ignoredDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$").String()
)

type typedFactorDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedFactorDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

type diskstatsCollector struct {
	ignoredDevicesPattern *regexp.Regexp
	fs                    blockdevice.FS
	infoDesc              typedFactorDesc
	descs                 []typedFactorDesc
	logger                log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	var diskLabelNames = []string{"device"}
	fs, err := blockdevice.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &diskstatsCollector{
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
		fs:                    fs,
		infoDesc: typedFactorDesc{
			desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "info"),
				"Info of /sys/block/<block_device>.",
				[]string{"device", "major", "minor"},
				nil,
			), valueType: prometheus.GaugeValue,
		},
		descs: []typedFactorDesc{
			{
				desc: readsCompletedDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "reads_merged_total"),
					"The total number of reads merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: readBytesDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: readTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: writesCompletedDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "writes_merged_total"),
					"The number of writes merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: writtenBytesDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: writeTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "io_now"),
					"The number of I/Os currently in progress.",
					diskLabelNames,
					nil,
				), valueType: prometheus.GaugeValue,
			},
			{
				desc: ioTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"),
					"The weighted # of seconds spent doing I/Os.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discards_completed_total"),
					"The total number of discards completed successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discards_merged_total"),
					"The total number of discards merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discarded_sectors_total"),
					"The total number of sectors discarded successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discard_time_seconds_total"),
					"This is the total number of seconds spent by all discards.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_total"),
					"The total number of flush requests completed successfully",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_time_seconds_total"),
					"This is the total number of seconds spent by all flush requests.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
		},
		logger: logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	diskStats, err := c.fs.ProcDiskstats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %w", err)
	}

	for _, stats := range diskStats {
		dev := stats.DeviceName
		if c.ignoredDevicesPattern.MatchString(dev) {
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev, "pattern", c.ignoredDevicesPattern)
			continue
		}

		diskSectorSize := 512.0
		blockQueue, err := c.fs.SysBlockDeviceQueueStats(dev)
		if err != nil {
			level.Debug(c.logger).Log("msg", "Error getting queue stats", "device", dev, "err", err)
		} else {
			diskSectorSize = float64(blockQueue.LogicalBlockSize)
		}

		ch <- c.infoDesc.mustNewConstMetric(1.0, dev, fmt.Sprint(stats.MajorNumber), fmt.Sprint(stats.MinorNumber))

		statCount := stats.IoStatsCount - 3 // Total diskstats record count, less MajorNumber, MinorNumber and DeviceName

		for i, val := range []float64{
			float64(stats.ReadIOs),
			float64(stats.ReadMerges),
			float64(stats.ReadSectors) * diskSectorSize,
			float64(stats.ReadTicks) * secondsPerTick,
			float64(stats.WriteIOs),
			float64(stats.WriteMerges),
			float64(stats.WriteSectors) * diskSectorSize,
			float64(stats.WriteTicks) * secondsPerTick,
			float64(stats.IOsInProgress),
			float64(stats.IOsTotalTicks) * secondsPerTick,
			float64(stats.WeightedIOTicks) * secondsPerTick,
			float64(stats.DiscardIOs),
			float64(stats.DiscardMerges),
			float64(stats.DiscardSectors),
			float64(stats.DiscardTicks) * secondsPerTick,
			float64(stats.FlushRequestsCompleted),
			float64(stats.TimeSpentFlushing) * secondsPerTick,
		} {
			if i >= statCount {
				break
			}
			ch <- c.descs[i].mustNewConstMetric(val, dev)
		}
	}
	return nil
}
