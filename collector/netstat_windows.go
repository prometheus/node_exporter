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

// +build !nonetstat

package collector

import (
	"fmt"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/net"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	netStatsSubsystem = "netstat"
)

var (
	netStatFields = kingpin.Flag("collector.netstat.fields", "Regexp of fields to return for netstat collector.").Default("(all|tcp|tcp4|tcp6|udp|udp4|udp6|inet|inet4|inet6)$").String()
)

type netStatCollector struct {
	fieldPattern *regexp.Regexp
	logger       log.Logger
}

func init() {
	registerCollector("netstat", defaultEnabled, NewNetStatCollector)
}

// NewNetStatCollector takes and returns
// a new Collector exposing network stats.
func NewNetStatCollector(logger log.Logger) (Collector, error) {
	pattern := regexp.MustCompile(*netStatFields)
	return &netStatCollector{
		fieldPattern: pattern,
		logger:       logger,
	}, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	netStats, err := getNetStats()
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	for protocol, protocolStats := range netStats {
		if !c.fieldPattern.MatchString(protocol) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, netStatsSubsystem, protocol),
				fmt.Sprintf("Statistic connections %s.", protocol),
				nil, nil,
			),
			prometheus.CounterValue, float64(protocolStats),
		)
	}
	return nil
}

func getNetStats() (map[string]int, error) {
	netStats := make(map[string]int)
	allNetStats, err := net.Connections("all")
	if err != nil {
		return nil, err
	}
	netStats["all"] = len(allNetStats)

	// Get connections tcp
	tcpNetStats, _ := net.Connections("tcp")
	netStats["tcp"] = len(tcpNetStats)

	// Get connections tcp4
	tcp4NetStats, _ := net.Connections("tcp4")
	netStats["tcp4"] = len(tcp4NetStats)

	// Get connections tcp6
	tcp6NetStats, _ := net.Connections("tcp6")
	netStats["tcp6"] = len(tcp6NetStats)

	// Get connections udp
	udpNetStats, _ := net.Connections("udp")
	netStats["udp"] = len(udpNetStats)

	// Get connections udp4
	udp4NetStats, _ := net.Connections("udp4")
	netStats["udp4"] = len(udp4NetStats)

	// Get connections udp6
	udp6NetStats, _ := net.Connections("udp6")
	netStats["udp6"] = len(udp6NetStats)

	// Get connections inet
	inetNetStats, _ := net.Connections("inet")
	netStats["inet"] = len(inetNetStats)

	// Get connections inet4
	inet4NetStats, _ := net.Connections("inet4")
	netStats["inet4"] = len(inet4NetStats)

	// Get connections inet6
	inet6NetStats, _ := net.Connections("inet6")
	netStats["inet6"] = len(inet6NetStats)

	return netStats, nil
}
