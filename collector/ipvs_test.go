// Copyright 2015 The Prometheus Authors
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
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	expectedIPVSStats = procfs.IPVSStats{
		Connections:     23765872,
		IncomingPackets: 3811989221,
		OutgoingPackets: 0,
		IncomingBytes:   89991519156915,
		OutgoingBytes:   0,
	}
	expectedIPVSBackendStatuses = []procfs.IPVSBackendStatus{
		{
			LocalAddress:  net.ParseIP("192.168.0.22"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.82.22"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    248,
			InactConn:     2,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.22"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.83.24"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    248,
			InactConn:     2,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.22"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.83.21"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    248,
			InactConn:     1,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.57"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.84.22"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        0,
			ActiveConn:    0,
			InactConn:     0,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.57"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.82.21"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    1499,
			InactConn:     0,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.57"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.50.21"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    1498,
			InactConn:     0,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.55"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.50.26"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        0,
			ActiveConn:    0,
			InactConn:     0,
		},
		{
			LocalAddress:  net.ParseIP("192.168.0.55"),
			LocalPort:     3306,
			RemoteAddress: net.ParseIP("192.168.49.32"),
			RemotePort:    3306,
			Proto:         "TCP",
			Weight:        100,
			ActiveConn:    0,
			InactConn:     0,
		},
	}
)

func TestIPVSCollector(t *testing.T) {
	if err := flag.Set("collector.procfs", "fixtures/proc"); err != nil {
		t.Fatal(err)
	}
	collector, err := newIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	sink := make(chan prometheus.Metric)
	go func() {
		for {
			<-sink
		}
	}()

	err = collector.Update(sink)
	if err != nil {
		t.Fatal(err)
	}

	for _, expect := range expectedIPVSBackendStatuses {
		labels := prometheus.Labels{
			"local_address":  expect.LocalAddress.String(),
			"local_port":     strconv.FormatUint(uint64(expect.LocalPort), 10),
			"remote_address": expect.RemoteAddress.String(),
			"remote_port":    strconv.FormatUint(uint64(expect.RemotePort), 10),
			"proto":          expect.Proto,
		}
		// TODO: Pending prometheus/client_golang#58, check the actual numbers
		_, err = collector.backendConnectionsActive.GetMetricWith(labels)
		if err != nil {
			t.Errorf("Missing active connections metric for label combination: %+v", labels)
		}
		_, err = collector.backendConnectionsInact.GetMetricWith(labels)
		if err != nil {
			t.Errorf("Missing inactive connections metric for label combination: %+v", labels)
		}
		_, err = collector.backendWeight.GetMetricWith(labels)
		if err != nil {
			t.Errorf("Missing weight metric for label combination: %+v", labels)
		}
	}
}

// mock collector
type miniCollector struct {
	c Collector
}

func (c miniCollector) Collect(ch chan<- prometheus.Metric) {
	c.c.Update(ch)
}

func (c miniCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "fake",
		Subsystem: "fake",
		Name:      "fake",
		Help:      "fake",
	}).Describe(ch)
}

func TestIPVSCollectorResponse(t *testing.T) {
	if err := flag.Set("collector.procfs", "fixtures/proc"); err != nil {
		t.Fatal(err)
	}
	collector, err := NewIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	prometheus.MustRegister(miniCollector{c: collector})

	rw := httptest.NewRecorder()
	prometheus.Handler().ServeHTTP(rw, &http.Request{})

	metricsFile := "fixtures/ip_vs_result.txt"
	wantMetrics, err := ioutil.ReadFile(metricsFile)
	if err != nil {
		t.Fatalf("unable to read input test file %s: %s", metricsFile, err)
	}

	wantLines := strings.Split(string(wantMetrics), "\n")
	gotLines := strings.Split(string(rw.Body.String()), "\n")
	gotLinesIdx := 0

	// Until the Prometheus Go client library offers better testability
	// (https://github.com/prometheus/client_golang/issues/58), we simply compare
	// verbatim text-format metrics outputs, but ignore any lines we don't have
	// in the fixture. Put differently, we are only testing that each line from
	// the fixture is present, in the order given.
wantLoop:
	for _, want := range wantLines {
		for _, got := range gotLines[gotLinesIdx:] {
			if want == got {
				// this is a line we are interested in, and it is correct
				continue wantLoop
			} else {
				gotLinesIdx++
			}
		}
		// if this point is reached, the line we want was missing
		t.Fatalf("Missing expected output line(s), first missing line is %s", want)
	}
}
