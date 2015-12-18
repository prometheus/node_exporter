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
	cpu          *prometheus.Desc
	intr         *prometheus.Desc
	ctxt         *prometheus.Desc
	forks        *prometheus.Desc
	btime        *prometheus.Desc
	procsRunning *prometheus.Desc
	procsBlocked *prometheus.Desc
}

func init() {
	Factories["stat"] = NewStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// kernel/system statistics.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "cpu"),
			"Seconds the cpus spent in each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		intr: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "intr"),
			"Total number of interrupts serviced.",
			nil, nil,
		),
		ctxt: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "context_switches"),
			"Total number of context switches.",
			nil, nil,
		),
		forks: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "forks"),
			"Total number of forks.",
			nil, nil,
		),
		btime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "boot_time"),
			"Node boot time, in unixtime.",
			nil, nil,
		),
		procsRunning: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "procs_running"),
			"Number of processes in runnable state.",
			nil, nil,
		),
		procsBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "procs_blocked"),
			"Number of processes blocked waiting for I/O to complete.",
			nil, nil,
		),
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
				ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, value, parts[0], cpuFields[i])
			}
		case parts[0] == "intr":
			// Only expose the overall number, use the 'interrupts' collector for more detail.
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.intr, prometheus.CounterValue, value)
		case parts[0] == "ctxt":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.ctxt, prometheus.CounterValue, value)
		case parts[0] == "processes":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.forks, prometheus.CounterValue, value)
		case parts[0] == "btime":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.btime, prometheus.GaugeValue, value)
		case parts[0] == "procs_running":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.procsRunning, prometheus.GaugeValue, value)
		case parts[0] == "procs_blocked":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return err
			}
			ch <- prometheus.MustNewConstMetric(c.procsBlocked, prometheus.GaugeValue, value)
		}
	}
	return err
}
