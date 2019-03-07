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

// +build !nonetdev
// +build freebsd dragonfly

package collector

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/prometheus/common/log"
)

/*
#cgo CFLAGS: -D_IFI_OQDROPS
#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <net/if.h>
*/
import "C"

func getNetDevStats(ignore *regexp.Regexp, accept *regexp.Regexp) (map[string]map[string]string, error) {
	netDev := map[string]map[string]string{}

	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family == C.AF_LINK {
			dev := C.GoString(ifa.ifa_name)
			if ignore != nil && ignore.MatchString(dev) {
				log.Debugf("Ignoring device: %s", dev)
				continue
			}
			if accept != nil && !accept.MatchString(dev) {
				log.Debugf("Ignoring device: %s", dev)
				continue
			}

			devStats := map[string]string{}
			data := (*C.struct_if_data)(ifa.ifa_data)

			devStats["receive_packets"] = convertFreeBSDCPUTime(uint64(data.ifi_ipackets))
			devStats["transmit_packets"] = convertFreeBSDCPUTime(uint64(data.ifi_opackets))
			devStats["receive_errs"] = convertFreeBSDCPUTime(uint64(data.ifi_ierrors))
			devStats["transmit_errs"] = convertFreeBSDCPUTime(uint64(data.ifi_oerrors))
			devStats["receive_bytes"] = convertFreeBSDCPUTime(uint64(data.ifi_ibytes))
			devStats["transmit_bytes"] = convertFreeBSDCPUTime(uint64(data.ifi_obytes))
			devStats["receive_multicast"] = convertFreeBSDCPUTime(uint64(data.ifi_imcasts))
			devStats["transmit_multicast"] = convertFreeBSDCPUTime(uint64(data.ifi_omcasts))
			devStats["receive_drop"] = convertFreeBSDCPUTime(uint64(data.ifi_iqdrops))
			devStats["transmit_drop"] = convertFreeBSDCPUTime(uint64(data.ifi_oqdrops))
			netDev[dev] = devStats
		}
	}

	return netDev, nil
}

func convertFreeBSDCPUTime(counter uint64) string {
	return strconv.FormatUint(counter, 10)
}
