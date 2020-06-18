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

// +build !notcpstat

package collector

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type tcpConnectionState int

const (
	// TCP_ESTABLISHED
	tcpEstablished tcpConnectionState = iota + 1
	// TCP_SYN_SENT
	tcpSynSent
	// TCP_SYN_RECV
	tcpSynRecv
	// TCP_FIN_WAIT1
	tcpFinWait1
	// TCP_FIN_WAIT2
	tcpFinWait2
	// TCP_TIME_WAIT
	tcpTimeWait
	// TCP_CLOSE
	tcpClose
	// TCP_CLOSE_WAIT
	tcpCloseWait
	// TCP_LAST_ACK
	tcpLastAck
	// TCP_LISTEN
	tcpListen
	// TCP_CLOSING
	tcpClosing
	// TCP_RX_BUFFER
	tcpRxQueuedBytes
	// TCP_TX_BUFFER
	tcpTxQueuedBytes
)

type tcpStatCollector struct {
	desc   typedDesc
	logger log.Logger
}

func init() {
	registerCollector("tcpstat", defaultDisabled, NewTCPStatCollector)
}

// NewTCPStatCollector returns a new Collector exposing network stats.
func NewTCPStatCollector(logger log.Logger) (Collector, error) {
	return &tcpStatCollector{
		desc: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tcp", "connection_states"),
			"Number of connection states.",
			[]string{"state"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

func (c *tcpStatCollector) Update(ch chan<- prometheus.Metric) error {
	tcpStats, err := getTCPStats(procFilePath("net/tcp"))
	if err != nil {
		return fmt.Errorf("couldn't get tcpstats: %w", err)
	}

	// if enabled ipv6 system
	tcp6File := procFilePath("net/tcp6")
	if _, hasIPv6 := os.Stat(tcp6File); hasIPv6 == nil {
		tcp6Stats, err := getTCPStats(tcp6File)
		if err != nil {
			return fmt.Errorf("couldn't get tcp6stats: %w", err)
		}

		for st, value := range tcp6Stats {
			tcpStats[st] += value
		}
	}

	for st, value := range tcpStats {
		ch <- c.desc.mustNewConstMetric(value, st.String())
	}
	return nil
}

func getTCPStats(statsFile string) (map[tcpConnectionState]float64, error) {
	file, err := os.Open(statsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseTCPStats(file)
}

func parseTCPStats(r io.Reader) (map[tcpConnectionState]float64, error) {
	tcpStats := map[tcpConnectionState]float64{}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(contents), "\n")[1:] {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) < 5 {
			return nil, fmt.Errorf("invalid TCP stats line: %q", line)
		}

		qu := strings.Split(parts[4], ":")
		if len(qu) < 2 {
			return nil, fmt.Errorf("cannot parse tx_queues and rx_queues: %q", line)
		}

		tx, err := strconv.ParseUint(qu[0], 16, 64)
		if err != nil {
			return nil, err
		}
		tcpStats[tcpConnectionState(tcpTxQueuedBytes)] += float64(tx)

		rx, err := strconv.ParseUint(qu[1], 16, 64)
		if err != nil {
			return nil, err
		}
		tcpStats[tcpConnectionState(tcpRxQueuedBytes)] += float64(rx)

		st, err := strconv.ParseInt(parts[3], 16, 8)
		if err != nil {
			return nil, err
		}

		tcpStats[tcpConnectionState(st)]++

	}

	return tcpStats, nil
}

func (st tcpConnectionState) String() string {
	switch st {
	case tcpEstablished:
		return "established"
	case tcpSynSent:
		return "syn_sent"
	case tcpSynRecv:
		return "syn_recv"
	case tcpFinWait1:
		return "fin_wait1"
	case tcpFinWait2:
		return "fin_wait2"
	case tcpTimeWait:
		return "time_wait"
	case tcpClose:
		return "close"
	case tcpCloseWait:
		return "close_wait"
	case tcpLastAck:
		return "last_ack"
	case tcpListen:
		return "listen"
	case tcpClosing:
		return "closing"
	case tcpRxQueuedBytes:
		return "rx_queued_bytes"
	case tcpTxQueuedBytes:
		return "tx_queued_bytes"
	default:
		return "unknown"
	}
}
