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
)

type netStatCollector struct {
	fs           procfs.FS
	fieldPattern *regexp.Regexp
	logger       *slog.Logger
}

func init() {
	registerCollector("netstat", defaultEnabled, NewNetStatCollector)
}

// NewNetStatCollector takes and returns
// a new Collector exposing network stats.
func NewNetStatCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	pattern := regexp.MustCompile(*netStatFields)
	return &netStatCollector{
		fs:           fs,
		fieldPattern: pattern,
		logger:       logger,
	}, nil
}

func (c *netStatCollector) Update(ch chan<- prometheus.Metric) error {
	netStats, err := c.getNetStats()
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %w", err)
	}

	for protocol, protocolStats := range netStats {
		for name, value := range protocolStats {
			key := protocol + "_" + name
			if !c.fieldPattern.MatchString(key) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, netStatsSubsystem, key),
					fmt.Sprintf("Statistic %s.", protocol+name),
					nil, nil,
				),
				prometheus.UntypedValue, value,
			)
		}
	}
	return nil
}

// getNetStats reads network statistics from the procfs using typed structs
// (procfs.ProcNetstat, procfs.ProcSnmp, procfs.ProcSnmp6) and returns a
// unified map of protocol → field → float64 value.
//
// On Linux /proc/net/netstat, /proc/net/snmp, and /proc/net/snmp6 are
// namespaced: they present the same data as /proc/1/net/... for the global
// (init) network namespace. The procfs library reads these via the
// per-process path; using PID 1 (init/systemd) is the standard way to
// access global network namespace statistics.
func (c *netStatCollector) getNetStats() (map[string]map[string]float64, error) {
	result := make(map[string]map[string]float64)

	proc, err := c.fs.Proc(1)
	if err != nil {
		return nil, fmt.Errorf("failed to open proc(1): %w", err)
	}

	// /proc/1/net/netstat → TcpExt and IpExt sections.
	netstat, err := proc.Netstat()
	if err != nil {
		return nil, fmt.Errorf("couldn't read netstat: %w", err)
	}
	addStructFields(result, "TcpExt", netstat.TcpExt)
	addStructFields(result, "IpExt", netstat.IpExt)

	// /proc/1/net/snmp → Ip, Icmp, IcmpMsg, Tcp, Udp, UdpLite sections.
	snmp, err := proc.Snmp()
	if err != nil {
		return nil, fmt.Errorf("couldn't read snmp: %w", err)
	}
	addStructFields(result, "Ip", snmp.Ip)
	addStructFields(result, "Icmp", snmp.Icmp)
	addStructFields(result, "IcmpMsg", snmp.IcmpMsg)
	addStructFields(result, "Tcp", snmp.Tcp)
	addStructFields(result, "Udp", snmp.Udp)
	addStructFields(result, "UdpLite", snmp.UdpLite)

	// /proc/1/net/snmp6 → Ip6, Icmp6, Udp6, UdpLite6 sections.
	// This file may not exist on systems with IPv6 disabled; absence is
	// treated as an empty result rather than an error (matching the prior
	// getSNMP6Stats behaviour).
	snmp6, err := proc.Snmp6()
	if err != nil {
		return nil, fmt.Errorf("couldn't read snmp6: %w", err)
	}
	addStructFields(result, "Ip6", snmp6.Ip6)
	addStructFields(result, "Icmp6", snmp6.Icmp6)
	addStructFields(result, "Udp6", snmp6.Udp6)
	addStructFields(result, "UdpLite6", snmp6.UdpLite6)

	return result, nil
}

// addStructFields uses reflection to iterate the exported *float64 fields of
// a procfs typed struct (e.g. procfs.TcpExt, procfs.Ip, procfs.Ip6) and
// stores non-nil values in result[protocol].
func addStructFields(result map[string]map[string]float64, protocol string, s interface{}) {
	if _, ok := result[protocol]; !ok {
		result[protocol] = make(map[string]float64)
	}

	v := reflect.ValueOf(s)
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fv := v.Field(i)
		// All stat fields in procfs network structs are *float64 pointers.
		if fv.Kind() == reflect.Ptr && !fv.IsNil() {
			result[protocol][field.Name] = fv.Elem().Float()
		}
	}
}
