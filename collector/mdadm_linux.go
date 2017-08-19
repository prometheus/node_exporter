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

// +build !nomdadm

package collector

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"syscall"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type (
	mdStatus struct {
		name string
		//active       bool
		state        float64
		disksActive  float64
		disksFailed  float64
		disksMissing float64
		disksSpare   float64
		bytesTotal   float64
		bytesSynced  float64
	}

	mdu_array_info struct {
		// Generic constant information
		major_version  int32
		minor_version  int32
		patch_version  int32
		ctime          uint32
		level          int32 //  RAID level
		size           int32 //  useless - only works for disk sizes that fit in an int32
		nr_disks       int32 //  Number of physically present disks. (working + failed)
		raid_disks     int32 //  Number of disks in the RAID.
		md_minor       int32 //     This may be larger than "nr" if disks are missing
		not_persistent int32 //     and smaller than "nr" when spare disks are around.
		// Generic state information
		utime         uint32 //  Superblock update time
		state         int32  //  State bits (clean, ...)
		active_disks  int32  //  Number of active (in sync) disks
		working_disks int32  //  Number of working disks (active + sync)
		failed_disks  int32  //  Number of failed disks
		spare_disks   int32  //  Number of spare (stand-by) disks
		// Personality information
		layout     int32 //  the array's physical layout
		chunk_size int32 //  chunk size in bytes
	}
)

const (
	// _IOR definition:
	// https://github.com/torvalds/linux/blob/e067eba5871c6922539dc1728699c14e6b22590f/include/uapi/asm-generic/ioctl.h#L85
	// from major.h
	// #define MD_MAJOR                9
	// from md_u.h
	// #define GET_ARRAY_INFO          _IOR (MD_MAJOR, 0x11, mdu_array_info_t)

	GET_ARRAY_INFO = (0x8048 << 16) | 0x0911
	// No idea what the 0x8048 means
	// I've just printed the const in a piece of C code on x86_64
	BLKGETSIZE64 = 0x80081272
)

type mdadmCollector struct{}

func init() {
	Factories["mdadm"] = NewMdadmCollector
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func ioctlArrayInfoReal(fd uintptr) (syscall.Errno, mdu_array_info) {
	var array mdu_array_info

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, GET_ARRAY_INFO, uintptr(unsafe.Pointer(&array)))
	return errno, array
}

var ioctlArrayInfo = ioctlArrayInfoReal

func ioctlBlockSizeReal(fd uintptr) (syscall.Errno, uint64) {
	var devsize uint64

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, BLKGETSIZE64, uintptr(unsafe.Pointer(&devsize)))
	return errno, devsize
}

var ioctlBlockSize = ioctlBlockSizeReal

func sysSyncCompletedFilenameReal(mdName string) string {
	return sysFilePath("block/" + mdName + "/md/sync_completed")
}

var sysSyncCompletedFilename = sysSyncCompletedFilenameReal

// getArrayInfo gets RAID info via the GET_ARRAY_INFO IOCTL syscall.
func getArrayInfo(md *mdStatus) (err error) {
	// needs md.name as input
	var (
		errno   syscall.Errno
		array   mdu_array_info
		devsize uint64
	)

	md.disksActive = math.NaN()
	md.disksFailed = math.NaN()
	md.disksMissing = math.NaN()
	md.disksSpare = math.NaN()
	md.bytesTotal = math.NaN()
	md.bytesSynced = math.NaN()
	md.state = math.NaN()

	mdDev := "/dev/" + md.name
	file, err := os.Open(mdDev)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", mdDev, err)
	}
	fd := file.Fd() // get the unix descriptor
	errno, array = ioctlArrayInfo(fd)
	if errno != 0 {
		//return fmt.Errorf("error getting RAID info via IOCTL syscall for %s: %s", mdDev, errno)
		log.Debugf("error getting RAID info via IOCTL syscall for %s: %s", mdDev, errno)
		return nil
	}
	errno, devsize = ioctlBlockSize(fd)
	if errno != 0 {
		//return fmt.Errorf("error getting RAID size via IOCTL syscall for %s: %s", mdDev, errno)
		log.Debugf("error getting RAID size via IOCTL syscall for %s: %s", mdDev, errno)
		return nil
	}
	file.Close()

	md.disksActive = float64(array.active_disks)
	md.disksFailed = float64(array.failed_disks)
	md.disksMissing = float64(max(0, array.raid_disks-array.nr_disks))
	md.disksSpare = float64(array.spare_disks)
	md.bytesTotal = float64(devsize)
	//TODO redo md.active as bool and document state (md.state)
	md.state = float64(array.state)

	sys_sync_completed := sysSyncCompletedFilename(md.name)
	log.Debugf("sys_sync_completed %s", sys_sync_completed)
	content, err := ioutil.ReadFile(sys_sync_completed)
	if err == nil {
		sync_completed := strings.Split(strings.Trim(string(content), " \t\n"), " / ")
		if len(sync_completed) == 2 {
			sync_sectors_completed, _ := strconv.ParseFloat(sync_completed[0], 64)
			sync_sectors, _ := strconv.ParseFloat(sync_completed[1], 64)
			md.bytesSynced = md.bytesTotal * sync_sectors_completed / sync_sectors
		} else {
			md.bytesSynced = md.bytesTotal
		}
	}
	return nil
}

// NewMdadmCollector returns a new Collector exposing raid statistics.
func NewMdadmCollector() (Collector, error) {
	return &mdadmCollector{}, nil
}

var (
	stateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "is_state"),
		"Indicator of the md-device state.",
		[]string{"device"},
		nil,
	)

	disksDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "disks"),
		"Number of disks in RAID. (total = active + failed + missing)",
		[]string{"device", "state"},
		nil,
	)

	bytesTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "bytes"),
		"Total number of bytes on device.",
		[]string{"device"},
		nil,
	)

	bytesSyncedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "bytes_synced"),
		"Number of bytes synced on device.",
		[]string{"device"},
		nil,
	)
)

func getMdDevices(procDiskStats string) ([]mdStatus, error) {
	content, err := ioutil.ReadFile(procDiskStats)
	if err != nil {
		return []mdStatus{}, err
	}

	lines := strings.Split(string(content), "\n")

	// Each md dev has probably at least two parent devices
	// therefore we will have something of the order len(lines)/3 devices
	// so we use that for preallocation.
	mdStates := make([]mdStatus, 0, len(lines)/3)

	for _, line := range lines {

		// major device number 9 indicates an md device.
		if len(line) > 0 && strings.Fields(line)[0] == "9" {
			md := mdStatus{
				name: strings.Fields(line)[2],
			}
			err := getArrayInfo(&md)
			if err != nil {
				return mdStates, err
			}
			mdStates = append(mdStates, md)
		}
	}
	return mdStates, nil
}

func (c *mdadmCollector) Update(ch chan<- prometheus.Metric) error {
	procDiskStats := procFilePath("diskstats")
	mdstate, err := getMdDevices(procDiskStats)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting md disks, file does not exist: %s", procDiskStats)
			return nil
		}
		return fmt.Errorf("error parsing diskstats: %s", err)
	}

	for _, mds := range mdstate {
		log.Debugf("collecting metrics for device %s", mds.name)

		ch <- prometheus.MustNewConstMetric(
			stateDesc,
			prometheus.GaugeValue,
			mds.state,
			mds.name,
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			mds.disksActive,
			mds.name, "active",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			mds.disksFailed,
			mds.name, "failed",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			mds.disksMissing,
			mds.name, "missing",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			mds.disksSpare,
			mds.name, "spare",
		)
		ch <- prometheus.MustNewConstMetric(
			bytesTotalDesc,
			prometheus.GaugeValue,
			mds.bytesTotal,
			mds.name,
		)
		ch <- prometheus.MustNewConstMetric(
			bytesSyncedDesc,
			prometheus.GaugeValue,
			mds.bytesSynced,
			mds.name,
		)
	}

	return nil
}
