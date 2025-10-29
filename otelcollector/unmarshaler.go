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

import "github.com/prometheus/exporter-toolkit/otlpreceiver"

type ConfigUnmarshaler struct{}

func (u *ConfigUnmarshaler) UnmarshalExporterConfig(data map[string]interface{}) (otlpreceiver.Config, error) {
	cfg := &Config{}

	if disableDefaults, ok := data["disable_defaults"].(bool); ok {
		cfg.DisableDefaults = disableDefaults
	}

	if enabledCollector, ok := data["enable_collectors"].([]string); ok {
		cfg.EnableCollectors = enabledCollector
	}

	if excludedCollectors, ok := data["exclude_collectors"].([]string); ok {
		cfg.ExcludeCollectors = excludedCollectors
	}
	return cfg, nil
}
