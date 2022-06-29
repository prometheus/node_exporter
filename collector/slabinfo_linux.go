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

//go:build linux && !noslabinfo
// +build linux,!noslabinfo

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type slabinfoCollector struct {
	fs        procfs.FS
	logger    log.Logger
	subsystem string
	labels    []string
}

func init() {
	registerCollector("slabinfo", defaultDisabled, NewSlabinfoCollector)
}

func NewSlabinfoCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &slabinfoCollector{logger: logger,
		fs:        fs,
		subsystem: "slabinfo",
		labels:    []string{"slab"},
	}, nil
}

func (c *slabinfoCollector) Update(ch chan<- prometheus.Metric) error {
	slabinfo, err := c.fs.SlabInfo()
	if err != nil {
		return fmt.Errorf("couldn't get %s: %w", c.subsystem, err)
	}

	for _, slab := range slabinfo.Slabs {
		ch <- c.activeObjects(slab.Name, slab.ObjActive)
		ch <- c.objects(slab.Name, slab.ObjNum)
		ch <- c.objectSizeBytes(slab.Name, slab.ObjSize)
		ch <- c.objectsPerSlab(slab.Name, slab.ObjPerSlab)
		ch <- c.pagesPerSlab(slab.Name, slab.PagesPerSlab)
	}

	return nil
}

func (c *slabinfoCollector) activeObjects(label string, val int64) prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, c.subsystem, "active_objects"),
		"The number of objects that are currently active (i.e., in use).",
		c.labels, nil)

	return prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(val), label,
	)
}

func (c *slabinfoCollector) objects(label string, val int64) prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, c.subsystem, "objects"),
		"The total number of allocated objects (i.e., objects that are both in use and not in use).",
		c.labels, nil)

	return prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(val), label,
	)
}

func (c *slabinfoCollector) objectSizeBytes(label string, val int64) prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, c.subsystem, "object_size_bytes"),
		"The size of objects in this slab, in bytes.",
		c.labels, nil)

	return prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(val), label,
	)
}

func (c *slabinfoCollector) objectsPerSlab(label string, val int64) prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, c.subsystem, "objects_per_slab"),
		"The number of objects stored in each slab.",
		c.labels, nil)

	return prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(val), label,
	)
}

func (c *slabinfoCollector) pagesPerSlab(label string, val int64) prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, c.subsystem, "pages_per_slab"),
		"The number of pages allocated for each slab.",
		c.labels, nil)

	return prometheus.MustNewConstMetric(
		desc, prometheus.GaugeValue, float64(val), label,
	)
}
