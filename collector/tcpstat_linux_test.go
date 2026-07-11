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
	"syscall"
	"testing"

	"github.com/mdlayher/netlink"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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

	stats, err := parseTCPStats(msg)
	if err != nil {
		t.Fatal(err)
	}

	assertStat(t, stats, tcpEstablished, 1)
	assertStat(t, stats, tcpListen, 1)
	assertStat(t, stats, tcpTxQueuedBytes, 42)
	assertStat(t, stats, tcpRxQueuedBytes, 22)
}

func assertStat(t *testing.T, stats map[tcpConnectionState]float64, state tcpConnectionState, expected int) {
	t.Helper()
	if got := int(stats[state]); got != expected {
		t.Errorf("expected %s = %d, got %d", state.String(), expected, got)
	}
}

func Test_emitTCPStatsPerPort(t *testing.T) {
	msg := []netlink.Message{
		{
			Data: encodeDiagMsg(InetDiagMsg{
				State: uint8(tcpEstablished),
				ID:    InetDiagSockID{SourcePort: [2]byte{0, 80}},
			}),
		},
		{
			Data: encodeDiagMsg(InetDiagMsg{
				State: uint8(tcpListen),
				ID:    InetDiagSockID{DestPort: [2]byte{0, 123}},
			}),
		},
		{
			Data: encodeDiagMsg(InetDiagMsg{
				State: uint8(tcpTimeWait),
				ID:    InetDiagSockID{DestPort: [2]byte{0, 123}},
			}),
		},
	}

	var metrics []string

	collector := &tcpStatCollector{
		desc: typedDesc{
			desc:      prometheus.NewDesc("test_tcp_stat", "Test metric", []string{"state", "port", "direction"}, nil),
			valueType: prometheus.GaugeValue,
		},
	}

	ch := make(chan prometheus.Metric, 10)

	emitTCPStatsPerPort(collector, ch, msg, []string{"80"}, "source", true)
	emitTCPStatsPerPort(collector, ch, msg, []string{"123"}, "dest", false)

	close(ch)
	for m := range ch {
		d := &dto.Metric{}
		if err := m.Write(d); err != nil {
			t.Fatalf("failed to write metric: %v", err)
		}

		var state, port, direction string
		for _, label := range d.Label {
			switch label.GetName() {
			case "state":
				state = label.GetValue()
			case "port":
				port = label.GetValue()
			case "direction":
				direction = label.GetValue()
			}
		}

		metrics = append(metrics, state+"_"+port+"_"+direction)
	}

	expected := map[string]bool{
		"established_80_source": true,
		"listen_123_dest":       true,
		"time_wait_123_dest":    true,
	}

	for _, metric := range metrics {
		if !expected[metric] {
			t.Errorf("unexpected metric emitted: %s", metric)
		}
		delete(expected, metric)
	}

	for k := range expected {
		t.Errorf("expected metric missing: %s", k)
	}
}
