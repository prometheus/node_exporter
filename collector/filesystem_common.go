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

// +build !nofilesystem
// +build linux freebsd openbsd darwin,amd64 dragonfly

package collector

import (
	"flag"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

// Arch-dependent implementation must define:
// * defIgnoredMountPoints
// * defIgnoredFSTypes
// * filesystemLabelNames
// * filesystemCollector.GetStats

var (
	ignoredMountPoints = flag.String(
		"collector.filesystem.ignored-mount-points",
		defIgnoredMountPoints,
		"Regexp of mount points to ignore for filesystem collector.")

	ignoredFSTypes = flag.String(
		"collector.filesystem.ignored-fs-types",
		defIgnoredFSTypes,
		"Regexp of filesystem types to ignore for filesystem collector.")

	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}
)

type filesystemCollector struct {
	ignoredMountPointsPattern *regexp.Regexp
	ignoredFSTypesPattern     *regexp.Regexp
	sizeDesc, freeDesc, availDesc,
	filesDesc, filesFreeDesc, roDesc *prometheus.Desc
}

type filesystemStats struct {
	labelValues                             []string
	size, free, avail, files, filesFree, ro float64
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Filesystems stats.
func NewFilesystemCollector() (Collector, error) {
	subsystem := "filesystem"
	mountPointPattern := regexp.MustCompile(*ignoredMountPoints)
	filesystemsTypesPattern := regexp.MustCompile(*ignoredFSTypes)

	sizeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "size"),
		"Filesystem size in bytes.",
		filesystemLabelNames, nil,
	)

	freeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "free"),
		"Filesystem free space in bytes.",
		filesystemLabelNames, nil,
	)

	availDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "avail"),
		"Filesystem space available to non-root users in bytes.",
		filesystemLabelNames, nil,
	)

	filesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files"),
		"Filesystem total file nodes.",
		filesystemLabelNames, nil,
	)

	filesFreeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files_free"),
		"Filesystem total free file nodes.",
		filesystemLabelNames, nil,
	)

	roDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "readonly"),
		"Filesystem read-only status.",
		filesystemLabelNames, nil,
	)

	return &filesystemCollector{
		ignoredMountPointsPattern: mountPointPattern,
		ignoredFSTypesPattern:     filesystemsTypesPattern,
		sizeDesc:                  sizeDesc,
		freeDesc:                  freeDesc,
		availDesc:                 availDesc,
		filesDesc:                 filesDesc,
		filesFreeDesc:             filesFreeDesc,
		roDesc:                    roDesc,
	}, nil
}

func (c *filesystemCollector) Update(ch chan<- prometheus.Metric) (err error) {
	stats, err := c.GetStats()
	if err != nil {
		return err
	}
	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.sizeDesc, prometheus.GaugeValue,
			s.size, s.labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.freeDesc, prometheus.GaugeValue,
			s.free, s.labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.availDesc, prometheus.GaugeValue,
			s.avail, s.labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.filesDesc, prometheus.GaugeValue,
			s.files, s.labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.filesFreeDesc, prometheus.GaugeValue,
			s.filesFree, s.labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.roDesc, prometheus.GaugeValue,
			s.ro, s.labelValues...,
		)
	}
	return nil
}
