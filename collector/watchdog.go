// Copyright 2023 The Prometheus Authors
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

//go:build linux && !nowatchdog
// +build linux,!nowatchdog

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type watchdogCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("watchdog", defaultEnabled, NewWatchdogCollector)
}

// NewWatchdogCollector returns a new Collector exposing watchdog stats.
func NewWatchdogCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &watchdogCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

var (
	watchdogBootstatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "bootstatus"),
		"Value of /sys/class/watchdog/<watchdog>/bootstatus",
		[]string{"name"}, nil,
	)
	watchdogFwVersionDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "fw_version"),
		"Value of /sys/class/watchdog/<watchdog>/fw_version",
		[]string{"name"}, nil,
	)
	watchdogNowayoutDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "nowayout"),
		"Value of /sys/class/watchdog/<watchdog>/nowayout",
		[]string{"name"}, nil,
	)
	watchdogTimeleftDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "timeleft_seconds"),
		"Value of /sys/class/watchdog/<watchdog>/timeleft",
		[]string{"name"}, nil,
	)
	watchdogTimeoutDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "timeout_seconds"),
		"Value of /sys/class/watchdog/<watchdog>/timeout",
		[]string{"name"}, nil,
	)
	watchdogPretimeoutDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "pretimeout_seconds"),
		"Value of /sys/class/watchdog/<watchdog>/pretimeout",
		[]string{"name"}, nil,
	)
	watchdogAccessCs0Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "access_cs0"),
		"Value of /sys/class/watchdog/<watchdog>/access_cs0",
		[]string{"name"}, nil,
	)
	watchdogInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "watchdog", "info"),
		"Info of /sys/class/watchdog/<watchdog>",
		[]string{"name", "options", "identity", "state", "status", "pretimeout_governor"}, nil,
	)
)

func toLabelValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (c *watchdogCollector) Update(ch chan<- prometheus.Metric) error {
	watchdogClass, err := c.fs.WatchdogClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrInvalid) {
			c.logger.Debug("Could not read watchdog stats", "err", err)
			return ErrNoData
		}
		return err
	}

	for _, wd := range watchdogClass {
		if wd.Bootstatus != nil {
			ch <- prometheus.MustNewConstMetric(watchdogBootstatusDesc, prometheus.GaugeValue, float64(*wd.Bootstatus), wd.Name)
		}
		if wd.FwVersion != nil {
			ch <- prometheus.MustNewConstMetric(watchdogFwVersionDesc, prometheus.GaugeValue, float64(*wd.FwVersion), wd.Name)
		}
		if wd.Nowayout != nil {
			ch <- prometheus.MustNewConstMetric(watchdogNowayoutDesc, prometheus.GaugeValue, float64(*wd.Nowayout), wd.Name)
		}
		if wd.Timeleft != nil {
			ch <- prometheus.MustNewConstMetric(watchdogTimeleftDesc, prometheus.GaugeValue, float64(*wd.Timeleft), wd.Name)
		}
		if wd.Timeout != nil {
			ch <- prometheus.MustNewConstMetric(watchdogTimeoutDesc, prometheus.GaugeValue, float64(*wd.Timeout), wd.Name)
		}
		if wd.Pretimeout != nil {
			ch <- prometheus.MustNewConstMetric(watchdogPretimeoutDesc, prometheus.GaugeValue, float64(*wd.Pretimeout), wd.Name)
		}
		if wd.AccessCs0 != nil {
			ch <- prometheus.MustNewConstMetric(watchdogAccessCs0Desc, prometheus.GaugeValue, float64(*wd.AccessCs0), wd.Name)
		}

		ch <- prometheus.MustNewConstMetric(watchdogInfoDesc, prometheus.GaugeValue, 1.0,
			wd.Name, toLabelValue(wd.Options), toLabelValue(wd.Identity), toLabelValue(wd.State), toLabelValue(wd.Status), toLabelValue(wd.PretimeoutGovernor))
	}

	return nil
}
