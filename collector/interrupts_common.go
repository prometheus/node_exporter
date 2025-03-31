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

//go:build (linux || openbsd) && !nointerrupts
// +build linux openbsd
// +build !nointerrupts

package collector

import (
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type interruptsCollector struct {
	desc         typedDesc
	logger       *slog.Logger
	nameFilter   deviceFilter
	includeZeros bool
}

func init() {
	registerCollector("interrupts", defaultDisabled, NewInterruptsCollector)
}

var (
	interruptsInclude      = kingpin.Flag("collector.interrupts.name-include", "Regexp of interrupts name to include (mutually exclusive to --collector.interrupts.name-exclude).").String()
	interruptsExclude      = kingpin.Flag("collector.interrupts.name-exclude", "Regexp of interrupts name to exclude (mutually exclusive to --collector.interrupts.name-include).").String()
	interruptsIncludeZeros = kingpin.Flag("collector.interrupts.include-zeros", "Include interrupts that have a zero value").Default("true").Bool()
)

// NewInterruptsCollector returns a new Collector exposing interrupts stats.
func NewInterruptsCollector(logger *slog.Logger) (Collector, error) {
	return &interruptsCollector{
		desc: typedDesc{prometheus.NewDesc(
			namespace+"_interrupts_total",
			"Interrupt details.",
			interruptLabelNames, nil,
		), prometheus.CounterValue},
		logger:       logger,
		nameFilter:   newDeviceFilter(*interruptsExclude, *interruptsInclude),
		includeZeros: *interruptsIncludeZeros,
	}, nil
}
