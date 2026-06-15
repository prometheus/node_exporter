// Copyright 2026 The Prometheus Authors
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

package collector

import (
	"log/slog"
	"sort"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/config"
)

// Runtime represents a single exporter runtime instance built from a reusable
// Config.
type Runtime struct {
	state     *collectorRuntimeState
	collector *NodeCollector
	logger    *slog.Logger
}

func NewRuntime(cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	state := newCollectorRuntimeState(cfg.CollectorDisableDefaults)

	return newRuntime(state, logger, cfg.EnabledCollectors)
}

func newRuntime(state *collectorRuntimeState, logger *slog.Logger, enabledCollectors []string) (*Runtime, error) {
	nodeCollector, err := state.NewNodeCollector(logger, enabledCollectors...)
	if err != nil {
		return nil, err
	}

	return &Runtime{state: state, collector: nodeCollector, logger: logger}, nil
}

func (r *Runtime) Filtered(enabledCollectors ...string) (*Runtime, error) {
	filtered, err := newRuntime(r.state, r.logger, enabledCollectors)
	if err != nil {
		return nil, err
	}

	return filtered, nil
}

func (r *Runtime) Collectors() []prometheus.Collector {
	return []prometheus.Collector{r.collector}
}

func (r *Runtime) Registry() (*prometheus.Registry, error) {
	registry := prometheus.NewRegistry()
	for _, c := range r.Collectors() {
		if err := registry.Register(c); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

func (r *Runtime) EnabledCollectors() []string {
	enabled := make([]string, 0, len(r.collector.Collectors))
	for name := range r.collector.Collectors {
		enabled = append(enabled, name)
	}
	sort.Strings(enabled)
	return enabled
}
