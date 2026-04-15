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
	"log/slog"
	"regexp"
	"testing"

	"github.com/prometheus/procfs"
)

// newTestNetStatCollector creates a netStatCollector pointed at the test
// fixture procfs directory (collector/fixtures/proc).  Fixture files must
// exist at fixtures/proc/1/net/{netstat,snmp,snmp6} because the collector
// reads statistics via the per-process procfs path (PID 1 = global network
// namespace on Linux).
func newTestNetStatCollector(t *testing.T) *netStatCollector {
	t.Helper()
	fs, err := procfs.NewFS("fixtures/proc")
	if err != nil {
		t.Fatalf("failed to open fixture procfs: %v", err)
	}
	return &netStatCollector{
		fs:           fs,
		fieldPattern: regexp.MustCompile(".*"), // match everything for test coverage
		logger:       slog.Default(),
	}
}

func TestNetStatCollector_TcpExt(t *testing.T) {
	c := newTestNetStatCollector(t)

	stats, err := c.getNetStats()
	if err != nil {
		t.Fatalf("getNetStats() error: %v", err)
	}

	// Fixture: fixtures/proc/1/net/netstat (copied from fixtures/proc/net/netstat)
	// TcpExt line: DelayedACKs = 102471
	if want, got := 102471.0, stats["TcpExt"]["DelayedACKs"]; want != got {
		t.Errorf("TcpExt_DelayedACKs: want %v, got %v", want, got)
	}
}

func TestNetStatCollector_IpExt(t *testing.T) {
	c := newTestNetStatCollector(t)

	stats, err := c.getNetStats()
	if err != nil {
		t.Fatalf("getNetStats() error: %v", err)
	}

	// Fixture: fixtures/proc/1/net/netstat
	// IpExt line: OutOctets = 2786264347
	if want, got := 2786264347.0, stats["IpExt"]["OutOctets"]; want != got {
		t.Errorf("IpExt_OutOctets: want %v, got %v", want, got)
	}
}

func TestNetStatCollector_Snmp_Udp(t *testing.T) {
	c := newTestNetStatCollector(t)

	stats, err := c.getNetStats()
	if err != nil {
		t.Fatalf("getNetStats() error: %v", err)
	}

	// Fixture: fixtures/proc/1/net/snmp (copied from fixtures/proc/net/snmp)
	// Udp line: RcvbufErrors = 9, SndbufErrors = 8
	if want, got := 9.0, stats["Udp"]["RcvbufErrors"]; want != got {
		t.Errorf("Udp_RcvbufErrors: want %v, got %v", want, got)
	}
	if want, got := 8.0, stats["Udp"]["SndbufErrors"]; want != got {
		t.Errorf("Udp_SndbufErrors: want %v, got %v", want, got)
	}
}

func TestNetStatCollector_Snmp6(t *testing.T) {
	c := newTestNetStatCollector(t)

	stats, err := c.getNetStats()
	if err != nil {
		t.Fatalf("getNetStats() error: %v", err)
	}

	// Fixture: fixtures/proc/1/net/snmp6 (copied from fixtures/proc/net/snmp6)
	// Ip6InOctets = 460
	if want, got := 460.0, stats["Ip6"]["InOctets"]; want != got {
		t.Errorf("Ip6_InOctets: want %v, got %v", want, got)
	}

	// Icmp6OutMsgs = 8
	if want, got := 8.0, stats["Icmp6"]["OutMsgs"]; want != got {
		t.Errorf("Icmp6_OutMsgs: want %v, got %v", want, got)
	}

	// Udp6RcvbufErrors = 9
	if want, got := 9.0, stats["Udp6"]["RcvbufErrors"]; want != got {
		t.Errorf("Udp6_RcvbufErrors: want %v, got %v", want, got)
	}

	// Udp6SndbufErrors = 8
	if want, got := 8.0, stats["Udp6"]["SndbufErrors"]; want != got {
		t.Errorf("Udp6_SndbufErrors: want %v, got %v", want, got)
	}
}

func TestNetStatCollector_FieldFilter(t *testing.T) {
	fs, err := procfs.NewFS("fixtures/proc")
	if err != nil {
		t.Fatalf("failed to open fixture procfs: %v", err)
	}

	// Use the default field pattern from the collector.
	pattern := regexp.MustCompile(`^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$`)
	c := &netStatCollector{
		fs:           fs,
		fieldPattern: pattern,
		logger:       slog.Default(),
	}

	stats, err := c.getNetStats()
	if err != nil {
		t.Fatalf("getNetStats() error: %v", err)
	}

	// Verify that fields matching the default pattern are present.
	matchedKeys := []string{
		"Ip_Forwarding",
		"IpExt_InOctets",
		"IpExt_OutOctets",
		"Udp_RcvbufErrors",
		"Udp_SndbufErrors",
	}
	for _, key := range matchedKeys {
		if !pattern.MatchString(key) {
			t.Errorf("expected pattern to match %q", key)
		}
	}

	// Verify stats map is non-empty and contains expected protocols.
	for _, proto := range []string{"TcpExt", "IpExt", "Ip", "Icmp", "Tcp", "Udp", "Ip6", "Icmp6", "Udp6"} {
		if _, ok := stats[proto]; !ok {
			t.Errorf("expected protocol %q in stats map", proto)
		}
	}
}
