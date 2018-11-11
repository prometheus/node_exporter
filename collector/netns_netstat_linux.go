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

// +build !nonetstat

package collector

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vishvananda/netns"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	netNsNetStatsSubsystem = "netns_netstat"
)

var (
	netNsNetStatFields = kingpin.Flag("collector.netns_netstat.fields", "Regexp of fields to return for netns_netstat collector.").Default("^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans)|Tcp_(ActiveOpens|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts))$").String()
)

type netNsNetStatCollector struct {
	fieldPattern *regexp.Regexp
}

func init() {
	registerCollector("netns_netstat", defaultDisabled, NewNetNsNetStatCollector)
}

func NewNetNsNetStatCollector() (Collector, error) {
	pattern := regexp.MustCompile(*netNsNetStatFields)
	return &netNsNetStatCollector{
		fieldPattern: pattern,
	}, nil
}

func (c *netNsNetStatCollector) Update(ch chan<- prometheus.Metric) error {
	f, err := os.Open("/var/run/netns")
	if err != nil {
		log.Debugf("couldn't open /var/run/netns (no ns): %s", err)
		return nil
	}
	fInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return fmt.Errorf("couldn't list dir /var/run/netns: %s", err)
	}

	for _, file := range fInfo {
		ns, err := netns.GetFromName(file.Name())
		if err != nil {
			return fmt.Errorf("couldn't get netns %s: %s", file.Name(), err)
		}
		err = netns.Set(ns)
		if err != nil {
			return fmt.Errorf("couldn't enter netns %s: %s", file.Name(), err)
		}
		hostname, _ := os.Hostname()
		if err != nil {
			hostname = ""
		}
		err = c.UpdateNetstat(ch, file.Name(), hostname)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *netNsNetStatCollector) UpdateNetstat(ch chan<- prometheus.Metric, ns string, hostname string) error {
	netStats, err := getNetStats(procFilePath("net/netstat"))
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	snmpStats, err := getNetStats(procFilePath("net/snmp"))
	if err != nil {
		return fmt.Errorf("couldn't get SNMP stats: %s", err)
	}
	snmp6Stats, err := getSNMP6Stats(procFilePath("net/snmp6"))
	if err != nil {
		return fmt.Errorf("couldn't get SNMP6 stats: %s", err)
	}
	// Merge the results of snmpStats into netStats (collisions are possible, but
	// we know that the keys are always unique for the given use case).
	for k, v := range snmpStats {
		netStats[k] = v
	}
	for k, v := range snmp6Stats {
		netStats[k] = v
	}
	for protocol, protocolStats := range netStats {
		for name, value := range protocolStats {
			key := protocol + "_" + name
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			if !c.fieldPattern.MatchString(key) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, netNsNetStatsSubsystem, key),
					fmt.Sprintf("Statistic %s.", protocol+name),
					[]string{"netns_ns", "netns_hostname"},
					nil,
				),
				prometheus.UntypedValue, v,
				ns, hostname,
			)
		}
	}
	return nil
}
