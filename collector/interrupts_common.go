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

import "github.com/prometheus/client_golang/prometheus"

type interruptsCollector struct {
	desc typedDesc
}

func init() {
	registerCollector("interrupts", defaultDisabled, NewInterruptsCollector)
}

// NewInterruptsCollector returns a new Collector exposing interrupts stats.
func NewInterruptsCollector() (Collector, error) {
	return &interruptsCollector{
		desc: typedDesc{prometheus.NewDesc(
			namespace+"_interrupts",
			"Interrupt details.",
			interruptLabelNames, nil,
		), prometheus.CounterValue},
	}, nil
}
