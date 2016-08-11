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

// +build !nonetstat

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	netStatsSubsystem = "netstat"
)

type netStatCollector struct{}

func init() {
	Factories["netstat"] = NewNetStatCollector
}

// NewNetStatCollector takes and returns
// a new Collector exposing network stats.
func NewNetStatCollector() (Collector, error) {
	return &netStatCollector{}, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netStats, err := getNetStats(procFilePath("net/netstat"))
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	snmpStats, err := getNetStats(procFilePath("net/snmp"))
	if err != nil {
		return fmt.Errorf("couldn't get SNMP stats: %s", err)
	}
	// Merge the results of snmpStats into netStats (collisions are possible, but
	// we know that the keys are always unique for the given use case).
	for k, v := range snmpStats {
		netStats[k] = v
	}
	for protocol, protocolStats := range netStats {
		for name, value := range protocolStats {
			key := protocol + "_" + name
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, netStatsSubsystem, key),
					fmt.Sprintf("Protocol %s statistic %s.", protocol, name),
					nil, nil,
				),
				prometheus.UntypedValue, v,
			)
		}
	}
	return nil
}

func getNetStats(fileName string) (map[string]map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetStats(file, fileName)
}

func parseNetStats(r io.Reader, fileName string) (map[string]map[string]string, error) {
	var (
		netStats = map[string]map[string]string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		nameParts := strings.Split(string(scanner.Text()), " ")
		scanner.Scan()
		valueParts := strings.Split(string(scanner.Text()), " ")
		// Remove trailing :.
		protocol := nameParts[0][:len(nameParts[0])-1]
		netStats[protocol] = map[string]string{}
		if len(nameParts) != len(valueParts) {
			return nil, fmt.Errorf("mismatch field count mismatch in %s: %s",
				fileName, protocol)
		}
		for i := 1; i < len(nameParts); i++ {
			netStats[protocol][nameParts[i]] = valueParts[i]
		}
	}

	return netStats, nil
}
