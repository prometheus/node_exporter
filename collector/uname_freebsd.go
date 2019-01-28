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
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/sys/unix"
)

func (c unameCollector) Update(ch chan<- prometheus.Metric) error {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return err
	}

	// We do a little bit of work here to emulate what happens in the Linux
	// uname calls since FreeBSD uname doesn't have a Domainname.
	nodename := string(uname.Nodename[:bytes.IndexByte(uname.Nodename[:], 0)])
	split := strings.SplitN(nodename, ".", 2)

	// We'll always have at least a single element in the array. We assume this
	// is the hostname.
	hostname := split[0]

	// If we have more than one element, we assume this is the domainname.
	// Otherwise leave it to "(none)" like Linux.
	domainname := "(none)"
	if len(split) > 1 {
		domainname = split[1]
	}

	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1,
		string(uname.Sysname[:bytes.IndexByte(uname.Sysname[:], 0)]),
		string(uname.Release[:bytes.IndexByte(uname.Release[:], 0)]),
		string(uname.Version[:bytes.IndexByte(uname.Version[:], 0)]),
		string(uname.Machine[:bytes.IndexByte(uname.Machine[:], 0)]),
		hostname,
		domainname,
	)
	return nil
}
