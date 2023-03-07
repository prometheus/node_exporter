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

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"
)

func TestDefaultProcPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", procfs.DefaultMountPoint}); err != nil {
		t.Fatal(err)
	}

	if got, want := procFilePath("somefile"), "/proc/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := procFilePath("some/file"), "/proc/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomProcPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}

	if got, want := procFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := procFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestDefaultSysPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "/sys"}); err != nil {
		t.Fatal(err)
	}

	if got, want := sysFilePath("somefile"), "/sys/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := sysFilePath("some/file"), "/sys/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomSysPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}

	if got, want := sysFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := sysFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}
