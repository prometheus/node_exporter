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
	"net"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vishvananda/netlink"
)

type networkRouteCollector struct {
	routeDesc       *prometheus.Desc
	routesTotalDesc *prometheus.Desc
	logger          log.Logger
}

func init() {
	registerCollector("network_route", defaultDisabled, NewNetworkRouteCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewNetworkRouteCollector(logger log.Logger) (Collector, error) {
	const subsystem = "network"

	routeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "route"),
		"network routing table", []string{"if", "src", "dest", "gw", "priority", "proto", "weight"}, nil,
	)
	routeTotalDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "routes_total"),
		"network total routes", []string{"if"}, nil,
	)

	return &networkRouteCollector{
		routeDesc:       routeDesc,
		routesTotalDesc: routeTotalDesc,
		logger:          logger,
	}, nil
}

func (n networkRouteCollector) Update(ch chan<- prometheus.Metric) error {
	deviceRoutes := make(map[string]int)

	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return fmt.Errorf("couldn't get route list: %w", err)
	}

	for _, route := range routes {
		if len(route.MultiPath) != 0 { // route has multipath
			for _, nexthop := range route.MultiPath {
				link, err := netlink.LinkByIndex(nexthop.LinkIndex)
				if err != nil {
					return fmt.Errorf("couldn't get link by index: %w", err)
				}
				labels := []string{
					link.Attrs().Name,                            // if
					networkRouteIPToString(route.Src),            // src
					networkRouteIPNetToString(route.Dst),         // dest
					networkRouteIPToString(nexthop.Gw),           // gw
					strconv.Itoa(route.Priority),                 // priority(metrics)
					networkRouteProtocolToString(route.Protocol), // proto
					strconv.Itoa(nexthop.Hops + 1),               // weight
				}
				ch <- prometheus.MustNewConstMetric(n.routeDesc, prometheus.GaugeValue, 1, labels...)
				deviceRoutes[link.Attrs().Name]++
			}
		} else {
			link, err := netlink.LinkByIndex(route.LinkIndex)
			if err != nil {
				return fmt.Errorf("couldn't get link by index: %w", err)
			}
			labels := []string{
				link.Attrs().Name,                            // if
				networkRouteIPToString(route.Src),            // src
				networkRouteIPNetToString(route.Dst),         // dest
				networkRouteIPToString(route.Gw),             // gw
				strconv.Itoa(route.Priority),                 // priority(metrics)
				networkRouteProtocolToString(route.Protocol), // proto
				"", // weight
			}
			ch <- prometheus.MustNewConstMetric(n.routeDesc, prometheus.GaugeValue, 1, labels...)
			deviceRoutes[link.Attrs().Name]++
		}
	}

	for dev, total := range deviceRoutes {
		ch <- prometheus.MustNewConstMetric(n.routesTotalDesc, prometheus.GaugeValue, float64(total), dev)
	}

	return nil
}

func networkRouteIPNetToString(ip *net.IPNet) string {
	if ip == nil {
		return "default"
	}
	return ip.String()
}

func networkRouteIPToString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

func networkRouteProtocolToString(protocol int) string {
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
