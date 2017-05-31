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
	fs                                                                                            procfs.FS
	backendConnectionsActive, backendConnectionsInact, backendWeight,
	backendIncomingConnectionsPerSecond, backendIncomingPackgesPerSecond,
	backendOutgoingPackgesPerSecond, backendIncomingBytesPerSecond, backendOutgoingBytesPerSecond typedDesc
	connections, incomingPackets, outgoingPackets, incomingBytes, outgoingBytes                   typedDesc
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

	c.connections = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "connections_total"),
		"The total number of connections made.",
		nil, nil,
	), prometheus.CounterValue}
	c.incomingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "incoming_packets_total"),
		"The total number of incoming packets.",
		nil, nil,
	), prometheus.CounterValue}
	c.outgoingPackets = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "outgoing_packets_total"),
		"The total number of outgoing packets.",
		nil, nil,
	), prometheus.CounterValue}
	c.incomingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "incoming_bytes_total"),
		"The total amount of incoming data.",
		nil, nil,
	), prometheus.CounterValue}
	c.outgoingBytes = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "outgoing_bytes_total"),
		"The total amount of outgoing data.",
		nil, nil,
	), prometheus.CounterValue}
	c.backendConnectionsActive = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_connections_active"),
		"The current active connections by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendConnectionsInact = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_connections_inactive"),
		"The current inactive connections by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendWeight = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_weight"),
		"The current backend weight by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	// For cps, pps, Bps
	c.backendIncomingConnectionsPerSecond = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_cps"),
		"The current backend cps by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendIncomingPackgesPerSecond = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_inpps"),
		"The current backend inpps by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendOutgoingPackgesPerSecond = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_outpps"),
		"The current backend outpps by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendIncomingBytesPerSecond = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_inbps"),
		"The current backend inbps by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}
	c.backendOutgoingBytesPerSecond = typedDesc{prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "backend_outbps"),
		"The current backend outbps by local and remote address.",
		ipvsBackendLabelNames, nil,
	), prometheus.GaugeValue}

	return &c, nil
}

func (c *ipvsCollector) Update(ch chan<- prometheus.Metric) error {
	ipvsStats, err := c.fs.NewIPVSStats()
	if err != nil {
		return fmt.Errorf("could not get IPVS stats: %s", err)
	}
	ch <- c.connections.mustNewConstMetric(float64(ipvsStats.Connections))
	ch <- c.incomingPackets.mustNewConstMetric(float64(ipvsStats.IncomingPackets))
	ch <- c.outgoingPackets.mustNewConstMetric(float64(ipvsStats.OutgoingPackets))
	ch <- c.incomingBytes.mustNewConstMetric(float64(ipvsStats.IncomingBytes))
	ch <- c.outgoingBytes.mustNewConstMetric(float64(ipvsStats.OutgoingBytes))

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
		ch <- c.backendConnectionsActive.mustNewConstMetric(float64(backend.ActiveConn), labelValues...)
		ch <- c.backendConnectionsInact.mustNewConstMetric(float64(backend.InactConn), labelValues...)
		ch <- c.backendWeight.mustNewConstMetric(float64(backend.Weight), labelValues...)
		// For cps, pps, Bps
		ch <- c.backendIncomingConnectionsPerSecond.mustNewConstMetric(float64(backend.IncomingConnectionsPerSecond), labelValues...)
		ch <- c.backendIncomingPackgesPerSecond.mustNewConstMetric(float64(backend.IncomingPackgesPerSecond), labelValues...)
		ch <- c.backendOutgoingPackgesPerSecond.mustNewConstMetric(float64(backend.OutgoingPackgesPerSecond), labelValues...)
		ch <- c.backendIncomingBytesPerSecond.mustNewConstMetric(float64(backend.IncomingBytesPerSecond), labelValues...)
		ch <- c.backendOutgoingBytesPerSecond.mustNewConstMetric(float64(backend.OutgoingBytesPerSecond), labelValues...)
	}
	return nil
}
