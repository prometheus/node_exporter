// Copyright 2024 The Prometheus Authors
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

//go:build !nopcidevice

package collector

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

func TestPCICollectorWithNameResolution(t *testing.T) {
	// Test the PCI collector with name resolution enabled and compare against expected output
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := fixtureConfig()
	cfg.PCIDevice = config.PCIDeviceConfig{Names: true, IDsFile: "fixtures/pci.ids"}
	c, err := NewPcideviceCollector(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testPCICollector{pc: c})

	// Read expected output from fixture file
	expectedOutput, err := os.ReadFile("fixtures/pcidevice-names-output.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = testutil.GatherAndCompare(reg, strings.NewReader(string(expectedOutput)))
	if err != nil {
		t.Fatal(err)
	}
}

// testPCICollector wraps the PCI collector for testing
type testPCICollector struct {
	pc Collector
}

func (tc *testPCICollector) Collect(ch chan<- prometheus.Metric) {
	sink := make(chan prometheus.Metric)
	go func() {
		err := tc.pc.Update(context.Background(), sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for m := range sink {
		ch <- m
	}
}

func (tc *testPCICollector) Describe(_ chan<- *prometheus.Desc) {
	// No-op for testing
}
