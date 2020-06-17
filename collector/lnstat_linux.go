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

	lnstats, err := procfs.Lnstat()
	if err != nil {
		return fmt.Errorf("Lnstat error: %s", err)
	}

	for _, lnstat := range lnstats {
		labelNames := []string{"subsystem", "cpu"}
		var cpu uint64 = 0
		for _, v := range lnstat.Value {
			labelValues := []string{lnstat.Filename, strconv.FormatUint(cpu, 10)}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, lnstat.Name),
					fmt.Sprintf("linux network cache stats"),
					labelNames, nil,
				),
				prometheus.CounterValue, float64(v), labelValues...,
			)
			cpu++
		}
	}
	return nil
}
