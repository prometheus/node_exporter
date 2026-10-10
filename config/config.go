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

package config

import "fmt"

const (
	DefaultWebTelemetryPath  = "/metrics"
	DefaultWebMaxRequests    = 40
	DefaultRuntimeGoMaxProcs = 1
)

// Config contains the user-facing configuration currently adapted by the
// node_exporter binary into the reusable runtime and HTTP handler layers.
type Config struct {
	WebTelemetryPath          string
	WebDisableExporterMetrics bool
	WebMaxRequests            int
	CollectorDisableDefaults  bool
	RuntimeGoMaxProcs         int
	EnabledCollectors         []string
}

func NewConfigWithDefaults() Config {
	return Config{
		WebTelemetryPath:  DefaultWebTelemetryPath,
		WebMaxRequests:    DefaultWebMaxRequests,
		RuntimeGoMaxProcs: DefaultRuntimeGoMaxProcs,
	}
}

func (c Config) Validate() error {
	if c.WebTelemetryPath == "" {
		return fmt.Errorf("web telemetry path must not be empty")
	}
	if c.WebMaxRequests < 0 {
		return fmt.Errorf("web max requests must be greater than or equal to zero")
	}
	if c.RuntimeGoMaxProcs <= 0 {
		return fmt.Errorf("runtime gomaxprocs must be greater than zero")
	}
	return nil
}
