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

// +build !noentropy

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type entropyCollector struct {
	entropyAvail *prometheus.Desc
}

func init() {
	registerCollector("entropy", defaultEnabled, NewEntropyCollector)
}

// NewEntropyCollector returns a new Collector exposing entropy stats.
func NewEntropyCollector() (Collector, error) {
	return &entropyCollector{
		entropyAvail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "entropy_available_bits"),
			"Bits of available entropy.",
			nil, nil,
		),
	}, nil
}

func (c *entropyCollector) Update(ch chan<- prometheus.Metric) error {
	value, err := readUintFromFile(procFilePath("sys/kernel/random/entropy_avail"))
	if err != nil {
		return fmt.Errorf("couldn't get entropy_avail: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(
		c.entropyAvail, prometheus.GaugeValue, float64(value))

	return nil
}
