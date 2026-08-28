// Copyright 2023 The Prometheus Authors
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

package collector

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

const cpuFreqCollectorSubsystem = "cpufreq"

var useCPUFreqPrefix = kingpin.Flag("collector.cpufreq.enable-cpufreq-prefix",
	"Expose cpufreq metrics with the node_cpufreq_ prefix instead of node_cpu_. This avoids a metric name collision with the cpu collector and will be the default behavior in 2.x.").Bool()

type cpuFreqDescs struct {
	hertz           *prometheus.Desc
	avgHertz        *prometheus.Desc
	minHertz        *prometheus.Desc
	maxHertz        *prometheus.Desc
	scalingHertz    *prometheus.Desc
	scalingMinHertz *prometheus.Desc
	scalingMaxHertz *prometheus.Desc
	scalingGovernor *prometheus.Desc
}

// newCPUFreqDescs builds the cpufreq metric descriptions. It must run after
// flag parsing because the metric prefix depends on
// --collector.cpufreq.enable-cpufreq-prefix.
func newCPUFreqDescs() cpuFreqDescs {
	subsystem := cpuCollectorSubsystem
	if *useCPUFreqPrefix {
		subsystem = cpuFreqCollectorSubsystem
	}
	return cpuFreqDescs{
		hertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_hertz"),
			"Current CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		avgHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_avg_hertz"),
			"Average CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		minHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_min_hertz"),
			"Minimum CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		maxHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "frequency_max_hertz"),
			"Maximum CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		scalingHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "scaling_frequency_hertz"),
			"Current scaled CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		scalingMinHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "scaling_frequency_min_hertz"),
			"Minimum scaled CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		scalingMaxHertz: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "scaling_frequency_max_hertz"),
			"Maximum scaled CPU thread frequency in hertz.",
			[]string{"cpu"}, nil,
		),
		scalingGovernor: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "scaling_governor"),
			"Current enabled CPU frequency governor.",
			[]string{"cpu", "governor"}, nil,
		),
	}
}
