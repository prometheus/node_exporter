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
// +build freebsd

package collector

import (
	"encoding/binary"
	"errors"
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
#include <netinet/ip_var.h>
#include <netinet6/ip6_var.h>
*/
import "C"

var (
	sysctlRaw              = unix.SysctlRaw
	tcpSendTotal           = "bsdNetstatTcpSendPacketsTotal"
	tcpRecvTotal           = "bsdNetstatTcpRecvPacketsTotal"
	udpSendTotal           = "bsdNetstatUdpSendPacketsTotal"
	udpRecvTotal           = "bsdNetstatUdpRecvPacketsTotal"
	ipv4SendTotal          = "bsdNetstatIPv4SendPacketsTotal"
	ipv4RawSendTotal       = "bsdNetstatIPv4RawSendPacketsTotal"
	ipv4RecvTotal          = "bsdNetstatIPv4RecvPacketsTotal"
	ipv4RecvFragmentsTotal = "bsdNetstatIPv4RecvFragmentsTotal"
	ipv4ForwardTotal       = "bsdNetstatIPv4ForwardTotal"
	ipv4FastForwardTotal   = "bsdNetstatIPv4FastForwardTotal"
	ipv4DeliveredTotal     = "bsdNetstatIPv4DeliveredTotal"
	ipv6SendTotal          = "bsdNetstatIPv6SendPacketsTotal"
	ipv6RawSendTotal       = "bsdNetstatIPv6RawSendPacketsTotal"
	ipv6RecvTotal          = "bsdNetstatIPv6RecvPacketsTotal"
	ipv6RecvFragmentsTotal = "bsdNetstatIPv6RecvFragmentsTotal"
	ipv6ForwardTotal       = "bsdNetstatIPv6ForwardTotal"
	ipv6DeliveredTotal     = "bsdNetstatIPv6DeliveredTotal"

	tcpStates = []string{
		"CLOSED", "LISTEN", "SYN_SENT", "SYN_RCVD",
		"ESTABLISHED", "CLOSE_WAIT", "FIN_WAIT_1", "CLOSING",
		"LAST_ACK", "FIN_WAIT_2", "TIME_WAIT",
	}

	tcpStatesMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netstat", "tcp_connections"),
		"Number of TCP connections per state", []string{"state"}, nil)

	counterMetrics = map[string]*prometheus.Desc{
		// TCP stats
		tcpSendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "tcp_transmit_packets_total"),
			"TCP packets sent", nil, nil),
		tcpRecvTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "tcp_receive_packets_total"),
			"TCP packets received", nil, nil),

		// UDP stats
		udpSendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "udp_transmit_packets_total"),
			"UDP packets sent", nil, nil),
		udpRecvTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "udp_receive_packets_total"),
			"UDP packets received", nil, nil),

		// IPv4 stats
		ipv4SendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_transmit_packets_total"),
			"IPv4 packets sent from this host", nil, nil),
		ipv4RawSendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_transmit_raw_packets_total"),
			"IPv4 raw packets generated", nil, nil),
		ipv4RecvTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_receive_packets_total"),
			"IPv4 packets received", nil, nil),
		ipv4RecvFragmentsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_receive_fragments_total"),
			"IPv4 fragments received", nil, nil),
		ipv4ForwardTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_forward_total"),
			"IPv4 packets forwarded", nil, nil),
		ipv4FastForwardTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_fast_forward_total"),
			"IPv4 packets fast forwarded", nil, nil),
		ipv4DeliveredTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip4_delivered_total"),
			"IPv4 packets delivered to the upper layer (packets for this host)", nil, nil),

		// IPv6 stats
		ipv6SendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_transmit_packets_total"),
			"IPv6 packets sent from this host", nil, nil),
		ipv6RawSendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_transmit_raw_packets_total"),
			"IPv6 raw packets generated", nil, nil),
		ipv6RecvTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_receive_packets_total"),
			"IPv6 packets received", nil, nil),
		ipv6RecvFragmentsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_receive_fragments_total"),
			"IPv6 fragments received", nil, nil),
		ipv6ForwardTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_forward_total"),
			"IPv6 packets forwarded", nil, nil),
		ipv6DeliveredTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "ip6_delivered_total"),
			"IPv6 packets delivered to the upper layer (packets for this host)", nil, nil),
	}
)

type NetstatData struct {
	structSize int
	sysctl     string
}

type NetstatMetrics map[string]float64

type NetstatTCPData NetstatData

func NewTCPStat() *NetstatTCPData {
	return &NetstatTCPData{
		structSize: int(unsafe.Sizeof(C.struct_tcpstat{})),
		sysctl:     "net.inet.tcp.stats",
	}
}

func (netstatMetric *NetstatTCPData) GetData() (NetstatMetrics, error) {
	data, err := getData(netstatMetric.sysctl, netstatMetric.structSize)
	if err != nil {
		return nil, err
	}

	tcpStats := *(*C.struct_tcpstat)(unsafe.Pointer(&data[0]))

	return NetstatMetrics{
		tcpSendTotal: float64(tcpStats.tcps_sndtotal),
		tcpRecvTotal: float64(tcpStats.tcps_rcvtotal),
	}, nil
}

type NetstatUDPData NetstatData

func NewUDPStat() *NetstatUDPData {
	return &NetstatUDPData{
		structSize: int(unsafe.Sizeof(C.struct_udpstat{})),
		sysctl:     "net.inet.udp.stats",
	}
}

func (netstatMetric *NetstatUDPData) GetData() (NetstatMetrics, error) {
	data, err := getData(netstatMetric.sysctl, netstatMetric.structSize)
	if err != nil {
		return nil, err
	}

	udpStats := *(*C.struct_udpstat)(unsafe.Pointer(&data[0]))

	return NetstatMetrics{
		udpSendTotal: float64(udpStats.udps_opackets),
		udpRecvTotal: float64(udpStats.udps_ipackets),
	}, nil
}

type NetstatIPv4Data NetstatData

func NewIPv4Stat() *NetstatIPv4Data {
	return &NetstatIPv4Data{
		structSize: int(unsafe.Sizeof(C.struct_ipstat{})),
		sysctl:     "net.inet.ip.stats",
	}
}

func (netstatMetric *NetstatIPv4Data) GetData() (NetstatMetrics, error) {
	data, err := getData(netstatMetric.sysctl, netstatMetric.structSize)
	if err != nil {
		return nil, err
	}

	ipStats := *(*C.struct_ipstat)(unsafe.Pointer(&data[0]))

	return NetstatMetrics{
		ipv4SendTotal:          float64(ipStats.ips_localout),
		ipv4RawSendTotal:       float64(ipStats.ips_rawout),
		ipv4RecvTotal:          float64(ipStats.ips_total),
		ipv4RecvFragmentsTotal: float64(ipStats.ips_fragments),
		ipv4ForwardTotal:       float64(ipStats.ips_forward),
		ipv4FastForwardTotal:   float64(ipStats.ips_fastforward),
		ipv4DeliveredTotal:     float64(ipStats.ips_delivered),
	}, nil
}

type NetstatIPv6Data NetstatData

func NewIPv6Stat() *NetstatIPv6Data {
	return &NetstatIPv6Data{
		structSize: int(unsafe.Sizeof(C.struct_ipstat{})),
		sysctl:     "net.inet6.ip6.stats",
	}
}

func (netstatMetric *NetstatIPv6Data) GetData() (NetstatMetrics, error) {
	data, err := getData(netstatMetric.sysctl, netstatMetric.structSize)
	if err != nil {
		return nil, err
	}

	ipStats := *(*C.struct_ip6stat)(unsafe.Pointer(&data[0]))

	return NetstatMetrics{
		ipv6SendTotal:          float64(ipStats.ip6s_localout),
		ipv6RawSendTotal:       float64(ipStats.ip6s_rawout),
		ipv6RecvTotal:          float64(ipStats.ip6s_total),
		ipv6RecvFragmentsTotal: float64(ipStats.ip6s_fragments),
		ipv6ForwardTotal:       float64(ipStats.ip6s_forward),
		ipv6DeliveredTotal:     float64(ipStats.ip6s_delivered),
	}, nil
}

func getData(queryString string, expectedSize int) ([]byte, error) {
	data, err := sysctlRaw(queryString)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if len(data) < expectedSize {
		return nil, errors.New("Data Size mismatch")
	}
	return data, nil
}

func getTCPStates() ([]uint64, error) {

	// This sysctl returns an array of uint64
	data, err := sysctlRaw("net.inet.tcp.states")

	if err != nil {
		return nil, err
	}

	if len(data)/8 != len(tcpStates) {
		return nil, fmt.Errorf("invalid TCP states data: expected %d entries, found %d", len(tcpStates), len(data)/8)
	}

	states := make([]uint64, 0)

	offset := 0
	for range len(tcpStates) {
		s := data[offset : offset+8]
		offset += 8
		states = append(states, binary.NativeEndian.Uint64(s))
	}
	return states, nil
}

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

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	tcpStats, err := NewTCPStat().GetData()
	if err != nil {
		return err
	}

	udpStats, err := NewUDPStat().GetData()
	if err != nil {
		return err
	}

	ipv4Stats, err := NewIPv4Stat().GetData()
	if err != nil {
		return err
	}

	ipv6Stats, err := NewIPv6Stat().GetData()
	if err != nil {
		return err
	}

	allStats := make(map[string]float64)

	for k, v := range tcpStats {
		allStats[k] = v
	}

	for k, v := range udpStats {
		allStats[k] = v
	}

	for k, v := range ipv4Stats {
		allStats[k] = v
	}

	for k, v := range ipv6Stats {
		allStats[k] = v
	}

	for metricKey, metricData := range counterMetrics {
		ch <- prometheus.MustNewConstMetric(
			metricData,
			prometheus.CounterValue,
			allStats[metricKey],
		)
	}

	tcpConnsPerStates, err := getTCPStates()

	if err != nil {
		return err
	}

	for i, value := range tcpConnsPerStates {
		ch <- prometheus.MustNewConstMetric(tcpStatesMetric, prometheus.GaugeValue, float64(value), tcpStates[i])
	}
	return nil
}

// Used by tests to mock unix.SysctlRaw
func getFreeBSDDataMock(sysctl string) []byte {

	if sysctl == "net.inet.tcp.stats" {
		tcpStats := C.struct_tcpstat{
			tcps_sndtotal: 1234,
			tcps_rcvtotal: 4321,
		}
		size := int(unsafe.Sizeof(C.struct_tcpstat{}))

		return unsafe.Slice((*byte)(unsafe.Pointer(&tcpStats)), size)
	} else if sysctl == "net.inet.tcp.states" {
		tcpStatesSlice := make([]byte, 0, len(tcpStates)*8)
		tcpStatesValues := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

		for _, value := range tcpStatesValues {
			tcpStatesSlice = binary.NativeEndian.AppendUint64(tcpStatesSlice, value)
		}

		return tcpStatesSlice

	} else if sysctl == "net.inet.udp.stats" {
		udpStats := C.struct_udpstat{
			udps_opackets: 1234,
			udps_ipackets: 4321,
		}
		size := int(unsafe.Sizeof(C.struct_udpstat{}))

		return unsafe.Slice((*byte)(unsafe.Pointer(&udpStats)), size)
	} else if sysctl == "net.inet.ip.stats" {
		ipStats := C.struct_ipstat{
			ips_localout:    1234,
			ips_rawout:      1235,
			ips_total:       1236,
			ips_fragments:   1237,
			ips_forward:     1238,
			ips_fastforward: 1239,
			ips_delivered:   1240,
		}
		size := int(unsafe.Sizeof(C.struct_ipstat{}))

		return unsafe.Slice((*byte)(unsafe.Pointer(&ipStats)), size)
	} else if sysctl == "net.inet6.ip6.stats" {
		ipStats := C.struct_ip6stat{
			ip6s_localout:  1234,
			ip6s_rawout:    1235,
			ip6s_total:     1236,
			ip6s_fragments: 1237,
			ip6s_forward:   1238,
			ip6s_delivered: 1240,
		}
		size := int(unsafe.Sizeof(C.struct_ip6stat{}))

		return unsafe.Slice((*byte)(unsafe.Pointer(&ipStats)), size)
	}

	return make([]byte, 0, 0)
}
