// Copyright 2020 The Prometheus Authors
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

// +build !nonetworkroute

package collector

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/jsimonetti/rtnetlink"
	"github.com/prometheus/client_golang/prometheus"
)

type networkRouteCollector struct {
	routeInfoDesc *prometheus.Desc
	routesDesc    *prometheus.Desc
	logger        log.Logger
}

func init() {
	registerCollector("network_route", defaultDisabled, NewNetworkRouteCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewNetworkRouteCollector(logger log.Logger) (Collector, error) {
	const subsystem = "network"

	routeInfoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "route_info"),
		"network routing table information", []string{"device", "src", "dest", "gw", "priority", "proto", "weight"}, nil,
	)
	routesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "routes"),
		"network routes by interface", []string{"device"}, nil,
	)

	return &networkRouteCollector{
		routeInfoDesc: routeInfoDesc,
		routesDesc:    routesDesc,
		logger:        logger,
	}, nil
}

func (n networkRouteCollector) Update(ch chan<- prometheus.Metric) error {
	deviceRoutes := make(map[string]int)

	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("couldn't connect rtnetlink: %w", err)
	}
	defer conn.Close()

	links, err := conn.Link.List()
	if err != nil {
		return fmt.Errorf("couldn't get links: %w", err)
	}

	routes, err := conn.Route.List()
	if err != nil {
		return fmt.Errorf("couldn't get routes: %w", err)
	}

	for _, route := range routes {
		if route.Type != unix.RTA_DST {
			continue
		}
		if len(route.Attributes.Multipath) != 0 {
			for _, nextHop := range route.Attributes.Multipath {
				ifName := ""
				for _, link := range links {
					if link.Index == nextHop.Hop.IfIndex {
						ifName = link.Attributes.Name
						break
					}
				}

				labels := []string{
					ifName, // if
					networkRouteIPToString(route.Attributes.Src),                            // src
					networkRouteIPWithPrefixToString(route.Attributes.Dst, route.DstLength), // dest
					networkRouteIPToString(nextHop.Gateway),                                 // gw
					strconv.FormatUint(uint64(route.Attributes.Priority), 10),               // priority(metrics)
					networkRouteProtocolToString(route.Protocol),                            // proto
					strconv.Itoa(int(nextHop.Hop.Hops) + 1),                                 // weight
				}
				ch <- prometheus.MustNewConstMetric(n.routeInfoDesc, prometheus.GaugeValue, 1, labels...)
				deviceRoutes[ifName]++
			}
		} else {
			ifName := ""
			for _, link := range links {
				if link.Index == route.Attributes.OutIface {
					ifName = link.Attributes.Name
					break
				}
			}

			labels := []string{
				ifName, // if
				networkRouteIPToString(route.Attributes.Src),                            // src
				networkRouteIPWithPrefixToString(route.Attributes.Dst, route.DstLength), // dest
				networkRouteIPToString(route.Attributes.Gateway),                        // gw
				strconv.FormatUint(uint64(route.Attributes.Priority), 10),               // priority(metrics)
				networkRouteProtocolToString(route.Protocol),                            // proto
				"", // weight
			}
			ch <- prometheus.MustNewConstMetric(n.routeInfoDesc, prometheus.GaugeValue, 1, labels...)
			deviceRoutes[ifName]++
		}
	}

	for dev, total := range deviceRoutes {
		ch <- prometheus.MustNewConstMetric(n.routesDesc, prometheus.GaugeValue, float64(total), dev)
	}

	return nil
}

func networkRouteIPWithPrefixToString(ip net.IP, len uint8) string {
	if len == 0 {
		return "default"
	}
	iplen := net.IPv4len
	if ip.To4() == nil {
		iplen = net.IPv6len
	}
	network := &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(int(len), iplen*8),
	}
	return network.String()
}

func networkRouteIPToString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

func networkRouteProtocolToString(protocol uint8) string {
	// from linux kernel 'include/uapi/linux/rtnetlink.h'
	switch protocol {
	case 0:
		return "unspec"
	case 1:
		return "redirect"
	case 2:
		return "kernel"
	case 3:
		return "boot"
	case 4:
		return "static"
	case 8:
		return "gated"
	case 9:
		return "ra"
	case 10:
		return "mrt"
	case 11:
		return "zebra"
	case 12:
		return "bird"
	case 13:
		return "dnrouted"
	case 14:
		return "xorp"
	case 15:
		return "ntk"
	case 16:
		return "dhcp"
	case 17:
		return "mrouted"
	case 42:
		return "babel"
	case 186:
		return "bgp"
	case 187:
		return "isis"
	case 188:
		return "ospf"
	case 189:
		return "rip"
	case 192:
		return "eigrp"
	}
	return "unknown"
}
