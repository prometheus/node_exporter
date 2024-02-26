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
// +build !notcpstat

package collector

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"testing"

	"github.com/josharian/native"
	"github.com/mdlayher/netlink"
)

func Test_parseTCPStats(t *testing.T) {
	encode := func(m InetDiagMsg) []byte {
		var buf bytes.Buffer
		err := binary.Write(&buf, native.Endian, m)
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}

	msg := []netlink.Message{
		{
			Data: encode(InetDiagMsg{
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
			Data: encode(InetDiagMsg{
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

	tcpStats, err := parseTCPStats(msg)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 1, int(tcpStats[tcpEstablished]); want != got {
		t.Errorf("want tcpstat number of established state %d, got %d", want, got)
	}

	if want, got := 1, int(tcpStats[tcpListen]); want != got {
		t.Errorf("want tcpstat number of listen state %d, got %d", want, got)
	}

	if want, got := 42, int(tcpStats[tcpTxQueuedBytes]); want != got {
		t.Errorf("want tcpstat number of bytes in tx queue %d, got %d", want, got)
	}
	if want, got := 22, int(tcpStats[tcpRxQueuedBytes]); want != got {
		t.Errorf("want tcpstat number of bytes in rx queue %d, got %d", want, got)
	}
}

func Test_parseTCPStatsPerPort(t *testing.T) {
	encode := func(m InetDiagMsg) []byte {
		var buf bytes.Buffer
		err := binary.Write(&buf, native.Endian, m)
		if err != nil {
			panic(err)
		}
		return buf.Bytes()
	}

	msg := []netlink.Message{
		{
			Data: encode(InetDiagMsg{
				Family: syscall.AF_INET,
				State:  uint8(tcpEstablished),
				ID: InetDiagSockID{
					DestPort: [2]byte{0, 22},
				},
			}),
		},
		{
			Data: encode(InetDiagMsg{
				Family: syscall.AF_INET,
				State:  uint8(tcpEstablished),
				ID: InetDiagSockID{
					DestPort:   [2]byte{0, 22},
					SourcePort: [2]byte{0, 23},
				},
			}),
		},
		{
			Data: encode(InetDiagMsg{
				Family: syscall.AF_INET6,
				State:  uint8(tcpEstablished),
				ID: InetDiagSockID{
					SourcePort: [2]byte{0, 23},
				},
			}),
		},
	}

	tcpStatsPerDestPort, err := parseTCPStatsPerDestPort(msg, []string{"22"})
	if err != nil {
		t.Fatal(err)
	}

	sp1 := statePortPair{
		state: tcpEstablished,
		port:  22,
	}

	sp2 := statePortPair{
		state: tcpEstablished,
		port:  23,
	}

	tcpStatsPerSourcePort, err := parseTCPStatsPerSourcePort(msg, []string{"23"})
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 2, int(tcpStatsPerDestPort[sp1]); want != got {
		t.Errorf("tcpstat connection per %s states per dest port %d. want %d, got %d", sp1.state.String(), sp1.port, want, got)
	}

	if want, got := 2, int(tcpStatsPerSourcePort[sp2]); want != got {
		t.Errorf("tcpstat connection per %s states per source port %d. want %d, got %d", sp2.state.String(), sp2.port, want, got)
	}
}
