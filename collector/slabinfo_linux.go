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

//go:build !noslabinfo
// +build !noslabinfo

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const (
	slabSubsystem = "slab"
)

type slabinfoCollector struct {
	fs         procfs.FS
	activeObjs *prometheus.Desc
	numObjs    *prometheus.Desc
	objSize    *prometheus.Desc
	logger     log.Logger
}

func init() {
	registerCollector("slabinfo", defaultDisabled, NewSlabinfoCollector)
}

// NewSlabinfoCollector returns a new Collector exposing slabinfo stats.
// https://www.kernel.org/doc/Documentation/vm/slub.txt
// It requires CAP_SYS_ADMIN privilege to read /proc/slabinfo
func NewSlabinfoCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &slabinfoCollector{
		fs: fs,
		activeObjs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, slabSubsystem, "active_objs"),
			"Active Objects",
			[]string{"name"}, nil),
		numObjs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, slabSubsystem, "num_objs"),
			"Number of Objects",
			[]string{"name"}, nil),
		objSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, slabSubsystem, "obj_size"),
			"Object Size",
			[]string{"name"}, nil),
		logger: logger,
	}, nil
}

func (c *slabinfoCollector) Update(ch chan<- prometheus.Metric) error {
	slabinfo, err := c.fs.SlabInfo()
	if err != nil {
		return fmt.Errorf("couldn't get slabinfo: %w", err)
	}

	for _, slab := range slabinfo.Slabs {
		name := slab.Name

		ch <- prometheus.MustNewConstMetric(
			c.activeObjs,
			prometheus.GaugeValue,
			float64(slab.ObjActive),
			name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.numObjs,
			prometheus.GaugeValue,
			float64(slab.ObjNum),
			name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.objSize,
			prometheus.GaugeValue,
			float64(slab.ObjSize),
			name,
		)
	}
	return nil
}
