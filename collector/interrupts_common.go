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

// +build linux openbsd
// +build !nointerrupts

package collector

import (
	"flag"

	"github.com/prometheus/client_golang/prometheus"
)

type interruptsCollector struct {
	metric *prometheus.CounterVec
}

func init() {
	Factories["interrupts"] = NewInterruptsCollector
	CollectorsEnabledState["interrupts"] = flag.Bool("collectors.interrupts.enabled", false, "enables interrupts-collector")
}

// Takes a prometheus registry and returns a new Collector exposing
// interrupts stats
func NewInterruptsCollector() (Collector, error) {
	return &interruptsCollector{
		metric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "interrupts",
				Help:      "Interrupt details.",
			},
			interruptLabelNames,
		),
	}, nil
}
