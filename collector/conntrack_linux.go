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

// +build !noconntrack

package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

type conntrackCollector struct {
	current         *prometheus.Desc
	limit           *prometheus.Desc
	kernelStatistic *prometheus.Desc
}

var (
	enableConntrackKernelStats = kingpin.Flag("collector.conntrack.kernel-stats", "fetch conntrack stats from kernel (requires root or CAP_NET_ADMIN)").Bool()
)

func init() {
	registerCollector("conntrack", defaultEnabled, NewConntrackCollector)
}

type testType struct {
	StatLen  uint8
	_        uint8
	StatType uint8
	_        uint8
	StatVal  uint32
}

// NewConntrackCollector returns a new Collector exposing conntrack stats.
func NewConntrackCollector() (Collector, error) {
	c := &conntrackCollector{
		current: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_entries"),
			"Number of currently allocated flow entries for connection tracking.",
			nil, nil,
		),
		limit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_entries_limit"),
			"Maximum size of connection tracking table.",
			nil, nil,
		),
	}

	if *enableConntrackKernelStats {
		c.kernelStatistic = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_kernel_statistic"),
			"Conntrack Kernel counter",
			[]string{"cpu", "statistic"}, nil,
		)
	}

	return c, nil
}

func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	value, err := readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_count"))
	if err != nil {
		// Conntrack probably not loaded into the kernel.
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.current, prometheus.GaugeValue, float64(value))

	value, err = readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_max"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.limit, prometheus.GaugeValue, float64(value))

	if *enableConntrackKernelStats {
		err = c.updateConntrackKernelStats(ch)
	}

	return err
}

type conntrackAttributes uint8

// Maps to ctattr_stats_cpu
const (
	ctaUnspecified conntrackAttributes = iota

	ctaSearched // Not used
	ctaFound
	ctaNew // Not used
	ctaInvalid
	ctaIgnore
	ctaDelete     // Not used
	ctaDeleteList // Not used
	ctaInsert
	ctaInsertFailed
	ctaDrop
	ctaEarlyDrop
	ctaError
	ctaSearchRestart
)

var conntrackAttributeLabels = map[conntrackAttributes]string{
	ctaUnspecified:   "unspecified",
	ctaSearched:      "searched",
	ctaFound:         "found",
	ctaNew:           "new",
	ctaInvalid:       "invalid",
	ctaIgnore:        "ignore",
	ctaDelete:        "delete",
	ctaDeleteList:    "delete_list",
	ctaInsert:        "insert",
	ctaInsertFailed:  "insert_failed",
	ctaDrop:          "drop",
	ctaEarlyDrop:     "early_drop",
	ctaError:         "error",
	ctaSearchRestart: "search_restart",
}

type conntrackStatistic struct {
	StatLen  uint8
	_        uint8 // padding
	StatType conntrackAttributes
	_        uint8 // Padding
	StatVal  uint32
}

const (
	netlinkFamilyNetfilter = 12
	// From conntrack.c, nfct_mnl_nlmsghdr_put, nlmsg_type = (subsys << 8) | type
	// where subsys = NFNL_SUBSYS_CTNETLINK        (include/uapi/linux/netfilter/nfnetlink.h)
	// and   type   = IPCTNL_MSG_CT_GET_STATS_CPU  (include/uapi/linux/netfilter/nfnetlink_conntrack.h)
	// gives (1 << 8) | 4 = 260
	netlinkConntrackStatisticsType = 260
)

func (c *conntrackCollector) updateConntrackKernelStats(ch chan<- prometheus.Metric) error {
	conn, err := netlink.Dial(netlinkFamilyNetfilter, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Debug("Opening netlink connection")
	req := netlink.Message{
		Header: netlink.Header{
			Flags: netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump,
			Type:  netlinkConntrackStatisticsType,
		},
		Data: []byte{0, 0, 0, 0}, // Unclear why this is needed, but conntrack -S sends it, and NL hangs if we don't
	}

	msgs, err := conn.Execute(req)
	if err != nil {
		return err
	}

	for cpuIdx, m := range msgs {
		payload := m.Data[4:]
		if len(payload)%8 != 0 {
			return fmt.Errorf("Unexpected size of conntrack stats from kernel, got %d bytes, expected multiple of 8")
		}

		n := int(len(payload) / 8)
		r := bytes.NewReader(payload)
		sts := make([]conntrackStatistic, n)
		err := binary.Read(r, binary.BigEndian, &sts)
		if err != nil {
			log.Errorf("Couldn't deserialize conntrack kernel stats")
			return err
		}

		cpuLabel := fmt.Sprintf("%d", cpuIdx)
		for _, st := range sts {
			statisticLabel, ok := conntrackAttributeLabels[st.StatType]
			if !ok {
				log.Warnf("Unknown conntrack statistic type %d", st.StatType)
				continue
			}
			ch <- prometheus.MustNewConstMetric(
				c.kernelStatistic,
				prometheus.CounterValue,
				float64(st.StatVal),
				cpuLabel, statisticLabel,
			)
		}
	}

	return nil
}
