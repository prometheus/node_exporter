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
	"bytes"
	"encoding/binary"
	"slices"
	"syscall"
	"testing"

	"github.com/mdlayher/netlink"
)

func encodeDiagMsg(m InetDiagMsg) []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.NativeEndian, m); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Test_parseTCPStats(t *testing.T) {
	msg := []netlink.Message{
		{
			Data: encodeDiagMsg(InetDiagMsg{
				Family:  syscall.AF_INET,
				State:   uint8(tcpEstablished),
				Timer:   0,
				Retrans: 0,
				ID:      InetDiagSockID{},
				Expires: 0,
				RQueue:  11,
				WQueue:  21,
				UID:     0,
				Inode:   0,
			}),
		},
		{
			Data: encodeDiagMsg(InetDiagMsg{
				Family:  syscall.AF_INET,
				State:   uint8(tcpListen),
				Timer:   0,
				Retrans: 0,
				ID:      InetDiagSockID{},
				Expires: 0,
				RQueue:  11,
				WQueue:  21,
				UID:     0,
				Inode:   0,
			}),
		},
	}

	tcpStats, err := parseTCPStats(msg, tcpPortFilter{})
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 1, int(tcpStats[tcpStatKey{state: tcpEstablished}]); want != got {
		t.Errorf("want tcpstat number of established state %d, got %d", want, got)
	}

	if want, got := 1, int(tcpStats[tcpStatKey{state: tcpListen}]); want != got {
		t.Errorf("want tcpstat number of listen state %d, got %d", want, got)
	}

	if want, got := 42, int(tcpStats[tcpStatKey{state: tcpTxQueuedBytes}]); want != got {
		t.Errorf("want tcpstat number of bytes in tx queue %d, got %d", want, got)
	}
	if want, got := 22, int(tcpStats[tcpStatKey{state: tcpRxQueuedBytes}]); want != got {
		t.Errorf("want tcpstat number of bytes in rx queue %d, got %d", want, got)
	}
}

func Test_parseTCPStatsPerPort(t *testing.T) {
	sock := func(state tcpConnectionState, sourcePort, destPort uint16, rQueue uint32) netlink.Message {
		var id InetDiagSockID
		binary.BigEndian.PutUint16(id.SourcePort[:], sourcePort)
		binary.BigEndian.PutUint16(id.DestPort[:], destPort)
		return netlink.Message{Data: encodeDiagMsg(InetDiagMsg{State: uint8(state), ID: id, RQueue: rQueue})}
	}

	msg := []netlink.Message{
		sock(tcpListen, 80, 0, 0),
		sock(tcpEstablished, 80, 50000, 5),
		sock(tcpEstablished, 40000, 443, 7),
		// Matches both the source and the dest filter, must be counted once.
		sock(tcpEstablished, 80, 443, 0),
		sock(tcpEstablished, 40001, 22, 3),
	}

	ports := tcpPortFilter{
		source: map[uint16]struct{}{80: {}},
		dest:   map[uint16]struct{}{443: {}},
	}

	tcpStats, err := parseTCPStats(msg, ports)
	if err != nil {
		t.Fatal(err)
	}

	want := map[tcpStatKey]float64{
		{tcpListen, "80", ""}:           1,
		{tcpEstablished, "80", ""}:      1,
		{tcpEstablished, "", "443"}:     1,
		{tcpEstablished, "80", "443"}:   1,
		{tcpEstablished, "", ""}:        1,
		{tcpRxQueuedBytes, "80", ""}:    5,
		{tcpRxQueuedBytes, "", "443"}:   7,
		{tcpRxQueuedBytes, "80", "443"}: 0,
		{tcpRxQueuedBytes, "", ""}:      3,
		{tcpTxQueuedBytes, "80", ""}:    0,
		{tcpTxQueuedBytes, "", "443"}:   0,
		{tcpTxQueuedBytes, "80", "443"}: 0,
		{tcpTxQueuedBytes, "", ""}:      0,
	}

	if len(tcpStats) != len(want) {
		t.Errorf("want %d series, got %d: %v", len(want), len(tcpStats), tcpStats)
	}
	for key, value := range want {
		if got := tcpStats[key]; got != value {
			t.Errorf("want %+v = %v, got %v", key, value, got)
		}
	}

	// Summing all connection state series must equal the number of sockets.
	var sockets float64
	for key, value := range tcpStats {
		if key.state != tcpRxQueuedBytes && key.state != tcpTxQueuedBytes {
			sockets += value
		}
	}
	if want, got := float64(len(msg)), sockets; want != got {
		t.Errorf("want sum of connection states %v, got %v", want, got)
	}
}

func Test_tcpStatLabels(t *testing.T) {
	key := tcpStatKey{tcpEstablished, "80", "443"}

	for _, tc := range []struct {
		name  string
		ports tcpPortFilter
		want  []string
	}{
		{"no filter", tcpPortFilter{}, []string{"established"}},
		{"source", tcpPortFilter{source: map[uint16]struct{}{}}, []string{"established", "80"}},
		{"dest", tcpPortFilter{dest: map[uint16]struct{}{}}, []string{"established", "443"}},
		{"both", tcpPortFilter{source: map[uint16]struct{}{}, dest: map[uint16]struct{}{}}, []string{"established", "80", "443"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := &tcpStatCollector{ports: tc.ports}
			if got := c.labelValues(key); !slices.Equal(got, tc.want) {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func Test_parsePortSet(t *testing.T) {
	set, err := parsePortSet([]string{"80", "443", "80"})
	if err != nil {
		t.Fatal(err)
	}
	if len(set) != 2 {
		t.Errorf("want 2 ports, got %d", len(set))
	}

	for _, p := range []string{"http", "-1", "65536", ""} {
		if _, err := parsePortSet([]string{p}); err == nil {
			t.Errorf("want error for port %q", p)
		}
	}
}
