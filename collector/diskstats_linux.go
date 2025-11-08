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

//go:build !nodiskstats
// +build !nodiskstats

package collector

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/blockdevice"
)

const (
	secondsPerTick = 1.0 / 1000.0

	// Read sectors and write sectors are the "standard UNIX 512-byte sectors, not any device- or filesystem-specific block size."
	// See also https://www.kernel.org/doc/Documentation/block/stat.txt
	unixSectorSize = 512.0

	diskstatsDefaultIgnoredDevices = "^(z?ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"

	// See udevadm(8).
	udevDevicePropertyPrefix = "E:"

	// Udev device properties.
	udevDMLVLayer               = "DM_LV_LAYER"
	udevDMLVName                = "DM_LV_NAME"
	udevDMName                  = "DM_NAME"
	udevDMUUID                  = "DM_UUID"
	udevDMVGName                = "DM_VG_NAME"
	udevIDATA                   = "ID_ATA"
	udevIDATARotationRateRPM    = "ID_ATA_ROTATION_RATE_RPM"
	udevIDATASATA               = "ID_ATA_SATA"
	udevIDATASATASignalRateGen1 = "ID_ATA_SATA_SIGNAL_RATE_GEN1"
	udevIDATASATASignalRateGen2 = "ID_ATA_SATA_SIGNAL_RATE_GEN2"
	udevIDATAWriteCache         = "ID_ATA_WRITE_CACHE"
	udevIDATAWriteCacheEnabled  = "ID_ATA_WRITE_CACHE_ENABLED"
	udevIDFSType                = "ID_FS_TYPE"
	udevIDFSUsage               = "ID_FS_USAGE"
	udevIDFSUUID                = "ID_FS_UUID"
	udevIDFSVersion             = "ID_FS_VERSION"
	udevIDModel                 = "ID_MODEL"
	udevIDPath                  = "ID_PATH"
	udevIDRevision              = "ID_REVISION"
	udevIDSerial                = "ID_SERIAL"
	udevIDSerialShort           = "ID_SERIAL_SHORT"
	udevIDWWN                   = "ID_WWN"
	udevSCSIIdentSerial         = "SCSI_IDENT_SERIAL"
)

type udevInfo map[string]string

type diskstatsCollector struct {
	deviceFilter            deviceFilter
	fs                      blockdevice.FS
	infoDesc                typedDesc
	descs                   []typedDesc
	filesystemInfoDesc      typedDesc
	deviceMapperInfoDesc    typedDesc
	ataDescs                map[string]typedDesc
	logger                  *slog.Logger
	getUdevDeviceProperties func(uint32, uint32) (udevInfo, error)
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
func NewDiskstatsCollector(logger *slog.Logger) (Collector, error) {
	var diskLabelNames = []string{"device"}
	fs, err := blockdevice.NewFS(*procPath, *sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	deviceFilter, err := newDiskstatsDeviceFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device filter flags: %w", err)
	}

	collector := diskstatsCollector{
		deviceFilter: deviceFilter,
		fs:           fs,
		infoDesc: typedDesc{
			desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "info"),
				"Info of /sys/block/<block_device>.",
				[]string{"device", "major", "minor", "path", "wwn", "model", "serial", "revision", "rotational"},
				nil,
			), valueType: prometheus.GaugeValue,
		},
		descs: []typedDesc{
			{
				desc: readsCompletedDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "reads_merged_total"),
					"The total number of reads merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: readBytesDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: readTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: writesCompletedDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "writes_merged_total"),
					"The number of writes merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: writtenBytesDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: writeTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "io_now"),
					"The number of I/Os currently in progress.",
					diskLabelNames,
					nil,
				), valueType: prometheus.GaugeValue,
			},
			{
				desc: ioTimeSecondsDesc, valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"),
					"The weighted # of seconds spent doing I/Os.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discards_completed_total"),
					"The total number of discards completed successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discards_merged_total"),
					"The total number of discards merged.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discarded_sectors_total"),
					"The total number of sectors discarded successfully.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "discard_time_seconds_total"),
					"This is the total number of seconds spent by all discards.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_total"),
					"The total number of flush requests completed successfully",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_time_seconds_total"),
					"This is the total number of seconds spent by all flush requests.",
					diskLabelNames,
					nil,
				), valueType: prometheus.CounterValue,
			},
		},
		filesystemInfoDesc: typedDesc{
			desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "filesystem_info"),
				"Info about disk filesystem.",
				[]string{"device", "type", "usage", "uuid", "version"},
				nil,
			), valueType: prometheus.GaugeValue,
		},
		deviceMapperInfoDesc: typedDesc{
			desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "device_mapper_info"),
				"Info about disk device mapper.",
				[]string{"device", "name", "uuid", "vg_name", "lv_name", "lv_layer"},
				nil,
			), valueType: prometheus.GaugeValue,
		},
		ataDescs: map[string]typedDesc{
			udevIDATAWriteCache: {
				desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "ata_write_cache"),
					"ATA disk has a write cache.",
					[]string{"device"},
					nil,
				), valueType: prometheus.GaugeValue,
			},
			udevIDATAWriteCacheEnabled: {
				desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "ata_write_cache_enabled"),
					"ATA disk has its write cache enabled.",
					[]string{"device"},
					nil,
				), valueType: prometheus.GaugeValue,
			},
			udevIDATARotationRateRPM: {
				desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "ata_rotation_rate_rpm"),
					"ATA disk rotation rate in RPMs (0 for SSDs).",
					[]string{"device"},
					nil,
				), valueType: prometheus.GaugeValue,
			},
		},
		logger: logger,
	}

	// Only enable getting device properties from udev if the directory is readable.
	if stat, err := os.Stat(*udevDataPath); err != nil || !stat.IsDir() {
		logger.Error("Failed to open directory, disabling udev device properties", "path", *udevDataPath)
	} else {
		collector.getUdevDeviceProperties = getUdevDeviceProperties
	}

	return &collector, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	diskStats, err := c.fs.ProcDiskstats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %w", err)
	}

	for _, stats := range diskStats {
		dev := stats.DeviceName
		if c.deviceFilter.ignored(dev) {
			continue
		}

		info, err := getUdevDeviceProperties(stats.MajorNumber, stats.MinorNumber)
		if err != nil {
			c.logger.Debug("Failed to parse udev info", "err", err)
		}

		// This is usually the serial printed on the disk label.
		serial := info[udevSCSIIdentSerial]

		// If it's undefined, fallback to ID_SERIAL_SHORT instead.
		if serial == "" {
			serial = info[udevIDSerialShort]
		}

		// If still undefined, fallback to ID_SERIAL (used by virtio devices).
		if serial == "" {
			serial = info[udevIDSerial]
		}

		queueStats, err := c.fs.SysBlockDeviceQueueStats(dev)
		// Block Device Queue stats may not exist for all devices.
		if err != nil && !os.IsNotExist(err) {
			c.logger.Debug("Failed to get block device queue stats", "device", dev, "err", err)
		}

		ch <- c.infoDesc.mustNewConstMetric(1.0, dev,
			fmt.Sprint(stats.MajorNumber),
			fmt.Sprint(stats.MinorNumber),
			info[udevIDPath],
			info[udevIDWWN],
			info[udevIDModel],
			serial,
			info[udevIDRevision],
			strconv.FormatUint(queueStats.Rotational, 2),
		)

		statCount := stats.IoStatsCount - 3 // Total diskstats record count, less MajorNumber, MinorNumber and DeviceName

		for i, val := range []float64{
			float64(stats.ReadIOs),
			float64(stats.ReadMerges),
			float64(stats.ReadSectors) * unixSectorSize,
			float64(stats.ReadTicks) * secondsPerTick,
			float64(stats.WriteIOs),
			float64(stats.WriteMerges),
			float64(stats.WriteSectors) * unixSectorSize,
			float64(stats.WriteTicks) * secondsPerTick,
			float64(stats.IOsInProgress),
			float64(stats.IOsTotalTicks) * secondsPerTick,
			float64(stats.WeightedIOTicks) * secondsPerTick,
			float64(stats.DiscardIOs),
			float64(stats.DiscardMerges),
			float64(stats.DiscardSectors),
			float64(stats.DiscardTicks) * secondsPerTick,
			float64(stats.FlushRequestsCompleted),
			float64(stats.TimeSpentFlushing) * secondsPerTick,
		} {
			if i >= statCount {
				break
			}
			ch <- c.descs[i].mustNewConstMetric(val, dev)
		}

		if fsType := info[udevIDFSType]; fsType != "" {
			ch <- c.filesystemInfoDesc.mustNewConstMetric(1.0, dev,
				fsType,
				info[udevIDFSUsage],
				info[udevIDFSUUID],
				info[udevIDFSVersion],
			)
		}

		if name := info[udevDMName]; name != "" {
			ch <- c.deviceMapperInfoDesc.mustNewConstMetric(1.0, dev,
				name,
				info[udevDMUUID],
				info[udevDMVGName],
				info[udevDMLVName],
				info[udevDMLVLayer],
			)
		}

		if ata := info[udevIDATA]; ata != "" {
			for attr, desc := range c.ataDescs {
				str, ok := info[attr]
				if !ok {
					c.logger.Debug("Udev attribute does not exist", "attribute", attr)
					continue
				}

				if value, err := strconv.ParseFloat(str, 64); err == nil {
					ch <- desc.mustNewConstMetric(value, dev)
				} else {
					c.logger.Error("Failed to parse ATA value", "err", err)
				}
			}
		}
	}
	return nil
}

func getUdevDeviceProperties(major, minor uint32) (udevInfo, error) {
	filename := udevDataFilePath(fmt.Sprintf("b%d:%d", major, minor))

	data, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	info := make(udevInfo)

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()

		// We're only interested in device properties.
		if !strings.HasPrefix(line, udevDevicePropertyPrefix) {
			continue
		}

		line = strings.TrimPrefix(line, udevDevicePropertyPrefix)

		if name, value, found := strings.Cut(line, "="); found {
			info[name] = value
		}
	}

	return info, nil
}
