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

//go:build !notcpstat

package collector

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/alecthomas/kingpin/v2"
	"github.com/mdlayher/netlink"
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

var (
	tcpstatSourcePorts = kingpin.Flag("collector.tcpstat.port.source", "List of tcpstat source ports").Strings()
	tcpstatDestPorts   = kingpin.Flag("collector.tcpstat.port.dest", "List of tcpstat destination ports").Strings()
)

type tcpStatCollector struct {
	desc   typedDesc
	logger *slog.Logger
}

func init() {
	registerCollector("tcpstat", defaultDisabled, NewTCPStatCollector)
}

// NewTCPStatCollector returns a new Collector exposing network stats.
func NewTCPStatCollector(logger *slog.Logger) (Collector, error) {
	return &tcpStatCollector{
		desc: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "tcp", "connection_states"),
			"Number of connection states.",
			[]string{"state", "port", "direction"}, nil,
		), prometheus.GaugeValue},
		logger: logger,
	}, nil
}

// InetDiagSockID (inet_diag_sockid) contains the socket identity.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L13
type InetDiagSockID struct {
	SourcePort [2]byte
	DestPort   [2]byte
	SourceIP   [4][4]byte
	DestIP     [4][4]byte
	Interface  uint32
	Cookie     [2]uint32
}

// InetDiagReqV2 (inet_diag_req_v2) is used to request diagnostic data.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L37
type InetDiagReqV2 struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	ID       InetDiagSockID
}

const sizeOfDiagRequest = 0x38

func (req *InetDiagReqV2) Serialize() []byte {
	return (*(*[sizeOfDiagRequest]byte)(unsafe.Pointer(req)))[:]
}

func (req *InetDiagReqV2) Len() int {
	return sizeOfDiagRequest
}

type InetDiagMsg struct {
	Family  uint8
	State   uint8
	Timer   uint8
	Retrans uint8
	ID      InetDiagSockID
	Expires uint32
	RQueue  uint32
	WQueue  uint32
	UID     uint32
	Inode   uint32
}

func parseInetDiagMsg(b []byte) *InetDiagMsg {
	return (*InetDiagMsg)(unsafe.Pointer(&b[0]))
}

func (c *tcpStatCollector) Update(ch chan<- prometheus.Metric) error {
	messages, err := getMessagesFromSocket(syscall.AF_INET)
	if err != nil {
		return fmt.Errorf("couldn't get tcpstats: %w", err)
	}

	tcpStats, err := parseTCPStats(messages)
	if err != nil {
		return fmt.Errorf("couldn't parse tcpstats: %w", err)
	}

	if _, hasIPv6 := os.Stat(procFilePath("net/tcp6")); hasIPv6 == nil {
		messagesIPv6, err := getMessagesFromSocket(syscall.AF_INET6)
		if err != nil {
			return fmt.Errorf("couldn't get tcp6stats: %w", err)
		}

		tcp6Stats, err := parseTCPStats(messagesIPv6)
		if err != nil {
			return fmt.Errorf("couldn't parse tcp6stats: %w", err)
		}

		for st, value := range tcp6Stats {
			tcpStats[st] += value
		}

		messages = append(messages, messagesIPv6...)
	}

	emitTotalTCPStats(c, ch, tcpStats)
	emitTCPStatsPerPort(c, ch, messages, *tcpstatSourcePorts, "source", true)
	emitTCPStatsPerPort(c, ch, messages, *tcpstatDestPorts, "dest", false)

	return nil
}

func emitTotalTCPStats(c *tcpStatCollector, ch chan<- prometheus.Metric, stats map[tcpConnectionState]float64) {
	for st, value := range stats {
		ch <- c.desc.mustNewConstMetric(value, st.String(), "0", "total")
	}
}

func emitTCPStatsPerPort(
	c *tcpStatCollector,
	ch chan<- prometheus.Metric,
	messages []netlink.Message,
	ports []string,
	direction string,
	isSource bool,
) {
	if len(ports) == 0 {
		return
	}

	portSet := map[string]struct{}{}
	for _, p := range ports {
		portSet[p] = struct{}{}
	}

	counts := map[string]map[string]float64{}

	for _, m := range messages {
		msg := parseInetDiagMsg(m.Data)

		state := tcpConnectionState(msg.State).String()

		var rawPort uint16
		if isSource {
			rawPort = binary.BigEndian.Uint16(msg.ID.SourcePort[:])
		} else {
			rawPort = binary.BigEndian.Uint16(msg.ID.DestPort[:])
		}

		portStr := strconv.Itoa(int(rawPort))

		if _, ok := portSet[portStr]; ok {
			if _, ok := counts[state]; !ok {
				counts[state] = make(map[string]float64)
			}

			counts[state][portStr]++
		}
	}

	for state, portMap := range counts {
		for port, count := range portMap {
			ch <- c.desc.mustNewConstMetric(count, state, port, direction)
		}
	}
}

func getMessagesFromSocket(family uint8) ([]netlink.Message, error) {
	const TCPFAll = 0xFFF
	const InetDiagInfo = 2
	const SockDiagByFamily = 20

	conn, err := netlink.Dial(syscall.NETLINK_INET_DIAG, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect netlink: %w", err)
	}
	defer conn.Close()

	msg := netlink.Message{
		Header: netlink.Header{
			Type:  SockDiagByFamily,
			Flags: syscall.NLM_F_REQUEST | syscall.NLM_F_DUMP,
		},
		Data: (&InetDiagReqV2{
			Family:   family,
			Protocol: syscall.IPPROTO_TCP,
			States:   TCPFAll,
			Ext:      0 | 1<<(InetDiagInfo-1),
		}).Serialize(),
	}

	return conn.Execute(msg)
}

func parseTCPStats(msgs []netlink.Message) (map[tcpConnectionState]float64, error) {
	stats := make(map[tcpConnectionState]float64)

	for _, m := range msgs {
		msg := parseInetDiagMsg(m.Data)
		stats[tcpTxQueuedBytes] += float64(msg.WQueue)
		stats[tcpRxQueuedBytes] += float64(msg.RQueue)
		stats[tcpConnectionState(msg.State)]++
	}

	return stats, nil
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
