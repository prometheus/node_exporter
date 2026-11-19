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

//go:build !nokexec

package collector

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testKexecCollector struct {
	kc Collector
}

func (c testKexecCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.kc.Update(ch); err != nil {
		panic(err)
	}
}

func (c testKexecCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestKexecCollector(t *testing.T) {
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = "fixtures/sys"

	c, err := NewKexecCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	expected := `
# HELP node_kexec_crash_loaded Whether a crash (kdump) kernel is loaded, 1 is true, 0 is false.
# TYPE node_kexec_crash_loaded gauge
node_kexec_crash_loaded 1
# HELP node_kexec_crash_size_bytes Size in bytes of the memory reserved for the crash (kdump) kernel, 0 means no reservation.
# TYPE node_kexec_crash_size_bytes gauge
node_kexec_crash_size_bytes 2.68435456e+08
# HELP node_kexec_loaded Whether a kexec kernel is loaded and ready for a kexec reboot, 1 is true, 0 is false.
# TYPE node_kexec_loaded gauge
node_kexec_loaded 0
`

	if err := testutil.CollectAndCompare(testKexecCollector{c}, strings.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestKexecCollectorNoData(t *testing.T) {
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = "fixtures/sys/nonexistent"

	c, err := NewKexecCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 8)
	if err := c.Update(ch); !errors.Is(err, ErrNoData) {
		t.Fatalf("want ErrNoData, got %v", err)
	}
	if got := len(ch); got != 0 {
		t.Fatalf("want 0 metrics, got %d", got)
	}
}
