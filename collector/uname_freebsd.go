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

// +build !nouname

package collector

import (
	"bytes"
	"strings"

	"golang.org/x/sys/unix"
)

func getUname() (uname, error) {
	var utsname unix.Utsname
	if err := unix.Uname(&utsname); err != nil {
		return uname{}, err
	}

	// We do a little bit of work here to emulate what happens in the Linux
	// uname calls since FreeBSD uname doesn't have a Domainname.
	nodename := string(utsname.Nodename[:bytes.IndexByte(utsname.Nodename[:], 0)])
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

	output := uname{
		SysName:    string(utsname.Sysname[:bytes.IndexByte(utsname.Sysname[:], 0)]),
		Release:    string(utsname.Release[:bytes.IndexByte(utsname.Release[:], 0)]),
		Version:    string(utsname.Version[:bytes.IndexByte(utsname.Version[:], 0)]),
		Machine:    string(utsname.Machine[:bytes.IndexByte(utsname.Machine[:], 0)]),
		NodeName:   hostname,
		DomainName: domainname,
	}

	return output, nil
}
