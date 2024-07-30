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

//go:build !nonetdev
// +build !nonetdev

package collector

import (
	"log/slog"

	"github.com/power-devops/perfstat"
)

func getNetDevStats(filter *deviceFilter, logger *slog.Logger) (netDevStats, error) {
	netDev := netDevStats{}

	stats, err := perfstat.NetAdapterStat()
	if err != nil {
		return nil, err
	}

	for _, stat := range stats {
		netDev[stat.Name] = map[string]uint64{
			"receive_packets":    uint64(stat.RxPackets),
			"transmit_packets":   uint64(stat.TxPackets),
			"receive_bytes":      uint64(stat.RxBytes),
			"transmit_bytes":     uint64(stat.TxBytes),
			"receive_errors":     uint64(stat.RxErrors),
			"transmit_errors":    uint64(stat.TxErrors),
			"receive_dropped":    uint64(stat.RxPacketsDropped),
			"transmit_dropped":   uint64(stat.TxPacketsDropped),
			"receive_multicast":  uint64(stat.RxMulticastPackets),
			"transmit_multicast": uint64(stat.TxMulticastPackets),
		}
	}

	return netDev, nil
}

func getNetDevLabels() (map[string]map[string]string, error) {
	// to be implemented if needed
	return nil, nil
}
