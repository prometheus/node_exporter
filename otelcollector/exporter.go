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

package otelcollector

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/exporter-toolkit/otlpreceiver"
)

type NodeExporter struct {
	config   *Config
	registry *prometheus.Registry
}

func NewNodeExporter(config *Config) *NodeExporter {
	return &NodeExporter{
		config:   config,
		registry: prometheus.NewRegistry(),
	}
}

func (ne *NodeExporter) Initialize(ctx context.Context, cfg otlpreceiver.Config) (*prometheus.Registry, error) {
	var exporterCfg *Config

	if cfg != nil {
		var ok bool
		exporterCfg, ok = cfg.(*Config)
		if !ok {
			return nil, errors.New("error reading configuration")
		}
	} else {
		// Use default configuration when none is provided
		exporterCfg = &Config{
			DisableDefaults:   false,
			EnableCollectors:  []string{},
			ExcludeCollectors: []string{},
		}
	}

	ne.registry.MustRegister(versioncollector.NewCollector("node_exporter"))

	ne.config = exporterCfg
	return ne.registry, nil
}

func (ne *NodeExporter) Shutdown(_ context.Context) error {
	// There's nothing special needed to shutdown node-exporter.
	return nil
}
