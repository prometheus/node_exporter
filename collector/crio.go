// Copyright 2020 The Prometheus Authors
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

// +build !nocrio

package collector

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	crioURL         = kingpin.Flag("collector.crio.url", "Metrics endpoint.").Default("http://localhost:9537/metrics").String()
	crioHTTPTimeout = kingpin.Flag("collector.crio.http_timeout", "http timeout.").Default("30s").Duration()
	client          *http.Client
)

type crioCollector struct {
	logger log.Logger
}

func init() {
	registerCollector("crio", defaultDisabled, NewCrioMetricsCollector)
}

// NewCrioCollector returns a new Collector exposing supervisord statistics.
func NewCrioMetricsCollector(logger log.Logger) (Collector, error) {
	client = &http.Client{
		Timeout: *crioHTTPTimeout,
	}
	return &crioCollector{logger: logger}, nil
}

func (c *crioCollector) Update(ch chan<- prometheus.Metric) error {
	req, err := http.NewRequest("GET", *crioURL, nil)
	if err != nil {
		return fmt.Errorf("unable to make http request: %s", err)
	}

	req.Header.Set("User-Agent", "node_exporter/crawler")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("crio http request returned StatusCode: %v", resp.StatusCode)
	}

	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return fmt.Errorf("reading text failed: %v", err)
	}
	for _, mf := range metricFamilies {
		convertMetricFamily(mf, ch, c.logger)
	}
	return nil
}
