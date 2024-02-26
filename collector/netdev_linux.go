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
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jsimonetti/rtnetlink"
	"github.com/prometheus/procfs"
)

var (
	netDevNetlink = kingpin.Flag("collector.netdev.netlink", "Use netlink to gather stats instead of /proc/net/dev.").Default("true").Bool()
)

func getNetDevStats(filter *deviceFilter, logger log.Logger) (netDevStats, error) {
	if *netDevNetlink {
		return netlinkStats(filter, logger)
	}
	return procNetDevStats(filter, logger)
}

func netlinkStats(filter *deviceFilter, logger log.Logger) (netDevStats, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	links, err := conn.Link.List()
	if err != nil {
		return nil, err
	}

	return parseNetlinkStats(links, filter, logger), nil
}

func parseNetlinkStats(links []rtnetlink.LinkMessage, filter *deviceFilter, logger log.Logger) netDevStats {
	metrics := netDevStats{}

	for _, msg := range links {
		if msg.Attributes == nil {
			level.Debug(logger).Log("msg", "No netlink attributes, skipping")
			continue
		}
		name := msg.Attributes.Name
		stats := msg.Attributes.Stats64
		if stats32 := msg.Attributes.Stats; stats == nil && stats32 != nil {
			stats = &rtnetlink.LinkStats64{
				RXPackets:          uint64(stats32.RXPackets),
				TXPackets:          uint64(stats32.TXPackets),
				RXBytes:            uint64(stats32.RXBytes),
				TXBytes:            uint64(stats32.TXBytes),
				RXErrors:           uint64(stats32.RXErrors),
				TXErrors:           uint64(stats32.TXErrors),
				RXDropped:          uint64(stats32.RXDropped),
				TXDropped:          uint64(stats32.TXDropped),
				Multicast:          uint64(stats32.Multicast),
				Collisions:         uint64(stats32.Collisions),
				RXLengthErrors:     uint64(stats32.RXLengthErrors),
				RXOverErrors:       uint64(stats32.RXOverErrors),
				RXCRCErrors:        uint64(stats32.RXCRCErrors),
				RXFrameErrors:      uint64(stats32.RXFrameErrors),
				RXFIFOErrors:       uint64(stats32.RXFIFOErrors),
				RXMissedErrors:     uint64(stats32.RXMissedErrors),
				TXAbortedErrors:    uint64(stats32.TXAbortedErrors),
				TXCarrierErrors:    uint64(stats32.TXCarrierErrors),
				TXFIFOErrors:       uint64(stats32.TXFIFOErrors),
				TXHeartbeatErrors:  uint64(stats32.TXHeartbeatErrors),
				TXWindowErrors:     uint64(stats32.TXWindowErrors),
				RXCompressed:       uint64(stats32.RXCompressed),
				TXCompressed:       uint64(stats32.TXCompressed),
				RXNoHandler:        uint64(stats32.RXNoHandler),
				RXOtherhostDropped: 0,
			}
		}

		if filter.ignored(name) {
			level.Debug(logger).Log("msg", "Ignoring device", "device", name)
			continue
		}

		// Make sure we don't panic when accessing `stats` attributes below.
		if stats == nil {
			level.Debug(logger).Log("msg", "No netlink stats, skipping")
			continue
		}

		// https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_link.h#L42-L246
		metrics[name] = map[string]uint64{
			"receive_packets":  stats.RXPackets,
			"transmit_packets": stats.TXPackets,
			"receive_bytes":    stats.RXBytes,
			"transmit_bytes":   stats.TXBytes,
			"receive_errors":   stats.RXErrors,
			"transmit_errors":  stats.TXErrors,
			"receive_dropped":  stats.RXDropped,
			"transmit_dropped": stats.TXDropped,
			"multicast":        stats.Multicast,
			"collisions":       stats.Collisions,

			// detailed rx_errors
			"receive_length_errors": stats.RXLengthErrors,
			"receive_over_errors":   stats.RXOverErrors,
			"receive_crc_errors":    stats.RXCRCErrors,
			"receive_frame_errors":  stats.RXFrameErrors,
			"receive_fifo_errors":   stats.RXFIFOErrors,
			"receive_missed_errors": stats.RXMissedErrors,

			// detailed tx_errors
			"transmit_aborted_errors":   stats.TXAbortedErrors,
			"transmit_carrier_errors":   stats.TXCarrierErrors,
			"transmit_fifo_errors":      stats.TXFIFOErrors,
			"transmit_heartbeat_errors": stats.TXHeartbeatErrors,
			"transmit_window_errors":    stats.TXWindowErrors,

			// for cslip etc
			"receive_compressed":  stats.RXCompressed,
			"transmit_compressed": stats.TXCompressed,
			"receive_nohandler":   stats.RXNoHandler,
		}
	}

	return metrics
}

func procNetDevStats(filter *deviceFilter, logger log.Logger) (netDevStats, error) {
	metrics := netDevStats{}

	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return metrics, fmt.Errorf("failed to open procfs: %w", err)
	}

	netDev, err := fs.NetDev()
	if err != nil {
		return metrics, fmt.Errorf("failed to parse /proc/net/dev: %w", err)
	}

	for _, stats := range netDev {
		name := stats.Name

		if filter.ignored(name) {
			level.Debug(logger).Log("msg", "Ignoring device", "device", name)
			continue
		}

		metrics[name] = map[string]uint64{
			"receive_bytes":       stats.RxBytes,
			"receive_packets":     stats.RxPackets,
			"receive_errors":      stats.RxErrors,
			"receive_dropped":     stats.RxDropped,
			"receive_fifo":        stats.RxFIFO,
			"receive_frame":       stats.RxFrame,
			"receive_compressed":  stats.RxCompressed,
			"receive_multicast":   stats.RxMulticast,
			"transmit_bytes":      stats.TxBytes,
			"transmit_packets":    stats.TxPackets,
			"transmit_errors":     stats.TxErrors,
			"transmit_dropped":    stats.TxDropped,
			"transmit_fifo":       stats.TxFIFO,
			"transmit_colls":      stats.TxCollisions,
			"transmit_carrier":    stats.TxCarrier,
			"transmit_compressed": stats.TxCompressed,
		}
	}

	return metrics, nil
}
