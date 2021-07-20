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
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/prometheus/procfs/blockdevice"
)

const (
	diskSectorSize = 512
)

var (
	ignoredDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|(mmcblk|nvme\\d+n)\\d+p)\\d+$").String()
	preferSysFS    = kingpin.Flag("collector.diskstats.prefer-sysfs", "Using /sys automatically skips partition metrics.").Default("false").Bool()
)

type typedFactorDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	factor    float64
}

func (d *typedFactorDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	if d.factor != 0 {
		value *= d.factor
	}
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

type diskstatsCollector struct {
	fs                    blockdevice.FS
	ignoredDevicesPattern *regexp.Regexp
	preferSysFS           bool
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
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &diskstatsCollector{
		fs:                    fs,
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
		preferSysFS:           *preferSysFS,
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
				factor: diskSectorSize,
			},
			{
				desc: readTimeSecondsDesc, valueType: prometheus.CounterValue,
				factor: .001,
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
				factor: diskSectorSize,
			},
			{
				desc: writeTimeSecondsDesc, valueType: prometheus.CounterValue,
				factor: .001,
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
				factor: .001,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"),
					"The weighted # of seconds spent doing I/Os.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
				factor: .001,
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
				factor: .001,
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
				factor: .001,
			},
		},
		logger: logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	var (
		stats    = map[string]blockdevice.IOStats{}
		counts   = map[string]int{}
		useSysFS = c.preferSysFS
	)

RETRY:
	for {
		if useSysFS {
			devices, err := c.fs.SysBlockDevices()
			if err != nil {
				level.Warn(c.logger).Log("msg", "couldn't list devices from /sys. Retry from /proc", "err", err)
				useSysFS = false
				continue RETRY
			}

			for _, dev := range devices {
				if c.ignoredDevicesPattern.MatchString(dev) {
					level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev)
					continue
				}

				stat, count, err := c.fs.SysBlockDeviceStat(dev)
				if err != nil {
					level.Warn(c.logger).Log("msg", "couldn't get diskstats. Retry from /proc", "device", dev, "err", err)
					useSysFS = false
					continue RETRY
				}

				stats[dev] = stat
				counts[dev] = count
			}
		} else {
			diskstats, err := c.fs.ProcDiskstats()
			if err != nil {
				return fmt.Errorf("couldn't get diskstats: %w", err)
			}

			for _, diskstat := range diskstats {
				dev := diskstat.Info.DeviceName

				if c.ignoredDevicesPattern.MatchString(dev) {
					level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev)
					continue
				}

				stats[dev] = diskstat.IOStats
				// Do not count major, minor and device name
				counts[dev] = diskstat.IoStatsCount - 3
			}
		}
		break
	}

	for dev, stat := range stats {

		count := counts[dev]

		statValues := []uint64{
			stat.ReadIOs,
			stat.ReadMerges,
			stat.ReadSectors,
			stat.ReadTicks,
			stat.WriteIOs,
			stat.WriteMerges,
			stat.WriteSectors,
			stat.WriteTicks,
			stat.IOsInProgress,
			stat.IOsTotalTicks,
			stat.WeightedIOTicks,
			stat.DiscardIOs,
			stat.DiscardMerges,
			stat.DiscardSectors,
			stat.DiscardTicks,
			stat.FlushRequestsCompleted,
			stat.TimeSpentFlushing,
		}

		for i := 0; i < count; i++ {
			// ignore unrecognized additional stats
			if i >= len(c.descs) {
				break
			}
			ch <- c.descs[i].mustNewConstMetric(float64(statValues[i]), dev)
		}
	}
	return nil
}
