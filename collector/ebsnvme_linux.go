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
//
// The EBS NVMe log-page parsing logic in this file is derived from the Amazon
// EBS CSI Driver, which is also licensed under the Apache License, Version 2.0:
//
//	Copyright 2024 The Kubernetes Authors
//	https://github.com/kubernetes-sigs/aws-ebs-csi-driver
//	(pkg/metrics/nvme.go)
//
// Author of this node_exporter collector:
//
//	Allen Xie <weifeng.xie@qq.com>, Storage Solution Architect

//go:build linux && !noebsnvme
// +build linux,!noebsnvme

// Package collector's ebsnvme collector exposes the Amazon EBS detailed
// performance statistics that Nitro-based EC2 instances vend through the NVMe
// device log page. The statistics, their names, and their semantics are
// documented in the Amazon EBS User Guide:
//
//	https://docs.aws.amazon.com/ebs/latest/userguide/nvme-detailed-performance-stats.html
//
// The collector reads NVMe log page 0xD0 from each EBS-backed NVMe device via
// an ioctl, parses the binary EBS statistics structure, and exposes the values
// as Prometheus metrics labelled by volume_id, device, and mount_path.

package collector

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	ebsNVMeSubsystem = "ebs"

	// ebsMagic identifies a valid EBS NVMe statistics log page.
	ebsMagic = 0x3C23B510

	// ebsLogPageID is the NVMe log page identifier for EBS statistics.
	ebsLogPageID = 0xD0

	// nvmeAdminGetLogPage is the NVMe admin opcode for "Get Log Page".
	nvmeAdminGetLogPage = 0x02

	// nvmeLogPageSize is the length, in bytes, of the EBS log page.
	nvmeLogPageSize = 4096

	// nvmeIoctlAdminCmd is the ioctl request code for NVME_IOCTL_ADMIN_CMD.
	nvmeIoctlAdminCmd = 0xC0484E41

	// microsecondsInSeconds converts microsecond counters to seconds.
	microsecondsInSeconds = 1e6

	// notMounted is the mount_path label value used for NVMe devices that are
	// not mounted, or whose top-level device has no direct mount point (for
	// example a disk that is only mounted through one of its partitions).
	notMounted = "NotMounted"
)

var (
	errInvalidEBSMagic = errors.New("invalid EBS magic number")
	errParseLogPage    = errors.New("failed to parse log page")
)

// ebsMetrics maps the EBS NVMe log page layout. The field order and the
// reserved area size are dictated by the on-disk EBS statistics structure.
//
// Reference: Amazon EBS User Guide, "Amazon EBS detailed performance statistics"
// https://docs.aws.amazon.com/ebs/latest/userguide/nvme-detailed-performance-stats.html
type ebsMetrics struct {
	EBSMagic              uint64
	ReadOps               uint64 // total_read_ops
	WriteOps              uint64 // total_write_ops
	ReadBytes             uint64 // total_read_bytes
	WriteBytes            uint64 // total_write_bytes
	TotalReadTime         uint64 // total_read_time (microseconds)
	TotalWriteTime        uint64 // total_write_time (microseconds)
	EBSIOPSExceeded       uint64 // ebs_volume_performance_exceeded_iops (microseconds)
	EBSThroughputExceeded uint64 // ebs_volume_performance_exceeded_tp (microseconds)
	EC2IOPSExceeded       uint64 // ec2_instance_ebs_performance_exceeded_iops (microseconds)
	EC2ThroughputExceeded uint64 // ec2_instance_ebs_performance_exceeded_tp (microseconds)
	QueueLength           uint64 // volume_queue_length (point in time)
	ReservedArea          [416]byte
	ReadLatency           ebsHistogram // read_io_latency_histogram
	WriteLatency          ebsHistogram // write_io_latency_histogram
}

// ebsHistogram is the latency histogram embedded in the EBS log page.
type ebsHistogram struct {
	BinCount uint64
	Bins     [64]ebsHistogramBin
}

// ebsHistogramBin is a single latency bucket of an ebsHistogram. Bounds are in
// microseconds.
type ebsHistogramBin struct {
	Lower uint64
	Upper uint64
	Count uint64
}

// nvmePassthruCommand mirrors struct nvme_passthru_cmd from <linux/nvme_ioctl.h>.
type nvmePassthruCommand struct {
	opcode      uint8
	flags       uint8
	rsvd1       uint16
	nsid        uint32
	cdw2        uint32
	cdw3        uint32
	metadata    uint64
	addr        uint64
	metadataLen uint32
	dataLen     uint32
	cdw10       uint32
	cdw11       uint32
	cdw12       uint32
	cdw13       uint32
	cdw14       uint32
	cdw15       uint32
	timeoutMs   uint32
	result      uint32
}

// ebsNVMeCollector exposes Amazon EBS volume performance statistics read from
// the EBS NVMe device log page.
type ebsNVMeCollector struct {
	logger  *slog.Logger
	metrics map[string]*prometheus.Desc
}

func init() {
	registerCollector("ebsnvme", defaultDisabled, NewEBSNVMeCollector)
}

// NewEBSNVMeCollector returns a new Collector exposing Amazon EBS NVMe
// performance statistics.
func NewEBSNVMeCollector(logger *slog.Logger) (Collector, error) {
	labels := []string{"volume_id", "device", "mount_path"}
	desc := func(name, help string) *prometheus.Desc {
		return prometheus.NewDesc(prometheus.BuildFQName(namespace, ebsNVMeSubsystem, name), help, labels, nil)
	}

	return &ebsNVMeCollector{
		logger: logger,
		metrics: map[string]*prometheus.Desc{
			"read_ops_total":                  desc("read_ops_total", "The total number of completed read operations. (EBS statistic: total_read_ops)"),
			"write_ops_total":                 desc("write_ops_total", "The total number of completed write operations. (EBS statistic: total_write_ops)"),
			"read_bytes_total":                desc("read_bytes_total", "The total number of read bytes transferred. (EBS statistic: total_read_bytes)"),
			"write_bytes_total":               desc("write_bytes_total", "The total number of write bytes transferred. (EBS statistic: total_write_bytes)"),
			"read_seconds_total":              desc("read_seconds_total", "The total time spent, in seconds, by all completed read operations. (EBS statistic: total_read_time)"),
			"write_seconds_total":             desc("write_seconds_total", "The total time spent, in seconds, by all completed write operations. (EBS statistic: total_write_time)"),
			"exceeded_iops_seconds_total":     desc("exceeded_iops_seconds_total", "The total time, in seconds, that IOPS demand exceeded the volume's provisioned IOPS performance. (EBS statistic: ebs_volume_performance_exceeded_iops)"),
			"exceeded_tp_seconds_total":       desc("exceeded_tp_seconds_total", "The total time, in seconds, that throughput demand exceeded the volume's provisioned throughput performance. (EBS statistic: ebs_volume_performance_exceeded_tp)"),
			"ec2_exceeded_iops_seconds_total": desc("ec2_exceeded_iops_seconds_total", "The total time, in seconds, that the EBS volume exceeded the attached EC2 instance's maximum IOPS performance. (EBS statistic: ec2_instance_ebs_performance_exceeded_iops)"),
			"ec2_exceeded_tp_seconds_total":   desc("ec2_exceeded_tp_seconds_total", "The total time, in seconds, that the EBS volume exceeded the attached EC2 instance's maximum throughput performance. (EBS statistic: ec2_instance_ebs_performance_exceeded_tp)"),
			"volume_queue_length":             desc("volume_queue_length", "The number of read and write operations waiting to be completed. (EBS statistic: volume_queue_length)"),
			"read_io_latency_seconds":         desc("read_io_latency_seconds", "The number of read operations completed within each latency bin, in seconds. (EBS statistic: read_io_latency_histogram)"),
			"write_io_latency_seconds":        desc("write_io_latency_seconds", "The number of write operations completed within each latency bin, in seconds. (EBS statistic: write_io_latency_histogram)"),
		},
	}, nil
}

// Update reads the EBS log page from every EBS NVMe device and emits the parsed
// metrics. Devices that are not EBS volumes (no valid EBS magic) are skipped.
func (c *ebsNVMeCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := fetchDeviceMapping()
	if err != nil {
		return fmt.Errorf("error mapping NVMe devices to EBS volumes: %w", err)
	}

	found := false
	for devicePath, info := range devices {
		data, err := readEBSLogPage(devicePath)
		if err != nil {
			c.logger.Debug("skipping device", "device", devicePath, "err", err)
			continue
		}

		metrics, err := parseEBSLogPage(data)
		if err != nil {
			// Not an EBS volume (or an unexpected layout); skip silently at debug.
			c.logger.Debug("skipping non-EBS device", "device", devicePath, "err", err)
			continue
		}
		found = true

		// Drop the /dev/ prefix to align with the diskstats collector.
		device := strings.TrimPrefix(devicePath, "/dev/")
		volumeID := info.volumeID
		mountPath := info.mountPath

		emit := func(name string, vt prometheus.ValueType, v float64) {
			ch <- prometheus.MustNewConstMetric(c.metrics[name], vt, v, volumeID, device, mountPath)
		}

		emit("read_ops_total", prometheus.CounterValue, float64(metrics.ReadOps))
		emit("write_ops_total", prometheus.CounterValue, float64(metrics.WriteOps))
		emit("read_bytes_total", prometheus.CounterValue, float64(metrics.ReadBytes))
		emit("write_bytes_total", prometheus.CounterValue, float64(metrics.WriteBytes))
		emit("read_seconds_total", prometheus.CounterValue, float64(metrics.TotalReadTime)/microsecondsInSeconds)
		emit("write_seconds_total", prometheus.CounterValue, float64(metrics.TotalWriteTime)/microsecondsInSeconds)
		emit("exceeded_iops_seconds_total", prometheus.CounterValue, float64(metrics.EBSIOPSExceeded)/microsecondsInSeconds)
		emit("exceeded_tp_seconds_total", prometheus.CounterValue, float64(metrics.EBSThroughputExceeded)/microsecondsInSeconds)
		emit("ec2_exceeded_iops_seconds_total", prometheus.CounterValue, float64(metrics.EC2IOPSExceeded)/microsecondsInSeconds)
		emit("ec2_exceeded_tp_seconds_total", prometheus.CounterValue, float64(metrics.EC2ThroughputExceeded)/microsecondsInSeconds)
		emit("volume_queue_length", prometheus.GaugeValue, float64(metrics.QueueLength))

		readCount, readBuckets := convertEBSHistogram(metrics.ReadLatency)
		ch <- prometheus.MustNewConstHistogram(c.metrics["read_io_latency_seconds"], readCount, 0, readBuckets, volumeID, device, mountPath)

		writeCount, writeBuckets := convertEBSHistogram(metrics.WriteLatency)
		ch <- prometheus.MustNewConstHistogram(c.metrics["write_io_latency_seconds"], writeCount, 0, writeBuckets, volumeID, device, mountPath)
	}

	if !found {
		return ErrNoData
	}
	return nil
}

// convertEBSHistogram converts an ebsHistogram (bounds in microseconds) into the
// cumulative bucket form expected by prometheus.MustNewConstHistogram (seconds).
func convertEBSHistogram(hist ebsHistogram) (uint64, map[float64]uint64) {
	var count uint64
	buckets := make(map[float64]uint64)

	for i := uint64(0); i < hist.BinCount && i < uint64(len(hist.Bins)); i++ {
		count += hist.Bins[i].Count
		buckets[float64(hist.Bins[i].Upper)/microsecondsInSeconds] = count
	}

	return count, buckets
}

// readEBSLogPage reads the EBS statistics log page from the NVMe device at the
// given path. The device is opened read-only: a write handle is not required to
// issue the NVMe admin ioctl on Linux (matching the AWS EBS CSI driver).
func readEBSLogPage(devicePath string) ([]byte, error) {
	f, err := os.OpenFile(devicePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error opening device %s: %w", devicePath, err)
	}
	defer f.Close()

	return nvmeReadLogPage(f.Fd(), ebsLogPageID)
}

// nvmeReadLogPage reads an NVMe log page via an ioctl system call.
func nvmeReadLogPage(fd uintptr, logID uint8) ([]byte, error) {
	data := make([]byte, nvmeLogPageSize)
	if len(data) > math.MaxUint32 {
		return nil, errors.New("nvmeReadLogPage: bufferLen exceeds MaxUint32")
	}

	cmd := nvmePassthruCommand{
		opcode:  nvmeAdminGetLogPage,
		addr:    uint64(uintptr(unsafe.Pointer(&data[0]))),
		nsid:    1,
		dataLen: uint32(len(data)),
		cdw10:   uint32(logID) | (1024 << 16),
	}

	status, _, errno := unix.Syscall(unix.SYS_IOCTL, fd, nvmeIoctlAdminCmd, uintptr(unsafe.Pointer(&cmd)))
	if errno != 0 {
		return nil, fmt.Errorf("nvmeReadLogPage: ioctl error: %w", errno)
	}
	if status != 0 {
		return nil, fmt.Errorf("nvmeReadLogPage: ioctl command failed with status %d", status)
	}

	return data, nil
}

// parseEBSLogPage parses the binary data from an EBS log page into ebsMetrics.
func parseEBSLogPage(data []byte) (ebsMetrics, error) {
	var metrics ebsMetrics

	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &metrics); err != nil {
		return ebsMetrics{}, fmt.Errorf("%w: %w", errParseLogPage, err)
	}

	if metrics.EBSMagic != ebsMagic {
		return ebsMetrics{}, fmt.Errorf("%w: %x", errInvalidEBSMagic, metrics.EBSMagic)
	}

	return metrics, nil
}

// blockDevice is a single entry of lsblk JSON output.
type blockDevice struct {
	Name       string  `json:"name"`
	Serial     string  `json:"serial"`
	MountPoint *string `json:"mountpoint"`
}

// lsblkOutput is the top-level lsblk JSON document.
type lsblkOutput struct {
	BlockDevices []blockDevice `json:"blockdevices"`
}

// ebsDeviceInfo holds the per-device attributes exported as metric labels.
type ebsDeviceInfo struct {
	volumeID  string
	mountPath string
}

// mapDevices parses lsblk output and returns a map of device paths to their EBS
// volume ID and mount path. A device with no mount point (lsblk reports null,
// e.g. an unmounted disk or one mounted only through a partition) is reported
// with the mount path set to notMounted.
func mapDevices(raw []byte) (map[string]ebsDeviceInfo, error) {
	var parsed lsblkOutput
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("mapDevices: error unmarshaling lsblk JSON: %w", err)
	}

	m := make(map[string]ebsDeviceInfo, len(parsed.BlockDevices))
	for _, device := range parsed.BlockDevices {
		volumeID := device.Serial
		// EBS exposes the volume ID as the NVMe serial, e.g. "vol0abc..."; AWS
		// tooling renders it as "vol-0abc...".
		if strings.HasPrefix(volumeID, "vol") && !strings.HasPrefix(volumeID, "vol-") {
			volumeID = "vol-" + volumeID[3:]
		}

		mountPath := notMounted
		if device.MountPoint != nil && *device.MountPoint != "" {
			mountPath = *device.MountPoint
		}

		m["/dev/"+device.Name] = ebsDeviceInfo{volumeID: volumeID, mountPath: mountPath}
	}

	return m, nil
}

// fetchDeviceMapping returns a map of device paths to their EBS volume ID and
// mount path, derived from lsblk.
func fetchDeviceMapping() (map[string]ebsDeviceInfo, error) {
	output, err := exec.Command("lsblk", "-nd", "--json", "-o", "NAME,SERIAL,MOUNTPOINT").Output()
	if err != nil {
		return nil, fmt.Errorf("fetchDeviceMapping: error running lsblk: %w", err)
	}

	return mapDevices(output)
}
