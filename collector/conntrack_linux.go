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

//go:build !noconntrack

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type conntrackCollector struct {
	logger *slog.Logger
}

type conntrackStatistics struct {
	found         uint64 // Number of searched entries which were successful
	invalid       uint64 // Number of packets seen which can not be tracked
	ignore        uint64 // Number of packets seen which are already connected to a conntrack entry
	insert        uint64 // Number of entries inserted into the list
	insertFailed  uint64 // Number of entries for which list insertion was attempted but failed (happens if the same entry is already present)
	drop          uint64 // Number of packets dropped due to conntrack failure. Either new conntrack entry allocation failed, or protocol helper dropped the packet
	earlyDrop     uint64 // Number of dropped conntrack entries to make room for new ones, if maximum table size was reached
	searchRestart uint64 // Number of conntrack table lookups which had to be restarted due to hashtable resizes
}

func init() {
	registerCollector("conntrack", defaultEnabled, NewConntrackCollector)
}

var (
	conntrackCurrent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_entries"),
		"Number of currently allocated flow entries for connection tracking.",
		nil, nil,
	)
	conntrackLimit = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_entries_limit"),
		"Maximum size of connection tracking table.",
		nil, nil,
	)
	conntrackFound = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_found"),
		"Number of searched entries which were successful.",
		nil, nil,
	)
	conntrackInvalid = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_invalid"),
		"Number of packets seen which can not be tracked.",
		nil, nil,
	)
	conntrackIgnore = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_ignore"),
		"Number of packets seen which are already connected to a conntrack entry.",
		nil, nil,
	)
	conntrackInsert = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_insert"),
		"Number of entries inserted into the list.",
		nil, nil,
	)
	conntrackInsertFailed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_insert_failed"),
		"Number of entries for which list insertion was attempted but failed.",
		nil, nil,
	)
	conntrackDrop = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_drop"),
		"Number of packets dropped due to conntrack failure.",
		nil, nil,
	)
	conntrackEarlyDrop = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_early_drop"),
		"Number of dropped conntrack entries to make room for new ones, if maximum table size was reached.",
		nil, nil,
	)
	conntrackSearchRestart = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_search_restart"),
		"Number of conntrack table lookups which had to be restarted due to hashtable resizes.",
		nil, nil,
	)
)

// NewConntrackCollector returns a new Collector exposing conntrack stats.
func NewConntrackCollector(logger *slog.Logger) (Collector, error) {
	return &conntrackCollector{
		logger: logger,
	}, nil
}

func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	value, err := readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_count"))
	if err != nil {
		return c.handleErr(err)
	}
	ch <- prometheus.MustNewConstMetric(
		conntrackCurrent, prometheus.GaugeValue, float64(value))

	value, err = readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_max"))
	if err != nil {
		return c.handleErr(err)
	}
	ch <- prometheus.MustNewConstMetric(
		conntrackLimit, prometheus.GaugeValue, float64(value))

	conntrackStats, err := getConntrackStatistics()
	if err != nil {
		return c.handleErr(err)
	}

	ch <- prometheus.MustNewConstMetric(
		conntrackFound, prometheus.GaugeValue, float64(conntrackStats.found))
	ch <- prometheus.MustNewConstMetric(
		conntrackInvalid, prometheus.GaugeValue, float64(conntrackStats.invalid))
	ch <- prometheus.MustNewConstMetric(
		conntrackIgnore, prometheus.GaugeValue, float64(conntrackStats.ignore))
	ch <- prometheus.MustNewConstMetric(
		conntrackInsert, prometheus.GaugeValue, float64(conntrackStats.insert))
	ch <- prometheus.MustNewConstMetric(
		conntrackInsertFailed, prometheus.GaugeValue, float64(conntrackStats.insertFailed))
	ch <- prometheus.MustNewConstMetric(
		conntrackDrop, prometheus.GaugeValue, float64(conntrackStats.drop))
	ch <- prometheus.MustNewConstMetric(
		conntrackEarlyDrop, prometheus.GaugeValue, float64(conntrackStats.earlyDrop))
	ch <- prometheus.MustNewConstMetric(
		conntrackSearchRestart, prometheus.GaugeValue, float64(conntrackStats.searchRestart))
	return nil
}

func (c *conntrackCollector) handleErr(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		c.logger.Debug("conntrack probably not loaded")
		return ErrNoData
	}
	return fmt.Errorf("failed to retrieve conntrack stats: %w", err)
}

func getConntrackStatistics() (*conntrackStatistics, error) {
	s := conntrackStatistics{}

	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	connStats, err := fs.ConntrackStat()
	if err != nil {
		return nil, err
	}

	for _, connStat := range connStats {
		s.found += connStat.Found
		s.invalid += connStat.Invalid
		s.ignore += connStat.Ignore
		s.insert += connStat.Insert
		s.insertFailed += connStat.InsertFailed
		s.drop += connStat.Drop
		s.earlyDrop += connStat.EarlyDrop
		s.searchRestart += connStat.SearchRestart
	}

	return &s, nil
}
