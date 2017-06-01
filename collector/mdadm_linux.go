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
	"os"
	"regexp"
	"strconv"
	"strings"

	"syscall"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	statuslineRE = regexp.MustCompile(`(\d+) blocks`)
	buildlineRE  = regexp.MustCompile(`\((\d+)/\d+\)`)
)

type mdStatus struct {
	name         string
	active       bool
	disksActive  int64
	disksFailed  int64
	disksMissing int64
	disksSpare   int64
	bytesTotal   int64
	bytesSynced  int64
}

type (
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

// evalBuildline gets the size that has already been synced out of the sync-line.
func evalBuildline(buildline string) (bytesSynced int64, err error) {
	matches := buildlineRE.FindStringSubmatch(buildline)

	// +1 to make it more obvious that the whole string containing the info is also returned as matches[0].
	if len(matches) != 1+1 {
		return 0, fmt.Errorf("unable to find the rebuild block count in buildline: %s", buildline)
	}

	blocksSynced, err := strconv.ParseInt(matches[1], 10, 64)
	bytesSynced = blocksSynced * 1024

	if err != nil {
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	}

	return bytesSynced, nil
}

// getArrayInfo gets RAID info via the GET_ARRAY_INFO IOCTL syscall.
func getArrayInfo(md *mdStatus) (err error) {
	// needs md.name as input
	var (
		array   mdu_array_info
		devsize uint64
	)
	mdDev := "/dev/" + md.name
	file, err := os.Open(mdDev)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", mdDev, err)
	}
	fd := file.Fd() // get the unix descriptor
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, GET_ARRAY_INFO, uintptr(unsafe.Pointer(&array)))
	if errno != 0 {
		return fmt.Errorf("error getting RAID info via IOCTL syscall for %s: %s", mdDev, errno)
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, fd, BLKGETSIZE64, uintptr(unsafe.Pointer(&devsize)))
	if errno != 0 {
		return fmt.Errorf("error getting RAID size via IOCTL syscall for %s: %s", mdDev, errno)
	}
	file.Close()

	//	ioctl(3, GET_DISK_INFO, 0x7ffe85c820d0) = 0
	//	ioctl(4, BLKSSZGET, [512])              = 0
	//	ioctl(4, BLKGETSIZE64, [200048573952])  = 0

	md.disksActive = int64(array.active_disks)
	md.disksFailed = int64(array.failed_disks)
	md.disksMissing = int64(max(0, array.raid_disks-array.nr_disks))
	md.disksSpare = int64(array.spare_disks)
	md.bytesTotal = int64(devsize)
	return nil
}

// parseMdstat parses an mdstat-file and returns a struct with the relevant infos.
func parseMdstat(mdStatusFilePath string) ([]mdStatus, error) {
	// TODO
	// we should probably get rid of parsing the mdstat file at all,
	// but this file in /proc does not depend on the underlying architecture
	// while the ioctl could.
	content, err := ioutil.ReadFile(mdStatusFilePath)
	if err != nil {
		return []mdStatus{}, err
	}

	lines := strings.Split(string(content), "\n")
	// Each md has at least the deviceline, statusline and one empty line afterwards
	// so we will have probably something of the order len(lines)/3 devices
	// so we use that for preallocation.
	mdStates := make([]mdStatus, 0, len(lines)/3)
	for i, line := range lines {
		if line == "" {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			// Lines starting with white space are not the beginning of a md-section.
			continue
		}
		if strings.HasPrefix(line, "Personalities") || strings.HasPrefix(line, "unused") {
			// These lines contain general information.
			continue
		}

		mainLine := strings.Split(line, " ")
		if len(mainLine) < 4 {
			return mdStates, fmt.Errorf("error parsing mdline: %s", line)
		}

		md := mdStatus{
			name:   mainLine[0],
			active: mainLine[2] == "active",
		}
		err := getArrayInfo(&md)
		if err != nil {
			return mdStates, err
		}

		if len(lines) <= i+3 {
			return mdStates, fmt.Errorf("error parsing mdstat: entry for %s has fewer lines than expected", md.name)
		}

		// Now get the number of synced blocks.

		// Get the line number of the sync-line.
		var j int
		if strings.Contains(lines[i+2], "bitmap") { // then skip the bitmap line
			j = i + 3
		} else {
			j = i + 2
		}

		// If device is syncing at the moment, get the number of currently synced bytes,
		// otherwise that number equals the size of the device.
		if strings.Contains(lines[j], "recovery") ||
			strings.Contains(lines[j], "resync") &&
				!strings.Contains(lines[j], "\tresync=") {
			md.bytesSynced, err = evalBuildline(lines[j])
			if err != nil {
				return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
			}
		} else {
			md.bytesSynced = md.bytesTotal
		}

		mdStates = append(mdStates, md)
	}

	return mdStates, nil
}

// NewMdadmCollector returns a new Collector exposing raid statistics.
func NewMdadmCollector() (Collector, error) {
	return &mdadmCollector{}, nil
}

var (
	isActiveDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "is_active"),
		"Indicator whether the md-device is active or not.",
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

func (c *mdadmCollector) Update(ch chan<- prometheus.Metric) error {
	statusfile := procFilePath("mdstat")
	mdstate, err := parseMdstat(statusfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting mdstat, file does not exist: %s", statusfile)
			return nil
		}
		return fmt.Errorf("error parsing mdstatus: %s", err)
	}

	for _, mds := range mdstate {
		log.Debugf("collecting metrics for device %s", mds.name)

		var active float64
		if mds.active {
			active = 1
		}
		ch <- prometheus.MustNewConstMetric(
			isActiveDesc,
			prometheus.GaugeValue,
			active,
			mds.name,
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			float64(mds.disksActive),
			mds.name, "active",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			float64(mds.disksFailed),
			mds.name, "failed",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			float64(mds.disksMissing),
			mds.name, "missing",
		)
		ch <- prometheus.MustNewConstMetric(
			disksDesc,
			prometheus.GaugeValue,
			float64(mds.disksSpare),
			mds.name, "spare",
		)
		ch <- prometheus.MustNewConstMetric(
			bytesTotalDesc,
			prometheus.GaugeValue,
			float64(mds.bytesTotal),
			mds.name,
		)
		ch <- prometheus.MustNewConstMetric(
			bytesSyncedDesc,
			prometheus.GaugeValue,
			float64(mds.bytesSynced),
			mds.name,
		)
	}

	return nil
}
