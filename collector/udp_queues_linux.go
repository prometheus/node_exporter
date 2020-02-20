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

// +build !noudp_queues

package collector

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type (
	udpQueuesCollector struct {
		fs     procfs.FS
		desc   *prometheus.Desc
		logger log.Logger
		ip_v4  bool
		ip_v6  bool
	}
)

var (
	ipv4 = kingpin.Flag(
		"collector.udp_queues.ipv4",
		"Read ipv4 based udp queues from /proc/net/udp (default: true). Disable the collection via '--no-collector.udp_queues.ipv4'.",
	).Default("true").Bool()
	ipv6 = kingpin.Flag(
		"collector.udp_queues.ipv6",
		"Read ipv6 based udp queues from /proc/net/udp6 (default: true).Disable the collection via '--no-collector.udp_queues.ipv6'.",
	).Default("true").Bool()
)

func init() {
	registerCollector("udp_queues", defaultDisabled, NewUDPqueuesCollector)
}

// NewUDPqueuesCollector returns a new Collector exposing network udp queued bytes.
func NewUDPqueuesCollector(logger log.Logger) (Collector, error) {
	if !*ipv4 && !*ipv6 {
		return nil,
			fmt.Errorf("Both flags '--no-collector.udp_queues.ipv4' and '--no-collector.udp_queues.ipv6' are set to 'false'. So, nothing to collect.")
	}
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}
	return &udpQueuesCollector{
		fs: fs,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "udp", "queues"),
			"Number of allocated memory in the kernel for UDP datagrams in bytes.",
			[]string{"queue", "ip"}, nil,
		),
		logger: logger,
		ip_v4:  *ipv4,
		ip_v6:  *ipv6,
	}, nil
}

func (c *udpQueuesCollector) Update(ch chan<- prometheus.Metric) error {
	if c.ip_v4 {
		s, err := c.fs.NetUDPSummary()
		if err != nil {
			return fmt.Errorf("couldn't get upd queued bytes: %s", err)
		}
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s.TxQueueLength), "tx", "v4")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s.RxQueueLength), "rx", "v4")
	}
	if c.ip_v6 {
		s6, err := c.fs.NetUDP6Summary()
		if err != nil {
			return fmt.Errorf("couldn't get upd6 queued bytes: %s", err)
		}
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s6.TxQueueLength), "tx", "v6")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s6.RxQueueLength), "rx", "v6")
	}
	return nil
}
