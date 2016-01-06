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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Version  string
	Revision string
	Branch   string
)

type versionCollector struct {
	metric *prometheus.GaugeVec
}

func init() {
	Factories["version"] = NewVersionCollector
}

func NewVersionCollector() (Collector, error) {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_exporter_build_info",
			Help: "A metric with a constant '1' value labeled by version, revision, and branch from which the node_exporter was built.",
		},
		[]string{"version", "revision", "branch"},
	)
	metric.WithLabelValues(Version, Revision, Branch).Set(1)
	return &versionCollector{
		metric: metric,
	}, nil
}

func (c *versionCollector) Update(ch chan<- prometheus.Metric) (err error) {
	c.metric.Collect(ch)
	return err
}
