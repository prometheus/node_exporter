// Copyright 2021 The Prometheus Authors
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

//go:build !nofibrechannel

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestFibreChannelCollector(t *testing.T) {
	// The fixtures directory includes host2 which has
	// port_state="Online" but an empty statistics/ dir.
	// Before the fix, this would panic with a nil pointer
	// dereference when trying to read counter values.

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Override sysPath to point to our fixtures
	testSysPath := "fixtures/sys"
	oldSysPath := *sysPath
	*sysPath = testSysPath
	defer func() { *sysPath = oldSysPath }()

	collector, err := NewFibreChannelCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	// This should not panic even with missing statistics files
	ch := make(chan prometheus.Metric, 100)
	err = collector.Update(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
