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
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/exporter-toolkit/otlpreceiver"
	"github.com/prometheus/node_exporter/collector"
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
	if cfg != nil {
		var ok bool
		ne.config, ok = cfg.(*Config)
		if !ok {
			return nil, errors.New("error reading configuration")
		}
	} else {
		// Use default configuration when none is provided
		ne.config = &Config{
			DisableDefaults:   false,
			EnableCollectors:  []string{},
			ExcludeCollectors: []string{},
		}
	}

	// Create logger for collectors
	logger := slog.Default().With("component", "node_exporter")

	collector.InitializeCollectorStateForOTEL()

	// Apply disable defaults configuration
	if ne.config.DisableDefaults {
		collector.DisableDefaultCollectors()
	}

	// Determine which collectors to create based on configuration
	var filters []string

	// Validate that both EnableCollectors and ExcludeCollectors are not used together
	if len(ne.config.EnableCollectors) > 0 && len(ne.config.ExcludeCollectors) > 0 {
		return nil, fmt.Errorf("enable_collectors and exclude_collectors cannot be used together")
	}

	if len(ne.config.EnableCollectors) > 0 {
		// If specific collectors are enabled, use only those
		filters = ne.config.EnableCollectors
	} else if len(ne.config.ExcludeCollectors) > 0 {
		// If excludes are specified, we need to create a list of all enabled collectors
		// minus the excluded ones, similar to how the main node_exporter handles it
		filters = []string{}

		// First, create a temporary NodeCollector to get the list of enabled collectors
		tempNC, err := collector.NewNodeCollector(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to get available collectors: %w", err)
		}

		// Get all enabled collector names
		for collectorName := range tempNC.Collectors {
			// Check if this collector is not in the exclude list
			excluded := false
			for _, excludeName := range ne.config.ExcludeCollectors {
				if collectorName == excludeName {
					excluded = true
					break
				}
			}
			if !excluded {
				filters = append(filters, collectorName)
			}
		}
	}
	// If neither EnableCollectors nor ExcludeCollectors are specified,
	// filters remains empty and NewNodeCollector will use all enabled collectors

	// Create the node collector with all enabled collectors
	nc, err := collector.NewNodeCollector(logger, filters...)
	if err != nil {
		return nil, fmt.Errorf("failed to create node collector: %w", err)
	}

	// Register version collector
	ne.registry.MustRegister(versioncollector.NewCollector("node_exporter"))

	// Register the node collector
	if err := ne.registry.Register(nc); err != nil {
		return nil, fmt.Errorf("failed to register node collector: %w", err)
	}

	logger.Info("Node exporter initialized successfully",
		"disable_defaults", ne.config.DisableDefaults,
		"enable_collectors", ne.config.EnableCollectors,
		"exclude_collectors", ne.config.ExcludeCollectors)

	return ne.registry, nil
}

func (ne *NodeExporter) Shutdown(_ context.Context) error {
	// There's nothing special needed to shutdown node-exporter.
	return nil
}
