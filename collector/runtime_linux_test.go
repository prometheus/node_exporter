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

//go:build linux && !nocpu

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/node_exporter/config"
)

func TestRuntimeSnapshotsCPUCollectorConfig(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	originalGuestEnabled := *enableCPUGuest
	originalInfoEnabled := *enableCPUInfo
	originalFlagsInclude := *flagsInclude
	originalBugsInclude := *bugsInclude
	t.Cleanup(func() {
		*enableCPUGuest = originalGuestEnabled
		*enableCPUInfo = originalInfoEnabled
		*flagsInclude = originalFlagsInclude
		*bugsInclude = originalBugsInclude
	})

	*enableCPUGuest = false
	*enableCPUInfo = false
	*flagsInclude = "foo"
	*bugsInclude = ""

	cfg := config.NewConfigWithDefaults()
	cfg.EnabledCollectors = []string{"cpu"}

	runtime, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	*enableCPUGuest = true
	*enableCPUInfo = true
	*flagsInclude = "bar"

	filtered, err := runtime.Filtered("cpu")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	collector := filtered.collector.Collectors["cpu"].(*cpuCollector)
	if collector.guestEnabled {
		t.Fatal("expected filtered runtime to keep initial cpu guest setting")
	}
	if !collector.infoEnabled {
		t.Fatal("expected flags include to enable cpu info in filtered runtime")
	}
	if collector.cpuFlagsIncludeRegexp == nil || !collector.cpuFlagsIncludeRegexp.MatchString("foo") || collector.cpuFlagsIncludeRegexp.MatchString("bar") {
		t.Fatalf("expected filtered runtime to keep initial cpu flags include regexp, got %v", collector.cpuFlagsIncludeRegexp)
	}

	runtime2, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	collector2 := runtime2.collector.Collectors["cpu"].(*cpuCollector)
	if !collector2.guestEnabled {
		t.Fatal("expected new runtime to use updated cpu guest setting")
	}
	if collector2.cpuFlagsIncludeRegexp == nil || !collector2.cpuFlagsIncludeRegexp.MatchString("bar") || collector2.cpuFlagsIncludeRegexp.MatchString("foo") {
		t.Fatalf("expected new runtime to use updated cpu flags include regexp, got %v", collector2.cpuFlagsIncludeRegexp)
	}
}
