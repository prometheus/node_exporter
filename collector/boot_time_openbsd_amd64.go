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

//go:build !noboottime
// +build !noboottime

package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
	"unsafe"
)

type bootTimeCollector struct {
	name, description string
	logger            log.Logger
}

func init() {
	registerCollector("boottime", defaultEnabled, newBootTimeCollector)
}

// newBootTimeCollector returns a new Collector exposing system boot time on BSD systems.
func newBootTimeCollector(logger log.Logger) (Collector, error) {
	return &bootTimeCollector{
		name:        "boot_time_seconds",
		description: "Unix time of last boot, including microseconds.",
		logger:      logger,
	}, nil
}

// Update pushes boot time onto ch
func (c *bootTimeCollector) Update(ch chan<- prometheus.Metric) error {
	raw, err := unix.SysctlRaw("kern.boottime")
	if err != nil {
		return err
	}

	tv := *(*unix.Timeval)(unsafe.Pointer(&raw[0]))
	v := (float64(tv.Sec) + (float64(tv.Usec) / float64(1000*1000)))

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", c.name),
			c.description,
			nil, nil,
		), prometheus.GaugeValue, v)

	return nil
}
