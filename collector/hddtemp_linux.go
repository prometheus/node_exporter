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

// +build !nohddtemp

package collector

import (
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type hddtempCollector struct {
	temp *prometheus.GaugeVec
}

func init() {
	registerCollector("hddtemp", defaultDisabled, NewHddtempCollector)
}

func NewHddtempCollector() (Collector, error) {
	return &hddtempCollector{
		temp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: diskSubsystem,
				Name:      "temperature_celsius",
				Help:      "Disk temperature in celsius",
			},
			[]string{"device", "model"},
		),
	}, nil
}

func (c *hddtempCollector) Update(ch chan<- prometheus.Metric) (err error) {
	conn, err := net.DialTimeout("tcp", "localhost:7634", time.Second*3)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	fields := strings.Split(string(data), "|")
	for index := 0; index < len(fields)/5; index++ {
		offset := index * 5
		device := fields[offset+1]
		device = device[strings.LastIndex(device, "/")+1:]
		temperatureField := fields[offset+3]
		temperature, err := strconv.ParseFloat(temperatureField, 64)
		if err != nil {
			continue
		}
		c.temp.WithLabelValues(device, fields[offset+2]).Set(temperature)
	}
	c.temp.Collect(ch)
	return nil
}
