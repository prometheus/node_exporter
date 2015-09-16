// +build !nofilesystem

package collector

import (
	"flag"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

// Arch-dependent implementation must define:
// * defIgnoredMountPoints
// * filesystemLabelNames
// * filesystemCollector.GetStats

var (
	ignoredMountPoints = flag.String(
		"collector.filesystem.ignored-mount-points",
		defIgnoredMountPoints,
		"Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	ignoredMountPointsPattern *regexp.Regexp
	sizeDesc, freeDesc, availDesc,
	filesDesc, filesFreeDesc *prometheus.Desc
}

type filesystemStats struct {
	labelValues                         []string
	size, free, avail, files, filesFree float64
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Filesystems stats.
func NewFilesystemCollector() (Collector, error) {
	subsystem := "filesystem"
	pattern := regexp.MustCompile(*ignoredMountPoints)

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

	return &filesystemCollector{
		ignoredMountPointsPattern: pattern,
		sizeDesc:                  sizeDesc,
		freeDesc:                  freeDesc,
		availDesc:                 availDesc,
		filesDesc:                 filesDesc,
		filesFreeDesc:             filesFreeDesc,
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
	}
	return nil
}
