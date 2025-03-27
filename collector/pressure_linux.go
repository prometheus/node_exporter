// Copyright 2019 The Prometheus Authors
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

//go:build !nopressure
// +build !nopressure

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	psiResourceCPU    = "cpu"
	psiResourceIO     = "io"
	psiResourceMemory = "memory"
	psiResourceIRQ    = "irq"
)

var (
	psiResources = []string{psiResourceCPU, psiResourceIO, psiResourceMemory, psiResourceIRQ}
)

type pressureStatsCollector struct {
	cpu     *prometheus.Desc
	io      *prometheus.Desc
	ioFull  *prometheus.Desc
	mem     *prometheus.Desc
	memFull *prometheus.Desc
	irqFull *prometheus.Desc

	fs procfs.FS

	logger *slog.Logger
}

func init() {
	registerCollector("pressure", defaultEnabled, NewPressureStatsCollector)
}

// NewPressureStatsCollector returns a Collector exposing pressure stall information
func NewPressureStatsCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &pressureStatsCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "cpu_waiting_seconds_total"),
			"Total time in seconds that processes have waited for CPU time",
			nil, nil,
		),
		io: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "io_waiting_seconds_total"),
			"Total time in seconds that processes have waited due to IO congestion",
			nil, nil,
		),
		ioFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "io_stalled_seconds_total"),
			"Total time in seconds no process could make progress due to IO congestion",
			nil, nil,
		),
		mem: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "memory_waiting_seconds_total"),
			"Total time in seconds that processes have waited for memory",
			nil, nil,
		),
		memFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "memory_stalled_seconds_total"),
			"Total time in seconds no process could make progress due to memory congestion",
			nil, nil,
		),
		irqFull: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pressure", "irq_stalled_seconds_total"),
			"Total time in seconds no process could make progress due to IRQ congestion",
			nil, nil,
		),
		fs:     fs,
		logger: logger,
	}, nil
}

// Update calls procfs.NewPSIStatsForResource for the different resources and updates the values
func (c *pressureStatsCollector) Update(ch chan<- prometheus.Metric) error {
	foundResources := 0
	for _, res := range psiResources {
		c.logger.Debug("collecting statistics for resource", "resource", res)
		vals, err := c.fs.PSIStatsForResource(res)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) && res != psiResourceIRQ {
				c.logger.Debug("pressure information is unavailable, you need a Linux kernel >= 4.20 and/or CONFIG_PSI enabled for your kernel", "resource", res)
				continue
			}
			if errors.Is(err, os.ErrNotExist) && res == psiResourceIRQ {
				c.logger.Debug("pressure information is unavailable, you need a Linux kernel >= 6.1 and/or CONFIG_PSI enabled for your kernel", "resource", res)
				continue
			}
			if errors.Is(err, syscall.ENOTSUP) {
				c.logger.Debug("pressure information is disabled, add psi=1 kernel command line to enable it")
				return ErrNoData
			}
			return fmt.Errorf("failed to retrieve pressure stats: %w", err)
		}
		// IRQ pressure does not have 'some' data.
		// See https://github.com/torvalds/linux/blob/v6.9/include/linux/psi_types.h#L65
		if vals.Some == nil && res != "irq" {
			c.logger.Debug("pressure information returned no 'some' data")
			return ErrNoData
		}
		if vals.Full == nil && res != "cpu" {
			c.logger.Debug("pressure information returned no 'full' data")
			return ErrNoData
		}
		switch res {
		case "cpu":
			ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
		case "io":
			ch <- prometheus.MustNewConstMetric(c.io, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
			ch <- prometheus.MustNewConstMetric(c.ioFull, prometheus.CounterValue, float64(vals.Full.Total)/1000.0/1000.0)
		case "memory":
			ch <- prometheus.MustNewConstMetric(c.mem, prometheus.CounterValue, float64(vals.Some.Total)/1000.0/1000.0)
			ch <- prometheus.MustNewConstMetric(c.memFull, prometheus.CounterValue, float64(vals.Full.Total)/1000.0/1000.0)
		case "irq":
			ch <- prometheus.MustNewConstMetric(c.irqFull, prometheus.CounterValue, float64(vals.Full.Total)/1000.0/1000.0)
		default:
			c.logger.Debug("did not account for resource", "resource", res)
			continue
		}
		foundResources++
	}

	if foundResources == 0 {
		c.logger.Debug("pressure information is unavailable, you need a Linux kernel >= 4.20 and/or CONFIG_PSI enabled for your kernel")
		return ErrNoData
	}

	return nil
}
