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

//go:build !nosockstat
// +build !nosockstat

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	sockStatSubsystem = "sockstat"
)

// Used for calculating the total memory bytes on TCP and UDP.
var pageSize = os.Getpagesize()

type sockStatCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector(sockStatSubsystem, defaultEnabled, NewSockStatCollector)
}

// NewSockStatCollector returns a new Collector exposing socket stats.
func NewSockStatCollector(logger *slog.Logger) (Collector, error) {
	return &sockStatCollector{logger}, nil
}

func (c *sockStatCollector) Update(ch chan<- prometheus.Metric) error {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	// If IPv4 and/or IPv6 are disabled on this kernel, handle it gracefully.
	stat4, err := fs.NetSockstat()
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		c.logger.Debug("IPv4 sockstat statistics not found, skipping")
	default:
		return fmt.Errorf("failed to get IPv4 sockstat data: %w", err)
	}

	stat6, err := fs.NetSockstat6()
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		c.logger.Debug("IPv6 sockstat statistics not found, skipping")
	default:
		return fmt.Errorf("failed to get IPv6 sockstat data: %w", err)
	}

	stats := []struct {
		isIPv6 bool
		stat   *procfs.NetSockstat
	}{
		{
			stat: stat4,
		},
		{
			isIPv6: true,
			stat:   stat6,
		},
	}

	for _, s := range stats {
		c.update(ch, s.isIPv6, s.stat)
	}

	return nil
}

func (c *sockStatCollector) update(ch chan<- prometheus.Metric, isIPv6 bool, s *procfs.NetSockstat) {
	if s == nil {
		// IPv6 disabled or similar; nothing to do.
		return
	}

	// If sockstat contains the number of used sockets, export it.
	if !isIPv6 && s.Used != nil {
		// TODO: this must be updated if sockstat6 ever exports this data.
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, sockStatSubsystem, "sockets_used"),
				"Number of IPv4 sockets in use.",
				nil,
				nil,
			),
			prometheus.GaugeValue,
			float64(*s.Used),
		)
	}

	// A name and optional value for a sockstat metric.
	type ssPair struct {
		name string
		v    *int
	}

	// Previously these metric names were generated directly from the file output.
	// In order to keep the same level of compatibility, we must map the fields
	// to their correct names.
	for _, p := range s.Protocols {
		pairs := []ssPair{
			{
				name: "inuse",
				v:    &p.InUse,
			},
			{
				name: "orphan",
				v:    p.Orphan,
			},
			{
				name: "tw",
				v:    p.TW,
			},
			{
				name: "alloc",
				v:    p.Alloc,
			},
			{
				name: "mem",
				v:    p.Mem,
			},
			{
				name: "memory",
				v:    p.Memory,
			},
		}

		// Also export mem_bytes values for sockets which have a mem value
		// stored in pages.
		if p.Mem != nil {
			v := *p.Mem * pageSize
			pairs = append(pairs, ssPair{
				name: "mem_bytes",
				v:    &v,
			})
		}

		for _, pair := range pairs {
			if pair.v == nil {
				// This value is not set for this protocol; nothing to do.
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(
						namespace,
						sockStatSubsystem,
						fmt.Sprintf("%s_%s", p.Protocol, pair.name),
					),
					fmt.Sprintf("Number of %s sockets in state %s.", p.Protocol, pair.name),
					nil,
					nil,
				),
				prometheus.GaugeValue,
				float64(*pair.v),
			)
		}
	}
}
