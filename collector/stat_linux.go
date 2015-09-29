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

// +build !nostat

package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	userHz = 100
)

type statCollector struct {
	cpu          *prometheus.CounterVec
	intr         prometheus.Counter
	ctxt         prometheus.Counter
	forks        prometheus.Counter
	btime        prometheus.Gauge
	procsRunning prometheus.Gauge
	procsBlocked prometheus.Gauge
}

func init() {
	Factories["stat"] = NewStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// kernel/system statistics.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "cpu",
				Help:      "Seconds the cpus spent in each mode.",
			},
			[]string{"cpu", "mode"},
		),
		intr: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "intr",
			Help:      "Total number of interrupts serviced.",
		}),
		ctxt: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "context_switches",
			Help:      "Total number of context switches.",
		}),
		forks: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "forks",
			Help:      "Total number of forks.",
		}),
		btime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "boot_time",
			Help:      "Node boot time, in unixtime.",
		}),
		procsRunning: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "procs_running",
			Help:      "Number of processes in runnable state.",
		}),
		procsBlocked: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "procs_blocked",
			Help:      "Number of processes blocked waiting for I/O to complete.",
		}),
	}, nil
}

// Expose kernel and system statistics.
func (c *statCollector) Update(ch chan<- prometheus.Metric) (err error) {
	file, err := os.Open(procFilePath("stat"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(parts[0], "cpu"):
			// Export only per-cpu stats, it can be aggregated up in prometheus.
			if parts[0] == "cpu" {
				break
			}
			// Only some of these may be present, depending on kernel version.
			cpuFields := []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"}
			// OpenVZ guests lack the "guest" CPU field, which needs to be ignored.
			expectedFieldNum := len(cpuFields) + 1
			if expectedFieldNum > len(parts) {
				expectedFieldNum = len(parts)
			}
			for i, v := range parts[1:expectedFieldNum] {
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				// Convert from ticks to seconds
				value /= userHz
				c.cpu.With(prometheus.Labels{"cpu": parts[0], "mode": cpuFields[i]}).Set(value)
			}
		case parts[0] == "intr":
			// Only expose the overall number, use the 'interrupts' collector for more detail.
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.intr.Set(value)
		case parts[0] == "ctxt":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.ctxt.Set(value)
		case parts[0] == "processes":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.forks.Set(value)
		case parts[0] == "btime":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.btime.Set(value)
		case parts[0] == "procs_running":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.procsRunning.Set(value)
		case parts[0] == "procs_blocked":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			c.procsBlocked.Set(value)
		}
	}
	c.cpu.Collect(ch)
	c.ctxt.Collect(ch)
	c.intr.Collect(ch)
	c.forks.Collect(ch)
	c.btime.Collect(ch)
	c.procsRunning.Collect(ch)
	c.procsBlocked.Collect(ch)
	return err
}
