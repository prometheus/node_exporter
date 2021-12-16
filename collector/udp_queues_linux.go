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

//go:build !noudp_queues
// +build !noudp_queues

package collector

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type (
	udpQueuesCollector struct {
		fs     procfs.FS
		desc   *prometheus.Desc
		logger log.Logger
	}
)

func init() {
	registerCollector("udp_queues", defaultEnabled, NewUDPqueuesCollector)
}

// NewUDPqueuesCollector returns a new Collector exposing network udp queued bytes.
func NewUDPqueuesCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &udpQueuesCollector{
		fs: fs,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "udp", "queues"),
			"Number of allocated memory in the kernel for UDP datagrams in bytes.",
			[]string{"queue", "ip"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *udpQueuesCollector) Update(ch chan<- prometheus.Metric) error {

	s4, errIPv4 := c.fs.NetUDPSummary()
	if errIPv4 == nil {
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s4.TxQueueLength), "tx", "v4")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s4.RxQueueLength), "rx", "v4")
	} else {
		if errors.Is(errIPv4, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "not collecting ipv4 based metrics")
		} else {
			return fmt.Errorf("couldn't get udp queued bytes: %w", errIPv4)
		}
	}

	s6, errIPv6 := c.fs.NetUDP6Summary()
	if errIPv6 == nil {
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s6.TxQueueLength), "tx", "v6")
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, float64(s6.RxQueueLength), "rx", "v6")
	} else {
		if errors.Is(errIPv6, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "not collecting ipv6 based metrics")
		} else {
			return fmt.Errorf("couldn't get udp6 queued bytes: %w", errIPv6)
		}
	}

	if errors.Is(errIPv4, os.ErrNotExist) && errors.Is(errIPv6, os.ErrNotExist) {
		return ErrNoData
	}
	return nil
}
