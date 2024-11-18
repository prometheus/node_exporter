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

//go:build !nonetdev
// +build !nonetdev

package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"golang.org/x/sys/unix"
)

func getNetDevStats(filter *deviceFilter, logger *slog.Logger) (netDevStats, error) {
	netDev := netDevStats{}

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("net.Interfaces() failed: %w", err)
	}

	for _, iface := range ifs {
		if filter.ignored(iface.Name) {
			logger.Debug("Ignoring device", "device", iface.Name)
			continue
		}

		ifaceData, err := getIfaceData(iface.Index)
		if err != nil {
			logger.Debug("failed to load data for interface", "device", iface.Name, "err", err)
			continue
		}

		netDev[iface.Name] = map[string]uint64{
			"receive_packets":    ifaceData.Data.Ipackets,
			"transmit_packets":   ifaceData.Data.Opackets,
			"receive_bytes":      ifaceData.Data.Ibytes,
			"transmit_bytes":     ifaceData.Data.Obytes,
			"receive_errors":     ifaceData.Data.Ierrors,
			"transmit_errors":    ifaceData.Data.Oerrors,
			"receive_dropped":    ifaceData.Data.Iqdrops,
			"receive_multicast":  ifaceData.Data.Imcasts,
			"transmit_multicast": ifaceData.Data.Omcasts,
			"collisions":         ifaceData.Data.Collisions,
			"noproto":            ifaceData.Data.Noproto,
		}
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

// https://github.com/apple/darwin-xnu/blob/main/bsd/net/if_var.h#L199-L231
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

func getNetDevLabels() (map[string]map[string]string, error) {
	// to be implemented if needed
	return nil, nil
}
