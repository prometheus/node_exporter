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

// Package collector includes all individual collectors to gather and export system metrics.
package collector

import (
	"testing"

	"gopkg.in/alecthomas/kingpin.v2"
)

func TestDisableDefaults(t *testing.T) {
	_, err := kingpin.CommandLine.Parse([]string{"--collectors.disable-defaults"})
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewNodeCollector()
	if err != nil {
		t.Error(err)
	}
	if len(c.Collectors) != 0 {
		t.Errorf("Expected no collectors, got %d", len(c.Collectors))
	}
}

func TestDisableDefaultsEnableTextfile(t *testing.T) {
	_, err := kingpin.CommandLine.Parse([]string{"--collectors.disable-defaults", "--collector.textfile.enabled"})
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewNodeCollector()
	if err != nil {
		t.Error(err)
	}
	if len(c.Collectors) != 1 {
		t.Errorf("Expected one collector, got %d", len(c.Collectors))
	}
	if _, ok := c.Collectors["textfile"]; !ok {
		t.Error("Expected textfile collector to be enabled, but it is not")
	}
}
