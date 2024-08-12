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

//go:build (darwin || freebsd || openbsd || netbsd || aix) && !nouname
// +build darwin freebsd openbsd netbsd aix
// +build !nouname

package collector

import (
	"strings"

	"golang.org/x/sys/unix"
)

func getUname() (uname, error) {
	var utsname unix.Utsname
	if err := unix.Uname(&utsname); err != nil {
		return uname{}, err
	}

	nodeName, domainName := parseHostNameAndDomainName(utsname)

	output := uname{
		SysName:    unix.ByteSliceToString(utsname.Sysname[:]),
		Release:    unix.ByteSliceToString(utsname.Release[:]),
		Version:    unix.ByteSliceToString(utsname.Version[:]),
		Machine:    unix.ByteSliceToString(utsname.Machine[:]),
		NodeName:   nodeName,
		DomainName: domainName,
	}

	return output, nil
}

// parseHostNameAndDomainName for FreeBSD,OpenBSD,Darwin.
// Attempts to emulate what happens in the Linux uname calls since these OS doesn't have a Domainname.
func parseHostNameAndDomainName(utsname unix.Utsname) (hostname string, domainname string) {
	nodename := unix.ByteSliceToString(utsname.Nodename[:])
	split := strings.SplitN(nodename, ".", 2)

	// We'll always have at least a single element in the array. We assume this
	// is the hostname.
	hostname = split[0]

	// If we have more than one element, we assume this is the domainname.
	// Otherwise leave it to "(none)" like Linux.
	domainname = "(none)"
	if len(split) > 1 {
		domainname = split[1]
	}
	return hostname, domainname
}
