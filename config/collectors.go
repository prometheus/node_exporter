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

const (
	DefaultCPUEnableGuest = true
)

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
