// Copyright 2017-2019 The Prometheus Authors
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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// TestInfiniBandCollectorDeviceExcludeSkipsUnparsableDevice verifies that a
// device excluded via --collector.infiniband.device-exclude is never parsed,
// so that a broken/unsupported attribute on that device (e.g. an unparsable
// "rate" file, as seen on Azure MANA RDMA devices) does not fail the whole
// collector. Before the fix, all devices were fully parsed before the filter
// was applied, so this update would fail even with the device excluded.
func TestInfiniBandCollectorDeviceExcludeSkipsUnparsableDevice(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	testSysPath := "fixtures/sys_infiniband"
	oldSysPath := *sysPath
	*sysPath = testSysPath
	defer func() { *sysPath = oldSysPath }()

	oldExclude := *infinibandDeviceExclude
	*infinibandDeviceExclude = "^mana_.*"
	defer func() { *infinibandDeviceExclude = oldExclude }()

	collector, err := NewInfiniBandCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 100)
	if err := collector.Update(ch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	close(ch)

	sawInfo := false
	for m := range ch {
		var pb dto.Metric
		if err := m.Write(&pb); err != nil {
			t.Fatalf("failed to write metric: %v", err)
		}
		for _, l := range pb.GetLabel() {
			if l.GetName() == "device" && l.GetValue() == "mana_0" {
				t.Fatalf("excluded device mana_0 was parsed and emitted a metric")
			}
			if l.GetName() == "device" && l.GetValue() == "mlx5_0" {
				sawInfo = true
			}
		}
	}
	if !sawInfo {
		t.Fatal("expected metrics for non-excluded device mlx5_0")
	}
}

// TestInfiniBandCollectorFailsWithoutExclude verifies that, without an
// exclusion filter, a device with an unparsable "rate" file still fails the
// collector as before -- the fix only changes filter ordering, not error
// handling for devices that are not excluded.
func TestInfiniBandCollectorFailsWithoutExclude(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	testSysPath := "fixtures/sys_infiniband"
	oldSysPath := *sysPath
	*sysPath = testSysPath
	defer func() { *sysPath = oldSysPath }()

	collector, err := NewInfiniBandCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 100)
	defer close(ch)
	if err := collector.Update(ch); err == nil {
		t.Fatal("expected an error parsing the unfiltered mana_0 device, got nil")
	}
}
