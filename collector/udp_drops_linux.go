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

//go:build !noudp_drops
// +build !noudp_drops

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type (
	udpDropsCollector struct {
		fs     procfs.FS
		desc   *prometheus.Desc
		logger *slog.Logger
	}
)

func init() {
	registerCollector("udp_drops", defaultEnabled, NewUDPdropsCollector)
}

// NewUDPdropsCollector returns a new Collector exposing network udp drops count.
func NewUDPdropsCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &udpDropsCollector{
		fs: fs,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "udp", "drops_total"),
			"Total number of datagrams dropped.",
			[]string{"ip"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *udpDropsCollector) Update(ch chan<- prometheus.Metric) error {

	s4, errIPv4 := c.fs.NetUDPSummary()
	if errIPv4 == nil {
		if s4.Drops != nil {
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(*s4.Drops), "v4")
		}
	} else {
		if errors.Is(errIPv4, os.ErrNotExist) {
			c.logger.Debug("not collecting ipv4 based metrics")
		} else {
			return fmt.Errorf("couldn't get udp drops: %w", errIPv4)
		}
	}

	s6, errIPv6 := c.fs.NetUDP6Summary()
	if errIPv6 == nil {
		if s6.Drops != nil {
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(*s6.Drops), "v6")
		}
	} else {
		if errors.Is(errIPv6, os.ErrNotExist) {
			c.logger.Debug("not collecting ipv6 based metrics")
		} else {
			return fmt.Errorf("couldn't get udp6 drops: %w", errIPv6)
		}
	}

	if errors.Is(errIPv4, os.ErrNotExist) && errors.Is(errIPv6, os.ErrNotExist) {
		return ErrNoData
	}
	return nil
}
