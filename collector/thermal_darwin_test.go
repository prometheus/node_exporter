// Copyright The Prometheus Authors
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

//go:build !notherm && darwin && cgo

package collector

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// Apple Silicon does not implement IOPMCopyCPUPowerStatus, so fetchCPUPowerStatus
// reports errNoCPUPowerStatus there. That is an expected condition and must not
// abort the collector, otherwise the temperature sensors, which are read after
// the CPU power status, are never collected.
func TestThermalUpdateWithoutCPUPowerStatus(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	c, err := NewThermCollector(logger)
	if err != nil {
		t.Fatalf("failed to create collector: %v", err)
	}

	ch := make(chan prometheus.Metric, 1024)
	err = c.Update(ch)
	close(ch)

	if errors.Is(err, errNoCPUPowerStatus) {
		t.Fatal("Update returned errNoCPUPowerStatus; a system without CPU power status must still collect temperatures")
	}
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	for range ch {
	}
}
