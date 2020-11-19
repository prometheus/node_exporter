// Copyright 2020 The Prometheus Authors
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
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	diskSectorSize = 512
)

var (
	ignoredDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$").String()
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
	ignoredDevicesPattern *regexp.Regexp
	descs                 []typedFactorDesc
	descUsages            []typedFactorDesc
	logger                log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	var diskLabelNames = []string{"device"}

	return &diskstatsCollector{
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
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
		},
		descUsages: []typedFactorDesc{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "total_bytes"),
					"The total of disk.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "free_bytes"),
					"The free of disk",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "used_bytes"),
					"The used in disk.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "usedPercent"),
					"The percentage used in disk.",
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
	diskIOCounters, err := getDiskIOCounters()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %s", err)
	}

	for dev, stats := range diskIOCounters {
		if c.ignoredDevicesPattern.MatchString(dev) {
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev)
			continue
		}

		for i, value := range stats {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in diskstats: %s", value, err)
			}
			ch <- c.descs[i].mustNewConstMetric(v, dev)
		}
	}

	diskUsages, err := getDiskUsages()
	if err != nil {
		return fmt.Errorf("couldn't get disk usages: %s", err)
	}
	for dev, stats := range diskUsages {
		if c.ignoredDevicesPattern.MatchString(dev) {
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", dev)
			continue
		}

		for i, value := range stats {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in diskstats: %s", value, err)
			}
			ch <- c.descUsages[i].mustNewConstMetric(v, dev)
		}
	}
	return nil
}

func getDiskIOCounters() (map[string][]string, error) {
	diskIOCounters := map[string][]string{}
	diskIOCs, err := disk.IOCounters()
	if err != nil {
		return diskIOCounters, err
	}
	for dev, diskIOC := range diskIOCs {
		diskIOCounters[dev] = []string{
			fmt.Sprintf(`%v`, diskIOC.ReadCount),
			fmt.Sprintf(`%v`, diskIOC.MergedReadCount),
			fmt.Sprintf(`%v`, diskIOC.ReadBytes),
			fmt.Sprintf(`%v`, diskIOC.ReadTime),
			fmt.Sprintf(`%v`, diskIOC.WriteCount),
			fmt.Sprintf(`%v`, diskIOC.MergedWriteCount),
			fmt.Sprintf(`%v`, diskIOC.WriteBytes),
			fmt.Sprintf(`%v`, diskIOC.WriteTime),
			fmt.Sprintf(`%v`, diskIOC.IopsInProgress),
			fmt.Sprintf(`%v`, diskIOC.IoTime),
			fmt.Sprintf(`%v`, diskIOC.WeightedIO),
		}
	}
	return diskIOCounters, nil
}

func getDiskUsages() (map[string][]string, error) {
	diskUsages := map[string][]string{}
	partitions, err := disk.Partitions(true)
	if err != nil {
		return diskUsages, err
	}
	for _, partition := range partitions {
		devUsages, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return diskUsages, err
		}
		diskUsages[partition.Mountpoint] = []string{
			fmt.Sprintf(`%v`, devUsages.Total),
			fmt.Sprintf(`%v`, devUsages.Free),
			fmt.Sprintf(`%v`, devUsages.Used),
			fmt.Sprintf("%.6f", devUsages.UsedPercent),
		}
	}
	return diskUsages, nil
}
