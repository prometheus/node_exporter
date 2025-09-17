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

//go:build !nocpu
// +build !nocpu

package collector

/*
#include <unistd.h>  // Include the standard Unix header
#include <errno.h>   // For errno
*/
import "C"
import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	nodeCPUPhysicalSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "physical_seconds_total"),
		"Seconds the physical CPUs spent in each mode.",
		[]string{"cpu", "mode"}, nil,
	)
	nodeCPUSRunQueueDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "runqueue"),
		"Length of the run queue.", []string{"cpu"}, nil,
	)
	nodeCPUFlagsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "flags"),
		"CPU flags.",
		[]string{"cpu", "flag"}, nil,
	)
	nodeCPUContextSwitchDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "context_switches_total"),
		"Number of context switches.",
		[]string{"cpu"}, nil,
	)
)

type cpuCollector struct {
	cpu              typedDesc
	cpuPhysical      typedDesc
	cpuRunQueue      typedDesc
	cpuFlags         typedDesc
	cpuContextSwitch typedDesc

	logger             *slog.Logger
	tickPerSecond      float64
	purrTicksPerSecond float64
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func tickPerSecond() (float64, error) {
	ticks, err := C.sysconf(C._SC_CLK_TCK)
	if ticks == -1 || err != nil {
		return 0, fmt.Errorf("failed to get clock ticks per second: %v", err)
	}
	return float64(ticks), nil
}

func NewCpuCollector(logger *slog.Logger) (Collector, error) {
	ticks, err := tickPerSecond()
	if err != nil {
		return nil, err
	}

	pconfig, err := perfstat.PartitionStat()

	if err != nil {
		return nil, err
	}

	return &cpuCollector{
		cpu:                typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		cpuPhysical:        typedDesc{nodeCPUPhysicalSecondsDesc, prometheus.CounterValue},
		cpuRunQueue:        typedDesc{nodeCPUSRunQueueDesc, prometheus.GaugeValue},
		cpuFlags:           typedDesc{nodeCPUFlagsDesc, prometheus.GaugeValue},
		cpuContextSwitch:   typedDesc{nodeCPUContextSwitchDesc, prometheus.CounterValue},
		logger:             logger,
		tickPerSecond:      ticks,
		purrTicksPerSecond: float64(pconfig.ProcessorMhz * 1e6),
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.CpuStat()
	if err != nil {
		return err
	}

	for n, stat := range stats {
		// LPAR metrics
		ch <- c.cpu.mustNewConstMetric(float64(stat.User)/c.tickPerSecond, strconv.Itoa(n), "user")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Sys)/c.tickPerSecond, strconv.Itoa(n), "system")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Idle)/c.tickPerSecond, strconv.Itoa(n), "idle")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Wait)/c.tickPerSecond, strconv.Itoa(n), "wait")

		// Physical CPU metrics
		ch <- c.cpuPhysical.mustNewConstMetric(float64(stat.PIdle)/c.purrTicksPerSecond, strconv.Itoa(n), "pidle")
		ch <- c.cpuPhysical.mustNewConstMetric(float64(stat.PUser)/c.purrTicksPerSecond, strconv.Itoa(n), "puser")
		ch <- c.cpuPhysical.mustNewConstMetric(float64(stat.PSys)/c.purrTicksPerSecond, strconv.Itoa(n), "psys")
		ch <- c.cpuPhysical.mustNewConstMetric(float64(stat.PWait)/c.purrTicksPerSecond, strconv.Itoa(n), "pwait")

		// Run queue length
		ch <- c.cpuRunQueue.mustNewConstMetric(float64(stat.RunQueue), strconv.Itoa(n))

		// Flags
		ch <- c.cpuFlags.mustNewConstMetric(float64(stat.SpurrFlag), strconv.Itoa(n), "spurr")

		// Context switches
		ch <- c.cpuContextSwitch.mustNewConstMetric(float64(stat.CSwitches), strconv.Itoa(n))
	}
	return nil
}
