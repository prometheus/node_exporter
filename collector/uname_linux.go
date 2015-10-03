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
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
)

var unameDesc = prometheus.NewDesc(
	prometheus.BuildFQName(Namespace, "uname", "info"),
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
	Factories["uname"] = newUnameCollector
}

// NewUnameCollector returns new unameCollector.
func newUnameCollector() (Collector, error) {
	return &unameCollector{}, nil
}

func (c unameCollector) Update(ch chan<- prometheus.Metric) error {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1,
		unameToString(uname.Sysname),
		unameToString(uname.Release),
		unameToString(uname.Version),
		unameToString(uname.Machine),
		unameToString(uname.Nodename),
		unameToString(uname.Domainname),
	)
	return nil
}
