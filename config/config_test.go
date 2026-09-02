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

package config

import "testing"

func TestNewConfigWithDefaults(t *testing.T) {
	cfg := NewConfigWithDefaults()

	if got, want := cfg.WebTelemetryPath, DefaultWebTelemetryPath; got != want {
		t.Errorf("Expected: %q, Got: %q", want, got)
	}
	if got, want := cfg.WebMaxRequests, DefaultWebMaxRequests; got != want {
		t.Errorf("Expected: %d, Got: %d", want, got)
	}
	if got, want := cfg.RuntimeGoMaxProcs, DefaultRuntimeGoMaxProcs; got != want {
		t.Errorf("Expected: %d, Got: %d", want, got)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewConfigWithDefaults()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestConfigValidateRequiresTelemetryPath(t *testing.T) {
	cfg := NewConfigWithDefaults()
	cfg.WebTelemetryPath = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfigValidateRejectsNegativeMaxRequests(t *testing.T) {
	cfg := NewConfigWithDefaults()
	cfg.WebMaxRequests = -1

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfigValidateRejectsNonPositiveGoMaxProcs(t *testing.T) {
	cfg := NewConfigWithDefaults()
	cfg.RuntimeGoMaxProcs = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}
