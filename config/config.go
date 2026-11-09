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

// Config contains every option of a node_exporter runtime. Each field has a
// command-line flag of the same meaning, named in its comment.
type Config struct {
	// CollectorDisableDefaults turns off every collector that Collectors
	// does not name, like --collector.disable-defaults.
	CollectorDisableDefaults bool
	// Collectors enables or disables collectors by name, like
	// --collector.<name> and --no-collector.<name>. It takes precedence
	// over the command-line flags.
	Collectors map[string]bool

	Path PathConfig

	CPU      CPUConfig
	Textfile TextfileConfig
}

const (
	DefaultPathProcFS   = "/proc"
	DefaultPathSysFS    = "/sys"
	DefaultPathRootFS   = "/"
	DefaultPathUdevData = "/run/udev/data"
)

// PathConfig contains the mount points the collectors read from.
type PathConfig struct {
	// ProcFS is the procfs mount point, like --path.procfs.
	ProcFS string
	// SysFS is the sysfs mount point, like --path.sysfs.
	SysFS string
	// RootFS is the rootfs mount point, like --path.rootfs.
	RootFS string
	// UdevData is the udev data path, like --path.udev.data.
	UdevData string
}

// NewConfigWithDefaults returns a Config with the node_exporter defaults.
func NewConfigWithDefaults() Config {
	return Config{
		Path: PathConfig{
			ProcFS:   DefaultPathProcFS,
			SysFS:    DefaultPathSysFS,
			RootFS:   DefaultPathRootFS,
			UdevData: DefaultPathUdevData,
		},
		CPU: CPUConfig{
			EnableGuest: DefaultCPUEnableGuest,
		},
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
