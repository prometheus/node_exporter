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

//go:build !nousb

package collector

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestUSBCollector(t *testing.T) {
	*sysPath = "fixtures/sys"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewUSBCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testUSBCollector{c: c})

	expectedOutput, err := os.ReadFile("fixtures/usb-output.txt")
	if err != nil {
		t.Fatal(err)
	}

	if err := testutil.GatherAndCompare(reg, strings.NewReader(string(expectedOutput))); err != nil {
		t.Fatal(err)
	}
}

func TestUSBCollectorNoDevices(t *testing.T) {
	// A checked-in fixture cannot express an absent /sys/bus/usb/devices
	// directory, so an empty temporary root is the smallest way to exercise
	// the ErrNoData path.
	*sysPath = t.TempDir()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewUSBCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 200)
	err = c.Update(ch)
	close(ch)

	if !errors.Is(err, ErrNoData) {
		t.Fatalf("expected ErrNoData, got %v", err)
	}
}

// testUSBCollector wraps the USB collector for testing.
type testUSBCollector struct {
	c Collector
}

func (tc *testUSBCollector) Collect(ch chan<- prometheus.Metric) {
	sink := make(chan prometheus.Metric)
	go func() {
		if err := tc.c.Update(sink); err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for m := range sink {
		ch <- m
	}
}

func (tc *testUSBCollector) Describe(_ chan<- *prometheus.Desc) {
	// No-op for testing.
}
