// Copyright 2024 The Prometheus Authors
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

//go:build !nonetdev && netbsd
// +build !nonetdev,netbsd

package collector

import (
	"errors"

	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/sysctl.h>
#include <ifaddrs.h>
#include <net/if.h>
*/
import "C"

func getNetDevStats(filter *deviceFilter, logger log.Logger) (netDevStats, error) {
	netDev := netDevStats{}

	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed")
	}
	defer C.freeifaddrs(ifap)

	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family != C.AF_LINK {
			continue
		}

		dev := C.GoString(ifa.ifa_name)
		if filter.ignored(dev) {
			level.Debug(logger).Log("msg", "Ignoring device", "device", dev)
			continue
		}

		data := (*C.struct_if_data)(ifa.ifa_data)

		// https://github.com/NetBSD/src/blob/trunk/sys/net/if.h#L180-L202
		netDev[dev] = map[string]uint64{
			"receive_packets":    uint64(data.ifi_ipackets),
			"transmit_packets":   uint64(data.ifi_opackets),
			"receive_bytes":      uint64(data.ifi_ibytes),
			"transmit_bytes":     uint64(data.ifi_obytes),
			"receive_errors":     uint64(data.ifi_ierrors),
			"transmit_errors":    uint64(data.ifi_oerrors),
			"receive_dropped":    uint64(data.ifi_iqdrops),
			"receive_multicast":  uint64(data.ifi_imcasts),
			"transmit_multicast": uint64(data.ifi_omcasts),
			"collisions":         uint64(data.ifi_collisions),
			"noproto":            uint64(data.ifi_noproto),
		}

		/*
		 * transmit_dropped (ifi_oqdrops) is not available in
		 * struct if_data for backwards compatibility reasons,
		 * but can be read via sysctl.
		 */
		var ifi_oqdrops uint64
		var len C.size_t = C.size_t(unsafe.Sizeof(ifi_oqdrops))
		name := C.CString(fmt.Sprintf("net.interfaces.%s.sndq.drops", dev))
		defer C.free(unsafe.Pointer(name))
		if C.sysctlbyname(name, unsafe.Pointer(&ifi_oqdrops), (*C.size_t)(unsafe.Pointer(&len)), unsafe.Pointer(nil), 0) == 0 {
			netDev[dev]["transmit_dropped"] = ifi_oqdrops
		}
	}

	return netDev, nil
}
