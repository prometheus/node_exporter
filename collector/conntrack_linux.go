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

// +build !noconntrack

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type conntrackCollector struct {
	acct,
	buckets,
	checksum,
	count,
	events,
	expectMax,
	frag6HighTresh,
	frag6LowTresh,
	frag6Timeout,
	genericTimeout,
	helper,
	icmpTimeout,
	icmp6Timeout,
	logInvalid,
	max,
	tcpBeLiberal,
	tcpLoose,
	tcpMaxRestrans,
	tcpTimeoutClose,
	tcpTimeoutCloseWait,
	tcpTimeoutEstablished,
	tcpTimeoutFinWait,
	tcpTimeoutLastAck,
	tcpTimeoutMaxRetrans,
	tcpTimeoutSynRecv,
	tcpTimeoutSynSent,
	tcpTimeoutTimeWait,
	tcpTimeoutUnacknowledged,
	timestamp,
	udpTimeout,
	udpTimeoutStream *prometheus.Desc
}

func init() {
	registerCollector("conntrack", defaultDisabled, NewConntrackCollector)
}

// NewConntrackCollector returns a new Collector exposing conntrack stats.
func NewConntrackCollector() (Collector, error) {
	//a lot of the description have been taken from here (and modified for some): https://www.kernel.org/doc/Documentation/networking/nf_conntrack-sysctl.txt
	return &conntrackCollector{
		acct:           buildDesc("nf_conntrack_acct", "Are new connections counted (0 = disabled)"),
		buckets:        buildDesc("nf_conntrack_buckets", "Size of the hash table"),
		checksum:       buildDesc("nf_conntrack_checksum", "Is the checksum of incoming packets verified (0 = disabled)"),
		count:          buildDesc("nf_conntrack_entries", "Number of currently allocated flow entries"),
		events:         buildDesc("nf_conntrack_events", "Is the connection tracking code providing userspace with connection tracking events via ctnetlink"),
		expectMax:      buildDesc("nf_conntrack_expect_max", "Maximum size of expectation table."),
		max:            buildDesc("nf_conntrack_max", "Size of connection tracking table"),
		helper:         buildDesc("nf_conntrack_helper", "Enable automatic conntrack helper assignment (0 = disabled)"),
		timestamp:      buildDesc("nf_conntrack_timestamp", "Is connection tracking flow timestamping enabled"),
		logInvalid:     buildDesc("nf_conntrack_log_invalid", "Log invalid packets of a type specified by value (0 = disabled; 1 = ICMP; 6 = TCP; 17 = UDP; 33 = DCCP; 41 = ICMPv6: 136 = UDPLITE; 255 = any)"),
		genericTimeout: buildDesc("nf_conntrack_generic_timeout", "Default for generic timeout"),

		frag6HighTresh: buildDesc("nf_conntrack_frag6_high_thresh", "Maximum memory used to reassemble IPv6 fragments"),
		frag6LowTresh:  buildDesc("nf_conntrack_frag6_low_thresh", "Amount of memory the fragment handler will go back to when the high treshold is reached"),
		frag6Timeout:   buildDesc("nf_conntrack_frag6_timeout", "Time to keep an IPv6 fragment in memory"),

		icmpTimeout:  buildDesc("nf_conntrack_icmp_timeout", "Default for ICMP timeout"),
		icmp6Timeout: buildDesc("nf_conntrack_icmpv6_timeout", "Default for ICMP6 timeout"),

		udpTimeout: buildDesc("nf_conntrack_udp_timeout", "Timeout for UDP connections"),
		udpTimeoutStream: buildDesc("nf_conntrack_udp_timeout_stream", "The extended timeout that will be used in case there is an UDP stream	detected"),

		tcpBeLiberal:             buildDesc("nf_conntrack_tcp_be_liberal", "If it's non-zero, we mark only out of window RST segments as INVALID"),
		tcpLoose:                 buildDesc("nf_conntrack_tcp_loose", "If it is set to zero, we disable picking up already established connections"),
		tcpMaxRestrans:           buildDesc("nf_conntrack_tcp_max_retrans", "Maximum number of packets that can be retransmitted without received an (acceptable) ACK from the destination"),
		tcpTimeoutClose:          buildDesc("nf_conntrack_tcp_timeout_close", "Timeout, in seconds, for closing TCP connections"),
		tcpTimeoutCloseWait:      buildDesc("nf_conntrack_tcp_timeout_close_wait", "Timeout, in seconds, for waiting during the closing of a TCP connection"),
		tcpTimeoutEstablished:    buildDesc("nf_conntrack_tcp_timeout_established", "Timeout, in seconds, for established TCP connections"),
		tcpTimeoutFinWait:        buildDesc("nf_conntrack_tcp_timeout_fin_wait", "Timeout, in seconds, for FIN of a TCP connections"),
		tcpTimeoutLastAck:        buildDesc("nf_conntrack_tcp_timeout_last_ack", "Timeout, in seconds, for last ACK of a TCP connection"),
		tcpTimeoutMaxRetrans:     buildDesc("nf_conntrack_tcp_timeout_max_retrans", "Timeout, in seconds, for the maximum retransmission of a TCP connection"),
		tcpTimeoutSynRecv:        buildDesc("nf_conntrack_tcp_timeout_syn_recv", "Timeout, in seconds, for receiving the syn of a TCP connection"),
		tcpTimeoutSynSent:        buildDesc("nf_conntrack_tcp_timeout_syn_sent", "Timeout in seconds"),
		tcpTimeoutTimeWait:       buildDesc("nf_conntrack_tcp_timeout_time_wait", "Timeout in seconds"),
		tcpTimeoutUnacknowledged: buildDesc("nf_conntrack_tcp_timeout_unacknowledged", "Timeout in seconds"),
	}, nil
}

func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	if err := readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_count", c.count, prometheus.GaugeValue); err != nil {
		// Conntrack probably not loaded into the kernel.
		return nil
	}
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_acct", c.acct, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_buckets", c.buckets, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_checksum", c.checksum, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_events", c.events, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_expect_max", c.expectMax, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_frag6_high_thresh", c.frag6HighTresh, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_frag6_low_thresh", c.frag6LowTresh, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_frag6_timeout", c.frag6Timeout, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_generic_timeout", c.genericTimeout, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_max", c.max, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_helper", c.helper, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_log_invalid", c.logInvalid, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_timestamp", c.timestamp, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_icmp_timeout", c.icmpTimeout, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_icmpv6_timeout", c.icmp6Timeout, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_udp_timeout", c.udpTimeout, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_udp_timeout_stream", c.udpTimeoutStream, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_be_liberal", c.tcpBeLiberal, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_loose", c.tcpLoose, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_max_retrans", c.tcpMaxRestrans, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_close", c.tcpTimeoutClose, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_close_wait", c.tcpTimeoutCloseWait, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_established", c.tcpTimeoutEstablished, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_fin_wait", c.tcpTimeoutFinWait, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_last_ack", c.tcpTimeoutLastAck, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_max_retrans", c.tcpTimeoutMaxRetrans, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_syn_recv", c.tcpTimeoutSynRecv, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_syn_sent", c.tcpTimeoutSynSent, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_time_wait", c.tcpTimeoutTimeWait, prometheus.GaugeValue)
	_ = readAndSendValue(ch, "sys/net/netfilter/nf_conntrack_tcp_timeout_unacknowledged", c.tcpTimeoutUnacknowledged, prometheus.GaugeValue)
	return nil
}

func buildDesc(name, description string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", name),
		description,
		nil, nil,
	)
}

func readAndSendValue(ch chan<- prometheus.Metric, file string, desc *prometheus.Desc, valueType prometheus.ValueType) error {
	value, err := readUintFromFile(procFilePath(file))
	if err == nil {
		ch <- prometheus.MustNewConstMetric(desc, valueType, float64(value))
	} else {
		log.Warn(fmt.Sprintf("a problem occured while reading the file: %s", err))
	}
	return err
}
