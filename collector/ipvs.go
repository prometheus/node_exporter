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

// +build !noipvs

package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type ipvsCollector struct {
	Collector
	fs                                                                          procfs.FS
	backendConnectionsActive, backendConnectionsInact, backendWeight            *prometheus.GaugeVec
	connections, incomingPackets, outgoingPackets, incomingBytes, outgoingBytes prometheus.Counter
}

func init() {
	Factories["ipvs"] = NewIPVSCollector
}

// NewIPVSCollector sets up a new collector for IPVS metrics. It accepts the
// "procfs" config parameter to override the default proc location (/proc).
func NewIPVSCollector() (Collector, error) {
	return newIPVSCollector()
}

func newIPVSCollector() (*ipvsCollector, error) {
	var (
		ipvsBackendLabelNames = []string{
			"local_address",
			"local_port",
			"remote_address",
			"remote_port",
			"proto",
		}
		c         ipvsCollector
		err       error
		subsystem = "ipvs"
	)

	c.fs, err = procfs.NewFS(*procPath)
	if err != nil {
		return nil, err
	}

	c.connections = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "connections_total",
			Help:      "The total number of connections made.",
		},
	)
	c.incomingPackets = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "incoming_packets_total",
			Help:      "The total number of incoming packets.",
		},
	)
	c.outgoingPackets = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "outgoing_packets_total",
			Help:      "The total number of outgoing packets.",
		},
	)
	c.incomingBytes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "incoming_bytes_total",
			Help:      "The total amount of incoming data.",
		},
	)
	c.outgoingBytes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "outgoing_bytes_total",
			Help:      "The total amount of outgoing data.",
		},
	)

	c.backendConnectionsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "backend_connections_active",
			Help:      "The current active connections by local and remote address.",
		},
		ipvsBackendLabelNames,
	)
	c.backendConnectionsInact = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "backend_connections_inactive",
			Help:      "The current inactive connections by local and remote address.",
		},
		ipvsBackendLabelNames,
	)
	c.backendWeight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subsystem,
			Name:      "backend_weight",
			Help:      "The current backend weight by local and remote address.",
		},
		ipvsBackendLabelNames,
	)

	return &c, nil
}

func (c *ipvsCollector) Update(ch chan<- prometheus.Metric) error {
	ipvsStats, err := c.fs.NewIPVSStats()
	if err != nil {
		return fmt.Errorf("could not get IPVS stats: %s", err)
	}

	c.connections.Set(float64(ipvsStats.Connections))
	c.incomingPackets.Set(float64(ipvsStats.IncomingPackets))
	c.outgoingPackets.Set(float64(ipvsStats.OutgoingPackets))
	c.incomingBytes.Set(float64(ipvsStats.IncomingBytes))
	c.outgoingBytes.Set(float64(ipvsStats.OutgoingBytes))

	c.connections.Collect(ch)
	c.incomingPackets.Collect(ch)
	c.outgoingPackets.Collect(ch)
	c.incomingBytes.Collect(ch)
	c.outgoingBytes.Collect(ch)

	backendStats, err := c.fs.NewIPVSBackendStatus()
	if err != nil {
		return fmt.Errorf("could not get backend status: %s", err)
	}

	for _, backend := range backendStats {
		labelValues := []string{
			backend.LocalAddress.String(),
			strconv.FormatUint(uint64(backend.LocalPort), 10),
			backend.RemoteAddress.String(),
			strconv.FormatUint(uint64(backend.RemotePort), 10),
			backend.Proto,
		}
		c.backendConnectionsActive.WithLabelValues(labelValues...).Set(float64(backend.ActiveConn))
		c.backendConnectionsInact.WithLabelValues(labelValues...).Set(float64(backend.InactConn))
		c.backendWeight.WithLabelValues(labelValues...).Set(float64(backend.Weight))
	}

	c.backendConnectionsActive.Collect(ch)
	c.backendConnectionsInact.Collect(ch)
	c.backendWeight.Collect(ch)

	return nil
}
