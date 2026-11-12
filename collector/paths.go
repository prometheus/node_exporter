// Copyright 2015 The Prometheus Authors
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
	"path/filepath"
	"strings"

	"github.com/prometheus/node_exporter/config"
)

// pathConfig adds file path helpers to config.PathConfig.
type pathConfig config.PathConfig

func (p pathConfig) procFilePath(name string) string {
	return filepath.Join(p.ProcFS, name)
}

func (p pathConfig) sysFilePath(name string) string {
	return filepath.Join(p.SysFS, name)
}

func (p pathConfig) rootfsFilePath(name string) string {
	return filepath.Join(p.RootFS, name)
}

func (p pathConfig) udevDataFilePath(name string) string {
	return filepath.Join(p.UdevData, name)
}

func (p pathConfig) rootfsStripPrefix(path string) string {
	if p.RootFS == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, p.RootFS)
	if stripped == "" {
		return "/"
	}
	return stripped
}
