// Copyright 2022 The Prometheus Authors
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

//go:build !nonetclass && linux
// +build !nonetclass,linux

package collector

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/ethtool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

var (
	netclassRTNLWithStats = kingpin.Flag("collector.netclass_rtnl.with-stats", "Expose the statistics for each network device, replacing netdev collector.").Bool()
	operstateStr          = []string{
		"unknown", "notpresent", "down", "lowerlayerdown", "testing",
		"dormant", "up",
	}
)

func (c *netClassCollector) netClassRTNLUpdate(ch chan<- prometheus.Metric) error {
	linkModes := make(map[string]*ethtool.LinkMode)
	lms, err := c.getLinkModes()
	if err != nil {
		if !errors.Is(errors.Unwrap(err), fs.ErrNotExist) {
			return fmt.Errorf("could not get link modes: %w", err)
		}
		level.Info(c.logger).Log("msg", "ETHTOOL netlink interface unavailable, duplex and linkspeed are not scraped.")
	} else {
		for _, lm := range lms {
			if c.ignoredDevicesPattern.MatchString(lm.Interface.Name) {
				continue
			}
			if lm.SpeedMegabits >= 0 {
				speedBytes := uint64(lm.SpeedMegabits * 1000 * 1000 / 8)
				pushMetric(ch, c.getFieldDesc("speed_bytes"), "speed_bytes", speedBytes, prometheus.GaugeValue, lm.Interface.Name)
			}
			linkModes[lm.Interface.Name] = lm
		}
	}

	// Get most attributes from Netlink
	lMsgs, err := c.getNetClassInfoRTNL()
	if err != nil {
		return fmt.Errorf("could not get net class info: %w", err)
	}

	relevantLinks := make([]rtnetlink.LinkMessage, 0, len(lMsgs))
	for _, msg := range lMsgs {
		if !c.ignoredDevicesPattern.MatchString(msg.Attributes.Name) {
			relevantLinks = append(relevantLinks, msg)
		}
	}

	// Read sysfs for attributes that Netlink doesn't expose
	sysfsAttrs, err := getSysfsAttributes(relevantLinks)
	if err != nil {
		return fmt.Errorf("could not get sysfs device info: %w", err)
	}

	// Parse all the info and update metrics
	for _, msg := range relevantLinks {
		upDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "up"),
			"Value is 1 if operstate is 'up', 0 otherwise.",
			[]string{"device"},
			nil,
		)
		upValue := 0.0
		if msg.Attributes.OperationalState == rtnetlink.OperStateUp {
			upValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, upValue, msg.Attributes.Name)

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data of <iface>, value is always 1.",
			[]string{"device", "address", "broadcast", "duplex", "operstate", "ifalias"},
			nil,
		)
		infoValue := 1.0

		var ifalias = ""
		if msg.Attributes.Alias != nil {
			ifalias = *msg.Attributes.Alias
		}

		duplex := ""
		lm, lmExists := linkModes[msg.Attributes.Name]
		if lmExists {
			duplex = lm.Duplex.String()
		}

		ifaceInfo := sysfsAttrs[msg.Attributes.Name]

		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, msg.Attributes.Name, msg.Attributes.Address.String(), msg.Attributes.Broadcast.String(), duplex, operstateStr[int(msg.Attributes.OperationalState)], ifalias)

		pushMetric(ch, c.getFieldDesc("address_assign_type"), "address_assign_type", ifaceInfo.AddrAssignType, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier"), "carrier", msg.Attributes.Carrier, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_changes_total"), "carrier_changes_total", msg.Attributes.CarrierChanges, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_up_changes_total"), "carrier_up_changes_total", msg.Attributes.CarrierUpCount, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_down_changes_total"), "carrier_down_changes_total", msg.Attributes.CarrierDownCount, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("device_id"), "device_id", ifaceInfo.DevID, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("flags"), "flags", msg.Flags, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("iface_id"), "iface_id", msg.Index, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("iface_link_mode"), "iface_link_mode", msg.Attributes.LinkMode, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("dormant"), "dormant", msg.Attributes.LinkMode, prometheus.GaugeValue, msg.Attributes.Name)

		// kernel logic: IFLA_LINK attribute will be ignore when ifindex is the same as iflink
		// (dev->ifindex != dev_get_iflink(dev) && nla_put_u32(skb, IFLA_LINK, dev_get_iflink(dev)))
		// As interface ID is never 0, we assume msg.Attributes.Type 0  means iflink is omitted in RTM_GETLINK response.
		if msg.Attributes.Type > 0 {
			pushMetric(ch, c.getFieldDesc("iface_link"), "iface_link", msg.Attributes.Type, prometheus.GaugeValue, msg.Attributes.Name)
		} else {
			pushMetric(ch, c.getFieldDesc("iface_link"), "iface_link", msg.Index, prometheus.GaugeValue, msg.Attributes.Name)
		}

		pushMetric(ch, c.getFieldDesc("mtu_bytes"), "mtu_bytes", msg.Attributes.MTU, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("name_assign_type"), "name_assign_type", ifaceInfo.NameAssignType, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("net_dev_group"), "net_dev_group", msg.Attributes.NetDevGroup, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_queue_length"), "transmit_queue_length", msg.Attributes.TxQueueLen, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("protocol_type"), "protocol_type", msg.Type, prometheus.GaugeValue, msg.Attributes.Name)

		// Skip statistics if argument collector.netclass_rtnl.with-stats is false or statistics are unavailable.
		if netclassRTNLWithStats == nil || !*netclassRTNLWithStats || msg.Attributes.Stats64 == nil {
			continue
		}

		pushMetric(ch, c.getFieldDesc("receive_packets_total"), "receive_packets_total", msg.Attributes.Stats64.RXPackets, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_packets_total"), "transmit_packets_total", msg.Attributes.Stats64.TXPackets, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_bytes_total"), "receive_bytes_total", msg.Attributes.Stats64.RXBytes, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_bytes_total"), "transmit_bytes_total", msg.Attributes.Stats64.TXBytes, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_errors_total"), "receive_errors_total", msg.Attributes.Stats64.RXErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_errors_total"), "transmit_errors_total", msg.Attributes.Stats64.TXErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_dropped_total"), "receive_dropped_total", msg.Attributes.Stats64.RXDropped, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_dropped_total"), "transmit_dropped_total", msg.Attributes.Stats64.TXDropped, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("multicast_total"), "multicast_total", msg.Attributes.Stats64.Multicast, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("collisions_total"), "collisions_total", msg.Attributes.Stats64.Collisions, prometheus.GaugeValue, msg.Attributes.Name)

		// Detailed rx_errors.
		pushMetric(ch, c.getFieldDesc("receive_length_errors_total"), "receive_length_errors_total", msg.Attributes.Stats64.RXLengthErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_over_errors_total"), "receive_over_errors_total", msg.Attributes.Stats64.RXOverErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_crc_errors_total"), "receive_crc_errors_total", msg.Attributes.Stats64.RXCRCErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_frame_errors_total"), "receive_frame_errors_total", msg.Attributes.Stats64.RXFrameErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_fifo_errors_total"), "receive_fifo_errors_total", msg.Attributes.Stats64.RXFIFOErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_missed_errors_total"), "receive_missed_errors_total", msg.Attributes.Stats64.RXMissedErrors, prometheus.GaugeValue, msg.Attributes.Name)

		// Detailed tx_errors.
		pushMetric(ch, c.getFieldDesc("transmit_aborted_errors_total"), "transmit_aborted_errors_total", msg.Attributes.Stats64.TXAbortedErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_carrier_errors_total"), "transmit_carrier_errors_total", msg.Attributes.Stats64.TXCarrierErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_fifo_errors_total"), "transmit_fifo_errors_total", msg.Attributes.Stats64.TXFIFOErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_heartbeat_errors_total"), "transmit_heartbeat_errors_total", msg.Attributes.Stats64.TXHeartbeatErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_window_errors_total"), "transmit_window_errors_total", msg.Attributes.Stats64.TXWindowErrors, prometheus.GaugeValue, msg.Attributes.Name)

		// For cslip, etc.
		pushMetric(ch, c.getFieldDesc("receive_compressed_total"), "receive_compressed_total", msg.Attributes.Stats64.RXCompressed, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_compressed_total"), "transmit_compressed_total", msg.Attributes.Stats64.TXCompressed, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_nohandler_total"), "receive_nohandler_total", msg.Attributes.Stats64.RXNoHandler, prometheus.GaugeValue, msg.Attributes.Name)

	}

	return nil
}

func (c *netClassCollector) getNetClassInfoRTNL() ([]rtnetlink.LinkMessage, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lMsgs, err := conn.Link.List()

	return lMsgs, err

}

func (c *netClassCollector) getLinkModes() ([]*ethtool.LinkMode, error) {
	conn, err := ethtool.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lms, err := conn.LinkModes()

	return lms, err
}

// getSysfsAttributes reads attributes that are absent from netlink but provided
// by sysfs.
func getSysfsAttributes(links []rtnetlink.LinkMessage) (sysfs.NetClass, error) {
	netClass := sysfs.NetClass{}
	for _, msg := range links {
		interfaceClass := sysfs.NetClassIface{}
		ifName := msg.Attributes.Name
		devPath := filepath.Join("/sys", "class", "net", ifName)

		// These three attributes hold a device-specific lock when
		// accessed, not the RTNL lock, so they are much less impactful
		// than reading most of the other attributes from sysfs.
		for _, attr := range []string{"addr_assign_type", "dev_id", "name_assign_type"} {
			if err := sysfs.ParseNetClassAttribute(devPath, attr, &interfaceClass); err != nil {
				return nil, err
			}
		}
		netClass[ifName] = interfaceClass
	}
	return netClass, nil
}
