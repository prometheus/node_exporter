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
	blocksTotal  int64
	blocksSynced int64
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

// evalStatusline gets the number of blocks of the device out of the status-line.
func evalStatusline(statusline string) (blocksTotal int64, err error) {
	matches := statuslineRE.FindStringSubmatch(statusline)

	// +1 to make it more obvious that the whole string containing the info is also returned as matches[0].
	if len(matches) != 1+1 {
		return 0, fmt.Errorf("unable to find the number of blocks in statusline: %s", statusline)
	}

	blocksTotal, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}

	return blocksTotal, nil
}

// evalBuildline gets the size that has already been synced out of the sync-line.
func evalBuildline(buildline string) (blocksSynced int64, err error) {
	matches := buildlineRE.FindStringSubmatch(buildline)

	// +1 to make it more obvious that the whole string containing the info is also returned as matches[0].
	if len(matches) != 1+1 {
		return 0, fmt.Errorf("unable to find the rebuild block count in buildline: %s", buildline)
	}

	blocksSynced, err = strconv.ParseInt(matches[1], 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	}

	return blocksSynced, nil
}

// getArrayInfo gets RAID info via the GET_ARRAY_INFO IOCTL syscall.
func getArrayInfo(md *mdStatus) (err error) {
	// needs md.name as input
	var array mdu_array_info
	mdDev := "/dev/" + md.name
	file, err := os.Open(mdDev)
	if err != nil {
		return fmt.Errorf("error opening %s: %s", mdDev, err)
	}
	fd := file.Fd() // get the unix descriptor
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, GET_ARRAY_INFO, uintptr(unsafe.Pointer(&array)))
	file.Close()
	if errno != 0 {
		return fmt.Errorf("error getting RAID info via IOCTL syscall for %s: %s", mdDev, errno)
	}
	md.disksActive = int64(array.active_disks)
	md.disksFailed = int64(array.failed_disks)
	md.disksMissing = int64(max(0, array.raid_disks-array.nr_disks))
	md.disksSpare = int64(array.spare_disks)
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
		return []mdStatus{}, fmt.Errorf("error parsing mdstat: %s", err)
	}

	mdStatusFile := string(content)

	lines := strings.Split(mdStatusFile, "\n")

	// Each md has at least the deviceline, statusline and one empty line afterwards
	// so we will have probably something of the order len(lines)/3 devices
	// so we use that for preallocation.
	estimateMDs := len(lines) / 3
	mdStates := make([]mdStatus, 0, estimateMDs)

	for i, l := range lines {

		if l == "" {
			// Skip entirely empty lines.
			continue
		}

		if l[0] == ' ' || l[0] == '\t' {
			// Those lines are not the beginning of a md-section.
			continue
		}

		if strings.HasPrefix(l, "Personalities") || strings.HasPrefix(l, "unused") {
			// We aren't interested in lines with general info.
			continue
		}

		mainLine := strings.Split(l, " ")
		if len(mainLine) < 4 {
			return mdStates, fmt.Errorf("error parsing mdline: %s", l)
		}

		var md mdStatus

		// mainLine[0] is the name of the md-device.
		// we have to set name before calling getArrayInfo
		md.name = mainLine[0]

		err := getArrayInfo(&md)
		if err != nil {
			return mdStates, err
		}

		// mainLine[2], either "active" or "inactive"
		md.active = mainLine[2] == "active"

		if len(lines) <= i+3 {
			return mdStates, fmt.Errorf("error parsing mdstat: entry for %s has fewer lines than expected", md.name)
		}

		md.blocksTotal, err = evalStatusline(lines[i+1]) // Parse statusline, always present.
		if err != nil {
			return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
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
			md.blocksSynced, err = evalBuildline(lines[j])
			if err != nil {
				return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
			}
		} else {
			md.blocksSynced = md.blocksTotal
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

	blocksTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "blocks"),
		"Total number of blocks on device.",
		[]string{"device"},
		nil,
	)

	blocksSyncedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "blocks_synced"),
		"Number of blocks synced on device.",
		[]string{"device"},
		nil,
	)
)

func (c *mdadmCollector) Update(ch chan<- prometheus.Metric) error {
	statusfile := procFilePath("mdstat")
	if _, err := os.Stat(statusfile); err != nil {
		// Take care we don't crash on non-existent statusfiles.
		if os.IsNotExist(err) {
			// no such file or directory, nothing to do, just return
			log.Debugf("Not collecting mdstat, file does not exist: %s", statusfile)
			return nil
		}
		return err
	}

	// First parse mdstat-file...
	mdstate, err := parseMdstat(statusfile)
	if err != nil {
		return fmt.Errorf("error parsing mdstatus: %s", err)
	}

	// ... and then plug the result into the metrics to be exported.
	var isActiveFloat float64
	for _, mds := range mdstate {

		log.Debugf("collecting metrics for device %s", mds.name)

		if mds.active {
			isActiveFloat = 1
		} else {
			isActiveFloat = 0
		}

		ch <- prometheus.MustNewConstMetric(
			isActiveDesc,
			prometheus.GaugeValue,
			isActiveFloat,
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
			blocksTotalDesc,
			prometheus.GaugeValue,
			float64(mds.blocksTotal),
			mds.name,
		)

		ch <- prometheus.MustNewConstMetric(
			blocksSyncedDesc,
			prometheus.GaugeValue,
			float64(mds.blocksSynced),
			mds.name,
		)

	}

	return nil
}
