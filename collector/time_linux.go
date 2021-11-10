// Copyright 2021 The Prometheus Authors
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

//go:build linux && !notime
// +build linux,!notime

package collector

import (
	"fmt"
	"strconv"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

func (c *timeCollector) update(ch chan<- prometheus.Metric) error {
	if c.sysfs == nil {
		fs, err := sysfs.NewFS(*sysPath)
		if err != nil {
			return fmt.Errorf("failed to open procfs: %w", err)
		}
		c.sysfs = &fs
	}

	clocksources, err := c.sysfs.ClockSources()
	if err != nil {
		return fmt.Errorf("couldn't get clocksources: %w", err)
	}
	level.Debug(c.logger).Log("msg", "in Update", "clocksources", fmt.Sprintf("%v", clocksources))

	for i, clocksource := range clocksources {
		is := strconv.Itoa(i)
		for _, cs := range clocksource.Available {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					namespace+"_time_clocksource_available_info",
					"Available clocksources read from '/sys/devices/system/clocksource'",
					[]string{"device", "clocksource"},
					nil),
				prometheus.GaugeValue,
				1.0,
				is,
				cs,
			)
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				namespace+"_time_clocksource_current_info",
				"Current clocksource read from '/sys/devices/system/clocksource'",
				[]string{"device", "clocksource"},
				nil),
			prometheus.GaugeValue,
			1.0,
			is,
			clocksource.Current,
		)
	}
	return nil
}
