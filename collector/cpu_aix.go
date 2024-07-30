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

type cpuCollector struct {
	cpu           typedDesc
	logger        *slog.Logger
	tickPerSecond int64
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func tickPerSecond() (int64, error) {
	ticks, err := C.sysconf(C._SC_CLK_TCK)
	if ticks == -1 || err != nil {
		return 0, fmt.Errorf("failed to get clock ticks per second: %v", err)
	}
	return int64(ticks), nil
}

func NewCpuCollector(logger *slog.Logger) (Collector, error) {
	ticks, err := tickPerSecond()
	if err != nil {
		return nil, err
	}
	return &cpuCollector{
		cpu:           typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		logger:        logger,
		tickPerSecond: ticks,
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.CpuStat()
	if err != nil {
		return err
	}

	for n, stat := range stats {
		ch <- c.cpu.mustNewConstMetric(float64(stat.User/c.tickPerSecond), strconv.Itoa(n), "user")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Sys/c.tickPerSecond), strconv.Itoa(n), "system")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Idle/c.tickPerSecond), strconv.Itoa(n), "idle")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Wait/c.tickPerSecond), strconv.Itoa(n), "wait")
	}
	return nil
}
