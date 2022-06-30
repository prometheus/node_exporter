// Copyright 2022 The Prometheus Authors
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

//go:build linux && !noselinux
// +build linux,!noselinux

package collector

import (
	"github.com/go-kit/log"
	"github.com/opencontainers/selinux/go-selinux"
	"github.com/prometheus/client_golang/prometheus"
)

type selinuxCollector struct {
	avcAllocations          *prometheus.Desc
	avcFrees                *prometheus.Desc
	avcHashBucketsAvailable *prometheus.Desc
	avcHashBucketsUsed      *prometheus.Desc
	avcHashEntries          *prometheus.Desc
	avcHashLongestChain     *prometheus.Desc
	avcHits                 *prometheus.Desc
	avcMisses               *prometheus.Desc
	avcLookups              *prometheus.Desc
	avcReclaims             *prometheus.Desc
	avcThreshold            *prometheus.Desc
	configMode              *prometheus.Desc
	currentMode             *prometheus.Desc
	enabled                 *prometheus.Desc
	logger                  log.Logger
}

func init() {
	registerCollector("selinux", defaultEnabled, NewSelinuxCollector)
}

// NewSelinuxCollector returns a new Collector exposing SELinux statistics.
func NewSelinuxCollector(logger log.Logger) (Collector, error) {
	const subsystem = "selinux"

	return &selinuxCollector{
		configMode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "config_mode"),
			"Configured SELinux enforcement mode",
			nil, nil,
		),
		currentMode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "current_mode"),
			"Current SELinux enforcement mode",
			nil, nil,
		),
		enabled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "enabled"),
			"SELinux is enabled, 1 is true, 0 is false",
			nil, nil,
		),
		avcAllocations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_allocations_total"),
			"SELinux AVC allocations",
			nil, nil,
		),
		avcFrees: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_frees_total"),
			"SELinux AVC frees",
			nil, nil,
		),
		avcHashBucketsAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_hash_buckets_available"),
			"SELinux AVC hash buckets available",
			nil, nil,
		),
		avcHashBucketsUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_hash_buckets_used"),
			"SELinux AVC hash buckets used",
			nil, nil,
		),
		avcHashEntries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_hash_entries"),
			"SELinux AVC hash entries",
			nil, nil,
		),
		avcHashLongestChain: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_hash_longest_chain"),
			"SELinux AVC hash longest chain",
			nil, nil,
		),
		avcHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_hits_total"),
			"SELinux AVC hits",
			nil, nil,
		),
		avcMisses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_misses_total"),
			"SELinux AVC misses",
			nil, nil,
		),
		avcLookups: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_lookups_total"),
			"SELinux AVC lookups",
			nil, nil,
		),
		avcReclaims: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_reclaims_total"),
			"SELinux AVC reclaims",
			nil, nil,
		),
		avcThreshold: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "avc_threshold"),
			"SELinux AVC threshold",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *selinuxCollector) Update(ch chan<- prometheus.Metric) error {
	if !selinux.GetEnabled() {
		ch <- prometheus.MustNewConstMetric(
			c.enabled, prometheus.GaugeValue, 0)

		return nil
	}

	ch <- prometheus.MustNewConstMetric(
		c.enabled, prometheus.GaugeValue, 1)

	ch <- prometheus.MustNewConstMetric(
		c.configMode, prometheus.GaugeValue, float64(selinux.DefaultEnforceMode()))

	ch <- prometheus.MustNewConstMetric(
		c.currentMode, prometheus.GaugeValue, float64(selinux.EnforceMode()))

	avcStats, err := getAVCStats("fs/selinux/avc/cache_stats")

	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.avcLookups, prometheus.CounterValue, float64(avcStats["lookups"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcHits, prometheus.CounterValue, float64(avcStats["hits"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcMisses, prometheus.CounterValue, float64(avcStats["misses"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcAllocations, prometheus.CounterValue, float64(avcStats["allocations"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcReclaims, prometheus.CounterValue, float64(avcStats["reclaims"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcFrees, prometheus.CounterValue, float64(avcStats["frees"]))

	avcHashStats, err := getAVCHashStats("fs/selinux/avc/hash_stats")

	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.avcHashEntries, prometheus.GaugeValue, float64(avcHashStats["entries"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcHashBucketsAvailable, prometheus.GaugeValue, float64(avcHashStats["buckets_available"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcHashBucketsUsed, prometheus.GaugeValue, float64(avcHashStats["buckets_used"]))

	ch <- prometheus.MustNewConstMetric(
		c.avcHashLongestChain, prometheus.GaugeValue, float64(avcHashStats["longest_chain"]))

	avcThreshold, err := readUintFromFile(sysFilePath("fs/selinux/avc/cache_threshold"))

	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.avcThreshold, prometheus.GaugeValue, float64(avcThreshold))

	return nil
}
