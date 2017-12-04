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

// +build linux
// +build !nomeminfo

package collector

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"os"
	"strconv"
	"strings"
)

const (
	netConnSubsystem = "netconns"
)

type netconnCollector struct{}

func init() {
	registerCollector(netConnSubsystem, defaultEnabled, NewNetconnCollector)
}

// NewNetconnCollector returns a new Collector exposing tcp connections stats.
func NewNetconnCollector() (Collector, error) {
	return &netconnCollector{}, nil
}

func hex2dec(hexstr string) string {
	i, _ := strconv.ParseInt(hexstr, 16, 0)
	return strconv.FormatInt(i, 10)
}

func hex_to_ip(hexstr string) (string, string) {
	var ip string
	if len(hexstr) != 8 {
		err := "parse error"
		return ip, err
	}

	i1, _ := strconv.ParseInt(hexstr[6:8], 16, 0)
	i2, _ := strconv.ParseInt(hexstr[4:6], 16, 0)
	i3, _ := strconv.ParseInt(hexstr[2:4], 16, 0)
	i4, _ := strconv.ParseInt(hexstr[0:2], 16, 0)
	ip = fmt.Sprintf("%d.%d.%d.%d", i1, i2, i3, i4)

	return ip, ""
}

func CountConns() map[string]int {
	fd, _ := os.Open("/proc/net/tcp")
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	m := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(strings.Trim(line, " "), " ")
		// Skip non-connection info
		if len(tokens[1]) != 13 || len(tokens[2]) != 13 {
			continue
		}
		localIP := strings.Split(tokens[1], ":")[0]
		remoteIP := strings.Split(tokens[2], ":")[0]

		key := fmt.Sprintf("%s-%s", localIP, remoteIP)
		c, ok := m[key]
		if ok {
			m[key] = c + 1
		} else {
			m[key] = 1
		}
	}
	return m
}

// Update calls to report information from map.
func (c *netconnCollector) Update(ch chan<- prometheus.Metric) error {
	log.Debugf("Collect tcp connctions info")
	connsMap := CountConns()
	variableLables := []string{"localIP", "remoteIP"}
	for k, v := range connsMap {
		ips := strings.Split(k, "-")
		localIP, _ := hex_to_ip(ips[0])
		remoteIP, _ := hex_to_ip(ips[1])
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, netConnSubsystem, "netconns"),
				"Connection information field ",
				variableLables, nil,
			),
			prometheus.GaugeValue, float64(v), localIP, remoteIP,
		)
	}
	return nil
}
