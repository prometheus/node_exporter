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

package collector

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/config"
)

// Arch-dependent implementation must define:
// * filesystemLabelNames
// * filesystemCollector.GetStats

var filesystemLabelNames = []string{"device", "mountpoint", "fstype", "device_error"}

type filesystemCollector struct {
	mountPointFilter              deviceFilter
	fsTypeFilter                  deviceFilter
	sizeDesc, freeDesc, availDesc *prometheus.Desc
	filesDesc, filesFreeDesc      *prometheus.Desc
	purgeableDesc                 *prometheus.Desc
	roDesc, deviceErrorDesc       *prometheus.Desc
	mountInfoDesc                 *prometheus.Desc
	logger                        *slog.Logger
	config                        config.FilesystemConfig
	paths                         pathConfig
}

type filesystemLabels struct {
	device, mountPoint, fsType, mountOptions, superOptions, deviceError, major, minor string
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
func NewFilesystemCollector(cfg config.Config, logger *slog.Logger) (Collector, error) {
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

	mountPointFilter, err := newMountPointsFilter(cfg.Filesystem, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to parse mount points filter flags: %w", err)
	}

	fsTypeFilter, err := newFSTypeFilter(cfg.Filesystem, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to parse fs types filter flags: %w", err)
	}

	return &filesystemCollector{
		paths:            pathConfig(cfg.Path),
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
		config:           cfg.Filesystem,
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

func newMountPointsFilter(cfg config.FilesystemConfig, logger *slog.Logger) (deviceFilter, error) {
	if cfg.MountPointsExclude != "" && cfg.MountPointsInclude != "" {
		return deviceFilter{}, errors.New("--collector.filesystem.mount-points-exclude and --collector.filesystem.mount-points-include are mutually exclusive")
	}

	if cfg.MountPointsExclude != "" {
		logger.Info("Parsed flag --collector.filesystem.mount-points-exclude", "flag", cfg.MountPointsExclude)
	}
	if cfg.MountPointsInclude != "" {
		logger.Info("Parsed flag --collector.filesystem.mount-points-include", "flag", cfg.MountPointsInclude)
	}

	return newDeviceFilter(cfg.MountPointsExclude, cfg.MountPointsInclude), nil
}

func newFSTypeFilter(cfg config.FilesystemConfig, logger *slog.Logger) (deviceFilter, error) {
	if cfg.FSTypesExclude != "" && cfg.FSTypesInclude != "" {
		return deviceFilter{}, errors.New("--collector.filesystem.fs-types-exclude and --collector.filesystem.fs-types-include are mutually exclusive")
	}

	if cfg.FSTypesExclude != "" {
		logger.Info("Parsed flag --collector.filesystem.fs-types-exclude", "flag", cfg.FSTypesExclude)
	}
	if cfg.FSTypesInclude != "" {
		logger.Info("Parsed flag --collector.filesystem.fs-types-include", "flag", cfg.FSTypesInclude)
	}

	return newDeviceFilter(cfg.FSTypesExclude, cfg.FSTypesInclude), nil
}
