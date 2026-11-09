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

import (
	"fmt"
	"regexp"
)

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

	CPU      CPUConfig
	Textfile TextfileConfig
}

// CPUConfig contains the options for the cpu collector.
type CPUConfig struct {
	// EnableGuest exposes node_cpu_guest_seconds_total, like
	// --collector.cpu.guest.
	EnableGuest bool
	// EnableInfo exposes node_cpu_info, like --collector.cpu.info.
	EnableInfo bool
	// FlagsInclude and BugsInclude are regular expressions that filter the
	// flags and bugs fields of node_cpu_info, like
	// --collector.cpu.info.flags-include and --collector.cpu.info.bugs-include.
	// Setting either one also enables node_cpu_info.
	FlagsInclude string
	BugsInclude  string
}

// TextfileConfig contains the options for the textfile collector.
type TextfileConfig struct {
	// Directories to read text files with metrics from, like
	// --collector.textfile.directory. Glob patterns are supported.
	Directories []string
}

// NewConfigWithDefaults returns a Config with the node_exporter defaults.
func NewConfigWithDefaults() Config {
	return Config{
		CPU: CPUConfig{EnableGuest: true},
	}
}

// Validate returns an error if the Config is invalid.
func (c Config) Validate() error {
	if _, err := regexp.Compile(c.CPU.FlagsInclude); err != nil {
		return fmt.Errorf("invalid CPU.FlagsInclude: %w", err)
	}
	if _, err := regexp.Compile(c.CPU.BugsInclude); err != nil {
		return fmt.Errorf("invalid CPU.BugsInclude: %w", err)
	}
	return nil
}
