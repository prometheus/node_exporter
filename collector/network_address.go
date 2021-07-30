// Copyright 2021 The Prometheus Authors
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

package collector

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"strconv"
)

type networkAddrCollector struct {
	entries *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector("node_network_address", defaultDisabled, NewNetworkAddressCollector)
}

// NewNetworkAddressCollector returns a new Collector exposing network address labels.
func NewNetworkAddressCollector(logger log.Logger) (Collector, error) {
	return &networkAddrCollector{
		entries: prometheus.NewDesc(prometheus.BuildFQName(namespace, "node_network_address",
			"entries"), "node network address by interface",
			[]string{"interface", "addr", "netmask", "scope"}, nil),
		logger: logger,
	}, nil
}

func (c *networkAddrCollector) Update(ch chan<- prometheus.Metric) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("could not get network interfaces: %w", err)
	}

	for _, addr := range getAddrsInfo(interfaces) {
		ch <- prometheus.MustNewConstMetric(c.entries, prometheus.GaugeValue, 1,
			addr.ifs, addr.addr, addr.netmask, addr.scope)
	}

	return nil
}

type addrInfo struct {
	ifs     string
	addr    string
	scope   string
	netmask string
}

func scope(ip net.IP) string {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return "link-local"
	}

	if ip.IsInterfaceLocalMulticast() {
		return "interface-local"
	}

	if ip.IsGlobalUnicast() {
		return "global"
	}

	return ""
}

// getAddrsInfo returns interface name, address, scope and netmask for all interfaces.
func getAddrsInfo(interfaces []net.Interface) []addrInfo {
	var res []addrInfo

	for _, ifs := range interfaces {
		addrs, _ := ifs.Addrs()
		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			size, _ := ipNet.Mask.Size()

			res = append(res, addrInfo{
				ifs:     ifs.Name,
				addr:    ip.String(),
				scope:   scope(ip),
				netmask: strconv.Itoa(size),
			})
		}
	}

	return res
}
