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

//go:build !noinfiniband

package collector

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInfiniBandDeviceFilterBeforeRead(t *testing.T) {
	sys := t.TempDir()
	included := filepath.Join(sys, "class", "infiniband", "mlx5_7")
	excluded := filepath.Join(sys, "class", "infiniband", "mlx5_5")

	if err := os.MkdirAll(filepath.Join(included, "ports"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(included, "fw_ver"), []byte("28.47.1088\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The excluded device deliberately lacks fw_ver and ports. Reading it would fail.
	if err := os.MkdirAll(excluded, 0o755); err != nil {
		t.Fatal(err)
	}

	oldSysPath := *sysPath
	oldInclude := *infinibandDeviceInclude
	oldExclude := *infinibandDeviceExclude
	t.Cleanup(func() {
		*sysPath = oldSysPath
		*infinibandDeviceInclude = oldInclude
		*infinibandDeviceExclude = oldExclude
	})
	*sysPath = sys
	*infinibandDeviceInclude = ""
	*infinibandDeviceExclude = "^mlx5_5$"

	collector, err := NewInfiniBandCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 1)
	if err := collector.Update(ch); err != nil {
		t.Fatalf("excluded device was read: %v", err)
	}
	if len(ch) != 1 {
		t.Fatalf("expected one metric from the included device, got %d", len(ch))
	}
}
