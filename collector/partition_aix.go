// Copyright 2025 The Prometheus Authors
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

//go:build !nopartition
// +build !nopartition

package collector

import (
	"log/slog"

	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

type partitionCollector struct {
	logger           *slog.Logger
	entitledCapacity *prometheus.Desc
	memoryMax        *prometheus.Desc
	memoryOnline     *prometheus.Desc
	cpuOnline        *prometheus.Desc
	cpuSys           *prometheus.Desc
	cpuPool          *prometheus.Desc
	powerSaveMode    *prometheus.Desc
	smtThreads       *prometheus.Desc
}

const (
	partitionCollectorSubsystem = "partition"
)

func init() {
	registerCollector("partition", defaultEnabled, NewPartitionCollector)
}

func NewPartitionCollector(logger *slog.Logger) (Collector, error) {
	return &partitionCollector{
		logger: logger,
		entitledCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "entitled_capacity"),
			"Entitled processor capacity of the partition in CPU units (e.g. 1.0 = one core).",
			nil, nil,
		),
		memoryMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "memory_max"),
			"Maximum memory of the partition in bytes.",
			nil, nil,
		),
		memoryOnline: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "memory_online"),
			"Online memory of the partition in bytes.",
			nil, nil,
		),
		cpuOnline: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "cpus_online"),
			"Number of online CPUs in the partition.",
			nil, nil,
		),
		cpuSys: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "cpus_sys"),
			"Number of physical CPUs in the system.",
			nil, nil,
		),
		cpuPool: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "cpus_pool"),
			"Number of physical CPUs in the pool.",
			nil, nil,
		),
		powerSaveMode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "power_save_mode"),
			"Power save mode of the partition (1 for enabled, 0 for disabled).",
			nil, nil,
		),
		smtThreads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, partitionCollectorSubsystem, "smt_threads"),
			"Number of SMT threads per core.",
			nil, nil,
		),
	}, nil
}

func (c *partitionCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.PartitionStat()
	if err != nil {
		return err
	}

	powerSaveMode := 0.0
	if stats.Conf.PowerSave {
		powerSaveMode = 1.0
	}

	ch <- prometheus.MustNewConstMetric(c.entitledCapacity, prometheus.GaugeValue, float64(stats.EntCapacity)/100.0)

	ch <- prometheus.MustNewConstMetric(c.memoryMax, prometheus.GaugeValue, float64(stats.Mem.Max)*1024*1024)
	ch <- prometheus.MustNewConstMetric(c.memoryOnline, prometheus.GaugeValue, float64(stats.Mem.Online)*1024*1024)

	ch <- prometheus.MustNewConstMetric(c.cpuOnline, prometheus.GaugeValue, float64(stats.VCpus.Online))

	ch <- prometheus.MustNewConstMetric(c.cpuSys, prometheus.GaugeValue, float64(stats.NumProcessors.Online))

	ch <- prometheus.MustNewConstMetric(c.cpuPool, prometheus.GaugeValue, float64(stats.ActiveCpusInPool))

	ch <- prometheus.MustNewConstMetric(c.powerSaveMode, prometheus.GaugeValue, powerSaveMode)
	ch <- prometheus.MustNewConstMetric(c.smtThreads, prometheus.GaugeValue, float64(stats.SmtThreads))

	return nil
}
