// Copyright 2019 The Prometheus Authors
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

package perconacollector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	cl "github.com/prometheus/node_exporter/collector"
)

func init() {
	fmt.Println("init standard")
	cl.RegisterCollectorPublic("standard.go", false, NewStandardGoCollector)
	cl.RegisterCollectorPublic("standard.process", false, NewStandardProcessCollector)
}

type standardGoCollector struct {
	origin prometheus.Collector
	logger log.Logger
}

// NewStandardGoCollector creates standard go collector.
func NewStandardGoCollector(logger log.Logger) (cl.Collector, error) {
	c := collectors.NewGoCollector()
	return &standardGoCollector{origin: c}, nil
}

func (c *standardGoCollector) Update(ch chan<- prometheus.Metric) error {
	c.origin.Collect(ch)
	return nil
}

type standardProcessCollector struct {
	origin prometheus.Collector
}

// NewStandardProcessCollector creates standard process collector.
func NewStandardProcessCollector(logger log.Logger) (cl.Collector, error) {
	c := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	return &standardProcessCollector{origin: c}, nil
}

func (c *standardProcessCollector) Update(ch chan<- prometheus.Metric) error {
	c.origin.Collect(ch)
	return nil
}
