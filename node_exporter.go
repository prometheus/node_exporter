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

package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/node_exporter/collector"
)

const (
	defaultCollectors = "conntrack,cpu,diskstats,entropy,filefd,filesystem,loadavg,mdadm,meminfo,netdev,netstat,sockstat,stat,textfile,time,uname,version,vmstat"
)

var (
	// Version of node_exporter. Set at build time.
	Version = "0.0.0.dev"

	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "node_exporter: Duration of a scrape job.",
		},
		[]string{"collector", "result"},
	)
)

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	collectors map[string]collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.collectors))
	for name, c := range n.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func filterAvailableCollectors(collectors string) string {
	availableCollectors := make([]string, 0)
	for _, c := range strings.Split(collectors, ",") {
		_, ok := collector.Factories[c]
		if ok {
			availableCollectors = append(availableCollectors, c)
		}
	}
	return strings.Join(availableCollectors, ",")
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func loadCollectors(list string) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	for _, name := range strings.Split(list, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn()
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func main() {
	var (
		listenAddress     = flag.String("web.listen-address", ":9100", "Address on which to expose metrics and web interface.")
		metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		enabledCollectors = flag.String("collectors.enabled", filterAvailableCollectors(defaultCollectors), "Comma-separated list of collectors to use.")
		printCollectors   = flag.Bool("collectors.print", false, "If true, print available collectors and exit.")
	)
	flag.Parse()

	if *printCollectors {
		collectorNames := make(sort.StringSlice, 0, len(collector.Factories))
		for n := range collector.Factories {
			collectorNames = append(collectorNames, n)
		}
		collectorNames.Sort()
		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}
		return
	}
	collectors, err := loadCollectors(*enabledCollectors)
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}

	log.Infof("Enabled collectors:")
	for n := range collectors {
		log.Infof(" - %s", n)
	}

	nodeCollector := NodeCollector{collectors: collectors}
	prometheus.MustRegister(nodeCollector)

	handler := prometheus.Handler()

	http.Handle(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infof("Starting node_exporter v%s at %s", Version, *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
