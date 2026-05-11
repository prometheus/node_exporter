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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	
	"github.com/prometheus/exporter-toolkit/otlpreceiver"
)

func NewFactory() receiver.Factory {
	defaults := map[string]interface{}{
		"disable_defaults":   false,
		"enable_collectors":  []string{},
		"exclude_collectors": []string{},
	}

	return otlpreceiver.NewFactory(
		otlpreceiver.WithType(component.MustNewType("prometheus_node_exporter")),
		otlpreceiver.WithInitializer(NewNodeExporter(&Config{})),
		otlpreceiver.WithConfigUnmarshaler(&ConfigUnmarshaler{}),
		otlpreceiver.WithComponentDefaults(defaults),
	)
}