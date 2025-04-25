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
*/
import "C"

var (
	sysctlRaw    = unix.SysctlRaw
	tcpSendTotal = "bsdNetstatTcpSendPacketsTotal"
	tcpRecvTotal = "bsdNetstatTcpRecvPacketsTotal"

	counterMetrics = map[string]*prometheus.Desc{
		tcpSendTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "tcp_transmit_packets_total"),
			"TCP packets sent", nil, nil),
		tcpRecvTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "netstat", "tcp_receive_packets_total"),
			"TCP packets received", nil, nil),
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

	allStats := make(map[string]float64)

	for k, v := range tcpStats {
		allStats[k] = v
	}

	for metricKey, metricData := range counterMetrics {
		ch <- prometheus.MustNewConstMetric(
			metricData,
			prometheus.CounterValue,
			allStats[metricKey],
		)
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
	}

	return make([]byte, 0, 0)
}
