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
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type timeCollector struct {
	now                   typedDesc
	zone                  typedDesc
	clocksourcesAvailable typedDesc
	clocksourceCurrent    typedDesc
	logger                *slog.Logger
}

func init() {
	registerCollector("time", defaultEnabled, NewTimeCollector)
}

// NewTimeCollector returns a new Collector exposing the current system time in
// seconds since epoch.
func NewTimeCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "time"
	return &timeCollector{
		now: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "seconds"),
			"System time in seconds since epoch (1970).",
			nil, nil,
		), prometheus.GaugeValue},
		zone: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "zone_offset_seconds"),
			"System time zone offset in seconds.",
			[]string{"time_zone"}, nil,
		), prometheus.GaugeValue},
		clocksourcesAvailable: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "clocksource_available_info"),
			"Available clocksources read from '/sys/devices/system/clocksource'.",
			[]string{"device", "clocksource"}, nil,
		), prometheus.GaugeValue},
		clocksourceCurrent: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "clocksource_current_info"),
			"Current clocksource read from '/sys/devices/system/clocksource'.",
			[]string{"device", "clocksource"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *timeCollector) Update(ch chan<- prometheus.Metric) error {
	now := time.Now()
	nowSec := float64(now.UnixNano()) / 1e9
	zone, zoneOffset := now.Zone()

	c.logger.Debug("Return time", "now", nowSec)
	ch <- c.now.mustNewConstMetric(nowSec)
	c.logger.Debug("Zone offset", "offset", zoneOffset, "time_zone", zone)
	ch <- c.zone.mustNewConstMetric(float64(zoneOffset), zone)
	return c.update(ch)
}
