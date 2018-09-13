// Copyright 2018 The Prometheus Authors
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

// +build linux,amd64
// +build !noiptables

package collector

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

type ipTablesCollector struct {
	entries *prometheus.Desc
}

type protocol int

const (
	ipV4 protocol = iota
	ipV6
)

var (
	_ipTablesTables = []string{
		"filter",
		"nat",
		"mangle",
		"raw",
		"security",
	}

	_ipTablesProto = []protocol{
		ipV4,
		ipV6,
	}
)

func init() {
	registerCollector("ipTables", defaultDisabled, NewIPTablesCollector)
}

// NewIPTablesCollector returns a new Collector exposing IpTables stats.
func NewIPTablesCollector() (Collector, error) {
	return &ipTablesCollector{
		entries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "iptables", "rules"),
			"iptables number of rules by tables",
			[]string{"proto", "table"}, nil,
		),
	}, nil
}

func (c *ipTablesCollector) Update(ch chan<- prometheus.Metric) (err error) {
	for _, table := range _ipTablesTables {
		for _, proto := range _ipTablesProto {
			res, err := ipTablesRules(proto, table)
			if err != nil {
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				c.entries,
				prometheus.GaugeValue,
				float64(res),
				proto.String(),
				table)
		}
	}
	return nil
}

// ipTablesRules returns the number of entries in a table
// entries are policies, rules, targets, ...
func ipTablesRules(proto protocol, table string) (int, error) {

	var domain, level int

	switch proto {
	case ipV6:
		domain = syscall.AF_INET6
		level = syscall.IPPROTO_IPV6
	default:
		domain = syscall.AF_INET
		level = syscall.IPPROTO_IP
	}

	sockFD, err := syscall.Socket(domain, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 84)
	size := uint32(len(buf))

	copy(buf, table)

	err = _getsockopt(sockFD, level, 64, unsafe.Pointer(&buf[0]), &size)
	if err != nil {
		return 0, err
	}

	p := bytes.NewBuffer(buf)

	p.Next(76)
	binary.Read(p, binary.LittleEndian, &size)

	return int(size), nil
}

// protocol Stringer
const _protocolName = "IPV4IPV6"

var _protocolIndex = [...]uint8{0, 4, 8}

func (i protocol) String() string {
	if i < 0 || i >= protocol(len(_protocolIndex)-1) {
		return "protocol(Unknown)"
	}
	return _protocolName[_protocolIndex[i]:_protocolIndex[i+1]]
}

// getsockopt implementation using syscall.Syscall6
func _getsockopt(fd int, level int, name int, val unsafe.Pointer, vallen *uint32) (err error) {
	_, _, eval := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		uintptr(level),
		uintptr(name),
		uintptr(val),
		uintptr(unsafe.Pointer(vallen)),
		0,
	)
	if eval != 0 {
		return eval
	}
	return nil
}
