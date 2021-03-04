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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestIPVSCollector(t *testing.T) {
	testcases := []struct {
		labels  string
		expects []string
		err     error
	}{
		{
			"<none>",
			[]string{
				prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto", "local_mark"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto", "local_mark"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto", "local_mark"}, nil).String(),
			},
			nil,
		},
		{
			"",
			[]string{
				prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", nil, nil).String(),
			},
			nil,
		},
		{
			"local_port",
			[]string{
				prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", []string{"local_port"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", []string{"local_port"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", []string{"local_port"}, nil).String(),
			},
			nil,
		},
		{
			"local_address,local_port",
			[]string{
				prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", []string{"local_address", "local_port"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", []string{"local_address", "local_port"}, nil).String(),
				prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", []string{"local_address", "local_port"}, nil).String(),
			},
			nil,
		},
		{
			"invalid_label",
			nil,
			errors.New(`unknown IPVS backend labels: "invalid_label"`),
		},
		{
			"invalid_label,bad_label",
			nil,
			errors.New(`unknown IPVS backend labels: "bad_label, invalid_label"`),
		},
	}
	for _, test := range testcases {
		t.Run(test.labels, func(t *testing.T) {
			args := []string{"--path.procfs", "fixtures/proc"}
			if test.labels != "<none>" {
				args = append(args, "--collector.ipvs.backend-labels="+test.labels)
			}
			if _, err := kingpin.CommandLine.Parse(args); err != nil {
				t.Fatal(err)
			}
			collector, err := newIPVSCollector(log.NewNopLogger())
			if err != nil {
				if test.err == nil {
					t.Fatal(err)
				}
				if !strings.Contains(err.Error(), test.err.Error()) {
					t.Fatalf("expect error: %v contains %v", err, test.err)
				}
				return
			}
			if test.err != nil {
				t.Fatalf("expect error: %v but got no error", test.err)
			}

			sink := make(chan prometheus.Metric)
			go func() {
				err = collector.Update(sink)
				if err != nil {
					panic(fmt.Sprintf("failed to update collector: %v", err))
				}
			}()
			for _, expected := range test.expects {
				got := (<-sink).Desc().String()
				if expected != got {
					t.Fatalf("Expected '%s' but got '%s'", expected, got)
				}
			}
		})
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
	testcases := []struct {
		labels      string
		metricsFile string
	}{
		{"<none>", "fixtures/ip_vs_result.txt"},
		{"", "fixtures/ip_vs_result_lbs_none.txt"},
		{"local_port", "fixtures/ip_vs_result_lbs_local_port.txt"},
		{"local_address,local_port", "fixtures/ip_vs_result_lbs_local_address_local_port.txt"},
	}
	for _, test := range testcases {
		t.Run(test.labels, func(t *testing.T) {
			args := []string{"--path.procfs", "fixtures/proc"}
			if test.labels != "<none>" {
				args = append(args, "--collector.ipvs.backend-labels="+test.labels)
			}
			if _, err := kingpin.CommandLine.Parse(args); err != nil {
				t.Fatal(err)
			}
			collector, err := NewIPVSCollector(log.NewNopLogger())
			if err != nil {
				t.Fatal(err)
			}
			registry := prometheus.NewRegistry()
			registry.MustRegister(miniCollector{c: collector})

			rw := httptest.NewRecorder()
			promhttp.InstrumentMetricHandler(registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})).ServeHTTP(rw, &http.Request{})

			wantMetrics, err := ioutil.ReadFile(test.metricsFile)
			if err != nil {
				t.Fatalf("unable to read input test file %s: %s", test.metricsFile, err)
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
		})
	}
}
