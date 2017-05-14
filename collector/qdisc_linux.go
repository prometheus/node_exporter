// Copyright 2017 The Prometheus Authors
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

// +build !noqdisc

package collector

import (
	"github.com/ema/qdisc"
	"github.com/prometheus/client_golang/prometheus"
)

type qdiscStatCollector struct {
	bytes      typedDesc
	packets    typedDesc
	drops      typedDesc
	requeues   typedDesc
	overlimits typedDesc
}

func init() {
	Factories["qdisc"] = NewQdiscStatCollector
}

func NewQdiscStatCollector() (Collector, error) {
	return &qdiscStatCollector{
		bytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "bytes_total"),
			"Number of bytes sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		packets: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "packets_total"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		drops: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "drops_total"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		requeues: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "requeues_total"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		overlimits: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "overlimits_total"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
	}, nil
}

func (c *qdiscStatCollector) Update(ch chan<- prometheus.Metric) error {
	msgs, err := qdisc.Get()
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		// Only report root qdisc info
		if msg.Parent != 0 {
			continue
		}

		ch <- c.bytes.mustNewConstMetric(float64(msg.Bytes), msg.IfaceName, msg.Kind)
		ch <- c.packets.mustNewConstMetric(float64(msg.Packets), msg.IfaceName, msg.Kind)
		ch <- c.drops.mustNewConstMetric(float64(msg.Drops), msg.IfaceName, msg.Kind)
		ch <- c.requeues.mustNewConstMetric(float64(msg.Requeues), msg.IfaceName, msg.Kind)
		ch <- c.overlimits.mustNewConstMetric(float64(msg.Overlimits), msg.IfaceName, msg.Kind)
		//fmt.Printf("%+v\n", m)
	}

	return nil
}
