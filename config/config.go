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

package config

// Config contains the options that control which collectors a runtime
// enables.
type Config struct {
	// CollectorDisableDefaults turns off every collector that Collectors
	// does not name, like --collector.disable-defaults.
	CollectorDisableDefaults bool
	// Collectors enables or disables collectors by name, like
	// --collector.<name> and --no-collector.<name>. It takes precedence
	// over the command-line flags.
	Collectors map[string]bool
}

// NewConfigWithDefaults returns a Config with the node_exporter defaults.
func NewConfigWithDefaults() Config {
	return Config{}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	return nil
}
