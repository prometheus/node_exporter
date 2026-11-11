// Copyright The Prometheus Authors
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
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/config"
)

// Runtime represents a single exporter runtime instance built from a reusable
// Config.
type Runtime struct {
	state     *collectorRuntimeState
	collector *NodeCollector
	logger    *slog.Logger
	cancel    context.CancelFunc
}

func NewRuntime(cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	return NewRuntimeWithContext(context.Background(), cfg, logger)
}

// NewRuntimeWithContext creates a Runtime whose lifetime is also bounded by ctx.
// Long-lived embedded consumers should use NewRuntime and call Runtime.Shutdown.
func NewRuntimeWithContext(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Runtime, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	runtimeCtx, cancel := context.WithCancel(ctx)
	state, err := newCollectorRuntimeState(runtimeCtx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	r, err := newRuntime(state, logger)
	if err != nil {
		cancel()
		return nil, err
	}
	r.cancel = cancel
	return r, nil
}

func newRuntime(state *collectorRuntimeState, logger *slog.Logger, filters ...string) (*Runtime, error) {
	nodeCollector, err := state.NewNodeCollector(logger, filters...)
	if err != nil {
		return nil, err
	}

	return &Runtime{state: state, collector: nodeCollector, logger: logger}, nil
}

// Filtered returns a Runtime with a subset of the enabled collectors. It
// shares the collectors and the lifetime of r, so only r needs Shutdown.
func (r *Runtime) Filtered(enabledCollectors ...string) (*Runtime, error) {
	for _, name := range enabledCollectors {
		if _, ok := r.collector.Collectors[name]; !ok {
			return nil, fmt.Errorf("collector not enabled in this runtime: %s", name)
		}
	}
	return newRuntime(r.state, r.logger, enabledCollectors...)
}

// Shutdown cancels in-flight collection work owned by the Runtime.
func (r *Runtime) Shutdown(context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *Runtime) Collectors() []prometheus.Collector {
	return []prometheus.Collector{r.collector}
}

func (r *Runtime) Registry() (*prometheus.Registry, error) {
	registry := prometheus.NewRegistry()
	if err := registry.Register(r.collector); err != nil {
		return nil, err
	}
	return registry, nil
}

func (r *Runtime) EnabledCollectors() []string {
	return slices.Sorted(maps.Keys(r.collector.Collectors))
}
