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
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/config"
)

// Namespace defines the common namespace to be used by all metrics.
const namespace = "node"

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
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	factories        = make(map[string]func(logger *slog.Logger) (Collector, error))
	collectorState   = make(map[string]*bool)
	defaultState     = make(map[string]bool)
	forcedCollectors = map[string]bool{} // collectors which have been explicitly enabled or disabled
)

type collectorRuntimeState struct {
	config     collectorConfig
	enabled    map[string]bool
	collectors map[string]Collector
	mtx        sync.Mutex
}

type runtimeConfiguredCollector interface {
	configureRuntimeState(*collectorRuntimeState) error
}

type collectorConfig struct {
	paths    collectorPathConfig
	cpu      config.CPUConfig
	textfile config.TextfileConfig
}

type collectorPathConfig struct {
	procMountPoint string
	sysMountPoint  string
	rootfsPath     string
	udevDataPath   string
}

func registerCollector(collector string, isDefaultEnabled bool, factory func(logger *slog.Logger) (Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Action(collectorFlagAction(collector)).Bool()
	collectorState[collector] = flag
	defaultState[collector] = isDefaultEnabled

	factories[collector] = factory
}

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	Collectors map[string]Collector
	logger     *slog.Logger
}

func newCollectorRuntimeState(cfg config.Config) (*collectorRuntimeState, error) {
	overrides := cfg.Collectors
	for collector := range overrides {
		if _, ok := collectorState[collector]; !ok {
			return nil, fmt.Errorf("unknown collector: %s", collector)
		}
	}

	enabled := make(map[string]bool, len(collectorState))
	for collector, state := range collectorState {
		// Use the flag value only when the flag was set on the command line,
		// so a runtime built without kingpin.Parse still gets the defaults.
		override, overridden := overrides[collector]
		switch {
		case overridden:
			enabled[collector] = override
		case forcedCollectors[collector]:
			enabled[collector] = *state
		case cfg.CollectorDisableDefaults:
			enabled[collector] = false
		default:
			enabled[collector] = defaultState[collector]
		}
	}

	return &collectorRuntimeState{
		config:     newCollectorConfig(cfg),
		enabled:    enabled,
		collectors: make(map[string]Collector),
	}, nil
}

func newCollectorConfig(cfg config.Config) collectorConfig {
	return collectorConfig{
		paths:    newCollectorPathConfig(),
		cpu:      cfg.CPU,
		textfile: cfg.Textfile,
	}
}

// ConfigFromFlags returns cfg with the collector options set from the parsed
// command-line flags.
func ConfigFromFlags(cfg config.Config) config.Config {
	cfg.CPU = cpuConfigFromFlags(cfg.CPU)
	cfg.Textfile = textfileConfigFromFlags(cfg.Textfile)
	return cfg
}

// collectorFlagAction generates a new action function for the given collector
// to track whether it has been explicitly enabled or disabled from the command line.
// A new action function is needed for each collector flag because the ParseContext
// does not contain information about which flag called the action.
// See: https://github.com/alecthomas/kingpin/issues/294
func collectorFlagAction(collector string) func(ctx *kingpin.ParseContext) error {
	return func(_ *kingpin.ParseContext) error {
		forcedCollectors[collector] = true
		return nil
	}
}

// NewNodeCollector creates a new NodeCollector.
func (s *collectorRuntimeState) NewNodeCollector(logger *slog.Logger, filters ...string) (*NodeCollector, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	f := make(map[string]bool)
	for _, filter := range filters {
		enabled, exist := s.enabled[filter]
		if !exist {
			return nil, fmt.Errorf("missing collector: %s", filter)
		}
		if !enabled {
			return nil, fmt.Errorf("disabled collector: %s", filter)
		}
		f[filter] = true
	}
	collectors := make(map[string]Collector)
	for key, enabled := range s.enabled {
		if !enabled || (len(f) > 0 && !f[key]) {
			continue
		}
		if collector, ok := s.collectors[key]; ok {
			collectors[key] = collector
		} else {
			collector, err := factories[key](logger.With("collector", key))
			if err != nil {
				return nil, err
			}
			if err := s.configureCollector(collector); err != nil {
				return nil, err
			}
			collectors[key] = collector
			s.collectors[key] = collector
		}
	}
	return &NodeCollector{Collectors: collectors, logger: logger}, nil
}

var legacyState = &collectorRuntimeState{collectors: make(map[string]Collector)}

// DisableDefaultCollectors sets the collector flags to disabled, except the
// flags set on the command line.
//
// Deprecated: Use NewRuntime with config.Config.CollectorDisableDefaults.
func DisableDefaultCollectors() {
	for c := range collectorState {
		if _, ok := forcedCollectors[c]; !ok {
			*collectorState[c] = false
		}
	}
}

// NewNodeCollector creates a NodeCollector from the collector flags. All
// callers share one set of collector instances.
//
// Deprecated: Use NewRuntime.
func NewNodeCollector(logger *slog.Logger, filters ...string) (*NodeCollector, error) {
	enabled := make(map[string]bool, len(collectorState))
	for c, state := range collectorState {
		enabled[c] = *state
	}
	legacyState.mtx.Lock()
	legacyState.enabled = enabled
	legacyState.config = newCollectorConfig(ConfigFromFlags(config.NewConfigWithDefaults()))
	legacyState.mtx.Unlock()
	return legacyState.NewNodeCollector(logger, filters...)
}

func (s *collectorRuntimeState) configureCollector(collector Collector) error {
	if configuredCollector, ok := collector.(runtimeConfiguredCollector); ok {
		return configuredCollector.configureRuntimeState(s)
	}

	return nil
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
			execute(name, c, ch, n.logger)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c Collector, ch chan<- prometheus.Metric, logger *slog.Logger) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		if IsNoDataError(err) {
			logger.Debug("collector returned no data", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		} else {
			logger.Error("collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		}
		success = 0
	} else {
		logger.Debug("collector succeeded", "name", name, "duration_seconds", duration.Seconds())
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

// ErrNoData indicates the collector found no data to collect, but had no other error.
var ErrNoData = errors.New("collector returned no data")

func IsNoDataError(err error) bool {
	return err == ErrNoData
}

// pushMetric helps construct and convert a variety of value types into Prometheus float64 metrics.
func pushMetric(ch chan<- prometheus.Metric, fieldDesc *prometheus.Desc, value any, valueType prometheus.ValueType, labelValues ...string) {
	var fVal float64
	switch val := value.(type) {
	case uint8:
		fVal = float64(val)
	case uint16:
		fVal = float64(val)
	case uint32:
		fVal = float64(val)
	case uint64:
		fVal = float64(val)
	case int64:
		fVal = float64(val)
	case *uint8:
		if val == nil {
			return
		}
		fVal = float64(*val)
	case *uint16:
		if val == nil {
			return
		}
		fVal = float64(*val)
	case *uint32:
		if val == nil {
			return
		}
		fVal = float64(*val)
	case *uint64:
		if val == nil {
			return
		}
		fVal = float64(*val)
	case *int64:
		if val == nil {
			return
		}
		fVal = float64(*val)
	default:
		return
	}

	ch <- prometheus.MustNewConstMetric(fieldDesc, valueType, fVal, labelValues...)
}
