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

//go:build !nostat
// +build !nostat

package collector

import (
	"fmt"
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type statCollector struct {
	fs           procfs.FS
	intr         *prometheus.Desc
	ctxt         *prometheus.Desc
	forks        *prometheus.Desc
	btime        *prometheus.Desc
	procsRunning *prometheus.Desc
	procsBlocked *prometheus.Desc
	softIRQ      *prometheus.Desc
	logger       *slog.Logger
}

var statSoftirqFlag = kingpin.Flag("collector.stat.softirq", "Export softirq calls per vector").Default("false").Bool()

func init() {
	registerCollector("stat", defaultEnabled, NewStatCollector)
}

// NewStatCollector returns a new Collector exposing kernel/system statistics.
func NewStatCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &statCollector{
		fs: fs,
		intr: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "intr_total"),
			"Total number of interrupts serviced.",
			nil, nil,
		),
		ctxt: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "context_switches_total"),
			"Total number of context switches.",
			nil, nil,
		),
		forks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "forks_total"),
			"Total number of forks.",
			nil, nil,
		),
		btime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_time_seconds"),
			"Node boot time, in unixtime.",
			nil, nil,
		),
		procsRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "procs_running"),
			"Number of processes in runnable state.",
			nil, nil,
		),
		procsBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "procs_blocked"),
			"Number of processes blocked waiting for I/O to complete.",
			nil, nil,
		),
		softIRQ: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "softirqs_total"),
			"Number of softirq calls.",
			[]string{"vector"}, nil,
		),
		logger: logger,
	}, nil
}

// Update implements Collector and exposes kernel and system statistics.
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.Stat()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.intr, prometheus.CounterValue, float64(stats.IRQTotal))
	ch <- prometheus.MustNewConstMetric(c.ctxt, prometheus.CounterValue, float64(stats.ContextSwitches))
	ch <- prometheus.MustNewConstMetric(c.forks, prometheus.CounterValue, float64(stats.ProcessCreated))

	ch <- prometheus.MustNewConstMetric(c.btime, prometheus.GaugeValue, float64(stats.BootTime))

	ch <- prometheus.MustNewConstMetric(c.procsRunning, prometheus.GaugeValue, float64(stats.ProcessesRunning))
	ch <- prometheus.MustNewConstMetric(c.procsBlocked, prometheus.GaugeValue, float64(stats.ProcessesBlocked))

	if *statSoftirqFlag {
		si := stats.SoftIRQ

		for _, vec := range []struct {
			name  string
			value uint64
		}{
			{name: "hi", value: si.Hi},
			{name: "timer", value: si.Timer},
			{name: "net_tx", value: si.NetTx},
			{name: "net_rx", value: si.NetRx},
			{name: "block", value: si.Block},
			{name: "block_iopoll", value: si.BlockIoPoll},
			{name: "tasklet", value: si.Tasklet},
			{name: "sched", value: si.Sched},
			{name: "hrtimer", value: si.Hrtimer},
			{name: "rcu", value: si.Rcu},
		} {
			ch <- prometheus.MustNewConstMetric(c.softIRQ, prometheus.CounterValue, float64(vec.value), vec.name)
		}
	}

	return nil
}
