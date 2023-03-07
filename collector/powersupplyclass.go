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

//go:build !nopowersupplyclass && (linux || darwin)
// +build !nopowersupplyclass
// +build linux darwin

package collector

import (
	"regexp"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	powerSupplyClassIgnoredPowerSupplies = kingpin.Flag("collector.powersupply.ignored-supplies", "Regexp of power supplies to ignore for powersupplyclass collector.").Default("^$").String()
)

type powerSupplyClassCollector struct {
	subsystem      string
	ignoredPattern *regexp.Regexp
	metricDescs    map[string]*prometheus.Desc
	logger         log.Logger
}

func init() {
	registerCollector("powersupplyclass", defaultEnabled, NewPowerSupplyClassCollector)
}

func NewPowerSupplyClassCollector(logger log.Logger) (Collector, error) {
	pattern := regexp.MustCompile(*powerSupplyClassIgnoredPowerSupplies)
	return &powerSupplyClassCollector{
		subsystem:      "power_supply",
		ignoredPattern: pattern,
		metricDescs:    map[string]*prometheus.Desc{},
		logger:         logger,
	}, nil
}
