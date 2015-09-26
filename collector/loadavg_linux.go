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
	metric prometheus.Gauge
}

func init() {
	Factories["loadavg"] = NewLoadavgCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewLoadavgCollector() (Collector, error) {
	return &loadavgCollector{
		metric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "load1",
			Help:      "1m load average.",
		}),
	}, nil
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) (err error) {
	load, err := getLoad1()
	if err != nil {
		return fmt.Errorf("Couldn't get load: %s", err)
	}
	log.Debugf("Set node_load: %f", load)
	c.metric.Set(load)
	c.metric.Collect(ch)
	return err
}

func getLoad1() (float64, error) {
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return 0, err
	}
	return parseLoad(string(data))
}

func parseLoad(data string) (float64, error) {
	parts := strings.Fields(data)
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}
