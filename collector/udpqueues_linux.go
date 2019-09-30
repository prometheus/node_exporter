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

// +build !noudp_queues

package collector

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type udpQueuesCollector struct {
	desc typedDesc
}

func init() {
	registerCollector("udp_queues", defaultDisabled, NewUDPqueuesCollector)
}

// NewUDPqueuesCollector returns a new Collector exposing network udp queued bytes.
func NewUDPqueuesCollector() (Collector, error) {
	return &udpQueuesCollector{
		desc: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "udp", "queues"),
			"Number of allocated memory in the kernel for UDP datagrams in bytes.",
			[]string{"queue"}, nil,
		), prometheus.GaugeValue},
	}, nil
}

func (c *udpQueuesCollector) Update(ch chan<- prometheus.Metric) error {
	updQueues, err := getUDPqueues(procFilePath("net/udp"))
	if err != nil {
		return fmt.Errorf("couldn't get upd queued bytes: %s", err)
	}

	// if enabled ipv6 system
	udp6File := procFilePath("net/udp6")
	if _, hasIPv6 := os.Stat(udp6File); hasIPv6 == nil {
		udp6Queues, err := getUDPqueues(udp6File)
		if err != nil {
			return fmt.Errorf("couldn't get udp6 queued bytes: %s", err)
		}

		for qu, value := range udp6Queues {
			updQueues[qu] += value
		}
	}

	for qu, value := range updQueues {
		ch <- c.desc.mustNewConstMetric(value, qu)
	}
	return nil
}

func getUDPqueues(statsFile string) (map[string]float64, error) {
	file, err := os.Open(statsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseUDPqueues(file)
}

func parseUDPqueues(r io.Reader) (map[string]float64, error) {
	updQueues := map[string]float64{}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(contents), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) < 5 {
			return nil, fmt.Errorf("invalid line in file: %q", line)
		}

		qu := strings.Split(fields[4], ":")
		if len(qu) < 2 {
			return nil, fmt.Errorf("cannot parse tx_queues and rx_queues: %q", line)
		}

		tx, err := strconv.ParseUint(qu[0], 16, 64)
		if err != nil {
			return nil, err
		}
		updQueues["tx_queue"] += float64(tx)

		rx, err := strconv.ParseUint(qu[1], 16, 64)
		if err != nil {
			return nil, err
		}
		updQueues["rx_queue"] += float64(rx)
	}

	return updQueues, nil
}
