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

//go:build !nonetstat

package collector

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNetStats(t *testing.T) {
	testNetStats(t, "fixtures/proc/net/netstat")
	testSNMPStats(t, "fixtures/proc/net/snmp")
	testSNMP6Stats(t, "fixtures/proc/net/snmp6")
}

func testNetStats(t *testing.T, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	netStats, err := parseNetStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "102471", netStats["TcpExt"]["DelayedACKs"]; want != got {
		t.Errorf("want netstat TCP DelayedACKs %s, got %s", want, got)
	}

	if want, got := "2786264347", netStats["IpExt"]["OutOctets"]; want != got {
		t.Errorf("want netstat IP OutOctets %s, got %s", want, got)
	}
}

func testSNMPStats(t *testing.T, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	snmpStats, err := parseNetStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "9", snmpStats["Udp"]["RcvbufErrors"]; want != got {
		t.Errorf("want netstat Udp RcvbufErrors %s, got %s", want, got)
	}

	if want, got := "8", snmpStats["Udp"]["SndbufErrors"]; want != got {
		t.Errorf("want netstat Udp SndbufErrors %s, got %s", want, got)
	}
}

func testSNMP6Stats(t *testing.T, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	snmp6Stats, err := parseSNMP6Stats(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "460", snmp6Stats["Ip6"]["InOctets"]; want != got {
		t.Errorf("want netstat IPv6 InOctets %s, got %s", want, got)
	}

	if want, got := "8", snmp6Stats["Icmp6"]["OutMsgs"]; want != got {
		t.Errorf("want netstat ICPM6 OutMsgs %s, got %s", want, got)
	}

	if want, got := "9", snmp6Stats["Udp6"]["RcvbufErrors"]; want != got {
		t.Errorf("want netstat Udp6 RcvbufErrors %s, got %s", want, got)
	}

	if want, got := "8", snmp6Stats["Udp6"]["SndbufErrors"]; want != got {
		t.Errorf("want netstat Udp6 SndbufErrors %s, got %s", want, got)
	}
}

func TestNetstatValueType(t *testing.T) {
	if got := netstatValueType("Tcp_CurrEstab"); got != prometheus.GaugeValue {
		t.Fatalf("Tcp_CurrEstab: got %v want gauge", got)
	}
	if got := netstatValueType("Ip_Forwarding"); got != prometheus.GaugeValue {
		t.Fatalf("Ip_Forwarding: got %v want gauge", got)
	}
	if got := netstatValueType("TcpExt_ListenDrops"); got != prometheus.CounterValue {
		t.Fatalf("TcpExt_ListenDrops: got %v want counter", got)
	}
}
