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

//go:build (freebsd || dragonfly || openbsd || netbsd || darwin) && !noboottime
// +build freebsd dragonfly openbsd netbsd darwin
// +build !noboottime

package collector

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

type bootTimeCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector("boottime", defaultEnabled, newBootTimeCollector)
}

// newBootTimeCollector returns a new Collector exposing system boot time on BSD systems.
func newBootTimeCollector(logger *slog.Logger) (Collector, error) {
	return &bootTimeCollector{
		logger: logger,
	}, nil
}

// Update pushes boot time onto ch
func (c *bootTimeCollector) Update(ch chan<- prometheus.Metric) error {
	tv, err := unix.SysctlTimeval("kern.boottime")
	if err != nil {
		return err
	}

	// This conversion maintains the usec precision.  Using the time
	// package did not.
	v := float64(tv.Sec) + (float64(tv.Usec) / float64(1000*1000))

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_time_seconds"),
			"Unix time of last boot, including microseconds.",
			nil, nil,
		), prometheus.GaugeValue, v)

	return nil
}
