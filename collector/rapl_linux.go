// Copyright 2019 The Prometheus Authors
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

//go:build !norapl
// +build !norapl

package collector

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const raplCollectorSubsystem = "rapl"

type raplCollector struct {
	fs     sysfs.FS
	logger log.Logger

	joulesMetricDesc *prometheus.Desc
}

func init() {
	registerCollector(raplCollectorSubsystem, defaultEnabled, NewRaplCollector)
}

var (
	raplZoneLabel = kingpin.Flag("collector.rapl.enable-zone-label", "Enables service unit metric unit_start_time_seconds").Bool()
)

// NewRaplCollector returns a new Collector exposing RAPL metrics.
func NewRaplCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)

	if err != nil {
		return nil, err
	}

	joulesMetricDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, raplCollectorSubsystem, "joules_total"),
		"Current RAPL value in joules",
		[]string{"index", "path", "rapl_zone"}, nil,
	)

	collector := raplCollector{
		fs:               fs,
		logger:           logger,
		joulesMetricDesc: joulesMetricDesc,
	}
	return &collector, nil
}

// Update implements Collector and exposes RAPL related metrics.
func (c *raplCollector) Update(ch chan<- prometheus.Metric) error {
	// nil zones are fine when platform doesn't have powercap files present.
	zones, err := sysfs.GetRaplZones(c.fs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "Platform doesn't have powercap files present", "err", err)
			return ErrNoData
		}
		if errors.Is(err, os.ErrPermission) {
			level.Debug(c.logger).Log("msg", "Can't access powercap files", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve rapl stats: %w", err)
	}

	for _, rz := range zones {
		microJoules, err := rz.GetEnergyMicrojoules()
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				level.Debug(c.logger).Log("msg", "Can't access energy_uj file", "zone", rz, "err", err)
				return ErrNoData
			}
			return err
		}

		joules := float64(microJoules) / 1000000.0

		if *raplZoneLabel {
			ch <- c.joulesMetricWithZoneLabel(rz, joules)
		} else {
			ch <- c.joulesMetric(rz, joules)
		}
	}
	return nil
}

func (c *raplCollector) joulesMetric(z sysfs.RaplZone, v float64) prometheus.Metric {
	index := strconv.Itoa(z.Index)
	descriptor := prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			raplCollectorSubsystem,
			fmt.Sprintf("%s_joules_total", SanitizeMetricName(z.Name)),
		),
		fmt.Sprintf("Current RAPL %s value in joules", z.Name),
		[]string{"index", "path"}, nil,
	)

	return prometheus.MustNewConstMetric(
		descriptor,
		prometheus.CounterValue,
		v,
		index,
		z.Path,
	)
}

func (c *raplCollector) joulesMetricWithZoneLabel(z sysfs.RaplZone, v float64) prometheus.Metric {
	index := strconv.Itoa(z.Index)

	return prometheus.MustNewConstMetric(
		c.joulesMetricDesc,
		prometheus.CounterValue,
		v,
		index,
		z.Path,
		z.Name,
	)
}
