// Copyright 2025 The Prometheus Authors
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

//go:build !nonetinterface
// +build !nonetinterface

package collector

import (
	"log/slog"

	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

type netinterfaceCollector struct {
	logger     *slog.Logger
	collisions *prometheus.Desc
	ibytes     *prometheus.Desc
	ipackets   *prometheus.Desc
	obytes     *prometheus.Desc
	opackets   *prometheus.Desc
}

const (
	netinterfaceSubsystem = "netinterface"
)

func init() {
	registerCollector("netinterface", defaultEnabled, NewNetinterfaceCollector)
}

func NewNetinterfaceCollector(logger *slog.Logger) (Collector, error) {
	labels := []string{"interface"}
	return &netinterfaceCollector{
		logger: logger,
		collisions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netinterfaceSubsystem, "collisions_total"),
			"Total number of CSMA collisions on the interface.", labels, nil,
		),
		ibytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netinterfaceSubsystem, "receive_bytes_total"),
			"Total number of bytes received on the interface.", labels, nil,
		),
		ipackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netinterfaceSubsystem, "receive_packets_total"),
			"Total number of packets received on the interface.", labels, nil,
		),
		obytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netinterfaceSubsystem, "transmit_bytes_total"),
			"Total number of bytes transmitted on the interface.", labels, nil,
		),
		opackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netinterfaceSubsystem, "transmit_packets_total"),
			"Total number of packets transmitted on the interface.", labels, nil,
		),
	}, nil
}

func (c *netinterfaceCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.NetIfaceStat()
	if err != nil {
		return err
	}

	for _, stat := range stats {
		iface := stat.Name

		ch <- prometheus.MustNewConstMetric(c.collisions, prometheus.CounterValue, float64(stat.Collisions), iface)
		ch <- prometheus.MustNewConstMetric(c.ibytes, prometheus.CounterValue, float64(stat.IBytes), iface)
		ch <- prometheus.MustNewConstMetric(c.ipackets, prometheus.CounterValue, float64(stat.IPackets), iface)
		ch <- prometheus.MustNewConstMetric(c.obytes, prometheus.CounterValue, float64(stat.OBytes), iface)
		ch <- prometheus.MustNewConstMetric(c.opackets, prometheus.CounterValue, float64(stat.OPackets), iface)
	}
	return nil
}
