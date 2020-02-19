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

// +build solaris
// +build !nocpu

package collector

import (
	"fmt"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	kstat "github.com/siebenmann/go-kstat"
)

// #include <unistd.h>
import "C"

type cpuFreqCollector struct {
	cpuFreq    *prometheus.Desc
	cpuFreqMax *prometheus.Desc
	logger     log.Logger
}

func init() {
	registerCollector("cpufreq", defaultEnabled, NewCpuFreqCollector)
}

func NewCpuFreqCollector(logger log.Logger) (Collector, error) {
	return &cpuFreqCollector{
		cpuFreq: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz"),
			"Current cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		cpuFreqMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_max_hertz"),
			"Maximum cpu thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *cpuFreqCollector) Update(ch chan<- prometheus.Metric) error {
	ncpus := C.sysconf(C._SC_NPROCESSORS_ONLN)

	tok, err := kstat.Open()
	if err != nil {
		return err
	}

	defer tok.Close()

	for cpu := 0; cpu < int(ncpus); cpu++ {
		ksCPUInfo, err := tok.Lookup("cpu_info", cpu, fmt.Sprintf("cpu_info%d", cpu))
		if err != nil {
			return err
		}
		cpuFreqV, err := ksCPUInfo.GetNamed("current_clock_Hz")
		if err != nil {
			return err
		}

		cpuFreqMaxV, err := ksCPUInfo.GetNamed("clock_MHz")
		if err != nil {
			return err
		}

		lcpu := strconv.Itoa(cpu)
		ch <- prometheus.MustNewConstMetric(
			c.cpuFreq,
			prometheus.GaugeValue,
			float64(cpuFreqV.UintVal),
			lcpu,
		)
		// Multiply by 1e+6 to convert MHz to Hz.
		ch <- prometheus.MustNewConstMetric(
			c.cpuFreqMax,
			prometheus.GaugeValue,
			float64(cpuFreqMaxV.IntVal)*1e+6,
			lcpu,
		)
	}
	return nil
}
