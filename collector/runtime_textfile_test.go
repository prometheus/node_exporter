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

//go:build !notextfile

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/node_exporter/config"
)

func TestRuntimeSnapshotsTextfileCollectorConfig(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	originalDirectories := append([]string(nil), (*textFileDirectories)...)
	t.Cleanup(func() {
		*textFileDirectories = originalDirectories
	})

	*textFileDirectories = []string{"fixtures/textfile/first"}
	cfg := config.NewConfigWithDefaults()
	cfg.EnabledCollectors = []string{"textfile"}

	runtime, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	*textFileDirectories = []string{"fixtures/textfile/second"}

	filtered, err := runtime.Filtered("textfile")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	collector := filtered.collector.Collectors["textfile"].(*textFileCollector)
	if got, want := collector.paths, []string{"fixtures/textfile/first"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("expected filtered runtime to keep initial textfile paths, got %v", got)
	}

	runtime2, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	collector2 := runtime2.collector.Collectors["textfile"].(*textFileCollector)
	if got, want := collector2.paths, []string{"fixtures/textfile/second"}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("expected new runtime to use updated textfile paths, got %v", got)
	}
}
