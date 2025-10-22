// Copyright The Prometheus Authors
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

var swapLabelNames = []string{"device", "swap_type"}

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

	metrics := make([]SwapsEntry, 0, len(swaps))

	for _, swap := range swaps {
		metrics = append(metrics, SwapsEntry{Device: swap.Filename, Type: swap.Type,
			Priority: swap.Priority, Size: swap.Size, Used: swap.Used})
	}

	return metrics, nil
}

func (c *swapCollector) Update(ch chan<- prometheus.Metric) error {
	swaps, err := c.getSwapInfo()
	if err != nil {
		return fmt.Errorf("couldn't get swap information: %w", err)
	}

	for _, swap := range swaps {
		swapLabelValues := []string{swap.Device, swap.Type}

		// Export swap size in bytes
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, swapSubsystem, "size_bytes"),
				"Swap device size in bytes.",
				[]string{"device", "swap_type"}, nil,
			),
			prometheus.GaugeValue,
			// Size is provided in kbytes (not bytes), translate to bytes
			// see https://github.com/torvalds/linux/blob/fd94619c43360eb44d28bd3ef326a4f85c600a07/mm/swapfile.c#L3079-L3080
			float64(swap.Size*1024),
			swapLabelValues...,
		)

		// Export swap used in bytes
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, swapSubsystem, "used_bytes"),
				"Swap device used in bytes.",
				swapLabelNames, nil,
			),
			prometheus.GaugeValue,
			// Swap used is also provided in kbytes, translate to bytes
			float64(swap.Used*1024),
			swapLabelValues...,
		)

		// Export swap priority
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, swapSubsystem, "priority"),
				"Swap device priority.",
				swapLabelNames, nil,
			),
			prometheus.GaugeValue,
			float64(swap.Priority),
			swapLabelValues...,
		)

	}

	return nil
}
