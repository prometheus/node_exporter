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
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	sockStatSubsystem = "sockstat"
)

// Used for calculating the total memory bytes on TCP and UDP.
var pageSize = os.Getpagesize()

type sockStatCollector struct {
	logger log.Logger
}

// A name and optional value for a sockstat metric.
type ssPair struct {
	name   string
	desc   string
	v      *int
	shared bool // Whether /proc/net/sockstat value is sum of IPv4+IPv6 values
}

func init() {
	registerCollector(sockStatSubsystem, defaultEnabled, NewSockStatCollector)
}

// NewSockStatCollector returns a new Collector exposing socket stats.
func NewSockStatCollector(logger log.Logger) (Collector, error) {
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
		level.Debug(c.logger).Log("msg", "IPv4 sockstat statistics not found, skipping")
	default:
		return fmt.Errorf("failed to get IPv4 sockstat data: %w", err)
	}

	stat6, err := fs.NetSockstat6()
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		level.Debug(c.logger).Log("msg", "IPv6 sockstat statistics not found, skipping")
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
		err := c.validate(s.isIPv6, s.stat)
		if err != nil {
			return err
		}
	}

	for _, s := range stats {
		c.update(ch, s.isIPv6, s.stat)
	}

	return nil
}

func getSsPairs(p procfs.NetSockstatProtocol) []ssPair {
	pairs := []ssPair{
		{
			name:   "inuse",
			desc:   "Number of %s in use.",
			shared: false,
			v:      &p.InUse,
		},
		{
			name:   "orphan",
			desc:   "Number of orphaned %s.",
			shared: true,
			v:      p.Orphan,
		},
		{
			name:   "tw",
			desc:   "Number of %s in state TIME_WAIT.",
			shared: true,
			v:      p.TW,
		},
		{
			name:   "alloc",
			desc:   "Number of allocated %s.",
			shared: true,
			v:      p.Alloc,
		},
		{
			name:   "mem",
			desc:   "Number of pages allocated for %s.",
			shared: true,
			v:      p.Mem,
		},
		{
			name:   "memory",
			desc:   "Number of bytes allocated for %s.",
			shared: false,
			v:      p.Memory,
		},
	}

	// Also export mem_bytes values for sockets which have a mem value
	// stored in pages.
	if p.Mem != nil {
		v := *p.Mem * pageSize
		pairs = append(pairs, ssPair{
			name:   "mem_bytes",
			desc:   "Number of bytes allocated for %s.",
			shared: true,
			v:      &v,
		})
	}

	return pairs
}

func (c *sockStatCollector) validate(isIPv6 bool, s *procfs.NetSockstat) error {
	for _, p := range s.Protocols {
		if isIPv6 {
			for _, pair := range getSsPairs(p) {
				if pair.shared && pair.v != nil {
					return fmt.Errorf("Unexpected %s %s in sockstat6", p.Protocol, pair.name)
				}
			}
		} else {
			if strings.HasSuffix(p.Protocol, "6") {
				return fmt.Errorf("Unexpected %s in IPv4 sockstat", p.Protocol)
			}
		}
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

	// Previously these metric names were generated directly from the file output.
	// In order to keep the same level of compatibility, we must map the fields
	// to their correct names.
	for _, p := range s.Protocols {
		for _, pair := range getSsPairs(p) {
			if pair.v == nil {
				// This value is not set for this protocol; nothing to do.
				continue
			}

			var kind string
			if p.Protocol == "FRAG" || p.Protocol == "FRAG6" {
				kind = "fragments"
			} else {
				kind = strings.TrimSuffix(p.Protocol, "6") + " sockets"
			}

			if isIPv6 {
				kind = "IPv6 " + kind
			} else if !pair.shared {
				kind = "IPv4 " + kind
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(
						namespace,
						sockStatSubsystem,
						fmt.Sprintf("%s_%s", p.Protocol, pair.name),
					),
					fmt.Sprintf(pair.desc, kind),
					nil,
					nil,
				),
				prometheus.GaugeValue,
				float64(*pair.v),
			)
		}
	}
}
