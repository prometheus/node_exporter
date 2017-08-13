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

// +build !nosockstat

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
	sockStatSubsystem = "sockstat"
)

// Used for calculating the total memory bytes on TCP and UDP.
var pageSize = os.Getpagesize()

type sockStatCollector struct{}

func init() {
	registerCollector(sockStatSubsystem, defaultEnabled, NewSockStatCollector)
}

// NewSockStatCollector returns a new Collector exposing socket stats.
func NewSockStatCollector() (Collector, error) {
	return &sockStatCollector{}, nil
}

func (c *sockStatCollector) Update(ch chan<- prometheus.Metric) error {
	sockStats, err := getSockStats(procFilePath("net/sockstat"))
	if err != nil {
		return fmt.Errorf("couldn't get sockstats: %s", err)
	}
	for protocol, protocolStats := range sockStats {
		for name, value := range protocolStats {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in sockstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, sockStatSubsystem, protocol+"_"+name),
					fmt.Sprintf("Number of %s sockets in state %s.", protocol, name),
					nil, nil,
				),
				prometheus.GaugeValue, v,
			)
		}
	}
	return err
}

func getSockStats(fileName string) (map[string]map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseSockStats(file, fileName)
}

func parseSockStats(r io.Reader, fileName string) (map[string]map[string]string, error) {
	var (
		sockStat = map[string]map[string]string{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		// Remove trailing ':'.
		protocol := line[0][:len(line[0])-1]
		sockStat[protocol] = map[string]string{}

		for i := 1; i < len(line) && i+1 < len(line); i++ {
			sockStat[protocol][line[i]] = line[i+1]
			i++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// The mem metrics is the count of pages used. Multiply the mem metrics by
	// the page size from the kernel to get the number of bytes used.
	//
	// Update the TCP mem from page count to bytes.
	pageCount, err := strconv.Atoi(sockStat["TCP"]["mem"])
	if err != nil {
		return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["TCP"]["mem"], err)
	}
	sockStat["TCP"]["mem_bytes"] = strconv.Itoa(pageCount * pageSize)

	// Update the UDP mem from page count to bytes.
	if udpMem := sockStat["UDP"]["mem"]; udpMem != "" {
		pageCount, err = strconv.Atoi(udpMem)
		if err != nil {
			return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["UDP"]["mem"], err)
		}
		sockStat["UDP"]["mem_bytes"] = strconv.Itoa(pageCount * pageSize)
	}

	return sockStat, nil
}
