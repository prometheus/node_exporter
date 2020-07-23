// Copyright 2018 The Prometheus Authors
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

// +build freebsd dragonfly openbsd netbsd darwin
// +build !noboottime

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type bootTimeCollector struct{ boottime bsdSysctl }

func init() {
	registerCollector("boottime", defaultEnabled, newBootTimeCollector)
}

// newBootTimeCollector returns a new Collector exposing system boot time on BSD systems.
func newBootTimeCollector() (Collector, error) {
	return &bootTimeCollector{
		boottime: bsdSysctl{
			name:        "boot_time_seconds",
			description: "Unix time of last boot, including microseconds.",
			mib:         "kern.boottime",
			dataType:    bsdSysctlTypeStructTimeval,
		},
	}, nil
}

// Update pushes boot time onto ch
func (c *bootTimeCollector) Update(ch chan<- prometheus.Metric) error {
	v, err := c.boottime.Value()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", c.boottime.name),
			c.boottime.description,
			nil, nil,
		), prometheus.GaugeValue, v)

	return nil
}
