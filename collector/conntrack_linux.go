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
	"bufio"
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"strconv"
	"strings"
)

type conntrackCollector struct {
	current      *prometheus.Desc
	limit        *prometheus.Desc
	searched     *prometheus.Desc
	insert       *prometheus.Desc
	insertFailed *prometheus.Desc
	delete       *prometheus.Desc
	drop         *prometheus.Desc
}

type ConntrackStatistics struct {
	Searched     uint64 // Number of conntrack table lookups performed
	Delete       uint64 // Number of conntrack entries which were removed
	Insert       uint64 // Number of entries inserted into the list
	InsertFailed uint64 // Number of entries for which list insertion was attempted but failed (happens if the same entry is already present)
	Drop         uint64 // Number of packets dropped due to conntrack failure. Either new conntrack entry allocation failed, or protocol helper dropped the packet
}

func init() {
	registerCollector("conntrack", defaultEnabled, NewConntrackCollector)
}

// NewConntrackCollector returns a new Collector exposing conntrack stats.
func NewConntrackCollector() (Collector, error) {
	return &conntrackCollector{
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
		searched: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_searched"),
			"Number of conntrack table lookups performed.",
			nil, nil,
		),
		insert: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_insert"),
			"Number of entries inserted into the list.",
			nil, nil,
		),
		insertFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_insert_failed"),
			"Number of entries for which list insertion was attempted but failed.",
			nil, nil,
		),
		delete: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_delete"),
			"Number of conntrack entries which were removed.",
			nil, nil,
		),
		drop: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "nf_conntrack_stat_drop"),
			"Number of packets dropped due to conntrack failure.",
			nil, nil,
		),
	}, nil
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

	data, err := ioutil.ReadFile(procFilePath("net/stat/nf_conntrack"))
	if err != nil {
		return nil
	}
	connStat := parseNfConntrackStat(data)
	ch <- prometheus.MustNewConstMetric(
		c.searched, prometheus.GaugeValue, float64(connStat.Searched))
	ch <- prometheus.MustNewConstMetric(
		c.insert, prometheus.GaugeValue, float64(connStat.Insert))
	ch <- prometheus.MustNewConstMetric(
		c.insertFailed, prometheus.GaugeValue, float64(connStat.InsertFailed))
	ch <- prometheus.MustNewConstMetric(
		c.delete, prometheus.GaugeValue, float64(connStat.Delete))
	ch <- prometheus.MustNewConstMetric(
		c.drop, prometheus.GaugeValue, float64(connStat.Drop))
	return nil
}

func parseNfConntrackStat(data []byte) *ConntrackStatistics {
	connstat := ConntrackStatistics{}

	r := bytes.NewReader(data)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		parseFields(&connstat, fields)
	}

	return &connstat
}

func parseFields(cs *ConntrackStatistics, fields []string) {
	searched, _ := strconv.ParseUint(fields[1], 16, 64)
	cs.Searched += searched

	delete, _ := strconv.ParseUint(fields[6], 16, 64)
	cs.Delete += delete

	insert, _ := strconv.ParseUint(fields[8], 16, 64)
	cs.Insert += insert

	insertFailed, _ := strconv.ParseUint(fields[9], 16, 64)
	cs.InsertFailed += insertFailed

	drop, _ := strconv.ParseUint(fields[10], 16, 64)
	cs.Drop += drop
}
