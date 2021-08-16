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

// +build !nolnstat

package collector

import (
	"fmt"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type lnstatCollector struct {
	logger log.Logger
}

func init() {
	registerCollector("lnstat", defaultEnabled, NewLnstatCollector)
}

func NewLnstatCollector(logger log.Logger) (Collector, error) {
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

	lnstats, err := fs.Lnstat()
	if err != nil {
		return fmt.Errorf("Lnstat error: %s", err)
	}

	for _, lnstatFile := range lnstats {
		labelNames := []string{"subsystem", "cpu"}
		for header, stats := range lnstatFile.Stats {
			for cpu, value := range stats {
				labelValues := []string{lnstatFile.Filename, strconv.Itoa(cpu)}
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, subsystem, header),
						fmt.Sprintf("linux network cache stats"),
						labelNames, nil,
					),
					prometheus.CounterValue, float64(value), labelValues...,
				)
			}
		}
	}
	return nil
}
