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
)

func (p *PathConfig) procFilePath(name string) string {
	return filepath.Join(*p.ProcPath, name)
}

func (p *PathConfig) sysFilePath(name string) string {
	return filepath.Join(*p.SysPath, name)
}

func (p *PathConfig) rootfsFilePath(name string) string {
	return filepath.Join(*p.RootfsPath, name)
}

func (p *PathConfig) udevDataFilePath(name string) string {
	return filepath.Join(*p.UdevDataPath, name)
}

func (p *PathConfig) rootfsStripPrefix(path string) string {
	if *p.RootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, *p.RootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
