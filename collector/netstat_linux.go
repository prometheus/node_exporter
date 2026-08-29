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

// Regenerate the explicit descriptor table after a procfs upgrade:
//go:generate go run netstat_descs_gen.go

package collector

import (
	"fmt"
	"log/slog"
	"reflect"
	"regexp"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	netStatsSubsystem = "netstat"
)

var (
	netStatFields = kingpin.Flag("collector.netstat.fields", "Regexp of fields to return for netstat collector.").Default("^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$").String()

	// netStatDescs holds the explicit metric descriptors, allocated once at
	// package init (see netstat_descs_linux.go) so collection never calls
	// prometheus.NewDesc per scrape.
	netStatDescs = netStatMetricDescs()
)

type netStatCollector struct {
	proc         procfs.Proc
	fieldPattern *regexp.Regexp
	logger       *slog.Logger
}

func init() {
	registerCollector("netstat", defaultEnabled, NewNetStatCollector)
}

// NewNetStatCollector takes and returns
// a new Collector exposing network stats.
func NewNetStatCollector(logger *slog.Logger) (Collector, error) {
	pattern := regexp.MustCompile(*netStatFields)
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	// Network statistics in /proc/net are network namespace local. Reading
	// them via the current process' /proc/self/net keeps the same semantics
	// while allowing the use of the procfs parsers.
	proc, err := fs.Self()
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/self: %w", err)
	}

	return &netStatCollector{
		proc:         proc,
		fieldPattern: pattern,
		logger:       logger,
	}, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	netStats, err := c.proc.Netstat()
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %w", err)
	}
	snmpStats, err := c.proc.Snmp()
	if err != nil {
		return fmt.Errorf("couldn't get SNMP stats: %w", err)
	}
	snmp6Stats, err := c.proc.Snmp6()
	if err != nil {
		return fmt.Errorf("couldn't get SNMP6 stats: %w", err)
	}

	c.emitStruct(ch, netStats.TcpExt)
	c.emitStruct(ch, netStats.IpExt)
	c.emitStruct(ch, snmpStats.Ip)
	c.emitStruct(ch, snmpStats.Icmp)
	c.emitStruct(ch, snmpStats.IcmpMsg)
	c.emitStruct(ch, snmpStats.Tcp)
	c.emitStruct(ch, snmpStats.Udp)
	c.emitStruct(ch, snmpStats.UdpLite)
	c.emitStruct(ch, snmp6Stats.Ip6)
	c.emitStruct(ch, snmp6Stats.Icmp6)
	c.emitStruct(ch, snmp6Stats.Udp6)
	c.emitStruct(ch, snmp6Stats.UdpLite6)

	return nil
}

// emitStruct emits one metric per non-nil field of a procfs netstat/snmp
// statistics struct, using the struct's type name as the protocol name and
// looking up the pre-built descriptor by "<protocol>_<field>".
func (c *netStatCollector) emitStruct(ch chan<- prometheus.Metric, stats any) {
	v := reflect.ValueOf(stats)
	protocol := v.Type().Name()

	for i := 0; i < v.NumField(); i++ {
		value, ok := v.Field(i).Interface().(*float64)
		if !ok || value == nil {
			continue
		}

		name := v.Type().Field(i).Name
		key := protocol + "_" + name
		desc, ok := netStatDescs[key]
		if !ok {
			continue
		}
		if !c.fieldPattern.MatchString(key) {
			continue
		}

		ch <- prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, *value)
	}
}
