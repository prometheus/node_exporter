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
	"flag"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type entropyCollector struct {
	entropy_avail *prometheus.Desc
}

func init() {
	Factories["entropy"] = NewEntropyCollector
	CollectorsEnabledState["entropy"] = flag.Bool("collectors.entropy.enabled", true, "enables entropy-collector")
}

// Takes a prometheus registry and returns a new Collector exposing
// entropy stats
func NewEntropyCollector() (Collector, error) {
	return &entropyCollector{
		entropy_avail: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "entropy_available_bits"),
			"Bits of available entropy.",
			nil, nil,
		),
	}, nil
}

func (c *entropyCollector) Update(ch chan<- prometheus.Metric) (err error) {
	value, err := readUintFromFile(procFilePath("sys/kernel/random/entropy_avail"))
	if err != nil {
		return fmt.Errorf("couldn't get entropy_avail: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(
		c.entropy_avail, prometheus.GaugeValue, float64(value))

	return nil
}
