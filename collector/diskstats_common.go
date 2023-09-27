// Copyright 2019 The Prometheus Authors
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

//go:build !nodiskstats && (openbsd || linux || darwin)
// +build !nodiskstats
// +build openbsd linux darwin

package collector

import (
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	diskSubsystem = "disk"
)

var (
	diskLabelNames = []string{"device"}

	readsCompletedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"),
		"The total number of reads completed successfully.",
		diskLabelNames, nil,
	)

	readBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"),
		"The total number of bytes read successfully.",
		diskLabelNames, nil,
	)

	writesCompletedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"),
		"The total number of writes completed successfully.",
		diskLabelNames, nil,
	)

	writtenBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"),
		"The total number of bytes written successfully.",
		diskLabelNames, nil,
	)

	ioTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "io_time_seconds_total"),
		"Total seconds spent doing I/Os.",
		diskLabelNames, nil,
	)

	readTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"),
		"The total number of seconds spent by all reads.",
		diskLabelNames,
		nil,
	)

	writeTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"),
		"This is the total number of seconds spent by all writes.",
		diskLabelNames,
		nil,
	)
)

type DiskstatsDeviceFilterConfig struct {
	DiskstatsDeviceExclude    *string
	DiskstatsDeviceExcludeSet bool
	OldDiskstatsDeviceExclude *string
	DiskstatsDeviceInclude    *string
}

func newDiskstatsDeviceFilter(config DiskstatsDeviceFilterConfig, logger log.Logger) (deviceFilter, error) {
	if *config.OldDiskstatsDeviceExclude != "" {
		if !config.DiskstatsDeviceExcludeSet {
			level.Warn(logger).Log("msg", "--collector.diskstats.ignored-devices is DEPRECATED and will be removed in 2.0.0, use --collector.diskstats.device-exclude")
			*config.DiskstatsDeviceExclude = *config.OldDiskstatsDeviceExclude
		} else {
			return deviceFilter{}, errors.New("--collector.diskstats.ignored-devices and --collector.diskstats.device-exclude are mutually exclusive")
		}
	}

	if *config.DiskstatsDeviceExclude != "" && *config.DiskstatsDeviceInclude != "" {
		return deviceFilter{}, errors.New("device-exclude & device-include are mutually exclusive")
	}

	if *config.DiskstatsDeviceExclude != "" {
		level.Info(logger).Log("msg", "Parsed flag --collector.diskstats.device-exclude", "flag", *config.DiskstatsDeviceExclude)
	} else {
		*config.DiskstatsDeviceExclude = diskstatsDefaultIgnoredDevices
	}

	if *config.DiskstatsDeviceInclude != "" {
		level.Info(logger).Log("msg", "Parsed Flag --collector.diskstats.device-include", "flag", *config.DiskstatsDeviceInclude)
	}

	return newDeviceFilter(*config.DiskstatsDeviceExclude, *config.DiskstatsDeviceInclude), nil
}
