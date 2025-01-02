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

//go:build freebsd
// +build freebsd

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
	"testing"
	"unsafe"
)

func TestNetStatCollectorDescribe(t *testing.T) {
	ch := make(chan *prometheus.Desc, 1)
	collector := &netStatCollector{
		netStatMetric: prometheus.NewDesc("dummy_metric", "dummy", nil, nil),
	}
	collector.Describe(ch)
	desc := <-ch

	if want, got := "dummy_metric", desc.String(); want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestGetData(t *testing.T) {
	data, err := getData("net.inet.tcp.stats")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if got, want := len(data), int(unsafe.Sizeof(unix.TCPStats{})); got < want {
		t.Errorf("data length too small: want >= %d, got %d", want, got)
	}
}

func TestNetStatCollectorUpdate(t *testing.T) {
	ch := make(chan prometheus.Metric, len(metrics))
	collector := &netStatCollector{
		netStatMetric: prometheus.NewDesc("netstat_metric", "NetStat Metric", nil, nil),
	}
	err := collector.Update(ch)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if got, want := len(ch), len(metrics); got != want {
		t.Errorf("metric count mismatch: want %d, got %d", want, got)
	}

	for range metrics {
		<-ch
	}
}

func TestNewNetStatCollector(t *testing.T) {
	collector, err := NewNetStatCollector(nil)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if collector == nil {
		t.Fatal("collector is nil, want non-nil")
	}
}
