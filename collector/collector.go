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

// Package collector includes all individual collectors to gather and export system metrics.
package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Namespace defines the common namespace to be used by all metrics.
const namespace = "node"

// Factories contains the list of all available collectors.
var Factories = make(map[string]func() (Collector, error))

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"node_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"node_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)

	disableDefaultCollectors = kingpin.Flag("collectors.disable-defaults", "Do not use the default collectors, only those explicitly enabled on the commandline.").Bool()
)

func warnDeprecated(collector string) {
	log.Warnf("The %s collector is deprecated and will be removed in the future!", collector)
}

const (
	defaultEnabled  = true
	defaultDisabled = false
)

type collectorState struct {
	flagState    bool
	flagSet      bool
	defaultState bool
}

func (state *collectorState) set(c *kingpin.ParseContext) error {
	state.flagSet = true
	return nil
}

var collectorStates = make(map[string]*collectorState)

func registerCollector(collector string, defaultState bool, factory func() (Collector, error)) {
	flagName := fmt.Sprintf("collector.%s.enabled", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector.", collector)

	state := collectorState{
		defaultState: defaultState,
	}
	kingpin.Flag(flagName, flagHelp).Action(state.set).BoolVar(&state.flagState)
	collectorStates[collector] = &state

	Factories[collector] = factory
}

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	Collectors map[string]Collector
}

func NewNodeCollector() (*NodeCollector, error) {
	collectors := make(map[string]Collector)
	for key, state := range collectorStates {
		enable := false
		// Enable the collector if it has been enabled by a flag, OR if it is enabled by default, and the defaults are not disabled
		if state.flagSet {
			enable = state.flagState
		} else if state.defaultState {
			enable = !*disableDefaultCollectors
		}

		if enable {
			collector, err := Factories[key]()
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
		}
	}
	return &NodeCollector{Collectors: collectors}, nil
}

// Describe implements the prometheus.Collector interface.
func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}
