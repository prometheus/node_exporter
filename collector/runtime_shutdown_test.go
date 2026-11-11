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

package collector

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// closeCounter is a collector that holds a resource.
type closeCounter struct {
	closed int
	err    error
}

func (c *closeCounter) Update(chan<- prometheus.Metric) error { return nil }

func (c *closeCounter) Close() error {
	c.closed++
	return c.err
}

// withCloser adds a closeCounter to the collectors that r owns.
func withCloser(r *Runtime, err error) *closeCounter {
	c := &closeCounter{err: err}
	r.state.collectors["closer"] = c
	return c
}

func TestRuntimeShutdownClosesCollectorsOnce(t *testing.T) {
	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	closer := withCloser(runtime, nil)
	filtered, err := runtime.Filtered("textfile")
	if err != nil {
		t.Fatal(err)
	}

	if err := filtered.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	if closer.closed != 0 {
		t.Fatalf("Shutdown of a filtered runtime closed its parent's collectors %d times", closer.closed)
	}

	for range 2 {
		if err := runtime.Shutdown(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if closer.closed != 1 {
		t.Fatalf("want the collector closed once after two Shutdown calls, got %d", closer.closed)
	}
}

func TestRuntimeShutdownReturnsCloseErrors(t *testing.T) {
	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	want := errors.New("socket busy")
	withCloser(runtime, want)
	if err := runtime.Shutdown(t.Context()); !errors.Is(err, want) {
		t.Fatalf("want %v, got %v", want, err)
	}
}
