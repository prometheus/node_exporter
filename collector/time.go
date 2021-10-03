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

//go:build !notime
// +build !notime

package collector

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type timeCollector struct {
	nowDesc  *prometheus.Desc
	zoneDesc *prometheus.Desc
	logger   log.Logger
}

func init() {
	registerCollector("time", defaultEnabled, NewTimeCollector)
}

// NewTimeCollector returns a new Collector exposing the current system time in
// seconds since epoch.
func NewTimeCollector(logger log.Logger) (Collector, error) {
	const subsystem = "time"
	return &timeCollector{
		nowDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "seconds"),
			"System time in seconds since epoch (1970).",
			nil, nil,
		),
		zoneDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "zone_offset_seconds"),
			"System time zone offset in seconds.",
			[]string{"time_zone"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *timeCollector) Update(ch chan<- prometheus.Metric) error {
	now := time.Now()
	nowSec := float64(now.UnixNano()) / 1e9
	zone, zoneOffset := now.Zone()

	level.Debug(c.logger).Log("msg", "Return time", "now", nowSec)
	ch <- prometheus.MustNewConstMetric(c.nowDesc, prometheus.GaugeValue, nowSec)
	level.Debug(c.logger).Log("msg", "Zone offset", "offset", zoneOffset, "time_zone", zone)
	ch <- prometheus.MustNewConstMetric(c.zoneDesc, prometheus.GaugeValue, float64(zoneOffset), zone)
	return nil
}
