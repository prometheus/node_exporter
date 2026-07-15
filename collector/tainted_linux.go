// Copyright 2024 The Prometheus Authors
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

package collector

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type taintedCollector struct {
	logger *slog.Logger
	desc   *prometheus.Desc
}

func init() {
	registerCollector("tainted", defaultDisabled, NewTaintedCollector)
}

// NewTaintedCollector returns a Collector exposing kernel taint flags from
// /proc/sys/kernel/tainted as a labelled gauge.
// See https://www.kernel.org/doc/html/latest/admin-guide/tainted-kernels.html
func NewTaintedCollector(logger *slog.Logger) (Collector, error) {
	return &taintedCollector{
		logger: logger,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "kernel", "tainted"),
			"Taint flags set on the running Linux kernel, as reported by /proc/sys/kernel/tainted. "+
				"Value is 1 if the flag is set, 0 otherwise. "+
				"See https://www.kernel.org/doc/html/latest/admin-guide/tainted-kernels.html for flag meanings.",
			[]string{"bit", "flag"},
			nil,
		),
	}, nil
}

func (c *taintedCollector) Update(ch chan<- prometheus.Metric) error {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %w", err)
	}

	tainted, err := fs.KernelTainted()
	if err != nil {
		return fmt.Errorf("couldn't read kernel tainted state: %w", err)
	}

	for _, b := range tainted.Bits {
		var val float64
		if b.Set {
			val = 1.0
		}
		ch <- prometheus.MustNewConstMetric(
			c.desc,
			prometheus.GaugeValue,
			val,
			strconv.Itoa(b.Index),
			b.Flag,
		)
	}
	return nil
}
