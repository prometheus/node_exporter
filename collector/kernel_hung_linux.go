// Copyright 2018 The Prometheus Authors
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

//go:build !noprocesses

package collector

import (
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type kernelHungCollector struct {
	fs     procfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("kernel_hung", defaultDisabled, NewKernelHungCollector)
}

func NewKernelHungCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &kernelHungCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

var (
	taskDetectCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "kernel_hung", "task_detect_count"),
		"Total number of interrupts serviced.",
		nil, nil,
	)
)

func (c *kernelHungCollector) Update(ch chan<- prometheus.Metric) error {
	kernelHung, err := c.fs.KernelHung()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(taskDetectCount, prometheus.CounterValue, float64(*kernelHung.HungTaskDetectCount))

	return nil
}
