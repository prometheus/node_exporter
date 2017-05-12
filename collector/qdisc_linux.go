// Copyright 2017 The Prometheus Authors
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

// +build !noqdisc

package collector

/*
#include <stdlib.h>
#include <net/if.h>
*/
import "C"
import "unsafe"

import (
	"fmt"
	"math"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	TCA_UNSPEC = iota
	TCA_KIND
	TCA_OPTIONS
	TCA_STATS
	TCA_XSTATS
	TCA_RATE
	TCA_FCNT
	TCA_STATS2
	TCA_STAB
	__TCA_MAX
)

const (
	TCA_STATS_UNSPEC = iota
	TCA_STATS_BASIC
	TCA_STATS_RATE_EST
	TCA_STATS_QUEUE
	TCA_STATS_APP
	TCA_STATS_RATE_EST64
	__TCA_STATS_MAX
)

// See struct tc_stats in /usr/include/linux/pkt_sched.h
type TC_Stats struct {
	Bytes      uint64
	Packets    uint32
	Drops      uint32
	Overlimits uint32
	Bps        uint32
	Pps        uint32
	Qlen       uint32
	Backlog    uint32
}

// See /usr/include/linux/gen_stats.h
type TC_Stats2 struct {
	// struct gnet_stats_basic
	Bytes   uint64
	Packets uint32
	// struct gnet_stats_queue
	Qlen       uint32
	Backlog    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

type Metric struct {
	IfaceIdx   uint32
	Parent     uint32
	Handle     uint32
	Kind       string
	Bytes      uint64
	Packets    uint32
	Drops      uint32
	Requeues   uint32
	Overlimits uint32
}

// See if_indextoname(3)
func ifIndexToName(index uint32) string {
	var name string
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.if_indextoname(C.uint(index), cName)
	return C.GoString(ret)
}

func parseTCAStats(attr netlink.Attribute) TC_Stats {
	var stats TC_Stats
	stats.Bytes = nlenc.Uint64(attr.Data[0:8])
	stats.Packets = nlenc.Uint32(attr.Data[8:12])
	stats.Drops = nlenc.Uint32(attr.Data[12:16])
	stats.Overlimits = nlenc.Uint32(attr.Data[16:20])
	stats.Bps = nlenc.Uint32(attr.Data[20:24])
	stats.Pps = nlenc.Uint32(attr.Data[24:28])
	stats.Qlen = nlenc.Uint32(attr.Data[28:32])
	stats.Backlog = nlenc.Uint32(attr.Data[32:36])
	return stats
}

func parseTCAStats2(attr netlink.Attribute) TC_Stats2 {
	var stats TC_Stats2

	nested, _ := netlink.UnmarshalAttributes(attr.Data)

	for _, a := range nested {
		switch a.Type {
		case TCA_STATS_BASIC:
			stats.Bytes = nlenc.Uint64(a.Data[0:8])
			stats.Packets = nlenc.Uint32(a.Data[8:12])
		case TCA_STATS_QUEUE:
			stats.Qlen = nlenc.Uint32(a.Data[0:4])
			stats.Backlog = nlenc.Uint32(a.Data[4:8])
			stats.Drops = nlenc.Uint32(a.Data[8:12])
			stats.Requeues = nlenc.Uint32(a.Data[12:16])
			stats.Overlimits = nlenc.Uint32(a.Data[16:20])
		default:
			// TODO: TCA_STATS_APP
		}
	}

	return stats
}

func getQdiscMsgs() ([]netlink.Message, error) {
	const familyRoute = 0

	c, err := netlink.Dial(familyRoute, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial netlink: %v", err)
	}
	defer c.Close()

	req := netlink.Message{
		Header: netlink.Header{
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
			Type:  38, // RTM_GETQDISC
		},
		Data: []byte{0},
	}

	// Perform a request, receive replies, and validate the replies
	msgs, err := c.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	return msgs, nil
}

// See https://tools.ietf.org/html/rfc3549#section-3.1.3
func parseMessage(msg netlink.Message) (Metric, error) {
	var m Metric
	var s TC_Stats
	var s2 TC_Stats2

	/*
	   struct tcmsg {
	       unsigned char   tcm_family;
	       unsigned char   tcm__pad1;
	       unsigned short  tcm__pad2;
	       int     tcm_ifindex;
	       __u32       tcm_handle;
	       __u32       tcm_parent;
	       __u32       tcm_info;
	   };
	*/
	m.IfaceIdx = nlenc.Uint32(msg.Data[4:8])
	m.Handle = nlenc.Uint32(msg.Data[8:12])
	m.Parent = nlenc.Uint32(msg.Data[12:16])

	if m.Parent == math.MaxUint32 {
		m.Parent = 0
	}

	// The first 20 bytes are taken by tcmsg
	attrs, err := netlink.UnmarshalAttributes(msg.Data[20:])

	if err != nil {
		return m, fmt.Errorf("failed to unmarshal attributes: %v", err)
	}

	for _, attr := range attrs {
		switch attr.Type {
		case TCA_KIND:
			m.Kind = nlenc.String(attr.Data)
		case TCA_STATS2:
			s2 = parseTCAStats2(attr)
			m.Bytes = s2.Bytes
			m.Packets = s2.Packets
			m.Drops = s2.Drops
			// requeues only available in TCA_STATS2, not in TCA_STATS
			m.Requeues = s2.Requeues
			m.Overlimits = s2.Overlimits
		case TCA_STATS:
			// Legacy
			s = parseTCAStats(attr)
			m.Bytes = s.Bytes
			m.Packets = s.Packets
			m.Drops = s.Drops
			m.Overlimits = s.Overlimits
		default:
			// TODO: TCA_OPTIONS and TCA_XSTATS
		}
	}

	return m, nil
}

type qdiscStatCollector struct {
	bytes      typedDesc
	pkts       typedDesc
	drops      typedDesc
	requeues   typedDesc
	overlimits typedDesc
}

func init() {
	Factories["qdisc"] = NewQdiscStatCollector
}

func NewQdiscStatCollector() (Collector, error) {
	return &qdiscStatCollector{
		bytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "bytes"),
			"Number of bytes sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		pkts: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "pkts"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		drops: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "drops"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		requeues: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "requeues"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
		overlimits: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "qdisc", "overlimits"),
			"Number of packets sent.",
			[]string{"iface", "kind"}, nil,
		), prometheus.CounterValue},
	}, nil
}

func (c *qdiscStatCollector) Update(ch chan<- prometheus.Metric) error {
	msgs, err := getQdiscMsgs()
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		m, err := parseMessage(msg)
		if err != nil {
			return err
		}

		// Only report root qdisc info
		if m.Parent != 0 {
			continue
		}

		ifname := ifIndexToName(m.IfaceIdx)

		ch <- c.bytes.mustNewConstMetric(float64(m.Bytes), ifname, m.Kind)
		ch <- c.pkts.mustNewConstMetric(float64(m.Packets), ifname, m.Kind)
		ch <- c.drops.mustNewConstMetric(float64(m.Drops), ifname, m.Kind)
		ch <- c.requeues.mustNewConstMetric(float64(m.Requeues), ifname, m.Kind)
		ch <- c.overlimits.mustNewConstMetric(float64(m.Overlimits), ifname, m.Kind)
		//fmt.Printf("%+v\n", m)
	}

	return nil
}
