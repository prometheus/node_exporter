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
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"regexp"
	"strconv"

	"github.com/prometheus/common/log"
	"golang.org/x/sys/unix"
)

func getNetDevStats(ignore *regexp.Regexp, accept *regexp.Regexp) (map[string]map[string]string, error) {
	netDev := map[string]map[string]string{}

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("net.Interfaces() failed")
	}

	for _, iface := range ifs {
		ifaceData, err := getIfaceData(iface.Index)
		if err != nil {
			log.Debugf("failed to load data for interface %q: %v", iface.Name, err)
			continue
		}

		if ignore != nil && ignore.MatchString(iface.Name) {
			log.Debugf("Ignoring device: %s", iface.Name)
			continue
		}
		if accept != nil && !accept.MatchString(iface.Name) {
			log.Debugf("Ignoring device: %s", iface.Name)
			continue
		}

		devStats := map[string]string{}
		devStats["receive_packets"] = strconv.FormatUint(ifaceData.Data.Ipackets, 10)
		devStats["transmit_packets"] = strconv.FormatUint(ifaceData.Data.Opackets, 10)
		devStats["receive_errs"] = strconv.FormatUint(ifaceData.Data.Ierrors, 10)
		devStats["transmit_errs"] = strconv.FormatUint(ifaceData.Data.Oerrors, 10)
		devStats["receive_bytes"] = strconv.FormatUint(ifaceData.Data.Ibytes, 10)
		devStats["transmit_bytes"] = strconv.FormatUint(ifaceData.Data.Obytes, 10)
		devStats["receive_multicast"] = strconv.FormatUint(ifaceData.Data.Imcasts, 10)
		devStats["transmit_multicast"] = strconv.FormatUint(ifaceData.Data.Omcasts, 10)
		netDev[iface.Name] = devStats
	}

	return netDev, nil
}

func getIfaceData(index int) (*ifMsghdr2, error) {
	var data ifMsghdr2
	rawData, err := unix.SysctlRaw("net", unix.AF_ROUTE, 0, 0, unix.NET_RT_IFLIST2, index)
	if err != nil {
		return nil, err
	}
	err = binary.Read(bytes.NewReader(rawData), binary.LittleEndian, &data)
	return &data, err
}

type ifMsghdr2 struct {
	Msglen    uint16
	Version   uint8
	Type      uint8
	Addrs     int32
	Flags     int32
	Index     uint16
	_         [2]byte
	SndLen    int32
	SndMaxlen int32
	SndDrops  int32
	Timer     int32
	Data      ifData64
}

type ifData64 struct {
	Type       uint8
	Typelen    uint8
	Physical   uint8
	Addrlen    uint8
	Hdrlen     uint8
	Recvquota  uint8
	Xmitquota  uint8
	Unused1    uint8
	Mtu        uint32
	Metric     uint32
	Baudrate   uint64
	Ipackets   uint64
	Ierrors    uint64
	Opackets   uint64
	Oerrors    uint64
	Collisions uint64
	Ibytes     uint64
	Obytes     uint64
	Imcasts    uint64
	Omcasts    uint64
	Iqdrops    uint64
	Noproto    uint64
	Recvtiming uint32
	Xmittiming uint32
	Lastchange unix.Timeval32
}
