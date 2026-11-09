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
	"maps"
	"sync"
	"time"

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
	factories    = make(map[string]func(cfg config.Config, logger *slog.Logger) (Collector, error))
	defaultState = make(map[string]bool)
)

type collectorRuntimeState struct {
	config     config.Config
	enabled    map[string]bool
	collectors map[string]Collector
	mtx        sync.Mutex
}

type runtimeConfiguredCollector interface {
	configureRuntimeState(*collectorRuntimeState)
}

func registerCollector(collector string, isDefaultEnabled bool, factory func(cfg config.Config, logger *slog.Logger) (Collector, error)) {
	defaultState[collector] = isDefaultEnabled
	factories[collector] = factory
}

// AvailableCollectors returns the name of every collector in this build and
// whether it is enabled by default.
func AvailableCollectors() map[string]bool {
	return maps.Clone(defaultState)
}

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	Collectors map[string]Collector
	logger     *slog.Logger
}

func newCollectorRuntimeState(cfg config.Config) (*collectorRuntimeState, error) {
	for collector := range cfg.Collectors {
		if _, ok := defaultState[collector]; !ok {
			return nil, fmt.Errorf("unknown collector: %s", collector)
		}
	}

	enabled := make(map[string]bool, len(defaultState))
	for collector, isDefaultEnabled := range defaultState {
		override, overridden := cfg.Collectors[collector]
		switch {
		case overridden:
			enabled[collector] = override
		case cfg.CollectorDisableDefaults:
			enabled[collector] = false
		default:
			enabled[collector] = isDefaultEnabled
		}
	}

	return &collectorRuntimeState{
		config:     cfg,
		enabled:    enabled,
		collectors: make(map[string]Collector),
	}, nil
}

// flagConfigs copy the parsed command-line flags into a Config. Each file
// that defines collector flags registers one, so a build that leaves out a
// collector also leaves out its flags.
var flagConfigs []func(*config.Config)

func registerFlagConfig(f func(*config.Config)) {
	flagConfigs = append(flagConfigs, f)
}

// ConfigFromFlags returns cfg with the options set from the parsed
// command-line flags.
func ConfigFromFlags(cfg config.Config) config.Config {
	for _, f := range flagConfigs {
		f(&cfg)
	}
	return cfg
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
			collector, err := factories[key](s.config, logger.With("collector", key))
			if err != nil {
				return nil, err
			}
			s.configureCollector(collector)
			collectors[key] = collector
			s.collectors[key] = collector
		}
	}
	return &NodeCollector{Collectors: collectors, logger: logger}, nil
}

func (s *collectorRuntimeState) configureCollector(collector Collector) {
	if configuredCollector, ok := collector.(runtimeConfiguredCollector); ok {
		configuredCollector.configureRuntimeState(s)
	}
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
