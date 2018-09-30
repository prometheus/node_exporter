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

// +build !noconntrack

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ti-mo/conntrack"
	"gopkg.in/alecthomas/kingpin.v2"
)

type conntrackCollector struct {
	current         *prometheus.Desc
	limit           *prometheus.Desc
	kernelStatistic *prometheus.Desc
}

var (
	enableConntrackKernelStats = kingpin.Flag("collector.conntrack.kernel-stats", "fetch conntrack stats from kernel (requires root or CAP_NET_ADMIN)").Bool()
)

func init() {
	registerCollector("conntrack", defaultEnabled, NewConntrackCollector)
}

// NewConntrackCollector returns a new Collector exposing conntrack stats.
func NewConntrackCollector() (Collector, error) {
	c := &conntrackCollector{
		current: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_entries"),
			"Number of currently allocated flow entries for connection tracking.",
			nil, nil,
		),
		limit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_entries_limit"),
			"Maximum size of connection tracking table.",
			nil, nil,
		),
	}

	if *enableConntrackKernelStats {
		c.kernelStatistic = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_kernel_statistic"),
			"Conntrack Kernel counter",
			[]string{"cpu", "statistic"}, nil,
		)
	}

	return c, nil
}

func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	value, err := readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_count"))
	if err != nil {
		// Conntrack probably not loaded into the kernel.
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.current, prometheus.GaugeValue, float64(value))

	value, err = readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_max"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.limit, prometheus.GaugeValue, float64(value))

	if *enableConntrackKernelStats {
		err = c.updateConntrackKernelStats(ch)
	}

	return err
}

func (c *conntrackCollector) updateConntrackKernelStats(ch chan<- prometheus.Metric) error {
	conn, err := conntrack.Dial(nil)
	stats, err := conn.Stats()
	if err != nil {
		return err
	}
	for cpuIdx, s := range stats {
		cpuLabel := fmt.Sprintf("%d", cpuIdx)

		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Found),
			cpuLabel, "found",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Invalid),
			cpuLabel, "invalid",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Ignore),
			cpuLabel, "ignore",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Insert),
			cpuLabel, "insert",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.InsertFailed),
			cpuLabel, "insert_failed",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Drop),
			cpuLabel, "drop",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.EarlyDrop),
			cpuLabel, "early_drop",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.Error),
			cpuLabel, "error",
		)
		ch <- prometheus.MustNewConstMetric(
			c.kernelStatistic,
			prometheus.CounterValue,
			float64(s.SearchRestart),
			cpuLabel, "search_restart",
		)
	}

	return nil
}
