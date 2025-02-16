// Copyright 2022 The Prometheus Authors
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

//go:build !noselinux
// +build !noselinux

package collector

import (
	"log/slog"

	"github.com/opencontainers/selinux/go-selinux"
	"github.com/prometheus/client_golang/prometheus"
)

type selinuxCollector struct {
	configMode  *prometheus.Desc
	currentMode *prometheus.Desc
	enabled     *prometheus.Desc
	logger      *slog.Logger
}

func init() {
	registerCollector("selinux", defaultEnabled, NewSelinuxCollector)
}

// NewSelinuxCollector returns a new Collector exposing SELinux statistics.
func NewSelinuxCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "selinux"

	return &selinuxCollector{
		configMode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "config_mode"),
			"Configured SELinux enforcement mode",
			nil, nil,
		),
		currentMode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "current_mode"),
			"Current SELinux enforcement mode",
			nil, nil,
		),
		enabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "enabled"),
			"SELinux is enabled, 1 is true, 0 is false",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *selinuxCollector) Update(ch chan<- prometheus.Metric) error {
	if !selinux.GetEnabled() {
		ch <- prometheus.MustNewConstMetric(
			c.enabled, prometheus.GaugeValue, 0)

		return nil
	}

	ch <- prometheus.MustNewConstMetric(
		c.enabled, prometheus.GaugeValue, 1)

	ch <- prometheus.MustNewConstMetric(
		c.configMode, prometheus.GaugeValue, float64(selinux.DefaultEnforceMode()))

	ch <- prometheus.MustNewConstMetric(
		c.currentMode, prometheus.GaugeValue, float64(selinux.EnforceMode()))

	return nil
}
