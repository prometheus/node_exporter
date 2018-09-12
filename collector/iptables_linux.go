// Copyright 2018 The Prometheus Authors
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

// +build linux
// +build !noiptables

package collector

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/prometheus/client_golang/prometheus"
)

type IptablesCollector struct {
	entries       *prometheus.Desc
	ipTableIfaces map[string]*iptables.IPTables
}

func init() {
	registerCollector("iptables", defaultDisabled, NewIptablesCollector)
}

// NewIptablesCollector returns a new Collector exposing Iptables stats.
func NewIptablesCollector() (Collector, error) {

	ipTableIfaces := make(map[string]*iptables.IPTables)

	if ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4); err == nil {
		ipTableIfaces["IPV4"] = ipt
	}

	if ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv6); err == nil {
		ipTableIfaces["IPV6"] = ipt
	}

	return &IptablesCollector{
		ipTableIfaces: ipTableIfaces,
		entries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "iptables", "rules"),
			"iptables number of rules by chains",
			[]string{"proto", "table", "chain"}, nil,
		),
	}, nil
}

var (
	cstIptablesTables = []string{
		"filter",
		"nat",
		"mangle",
		"raw",
		"security",
	}
)

func (c *IptablesCollector) Update(ch chan<- prometheus.Metric) (err error) {

	for _, table := range cstIptablesTables {

		for proto, ipt := range c.ipTableIfaces {
			chains, err := ipt.ListChains(table)

			if err != nil {
				continue
			}

			for _, chain := range chains {

				res := 0

				rules, err := ipt.List(table, chain)

				if err == nil {
					res = len(rules) - 1
				}

				ch <- prometheus.MustNewConstMetric(c.entries, prometheus.GaugeValue, float64(res), proto, table, chain)
			}
		}
	}
	return nil
}
