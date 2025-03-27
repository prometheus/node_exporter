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

//go:build !nofilesystem && (linux || freebsd || netbsd || openbsd || darwin || dragonfly || aix)
// +build !nofilesystem
// +build linux freebsd netbsd openbsd darwin dragonfly aix

package collector

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// Arch-dependent implementation must define:
// * defMountPointsExcluded
// * defFSTypesExcluded
// * filesystemLabelNames
// * filesystemCollector.GetStats

var (
	mountPointsExcludeSet bool
	mountPointsExclude    = kingpin.Flag(
		"collector.filesystem.mount-points-exclude",
		"Regexp of mount points to exclude for filesystem collector. (mutually exclusive to mount-points-include)",
	).Default(defMountPointsExcluded).PreAction(func(c *kingpin.ParseContext) error {
		mountPointsExcludeSet = true
		return nil
	}).String()
	oldMountPointsExcluded = kingpin.Flag(
		"collector.filesystem.ignored-mount-points",
		"Regexp of mount points to ignore for filesystem collector.",
	).Hidden().String()
	mountPointsInclude = kingpin.Flag(
		"collector.filesystem.mount-points-include",
		"Regexp of mount points to include for filesystem collector. (mutually exclusive to mount-points-exclude)",
	).String()

	fsTypesExcludeSet bool
	fsTypesExclude    = kingpin.Flag(
		"collector.filesystem.fs-types-exclude",
		"Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-include)",
	).Default(defFSTypesExcluded).PreAction(func(c *kingpin.ParseContext) error {
		fsTypesExcludeSet = true
		return nil
	}).String()
	oldFSTypesExcluded = kingpin.Flag(
		"collector.filesystem.ignored-fs-types",
		"Regexp of filesystem types to ignore for filesystem collector.",
	).Hidden().String()
	fsTypesInclude = kingpin.Flag(
		"collector.filesystem.fs-types-include",
		"Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-exclude)",
	).String()

	filesystemLabelNames = []string{"device", "mountpoint", "fstype", "device_error"}
)

type filesystemCollector struct {
	mountPointFilter              deviceFilter
	fsTypeFilter                  deviceFilter
	sizeDesc, freeDesc, availDesc *prometheus.Desc
	filesDesc, filesFreeDesc      *prometheus.Desc
	purgeableDesc                 *prometheus.Desc
	roDesc, deviceErrorDesc       *prometheus.Desc
	mountInfoDesc                 *prometheus.Desc
	logger                        *slog.Logger
}

type filesystemLabels struct {
	device, mountPoint, fsType, options, deviceError, major, minor string
}

type filesystemStats struct {
	labels            filesystemLabels
	size, free, avail float64
	files, filesFree  float64
	purgeable         float64
	ro, deviceError   float64
}

func init() {
	registerCollector("filesystem", defaultEnabled, NewFilesystemCollector)
}

// NewFilesystemCollector returns a new Collector exposing filesystems stats.
func NewFilesystemCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "filesystem"

	sizeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "size_bytes"),
		"Filesystem size in bytes.",
		filesystemLabelNames, nil,
	)

	freeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "free_bytes"),
		"Filesystem free space in bytes.",
		filesystemLabelNames, nil,
	)

	availDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "avail_bytes"),
		"Filesystem space available to non-root users in bytes.",
		filesystemLabelNames, nil,
	)

	filesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "files"),
		"Filesystem total file nodes.",
		filesystemLabelNames, nil,
	)

	filesFreeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "files_free"),
		"Filesystem total free file nodes.",
		filesystemLabelNames, nil,
	)

	purgeableDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "purgeable_bytes"),
		"Filesystem space available including purgeable space (MacOS specific).",
		filesystemLabelNames, nil,
	)

	roDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "readonly"),
		"Filesystem read-only status.",
		filesystemLabelNames, nil,
	)

	deviceErrorDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "device_error"),
		"Whether an error occurred while getting statistics for the given device.",
		filesystemLabelNames, nil,
	)

	mountInfoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "mount_info"),
		"Filesystem mount information.",
		[]string{"device", "major", "minor", "mountpoint"},
		nil,
	)

	mountPointFilter, err := newMountPointsFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("unable to parse mount points filter flags: %w", err)
	}

	fsTypeFilter, err := newFSTypeFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("unable to parse fs types filter flags: %w", err)
	}

	return &filesystemCollector{
		mountPointFilter: mountPointFilter,
		fsTypeFilter:     fsTypeFilter,
		sizeDesc:         sizeDesc,
		freeDesc:         freeDesc,
		availDesc:        availDesc,
		filesDesc:        filesDesc,
		filesFreeDesc:    filesFreeDesc,
		purgeableDesc:    purgeableDesc,
		roDesc:           roDesc,
		deviceErrorDesc:  deviceErrorDesc,
		mountInfoDesc:    mountInfoDesc,
		logger:           logger,
	}, nil
}

func (c *filesystemCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.GetStats()
	if err != nil {
		return err
	}
	// Make sure we expose a metric once, even if there are multiple mounts
	seen := map[filesystemLabels]bool{}
	for _, s := range stats {
		if seen[s.labels] {
			continue
		}
		seen[s.labels] = true

		ch <- prometheus.MustNewConstMetric(
			c.deviceErrorDesc, prometheus.GaugeValue,
			s.deviceError, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.roDesc, prometheus.GaugeValue,
			s.ro, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)

		if s.deviceError > 0 {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.sizeDesc, prometheus.GaugeValue,
			s.size, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.freeDesc, prometheus.GaugeValue,
			s.free, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.availDesc, prometheus.GaugeValue,
			s.avail, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.filesDesc, prometheus.GaugeValue,
			s.files, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.filesFreeDesc, prometheus.GaugeValue,
			s.filesFree, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
		)
		ch <- prometheus.MustNewConstMetric(
			c.mountInfoDesc, prometheus.GaugeValue,
			1.0, s.labels.device, s.labels.major, s.labels.minor, s.labels.mountPoint,
		)
		if s.purgeable >= 0 {
			ch <- prometheus.MustNewConstMetric(
				c.purgeableDesc, prometheus.GaugeValue,
				s.purgeable, s.labels.device, s.labels.mountPoint, s.labels.fsType, s.labels.deviceError,
			)
		}
	}
	return nil
}

func newMountPointsFilter(logger *slog.Logger) (deviceFilter, error) {
	if *oldMountPointsExcluded != "" {
		if !mountPointsExcludeSet {
			logger.Warn("--collector.filesystem.ignored-mount-points is DEPRECATED and will be removed in 2.0.0, use --collector.filesystem.mount-points-exclude")
			*mountPointsExclude = *oldMountPointsExcluded
		} else {
			return deviceFilter{}, errors.New("--collector.filesystem.ignored-mount-points and --collector.filesystem.mount-points-exclude are mutually exclusive")
		}
	}

	if *mountPointsInclude != "" && !mountPointsExcludeSet {
		logger.Debug("mount-points-exclude flag not set when mount-points-include flag is set, assuming include is desired")
		*mountPointsExclude = ""
	}

	if *mountPointsExclude != "" && *mountPointsInclude != "" {
		return deviceFilter{}, errors.New("--collector.filesystem.mount-points-exclude and --collector.filesystem.mount-points-include are mutually exclusive")
	}

	if *mountPointsExclude != "" {
		logger.Info("Parsed flag --collector.filesystem.mount-points-exclude", "flag", *mountPointsExclude)
	}
	if *mountPointsInclude != "" {
		logger.Info("Parsed flag --collector.filesystem.mount-points-include", "flag", *mountPointsInclude)
	}

	return newDeviceFilter(*mountPointsExclude, *mountPointsInclude), nil
}

func newFSTypeFilter(logger *slog.Logger) (deviceFilter, error) {
	if *oldFSTypesExcluded != "" {
		if !fsTypesExcludeSet {
			logger.Warn("--collector.filesystem.ignored-fs-types is DEPRECATED and will be removed in 2.0.0, use --collector.filesystem.fs-types-exclude")
			*fsTypesExclude = *oldFSTypesExcluded
		} else {
			return deviceFilter{}, errors.New("--collector.filesystem.ignored-fs-types and --collector.filesystem.fs-types-exclude are mutually exclusive")
		}
	}

	if *fsTypesInclude != "" && !fsTypesExcludeSet {
		logger.Debug("fs-types-exclude flag not set when fs-types-include flag is set, assuming include is desired")
		*fsTypesExclude = ""
	}

	if *fsTypesExclude != "" && *fsTypesInclude != "" {
		return deviceFilter{}, errors.New("--collector.filesystem.fs-types-exclude and --collector.filesystem.fs-types-include are mutually exclusive")
	}

	if *fsTypesExclude != "" {
		logger.Info("Parsed flag --collector.filesystem.fs-types-exclude", "flag", *fsTypesExclude)
	}
	if *fsTypesInclude != "" {
		logger.Info("Parsed flag --collector.filesystem.fs-types-include", "flag", *fsTypesInclude)
	}

	return newDeviceFilter(*fsTypesExclude, *fsTypesInclude), nil
}
