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
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

type fnatCollector struct {
	Collector
	fs                                                                          procfs.FS
	backendConnectionsActive, backendConnectionsInact, backendWeight            typedDesc
	connections, incomingPackets, outgoingPackets, incomingBytes, outgoingBytes typedDesc
	stat                                                                        typedDesc
}

func init() {
	registerCollector("fnat", defaultDisabled, NewFNATCollector)
}

// NewFNATCollector sets up a new collector for FNAT metrics. It accepts the
// "procfs" config parameter to override the default proc location (/proc).
func NewFNATCollector() (Collector, error) {
	return newFNATCollector()
}

func newFNATCollector() (*fnatCollector, error) {
	var (
		fnatStatLabelNames = []string{
			"CPU",
		}
		fnatBackendLabelNames = []string{
			"local_address",
			"local_port",
			"remote_address",
			"remote_port",
			"proto",
		}
		c         fnatCollector
		err       error
		subsystem = "fnat"
	)

	c.fs, err = procfs.NewFS(*procPath)
	if err != nil {
		return nil, err
	}

	c.connections = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "connections_total"),
		"The total number of connections made.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.incomingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "incoming_packets_total"),
		"The total number of incoming packets.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.outgoingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "outgoing_packets_total"),
		"The total number of outgoing packets.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.incomingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "incoming_bytes_total"),
		"The total amount of incoming data.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.outgoingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "outgoing_bytes_total"),
		"The total amount of outgoing data.",
		fnatStatLabelNames, nil,
	), prometheus.CounterValue}
	c.backendConnectionsActive = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_connections_active"),
		"The current active connections by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendConnectionsInact = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_connections_inactive"),
		"The current inactive connections by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendWeight = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "backend_weight"),
		"The current backend weight by local and remote address.",
		fnatBackendLabelNames, nil,
	), prometheus.GaugeValue}

	return &c, nil
}

func (c *fnatCollector) Update(ch chan<- prometheus.Metric) error {
	fnatStats, err := c.fs.NewFNATStats()
	if err != nil {
		// Cannot access ipvs metrics, report no error.
		if os.IsNotExist(err) {
			log.Debug("fnat collector metrics are not available for this system")
			return nil
		}
		return fmt.Errorf("could not get FNAT stats: %s", err)
	}
	for _, statsFnat := range fnatStats.Stat {
		statsLabelValues := []string{
			statsFnat.Cpu,
		}
		ch <- c.connections.mustNewConstMetric(float64(statsFnat.Connections), statsLabelValues...)
		ch <- c.incomingPackets.mustNewConstMetric(float64(statsFnat.IncomingPackets), statsLabelValues...)
		ch <- c.outgoingPackets.mustNewConstMetric(float64(statsFnat.OutgoingPackets), statsLabelValues...)
		ch <- c.incomingBytes.mustNewConstMetric(float64(statsFnat.IncomingBytes), statsLabelValues...)
		ch <- c.outgoingBytes.mustNewConstMetric(float64(statsFnat.OutgoingBytes), statsLabelValues...)

	}

	backendStats, err := c.fs.NewFNATBackendStatus()
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
		ch <- c.backendConnectionsActive.mustNewConstMetric(float64(backend.ActiveConn), labelValues...)
		ch <- c.backendConnectionsInact.mustNewConstMetric(float64(backend.InactConn), labelValues...)
		ch <- c.backendWeight.mustNewConstMetric(float64(backend.Weight), labelValues...)
	}
	return nil
}
