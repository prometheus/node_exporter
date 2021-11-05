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

//go:build linux && !noclocksource
// +build linux,!noclocksource

package collector

import (
	"fmt"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type clockSourceCollector struct {
	fs                    sysfs.FS
	clocksourcesAvailable typedDesc
	clocksourceCurrent    typedDesc
	logger                log.Logger
}

func init() {
	registerCollector("clocksource", defaultEnabled, NewClockSourceCollector)
}

// NewClockSourceCollector returns a new Collector exposing clocksource info metrics from /sys/devices/system/clocksource.
func NewClockSourceCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &clockSourceCollector{
		fs:                    fs,
		clocksourcesAvailable: typedDesc{prometheus.NewDesc(namespace+"_clocksource_available_info", "Available clocksources read from '/sys/devices/system/clocksource'", []string{"device", "clocksource"}, nil), prometheus.GaugeValue},
		clocksourceCurrent:    typedDesc{prometheus.NewDesc(namespace+"_clocksource_current_info", "Current clocksource read from '/sys/devices/system/clocksource'", []string{"device", "clocksource"}, nil), prometheus.GaugeValue},
		logger:                logger,
	}, nil
}

func (c *clockSourceCollector) Update(ch chan<- prometheus.Metric) error {
	clocksources, err := c.fs.ClockSources()
	if err != nil {
		return fmt.Errorf("couldn't get clocksources: %w", err)
	}
	level.Debug(c.logger).Log("msg", "in Update", "clocksources", fmt.Sprintf("%v", clocksources))

	for i, clocksource := range clocksources {
		is := strconv.Itoa(i)
		for _, cs := range clocksource.Available {
			ch <- c.clocksourcesAvailable.mustNewConstMetric(1.0, is, cs)
		}
		ch <- c.clocksourceCurrent.mustNewConstMetric(1.0, is, clocksource.Current)
	}
	return err
}
