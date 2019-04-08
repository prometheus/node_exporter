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
		acct:           buildDesc("acct", "Are new connections counted (0 = disabled)"),
		buckets:        buildDesc("buckets", "Size of the hash table"),
		checksum:       buildDesc("checksum", "Is the checksum of incoming packets verified (0 = disabled)"),
		count:          buildDesc("entries", "Number of currently allocated flow entries"),
		events:         buildDesc("events", "Is the connection tracking code providing userspace with connection tracking events via ctnetlink"),
		expectMax:      buildDesc("expect_max", "Maximum size of expectation table."),
		max:            buildDesc("max", "Size of connection tracking table"),
		helper:         buildDesc("helper", "Enable automatic conntrack helper assignment (0 = disabled)"),
		timestamp:      buildDesc("timestamp", "Is connection tracking flow timestamping enabled"),
		logInvalid:     buildDesc("log_invalid", "Log invalid packets of a type specified by value (0 = disabled; 1 = ICMP; 6 = TCP; 17 = UDP; 33 = DCCP; 41 = ICMPv6: 136 = UDPLITE; 255 = any)"),
		genericTimeout: buildDesc("generic_timeout", "Default for generic timeout"),

		frag6HighTresh: buildDesc("frag6_high_thresh", "Maximum memory used to reassemble IPv6 fragments"),
		frag6LowTresh:  buildDesc("frag6_low_thresh", "Amount of memory the fragment handler will go back to when the high threshold is reached"),
		frag6Timeout:   buildDesc("frag6_timeout", "Time to keep an IPv6 fragment in memory"),

		icmpTimeout:  buildDesc("icmp_timeout", "Default for ICMP timeout"),
		icmp6Timeout: buildDesc("icmpv6_timeout", "Default for ICMP6 timeout"),

		udpTimeout: buildDesc("udp_timeout", "Timeout for UDP connections"),
		udpTimeoutStream: buildDesc("udp_timeout_stream", "The extended timeout that will be used in case there is an UDP stream	detected"),

		tcpBeLiberal:             buildDesc("tcp_be_liberal", "If it's non-zero, we mark only out of window RST segments as INVALID"),
		tcpLoose:                 buildDesc("tcp_loose", "If it is set to zero, we disable picking up already established connections"),
		tcpMaxRestrans:           buildDesc("tcp_max_retrans", "Maximum number of packets that can be retransmitted without received an (acceptable) ACK from the destination"),
		tcpTimeoutClose:          buildDesc("tcp_timeout_close", "Timeout, in seconds, for closing TCP connections"),
		tcpTimeoutCloseWait:      buildDesc("tcp_timeout_close_wait", "Timeout, in seconds, for waiting during the closing of a TCP connection"),
		tcpTimeoutEstablished:    buildDesc("tcp_timeout_established", "Timeout, in seconds, for established TCP connections"),
		tcpTimeoutFinWait:        buildDesc("tcp_timeout_fin_wait", "Timeout, in seconds, for FIN of a TCP connections"),
		tcpTimeoutLastAck:        buildDesc("tcp_timeout_last_ack", "Timeout, in seconds, for last ACK of a TCP connection"),
		tcpTimeoutMaxRetrans:     buildDesc("tcp_timeout_max_retrans", "Timeout, in seconds, for the maximum retransmission of a TCP connection"),
		tcpTimeoutSynRecv:        buildDesc("tcp_timeout_syn_recv", "Timeout, in seconds, for receiving the syn of a TCP connection"),
		tcpTimeoutSynSent:        buildDesc("tcp_timeout_syn_sent", "Timeout in seconds"),
		tcpTimeoutTimeWait:       buildDesc("tcp_timeout_time_wait", "Timeout in seconds"),
		tcpTimeoutUnacknowledged: buildDesc("tcp_timeout_unacknowledged", "Timeout in seconds"),
	}, nil
}

func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	for file, desc := range map[string]*prometheus.Desc{
		"count":                      c.count,
		"acct":                       c.acct,
		"buckets":                    c.buckets,
		"checksum":                   c.checksum,
		"events":                     c.events,
		"expect_max":                 c.expectMax,
		"frag6_high_thresh":          c.frag6HighTresh,
		"frag6_low_thresh":           c.frag6LowTresh,
		"frag6_timeout":              c.frag6Timeout,
		"generic_timeout":            c.genericTimeout,
		"max":                        c.max,
		"helper":                     c.helper,
		"log_invalid":                c.logInvalid,
		"timestamp":                  c.timestamp,
		"icmp_timeout":               c.icmpTimeout,
		"icmpv6_timeout":             c.icmp6Timeout,
		"udp_timeout":                c.udpTimeout,
		"udp_timeout_stream":         c.udpTimeoutStream,
		"tcp_be_liberal":             c.tcpBeLiberal,
		"tcp_loose":                  c.tcpLoose,
		"tcp_max_retrans":            c.tcpMaxRestrans,
		"tcp_timeout_close":          c.tcpTimeoutClose,
		"tcp_timeout_close_wait":     c.tcpTimeoutCloseWait,
		"tcp_timeout_established":    c.tcpTimeoutEstablished,
		"tcp_timeout_fin_wait":       c.tcpTimeoutFinWait,
		"tcp_timeout_last_ack":       c.tcpTimeoutLastAck,
		"tcp_timeout_max_retrans":    c.tcpTimeoutMaxRetrans,
		"tcp_timeout_syn_recv":       c.tcpTimeoutSynRecv,
		"tcp_timeout_syn_sent":       c.tcpTimeoutSynSent,
		"tcp_timeout_time_wait":      c.tcpTimeoutTimeWait,
		"tcp_timeout_unacknowledged": c.tcpTimeoutUnacknowledged,
	} {
		completePath := procFilePath("sys/net/netfilter/nf_conntrack_" + file)
		log.Debugf("reading from file %s", completePath)
		value, err := readUintFromFile(completePath)
		if err == nil {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(value))
		} else {
			log.Warn(fmt.Sprintf("an error (%s) occurred while reading the file %s", err.Error(), completePath))
		}
	}
	return nil
}

func buildDesc(name, description string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "netfilter", name),
		description,
		nil, nil,
	)
}
