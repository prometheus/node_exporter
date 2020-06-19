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

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type entropyCollector struct {
	fs              procfs.FS
	entropyAvail    *prometheus.Desc
	entropyPoolSize *prometheus.Desc
	logger          log.Logger
}

func init() {
	registerCollector("entropy", defaultEnabled, NewEntropyCollector)
}

// NewEntropyCollector returns a new Collector exposing entropy stats.
func NewEntropyCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &entropyCollector{
		fs: fs,
		entropyAvail: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "entropy_available_bits"),
			"Bits of available entropy.",
			nil, nil,
		),
		entropyPoolSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "entropy_pool_size_bits"),
			"Bits of entropy pool.",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *entropyCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.KernelRandom()
	if err != nil {
		return fmt.Errorf("failed to get kernel random stats: %w", err)
	}

	if stats.EntropyAvaliable == nil {
		return fmt.Errorf("couldn't get entropy_avail")
	}
	ch <- prometheus.MustNewConstMetric(
		c.entropyAvail, prometheus.GaugeValue, float64(*stats.EntropyAvaliable))

	if stats.PoolSize == nil {
		return fmt.Errorf("couldn't get entropy poolsize")
	}
	ch <- prometheus.MustNewConstMetric(
		c.entropyPoolSize, prometheus.GaugeValue, float64(*stats.PoolSize))

	return nil
}
