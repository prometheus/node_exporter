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

//go:build !nobtrfs
// +build !nobtrfs

package collector

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/btrfs"
)

// A btrfsCollector is a Collector which gathers metrics from Btrfs filesystems.
type btrfsCollector struct {
	fs     btrfs.FS
	logger log.Logger
}

func init() {
	registerCollector("btrfs", defaultEnabled, NewBtrfsCollector)
}

// NewBtrfsCollector returns a new Collector exposing Btrfs statistics.
func NewBtrfsCollector(logger log.Logger) (Collector, error) {
	fs, err := btrfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &btrfsCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update retrieves and exports Btrfs statistics.
// It implements Collector.
func (c *btrfsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.Stats()
	if err != nil {
		return fmt.Errorf("failed to retrieve Btrfs stats from procfs: %w", err)
	}

	ioctlStatsMap, err := c.getIoctlStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve Btrfs stats with ioctl: %w", err)
	}

	for _, s := range stats {
		// match up procfs and ioctl info by filesystem UUID
		ioctlStats := ioctlStatsMap[s.UUID]
		c.updateBtrfsStats(ch, s, ioctlStats)
	}

	return nil
}

type ioctlFsDeviceStats struct {
	path string
	uuid string

	bytesUsed  uint64
	totalBytes uint64

	errorStats map[string]uint64
}

type ioctlFsStats struct {
	uuid    string
	devices []ioctlFsDeviceStats
}

// Magic constants for ioctl
//nolint:revive
const (
	_BTRFS_IOC_FS_INFO       = 0x8400941F
	_BTRFS_IOC_DEV_INFO      = 0xD000941E
	_BTRFS_IOC_GET_DEV_STATS = 0x00c4089434
)

type _UuidBytes [16]byte

func (id _UuidBytes) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}

//name matches linux struct
//nolint:revive
type btrfs_ioctl_fs_info_args struct {
	maxID          uint64          // out
	numDevices     uint64          // out
	fsID           _UuidBytes      // out
	nodeSize       uint32          // out
	sectorSize     uint32          // out
	cloneAlignment uint32          // out
	_              [122*8 + 4]byte // pad to 1k
}

func (c *btrfsCollector) getIoctlStats() (map[string]*ioctlFsStats, error) {
	// Instead of introducing more ioctl calls to scan for all btrfs
	// filesytems re-use our mount point utils to find known mounts
	mountsList, err := mountPointDetails(c.logger)
	if err != nil {
		return nil, err
	}

	// track devices we have successfully scanned, by device path
	devicesDone := make(map[string]struct{})
	// filesystems scann results by UUID
	fsStats := make(map[string]*ioctlFsStats)

	for _, mount := range mountsList {
		if mount.fsType != "btrfs" {
			continue
		}

		if _, found := devicesDone[mount.device]; found {
			// We already found this filesystem by another mount point
			continue
		}

		fd, err := os.Open(mount.mountPoint)
		if err != nil {
			// failed to open this mount point, maybe we didn't have permission
			// maybe we'll find another mount point for this FS later
			continue
		}

		var fsInfo = btrfs_ioctl_fs_info_args{}
		_, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			fd.Fd(),
			uintptr(_BTRFS_IOC_FS_INFO),
			uintptr(unsafe.Pointer(&fsInfo)))
		if errno != 0 {
			// Failed to get the FS info for some reason,
			// perhaps it'll work with a different mount point
			continue
		}

		fsID := fsInfo.fsID.String()
		if _, found := fsStats[fsID]; found {
			// We already found this filesystem by another mount point
			continue
		}

		deviceStats, err := c.getIoctlDeviceStats(fd, &fsInfo)
		if err != nil {
			return nil, err
		}

		devicesDone[mount.device] = struct{}{}
		fsStats[fsID] = &ioctlFsStats{
			uuid:    fsID,
			devices: deviceStats,
		}
	}

	return fsStats, nil
}

// Known/supported device stats fields
//nolint:revive
const (
	// direct indicators of I/O failures:

	_BTRFS_DEV_STAT_WRITE_ERRS = iota
	_BTRFS_DEV_STAT_READ_ERRS
	_BTRFS_DEV_STAT_FLUSH_ERRS

	// indirect indicators of I/O failures:

	_BTRFS_DEV_STAT_CORRUPTION_ERRS // checksum error, bytenr error or contents is illegal
	_BTRFS_DEV_STAT_GENERATION_ERRS // an indication that blocks have not been written

	_BTRFS_DEV_STAT_VALUES_MAX // counter to indicate the number of known stats we support
)

var errorStatFields = []int{
	_BTRFS_DEV_STAT_WRITE_ERRS,
	_BTRFS_DEV_STAT_READ_ERRS,
	_BTRFS_DEV_STAT_FLUSH_ERRS,
	_BTRFS_DEV_STAT_CORRUPTION_ERRS,
	_BTRFS_DEV_STAT_GENERATION_ERRS,
}
var errorStatNames = []string{
	"write",
	"read",
	"flush",
	"corruption",
	"generation",
}

//name matches linux struct
//nolint:revive
type btrfs_ioctl_dev_info_args struct {
	deviceID   uint64      // in/out
	uuid       _UuidBytes  // in/out
	bytesUsed  uint64      // out
	totalBytes uint64      // out
	_          [379]uint64 // pad to 4k
	path       [1024]byte  // out
}

//name matches linux struct
//nolint:revive
type btrfs_ioctl_get_dev_stats struct {
	deviceID  uint64                                       // in
	itemCount uint64                                       // in/out
	flags     uint64                                       // in/out
	values    [_BTRFS_DEV_STAT_VALUES_MAX]uint64           // out values
	_         [128 - 2 - _BTRFS_DEV_STAT_VALUES_MAX]uint64 // pad to 1k
}

func (c *btrfsCollector) getIoctlDeviceStats(fd *os.File, fsInfo *btrfs_ioctl_fs_info_args) ([]ioctlFsDeviceStats, error) {
	devices := make([]ioctlFsDeviceStats, 0, fsInfo.numDevices)

	var deviceInfo btrfs_ioctl_dev_info_args
	var deviceStats btrfs_ioctl_get_dev_stats

	for i := uint64(0); i <= fsInfo.maxID; i++ {
		deviceInfo = btrfs_ioctl_dev_info_args{
			deviceID: i,
		}
		deviceStats = btrfs_ioctl_get_dev_stats{
			deviceID:  i,
			itemCount: _BTRFS_DEV_STAT_VALUES_MAX,
		}

		_, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			fd.Fd(),
			uintptr(_BTRFS_IOC_DEV_INFO),
			uintptr(unsafe.Pointer(&deviceInfo)))

		if errno == syscall.ENODEV {
			// device IDs do not consistently start at 0, so we expect this
			continue
		}
		if errno != 0 {
			return nil, errno
		}

		_, _, errno = syscall.Syscall(
			syscall.SYS_IOCTL,
			fd.Fd(),
			uintptr(_BTRFS_IOC_GET_DEV_STATS),
			uintptr(unsafe.Pointer(&deviceStats)))

		if errno != 0 {
			return nil, errno
		}

		errorStats := make(map[string]uint64, deviceStats.itemCount)
		for i, fieldIndex := range errorStatFields {
			if int(deviceStats.itemCount) >= fieldIndex {
				errorStats[errorStatNames[i]] = deviceStats.values[fieldIndex]
			}
		}

		devices = append(devices, ioctlFsDeviceStats{
			path:       string(bytes.Trim(deviceInfo.path[:], "\x00")),
			uuid:       deviceInfo.uuid.String(),
			bytesUsed:  deviceInfo.bytesUsed,
			totalBytes: deviceInfo.totalBytes,
			errorStats: errorStats,
		})

		if uint64(len(devices)) == fsInfo.numDevices {
			break
		}
	}

	return devices, nil
}

// btrfsMetric represents a single Btrfs metric that is converted into a Prometheus Metric.
type btrfsMetric struct {
	name            string
	metricType      prometheus.ValueType
	desc            string
	value           float64
	extraLabel      []string
	extraLabelValue []string
}

// updateBtrfsStats collects statistics for one bcache ID.
func (c *btrfsCollector) updateBtrfsStats(ch chan<- prometheus.Metric, s *btrfs.Stats, iocStats *ioctlFsStats) {
	const subsystem = "btrfs"

	// Basic information about the filesystem.
	devLabels := []string{"uuid"}

	// Retrieve the metrics.
	metrics := c.getMetrics(s, iocStats)

	// Convert all gathered metrics to Prometheus Metrics and add to channel.
	for _, m := range metrics {
		labels := append(devLabels, m.extraLabel...)

		desc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, m.name),
			m.desc,
			labels,
			nil,
		)

		labelValues := []string{s.UUID}
		if len(m.extraLabelValue) > 0 {
			labelValues = append(labelValues, m.extraLabelValue...)
		}

		ch <- prometheus.MustNewConstMetric(
			desc,
			m.metricType,
			m.value,
			labelValues...,
		)
	}
}

// getMetrics returns metrics for the given Btrfs statistics.
func (c *btrfsCollector) getMetrics(s *btrfs.Stats, iocStats *ioctlFsStats) []btrfsMetric {
	metrics := []btrfsMetric{
		{
			name:            "info",
			desc:            "Filesystem information",
			value:           1,
			metricType:      prometheus.GaugeValue,
			extraLabel:      []string{"label"},
			extraLabelValue: []string{s.Label},
		},
		{
			name:       "global_rsv_size_bytes",
			desc:       "Size of global reserve.",
			metricType: prometheus.GaugeValue,
			value:      float64(s.Allocation.GlobalRsvSize),
		},
	}

	// Information about devices.
	if iocStats == nil {
		for n, dev := range s.Devices {
			metrics = append(metrics, btrfsMetric{
				name:            "device_size_bytes",
				desc:            "Size of a device that is part of the filesystem.",
				metricType:      prometheus.GaugeValue,
				value:           float64(dev.Size),
				extraLabel:      []string{"device"},
				extraLabelValue: []string{n},
			})
		}
	} else {
		for _, dev := range iocStats.devices {
			// trim the /dev/ prefix from the device name so the value should match
			// the value used in the fallback branch above
			device := strings.TrimPrefix(dev.path, "/dev/")

			extraLabels := []string{"device", "device_uuid"}
			extraLabelValues := []string{device, dev.uuid}

			metrics = append(metrics,
				btrfsMetric{
					name:            "device_size_bytes",
					desc:            "Size of a device that is part of the filesystem.",
					metricType:      prometheus.GaugeValue,
					value:           float64(dev.totalBytes),
					extraLabel:      extraLabels,
					extraLabelValue: extraLabelValues,
				},
				// A bytes available metric is probably more useful than a
				// bytes used metric, because large numbers of bytes will
				// suffer from floating point representation issues
				// and we probably care more about the number when it's low anyway
				btrfsMetric{
					name:            "device_unused_bytes",
					desc:            "Unused bytes unused on a device that is part of the filesystem.",
					metricType:      prometheus.GaugeValue,
					value:           float64(dev.totalBytes - dev.bytesUsed),
					extraLabel:      extraLabels,
					extraLabelValue: extraLabelValues,
				})

			errorLabels := append([]string{"type"}, extraLabels...)
			for errorName, count := range dev.errorStats {
				errorLabelValues := append([]string{errorName}, extraLabelValues...)
				metrics = append(metrics,
					btrfsMetric{
						name:            "device_errors_total",
						desc:            "Errors reported for the device",
						metricType:      prometheus.CounterValue,
						value:           float64(count),
						extraLabel:      errorLabels,
						extraLabelValue: errorLabelValues,
					})
			}
		}
	}

	// Information about data, metadata and system data.
	metrics = append(metrics, c.getAllocationStats("data", s.Allocation.Data)...)
	metrics = append(metrics, c.getAllocationStats("metadata", s.Allocation.Metadata)...)
	metrics = append(metrics, c.getAllocationStats("system", s.Allocation.System)...)

	return metrics
}

// getAllocationStats returns allocation metrics for the given Btrfs Allocation statistics.
func (c *btrfsCollector) getAllocationStats(a string, s *btrfs.AllocationStats) []btrfsMetric {
	metrics := []btrfsMetric{
		{
			name:            "reserved_bytes",
			desc:            "Amount of space reserved for a data type",
			metricType:      prometheus.GaugeValue,
			value:           float64(s.ReservedBytes),
			extraLabel:      []string{"block_group_type"},
			extraLabelValue: []string{a},
		},
	}

	// Add all layout statistics.
	for layout, stats := range s.Layouts {
		metrics = append(metrics, c.getLayoutStats(a, layout, stats)...)
	}

	return metrics
}

// getLayoutStats returns metrics for a data layout.
func (c *btrfsCollector) getLayoutStats(a, l string, s *btrfs.LayoutUsage) []btrfsMetric {
	return []btrfsMetric{
		{
			name:            "used_bytes",
			desc:            "Amount of used space by a layout/data type",
			metricType:      prometheus.GaugeValue,
			value:           float64(s.UsedBytes),
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
		{
			name:            "size_bytes",
			desc:            "Amount of space allocated for a layout/data type",
			metricType:      prometheus.GaugeValue,
			value:           float64(s.TotalBytes),
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
		{
			name:            "allocation_ratio",
			desc:            "Data allocation ratio for a layout/data type",
			metricType:      prometheus.GaugeValue,
			value:           s.Ratio,
			extraLabel:      []string{"block_group_type", "mode"},
			extraLabelValue: []string{a, l},
		},
	}
}
