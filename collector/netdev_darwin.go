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

package collector

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/prometheus/common/log"
)

/*
#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <net/if.h>
*/
import "C"

func getNetDevStats(ignore *regexp.Regexp) (map[string]map[string]string, error) {
	netDev := map[string]map[string]string{}

	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family == C.AF_LINK {
			dev := C.GoString(ifa.ifa_name)
			if ignore.MatchString(dev) {
				log.Debugf("Ignoring device: %s", dev)
				continue
			}

			devStats := map[string]string{}
			data := (*C.struct_if_data)(ifa.ifa_data)
			devStats["receive_packets"] = strconv.FormatUint(uint64(data.ifi_ipackets), 10)
			devStats["transmit_packets"] = strconv.FormatUint(uint64(data.ifi_opackets), 10)
			devStats["receive_errs"] = strconv.FormatUint(uint64(data.ifi_ierrors), 10)
			devStats["transmit_errs"] = strconv.FormatUint(uint64(data.ifi_oerrors), 10)
			devStats["receive_bytes"] = strconv.FormatUint(uint64(data.ifi_ibytes), 10)
			devStats["transmit_bytes"] = strconv.FormatUint(uint64(data.ifi_obytes), 10)
			devStats["receive_multicast"] = strconv.FormatUint(uint64(data.ifi_imcasts), 10)
			devStats["transmit_multicast"] = strconv.FormatUint(uint64(data.ifi_omcasts), 10)
			netDev[dev] = devStats
		}
	}

	return netDev, nil
}
