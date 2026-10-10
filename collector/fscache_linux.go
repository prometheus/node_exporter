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

//go:build !nofscache
// +build !nofscache

package collector

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	fscacheSubsystem = "fscache"
)

type fscacheCollector struct {
	// Metrics
	objectsAllocated *prometheus.Desc
	objectsAvailable *prometheus.Desc

	acquireAttempts *prometheus.Desc
	acquireSuccess  *prometheus.Desc

	lookupsTotal    *prometheus.Desc
	lookupsPositive *prometheus.Desc
	lookupsNegative *prometheus.Desc

	invalidationsTotal *prometheus.Desc

	updatesTotal *prometheus.Desc

	relinquishesTotal *prometheus.Desc

	attributeChangesTotal   *prometheus.Desc
	attributeChangesSuccess *prometheus.Desc

	allocationsTotal   *prometheus.Desc
	allocationsSuccess *prometheus.Desc

	retrievalsTotal    *prometheus.Desc
	retrievalsSuccess  *prometheus.Desc
	retrievalsNoBuffer *prometheus.Desc

	storesTotal   *prometheus.Desc
	storesSuccess *prometheus.Desc

	cacheEventRetired *prometheus.Desc
	cacheEventCulled  *prometheus.Desc

	fs     procfs.FS
	logger *slog.Logger
}

var _ prometheus.Collector = (*fscacheCollector)(nil)

func init() {
	registerCollector("fscache", defaultEnabled, func(logger *slog.Logger) (Collector, error) {
		return NewFscacheCollector(logger)
	})
}

// NewFscacheCollector returns a new Collector exposing fscache stats.
func NewFscacheCollector(logger *slog.Logger) (*fscacheCollector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &fscacheCollector{
		objectsAllocated: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "objects_allocated_total"),
			"Number of index cookies allocated (Cookies: idx=allocated/available/unused).",
			nil, nil,
		),
		objectsAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "objects_available_total"),
			"Number of index cookies available (Cookies: idx=allocated/available/unused).",
			nil, nil,
		),
		acquireAttempts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "acquire_attempts_total"),
			"Number of acquire operations attempted (Acquire: n=attempts).",
			nil, nil,
		),
		acquireSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "acquire_success_total"),
			"Number of acquire operations successful (Acquire: ok=success).",
			nil, nil,
		),
		lookupsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "lookups_total"),
			"Number of lookup operations (Lookups: n=tot).",
			nil, nil,
		),
		lookupsPositive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "lookups_positive_total"),
			"Number of positive lookup operations (Lookups: pos=positive).",
			nil, nil,
		),
		lookupsNegative: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "lookups_negative_total"),
			"Number of negative lookup operations (Lookups: neg=negative).",
			nil, nil,
		),
		invalidationsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "invalidations_total"),
			"Number of invalidation operations (Invals: n=tot).",
			nil, nil,
		),
		updatesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "updates_total"),
			"Number of update operations (Updates: n=tot).",
			nil, nil,
		),
		relinquishesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "relinquishes_total"),
			"Number of relinquish operations (Relinqs: n=tot).",
			nil, nil,
		),
		attributeChangesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "attribute_changes_total"),
			"Number of attribute change operations attempted (AttrChg: n=attempts).",
			nil, nil,
		),
		attributeChangesSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "attribute_changes_success_total"),
			"Number of successful attribute change operations (AttrChg: ok=success).",
			nil, nil,
		),
		allocationsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "allocations_total"),
			"Number of allocation operations attempted (Allocs: n=attempts).",
			nil, nil,
		),
		allocationsSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "allocations_success_total"),
			"Number of successful allocation operations (Allocs: ok=success).",
			nil, nil,
		),
		retrievalsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "retrievals_total"),
			"Number of retrieval (read) operations attempted (Retrvls: n=attempts).",
			nil, nil,
		),
		retrievalsSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "retrievals_success_total"),
			"Number of successful retrieval (read) operations (Retrvls: ok=success).",
			nil, nil,
		),
		retrievalsNoBuffer: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "retrievals_nobuffer_total"),
			"Number of retrieval (read) operations failed due to no buffer (Retrvls: nbf=nobuff).",
			nil, nil,
		),
		storesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "stores_total"),
			"Number of store (write) operations attempted (Stores: n=attempts).",
			nil, nil,
		),
		storesSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "stores_success_total"),
			"Number of successful store (write) operations (Stores: ok=success).",
			nil, nil,
		),
		cacheEventRetired: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "objects_retired_total"),
			"Number of objects retired (CacheEv: rtr=retired).",
			nil, nil,
		),
		cacheEventCulled: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fscacheSubsystem, "objects_culled_total"),
			"Number of objects culled (CacheEv: cul=culled).",
			nil, nil,
		),
		fs:     fs,
		logger: logger,
	}, nil
}

// Describe implements prometheus.Collector.
func (c *fscacheCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.objectsAllocated
	ch <- c.objectsAvailable
	ch <- c.acquireAttempts
	ch <- c.acquireSuccess
	ch <- c.lookupsTotal
	ch <- c.lookupsPositive
	ch <- c.lookupsNegative
	ch <- c.invalidationsTotal
	ch <- c.updatesTotal
	ch <- c.relinquishesTotal
	ch <- c.attributeChangesTotal
	ch <- c.attributeChangesSuccess
	ch <- c.allocationsTotal
	ch <- c.allocationsSuccess
	ch <- c.retrievalsTotal
	ch <- c.retrievalsSuccess
	ch <- c.retrievalsNoBuffer
	ch <- c.storesTotal
	ch <- c.storesSuccess
	ch <- c.cacheEventRetired
	ch <- c.cacheEventCulled
}

// Collect implements prometheus.Collector.
func (c *fscacheCollector) Collect(ch chan<- prometheus.Metric) {
	// Let the collector helper handle scrape success/failure based on Update's error.
	if err := c.Update(ch); err != nil {
		// Optionally log the error, but don't send invalid metrics here.
		c.logger.Error("Error updating fscache stats", "err", err)
	}
}

// Update gathers metrics from the fscache subsystem.
func (c *fscacheCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.Fscacheinfo()
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Debug("Not collecting fscache statistics, as /proc/fs/fscache/stats is not available", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get fscache stats: %w", err)
	}

	ch <- prometheus.MustNewConstMetric(c.objectsAllocated, prometheus.CounterValue, float64(stats.IndexCookiesAllocated))
	ch <- prometheus.MustNewConstMetric(c.objectsAvailable, prometheus.CounterValue, float64(stats.ObjectsAvailable))
	ch <- prometheus.MustNewConstMetric(c.acquireAttempts, prometheus.CounterValue, float64(stats.AcquireCookiesRequestSeen))
	ch <- prometheus.MustNewConstMetric(c.acquireSuccess, prometheus.CounterValue, float64(stats.AcquireRequestsSucceeded))
	ch <- prometheus.MustNewConstMetric(c.lookupsTotal, prometheus.CounterValue, float64(stats.LookupsNumber))
	ch <- prometheus.MustNewConstMetric(c.lookupsPositive, prometheus.CounterValue, float64(stats.LookupsPositive))
	ch <- prometheus.MustNewConstMetric(c.lookupsNegative, prometheus.CounterValue, float64(stats.LookupsNegative))
	ch <- prometheus.MustNewConstMetric(c.invalidationsTotal, prometheus.CounterValue, float64(stats.InvalidationsNumber))
	ch <- prometheus.MustNewConstMetric(c.updatesTotal, prometheus.CounterValue, float64(stats.UpdateCookieRequestSeen))
	ch <- prometheus.MustNewConstMetric(c.relinquishesTotal, prometheus.CounterValue, float64(stats.RelinquishCookiesRequestSeen))
	ch <- prometheus.MustNewConstMetric(c.attributeChangesTotal, prometheus.CounterValue, float64(stats.AttributeChangedRequestsSeen))
	ch <- prometheus.MustNewConstMetric(c.attributeChangesSuccess, prometheus.CounterValue, float64(stats.AttributeChangedOps))
	ch <- prometheus.MustNewConstMetric(c.allocationsTotal, prometheus.CounterValue, float64(stats.AllocationRequestsSeen))
	ch <- prometheus.MustNewConstMetric(c.allocationsSuccess, prometheus.CounterValue, float64(stats.AllocationOkRequests))
	ch <- prometheus.MustNewConstMetric(c.retrievalsTotal, prometheus.CounterValue, float64(stats.RetrievalsReadRequests))
	ch <- prometheus.MustNewConstMetric(c.retrievalsSuccess, prometheus.CounterValue, float64(stats.RetrievalsOk))
	ch <- prometheus.MustNewConstMetric(c.retrievalsNoBuffer, prometheus.CounterValue, float64(stats.RetrievalsRejectedDueToEnobufs))
	ch <- prometheus.MustNewConstMetric(c.storesTotal, prometheus.CounterValue, float64(stats.StoreWriteRequests))
	ch <- prometheus.MustNewConstMetric(c.storesSuccess, prometheus.CounterValue, float64(stats.StoreSuccessfulRequests))
	ch <- prometheus.MustNewConstMetric(c.cacheEventRetired, prometheus.CounterValue, float64(stats.CacheevRetiredWhenReliquished))
	ch <- prometheus.MustNewConstMetric(c.cacheEventCulled, prometheus.CounterValue, float64(stats.CacheevObjectsCulled))

	return nil
}
