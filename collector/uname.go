// Copyright 2019 The Prometheus Authors
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

//go:build (darwin || freebsd || openbsd || netbsd || linux) && !nouname
// +build darwin freebsd openbsd netbsd linux
// +build !nouname

package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var unameDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "uname", "info"),
	"Labeled system information as provided by the uname system call.",
	[]string{
		"sysname",
		"release",
		"version",
		"machine",
		"nodename",
		"domainname",
	},
	nil,
)

type unameCollector struct {
	logger log.Logger
}
type uname struct {
	SysName    string
	Release    string
	Version    string
	Machine    string
	NodeName   string
	DomainName string
}

func init() {
	registerCollector("uname", defaultEnabled, newUnameCollector)
}

// NewUnameCollector returns new unameCollector.
func newUnameCollector(logger log.Logger) (Collector, error) {
	return &unameCollector{logger}, nil
}

func (c *unameCollector) Update(ch chan<- prometheus.Metric) error {
	uname, err := getUname()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1,
		uname.SysName,
		uname.Release,
		uname.Version,
		uname.Machine,
		uname.NodeName,
		uname.DomainName,
	)

	return nil
}
