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

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/node_exporter/config"
)

func TestNewRuntimeCollectors(t *testing.T) {
	runtime, err := NewRuntime(config.NewConfigWithDefaults(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got, want := len(runtime.Collectors()), 1; got != want {
		t.Errorf("Expected: %d, Got: %d", want, got)
	}

	if got := len(runtime.EnabledCollectors()); got == 0 {
		t.Fatal("expected at least one enabled collector")
	}
}

func TestNewRuntimeValidateConfig(t *testing.T) {
	cfg := config.NewConfigWithDefaults()
	cfg.RuntimeGoMaxProcs = 0

	if _, err := NewRuntime(cfg, slog.New(slog.NewTextHandler(io.Discard, nil))); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRuntimeRegistry(t *testing.T) {
	runtime, err := NewRuntime(config.NewConfigWithDefaults(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	registry, err := runtime.Registry()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got := len(metrics); got == 0 {
		t.Fatal("expected gathered metrics, got none")
	}
}
