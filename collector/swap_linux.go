// Copyright 2025 The Prometheus Authors
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

//go:build !noswap
// +build !noswap

package collector

import (
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	swapSubsystem = "swap"
)

type swapCollector struct {
	fs     procfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("swap", defaultDisabled, NewSwapCollector)
}

// NewSwapCollector returns a new Collector exposing swap device statistics.
func NewSwapCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &swapCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

type SwapsEntry struct {
	Device   string
	Type     string
	Priority int
	Size     int
	Used     int
}

func (c *swapCollector) getSwapInfo() ([]SwapsEntry, error) {
	swaps, err := c.fs.Swaps()
	if err != nil {
		return nil, fmt.Errorf("couldn't get proc/swap information: %w", err)
	}

	metrics := make([]SwapsEntry, 0)

	for _, swap := range swaps {
		metrics = append(metrics, SwapsEntry{swap.Filename, swap.Type, swap.Priority, swap.Size, swap.Used})
	}

	return metrics, nil
}

func (c *swapCollector) Update(ch chan<- prometheus.Metric) error {
	swaps, err := c.getSwapInfo()
	if err != nil {
		return fmt.Errorf("couldn't get swap information: %w", err)
	}

	for _, swap := range swaps {
		labels := []string{swap.Device, swap.Type, fmt.Sprintf("%d", swap.Priority)}

		// Export swap size in bytes
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, swapSubsystem, "size_bytes"),
				"Swap device size in bytes.",
				[]string{"device", "type", "priority"}, nil,
			),
			prometheus.GaugeValue, float64(swap.Size*1024), labels...,
		)

		// Export swap used in bytes
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, swapSubsystem, "used_bytes"),
				"Swap device used in bytes.",
				[]string{"device", "type", "priority"}, nil,
			),
			prometheus.GaugeValue, float64(swap.Used*1024), labels...,
		)
	}

	return nil
}
