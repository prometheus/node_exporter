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
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/jsimonetti/rtnetlink"
)

func getNetDevStats(filter *deviceFilter, logger log.Logger) (netDevStats, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	links, err := conn.Link.List()
	if err != nil {
		return nil, err
	}

	return netlinkStats(links, filter, logger), nil
}

func netlinkStats(links []rtnetlink.LinkMessage, filter *deviceFilter, logger log.Logger) netDevStats {
	metrics := netDevStats{}

	for _, msg := range links {
		name := msg.Attributes.Name
		stats := msg.Attributes.Stats64

		if filter.ignored(name) {
			level.Debug(logger).Log("msg", "Ignoring device", "device", name)
			continue
		}

		// https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_link.h#L42-L246
		// https://github.com/torvalds/linux/blob/master/net/core/net-procfs.c#L75-L97
		metrics[name] = map[string]uint64{
			"receive_packets":     stats.RXPackets,
			"transmit_packets":    stats.TXPackets,
			"receive_bytes":       stats.RXBytes,
			"transmit_bytes":      stats.TXBytes,
			"receive_errs":        stats.RXErrors,
			"transmit_errs":       stats.TXErrors,
			"receive_drop":        stats.RXDropped + stats.RXMissedErrors,
			"transmit_drop":       stats.TXDropped,
			"receive_multicast":   stats.Multicast,
			"transmit_colls":      stats.Collisions,
			"receive_frame":       stats.RXLengthErrors + stats.RXOverErrors + stats.RXCRCErrors + stats.RXFrameErrors,
			"receive_fifo":        stats.RXFIFOErrors,
			"transmit_carrier":    stats.TXAbortedErrors + stats.TXCarrierErrors + stats.TXHeartbeatErrors + stats.TXWindowErrors,
			"transmit_fifo":       stats.TXFIFOErrors,
			"receive_compressed":  stats.RXCompressed,
			"transmit_compressed": stats.TXCompressed,
		}
	}

	return metrics
}
