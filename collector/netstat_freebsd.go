// Copyright 2024 The Prometheus Authors
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

//go:build freebsd

package collector

import (
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

/*
#include <sys/types.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/tcp.h>
#include <netinet/tcp_var.h>
#include <netinet/udp.h>
#include <netinet/udp_var.h>
*/
import "C"

var (
	// TCP metrics
	bsdNetstatTcpSendPacketsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_transmit_packets_total"),
		"TCP packets sent",
		nil, nil,
	)
	bsdNetstatTcpRecvPacketsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_receive_packets_total"),
		"TCP packets received",
		nil, nil,
	)
	bsdNetstatTcpConnectionAttempts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_connection_attempts_total"),
		"Number of times TCP connections have been initiated",
		nil, nil,
	)
	bsdNetstatTcpConnectionAccepts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_connection_accepts_total"),
		"Number of times TCP connections have made it to established state",
		nil, nil,
	)
	bsdNetstatTcpConnectionDrops = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_connection_drops_total"),
		"Number of dropped TCP connections",
		nil, nil,
	)
	bsdNetstatTcpRetransmitPackets = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_retransmit_packets_total"),
		"Number of TCP data packets retransmitted",
		nil, nil,
	)
	// UDP metrics
	bsdNetstatUdpSendPacketsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "udp_transmit_packets_total"),
		"UDP packets sent",
		nil, nil,
	)
	bsdNetstatUdpRecvPacketsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "udp_receive_packets_total"),
		"UDP packets received",
		nil, nil,
	)
	bsdNetstatUdpHeaderDrops = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "udp_header_drops_total"),
		"Number of UDP packets dropped due to invalid header",
		nil, nil,
	)
	bsdNetstatUdpBadChecksum = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "udp_bad_checksum_total"),
		"Number of UDP packets dropped due to bad checksum",
		nil, nil,
	)
	bsdNetstatUdpNoPort = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "udp_no_port_total"),
		"Number of UDP packets to port with no listener",
		nil, nil,
	)
)

type netStatCollector struct {
	netStatMetric *prometheus.Desc
}

func init() {
	registerCollector("netstat", defaultEnabled, NewNetStatCollector)
}

func NewNetStatCollector(logger *slog.Logger) (Collector, error) {
	return &netStatCollector{}, nil
}

func (c *netStatCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.netStatMetric
}

func (c *netStatCollector) Collect(ch chan<- prometheus.Metric) {
	_ = c.Update(ch)
}

func getData(queryString string) ([]byte, error) {
	data, err := unix.SysctlRaw(queryString)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return data, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	tcpData, err := getData("net.inet.tcp.stats")
	if err != nil {
		return fmt.Errorf("failed to get TCP stats: %w", err)
	}
	if len(tcpData) < int(unsafe.Sizeof(C.struct_tcpstat{})) {
		return fmt.Errorf("TCP data size mismatch: got %d, want >= %d", len(tcpData), unsafe.Sizeof(C.struct_tcpstat{}))
	}

	tcpStats := *(*C.struct_tcpstat)(unsafe.Pointer(&tcpData[0]))

	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpSendPacketsTotal,
		prometheus.CounterValue,
		float64(tcpStats.tcps_sndtotal),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpRecvPacketsTotal,
		prometheus.CounterValue,
		float64(tcpStats.tcps_rcvtotal),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpConnectionAttempts,
		prometheus.CounterValue,
		float64(tcpStats.tcps_connattempt),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpConnectionAccepts,
		prometheus.CounterValue,
		float64(tcpStats.tcps_accepts),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpConnectionDrops,
		prometheus.CounterValue,
		float64(tcpStats.tcps_drops),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatTcpRetransmitPackets,
		prometheus.CounterValue,
		float64(tcpStats.tcps_sndrexmitpack),
	)

	udpData, err := getData("net.inet.udp.stats")
	if err != nil {
		return fmt.Errorf("failed to get UDP stats: %w", err)
	}
	if len(udpData) < int(unsafe.Sizeof(C.struct_udpstat{})) {
		return fmt.Errorf("UDP data size mismatch: got %d, want >= %d", len(udpData), unsafe.Sizeof(C.struct_udpstat{}))
	}

	udpStats := *(*C.struct_udpstat)(unsafe.Pointer(&udpData[0]))

	ch <- prometheus.MustNewConstMetric(
		bsdNetstatUdpSendPacketsTotal,
		prometheus.CounterValue,
		float64(udpStats.udps_opackets),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatUdpRecvPacketsTotal,
		prometheus.CounterValue,
		float64(udpStats.udps_ipackets),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatUdpHeaderDrops,
		prometheus.CounterValue,
		float64(udpStats.udps_hdrops),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatUdpBadChecksum,
		prometheus.CounterValue,
		float64(udpStats.udps_badsum),
	)
	ch <- prometheus.MustNewConstMetric(
		bsdNetstatUdpNoPort,
		prometheus.CounterValue,
		float64(udpStats.udps_noport),
	)

	return nil
}
