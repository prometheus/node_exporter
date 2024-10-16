// Copyright 2020 The Prometheus Authors
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

//go:build !nonetdev
// +build !nonetdev

package collector

import (
	"log/slog"

	"unsafe"

	"golang.org/x/sys/unix"
)

func getNetDevStats(filter *deviceFilter, logger *slog.Logger) (netDevStats, error) {
	netDev := netDevStats{}

	mib := [6]_C_int{unix.CTL_NET, unix.AF_ROUTE, 0, 0, unix.NET_RT_IFLIST, 0}
	buf, err := sysctl(mib[:])
	if err != nil {
		return nil, err
	}
	n := uintptr(len(buf))
	index := uintptr(unsafe.Pointer(&buf[0]))
	next := uintptr(0)

	var rtm *unix.RtMsghdr

	for next = index; next < (index + n); next += uintptr(rtm.Msglen) {
		rtm = (*unix.RtMsghdr)(unsafe.Pointer(next))
		if rtm.Version != unix.RTM_VERSION || rtm.Type != unix.RTM_IFINFO {
			continue
		}
		ifm := (*unix.IfMsghdr)(unsafe.Pointer(next))
		if ifm.Addrs&unix.RTA_IFP == 0 {
			continue
		}
		dl := (*unix.RawSockaddrDatalink)(unsafe.Pointer(next + uintptr(rtm.Hdrlen)))
		if dl.Family != unix.AF_LINK {
			continue
		}
		data := ifm.Data
		dev := int8ToString(dl.Data[:dl.Nlen])
		if filter.ignored(dev) {
			logger.Debug("Ignoring device", "device", dev)
			continue
		}

		// https://cs.opensource.google/go/x/sys/+/master:unix/ztypes_openbsd_amd64.go;l=292-316
		netDev[dev] = map[string]uint64{
			"receive_packets":    data.Ipackets,
			"transmit_packets":   data.Opackets,
			"receive_bytes":      data.Ibytes,
			"transmit_bytes":     data.Obytes,
			"receive_errors":     data.Ierrors,
			"transmit_errors":    data.Oerrors,
			"receive_dropped":    data.Iqdrops,
			"transmit_dropped":   data.Oqdrops,
			"receive_multicast":  data.Imcasts,
			"transmit_multicast": data.Omcasts,
			"collisions":         data.Collisions,
			"noproto":            data.Noproto,
		}
	}
	return netDev, nil
}

func getNetDevLabels() (map[string]map[string]string, error) {
	// to be implemented if needed
	return nil, nil
}
