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

// Exporter is a prometheus exporter using multiple Factories to collect and export system metrics.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "node"

var Factories = make(map[string]func() (Collector, error))

// Interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) (err error)
}

// TODO: Instead of periodically call Update, a Collector could be implemented
// as a real prometheus.Collector that only gathers metrics when
// scraped. (However, for metric gathering that takes very long, it might
// actually be better to do them proactively before scraping to minimize scrape
// time.)
