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

var metricDescs = []*prometheus.Desc{
	prometheus.NewDesc(
		"tcp_send_packet_total",
		"tcp_send_packet_total",
		nil, nil,
	),
	prometheus.NewDesc(
		"tcp_recv_packet_total",
		"tcp_recv_packet_total",
		nil, nil,
	),
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

func getData(queryString string) ([]byte, error) {
	data, err := unix.SysctlRaw(queryString)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	if len(data) < int(unsafe.Sizeof(C.struct_tcpstat{})) {
		return nil, errors.New("Data Size mismatch")
	}
	return data, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {

	var result []float64

	tcpData, err := getData("net.inet.tcp.stats")
	if err != nil {
		return err
	}

	tcpStats := *(*C.struct_tcpstat)(unsafe.Pointer(&tcpData[0]))

	result = append(result, float64(tcpStats.tcps_sndtotal))
	result = append(result, float64(tcpStats.tcps_rcvtotal))

	for index, value := range metricDescs {
		ch <- prometheus.MustNewConstMetric(
			value,
			prometheus.UntypedValue,
			result[index],
		)
	}

	return nil
}
