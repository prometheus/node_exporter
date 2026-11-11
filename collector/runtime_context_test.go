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
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// ctxRecorder records the context error seen by its last Update.
type ctxRecorder struct{ err error }

func (c *ctxRecorder) Update(ctx context.Context, _ chan<- prometheus.Metric) error {
	c.err = ctx.Err()
	return nil
}

// recordContext replaces the collectors of r with a ctxRecorder, gathers r
// once and returns the context error the collector saw.
func recordContext(t *testing.T, r *Runtime) error {
	t.Helper()
	rec := &ctxRecorder{}
	r.collector.Collectors = map[string]Collector{"recorder": rec}
	testutil.CollectAndCount(r.Collectors()[0])
	return rec.err
}

func TestRuntimeShutdownCancelsCollectorContext(t *testing.T) {
	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	filtered, err := runtime.Filtered("textfile")
	if err != nil {
		t.Fatal(err)
	}

	if err := filtered.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := recordContext(t, runtime); err != nil {
		t.Fatalf("Shutdown of a filtered runtime must not cancel its parent, got %v", err)
	}

	if err := runtime.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := recordContext(t, runtime); !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled after Shutdown, got %v", err)
	}
	if err := recordContext(t, filtered); !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled in the filtered runtime after Shutdown, got %v", err)
	}
}

func TestNewRuntimeWithContextFollowsParent(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	runtime, err := NewRuntimeWithContext(ctx, textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := recordContext(t, runtime); !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled after the parent is canceled, got %v", err)
	}
}

func TestNewRuntimeWithContextRequiresContext(t *testing.T) {
	var ctx context.Context
	if _, err := NewRuntimeWithContext(ctx, textfileOnlyConfig(), testLogger()); err == nil {
		t.Fatal("want error for a nil context")
	}
}
