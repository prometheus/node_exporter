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

// +build !noloadavg

package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

type loadavgCollector struct {
	metric []prometheus.Gauge
}

func init() {
	Factories["loadavg"] = NewLoadavgCollector
}

// Take a prometheus registry and return a new Collector exposing load average.
func NewLoadavgCollector() (Collector, error) {
	return &loadavgCollector{
		metric: []prometheus.Gauge{
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "load1",
				Help:      "1m load average.",
			}),
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "load5",
				Help:      "5m load average.",
			}),
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "load15",
				Help:      "15m load average.",
			}),
		},
	}, nil
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) (err error) {
	loads, err := getLoad()
	if err != nil {
		return fmt.Errorf("couldn't get load: %s", err)
	}
	for i, load := range loads {
		log.Debugf("Set load %d: %f", i, load)
		c.metric[i].Set(load)
		c.metric[i].Collect(ch)
	}
	return err
}

// Read loadavg from /proc.
func getLoad() (loads []float64, err error) {
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return nil, err
	}
	loads, err = parseLoad(string(data))
	if err != nil {
		return nil, err
	}
	return loads, nil
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func parseLoad(data string) (loads []float64, err error) {
	loads = make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %s", procFilePath("loadavg"))
	}
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%s': %s", load, err)
		}
	}
	return loads, nil
}
