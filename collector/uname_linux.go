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

// +build !nouname

package collector

import (
	"bytes"

	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/sys/unix"
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

type unameCollector struct{}

func init() {
	registerCollector("uname", defaultEnabled, newUnameCollector)
}

// NewUnameCollector returns new unameCollector.
func newUnameCollector() (Collector, error) {
	return &unameCollector{}, nil
}

func (c unameCollector) Update(ch chan<- prometheus.Metric) error {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1,
		string(uname.Sysname[:bytes.IndexByte(uname.Sysname[:], 0)]),
		string(uname.Release[:bytes.IndexByte(uname.Release[:], 0)]),
		string(uname.Version[:bytes.IndexByte(uname.Version[:], 0)]),
		string(uname.Machine[:bytes.IndexByte(uname.Machine[:], 0)]),
		string(uname.Nodename[:bytes.IndexByte(uname.Nodename[:], 0)]),
		string(uname.Domainname[:bytes.IndexByte(uname.Domainname[:], 0)]),
	)
	return nil
}
