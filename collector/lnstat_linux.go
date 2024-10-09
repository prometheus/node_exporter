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

//go:build !nolnstat
// +build !nolnstat

package collector

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type lnstatCollector struct {
	logger *slog.Logger
}

func init() {
	registerCollector("lnstat", defaultDisabled, NewLnstatCollector)
}

func NewLnstatCollector(logger *slog.Logger) (Collector, error) {
	return &lnstatCollector{logger}, nil
}

func (c *lnstatCollector) Update(ch chan<- prometheus.Metric) error {
	const (
		subsystem = "lnstat"
	)

	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	netStats, err := fs.NetStat()
	if err != nil {
		return fmt.Errorf("lnstat error: %s", err)
	}

	for _, netStatFile := range netStats {
		labelNames := []string{"subsystem", "cpu"}
		for header, stats := range netStatFile.Stats {
			for cpu, value := range stats {
				labelValues := []string{netStatFile.Filename, strconv.Itoa(cpu)}
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, subsystem, header+"_total"),
						"linux network cache stats",
						labelNames, nil,
					),
					prometheus.CounterValue, float64(value), labelValues...,
				)
			}
		}
	}
	return nil
}
