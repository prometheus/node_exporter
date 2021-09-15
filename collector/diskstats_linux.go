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
		descs: []typedFactorDesc{
			{
				desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "info"),
					"Info of /sys/block/<block_device>.",
					[]string{"device", "major", "minor"},
					nil,
				), valueType: prometheus.GaugeValue,
			},
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
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev)
			continue
		}
		blockQueue, err := c.fs.SysBlockDeviceQueueStats(dev)
		diskSectorSize := 512.0
		if err != nil {
			level.Debug(c.logger).Log("msg", "Error getting queue stats", "device", dev, "err", err)
			diskSectorSize = 512.0
		} else {
			diskSectorSize = float64(blockQueue.LogicalBlockSize)
		}

		scaleMilliseconds := 0.001

		ch <- c.descs[0].mustNewConstMetric(1.0, dev, fmt.Sprint(stats.MajorNumber), fmt.Sprint(stats.MinorNumber))
		ch <- c.descs[1].mustNewConstMetric(float64(stats.ReadIOs), dev)
		ch <- c.descs[2].mustNewConstMetric(float64(stats.ReadMerges), dev)
		ch <- c.descs[3].mustNewConstMetric(float64(stats.ReadSectors)*diskSectorSize, dev)
		ch <- c.descs[4].mustNewConstMetric(float64(stats.ReadTicks)*scaleMilliseconds, dev)
		ch <- c.descs[5].mustNewConstMetric(float64(stats.WriteIOs), dev)
		ch <- c.descs[6].mustNewConstMetric(float64(stats.WriteMerges), dev)
		ch <- c.descs[7].mustNewConstMetric(float64(stats.WriteSectors)*diskSectorSize, dev)
		ch <- c.descs[8].mustNewConstMetric(float64(stats.WriteTicks)*scaleMilliseconds, dev)
		ch <- c.descs[9].mustNewConstMetric(float64(stats.IOsInProgress), dev)
		ch <- c.descs[10].mustNewConstMetric(float64(stats.IOsTotalTicks), dev)
		ch <- c.descs[11].mustNewConstMetric(float64(stats.WeightedIOTicks), dev)
		ch <- c.descs[12].mustNewConstMetric(float64(stats.DiscardIOs), dev)
		ch <- c.descs[13].mustNewConstMetric(float64(stats.DiscardMerges), dev)
		ch <- c.descs[14].mustNewConstMetric(float64(stats.DiscardSectors)*diskSectorSize, dev)
		ch <- c.descs[15].mustNewConstMetric(float64(stats.DiscardTicks)*scaleMilliseconds, dev)
		ch <- c.descs[16].mustNewConstMetric(float64(stats.FlushRequestsCompleted), dev)
		ch <- c.descs[17].mustNewConstMetric(float64(stats.TimeSpentFlushing)*scaleMilliseconds, dev)
	}
	return nil
}
