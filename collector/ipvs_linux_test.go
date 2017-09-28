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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestIPVSCollector(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	collector, err := newIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Sprintf("failed to update collector: %v", err))
		}
	}()
	for expected, got := range map[string]string{
		prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String():                                                                                                                  (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String():                                                                                                             (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String():                                                                                                             (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String():                                                                                                                  (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String():                                                                                                                  (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String():     (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String(): (<-sink).Desc().String(),
		prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String():                     (<-sink).Desc().String(),
	} {
		if expected != got {
			t.Fatalf("Expected '%s' but got '%s'", expected, got)
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
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	collector, err := NewIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	prometheus.MustRegister(miniCollector{c: collector})

	rw := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(rw, &http.Request{})

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
