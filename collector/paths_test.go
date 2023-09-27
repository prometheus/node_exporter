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
	"testing"

	"github.com/prometheus/procfs"
)

func TestDefaultProcPath(t *testing.T) {
	config := newNodeCollectorWithPaths()
	path := procfs.DefaultMountPoint
	config.Path.ProcPath = &path

	if got, want := config.Path.procFilePath("somefile"), "/proc/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := config.Path.procFilePath("some/file"), "/proc/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomProcPath(t *testing.T) {
	config := newNodeCollectorWithPaths()
	path := "./../some/./place/"
	config.Path.ProcPath = &path

	if got, want := config.Path.procFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := config.Path.procFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestDefaultSysPath(t *testing.T) {
	config := newNodeCollectorWithPaths()
	path := "/sys"
	config.Path.SysPath = &path

	if got, want := config.Path.sysFilePath("somefile"), "/sys/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := config.Path.sysFilePath("some/file"), "/sys/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomSysPath(t *testing.T) {
	config := newNodeCollectorWithPaths()
	path := "./../some/./place/"
	config.Path.SysPath = &path

	if got, want := config.Path.sysFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := config.Path.sysFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func newNodeCollectorWithPaths() *NodeCollectorConfig {
	return &NodeCollectorConfig{Path: PathConfig{
		ProcPath:     new(string),
		SysPath:      new(string),
		RootfsPath:   new(string),
		UdevDataPath: new(string),
	}}
}
