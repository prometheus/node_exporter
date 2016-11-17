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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	statuslineRE             = regexp.MustCompile(`(\d+) blocks .*\[(\d+)/(\d+)\] \[[U_]+\]`)
	raid0lineRE              = regexp.MustCompile(`(\d+) blocks( super ([0-9\.])*)? \d+k chunks`)
	buildlineRE              = regexp.MustCompile(`\((\d+)/\d+\)`)
	unknownPersonalityLineRE = regexp.MustCompile(`(\d+) blocks (.*)`)
	raidPersonalityRE        = regexp.MustCompile(`raid[0-9]+`)
)

type mdStatus struct {
	mdName       string
	isActive     bool
	disksActive  int64
	disksTotal   int64
	blocksTotal  int64
	blocksSynced int64
}

type mdadmCollector struct{}

func init() {
	Factories["mdadm"] = NewMdadmCollector
}

func evalStatusline(statusline string) (active, total, size int64, err error) {
	matches := statuslineRE.FindStringSubmatch(statusline)

	// +1 to make it more obvious that the whole string containing the info is also returned as matches[0].
	if len(matches) < 3+1 {
		return 0, 0, 0, fmt.Errorf("too few matches found in statusline: %s", statusline)
	} else {
		if len(matches) > 3+1 {
			return 0, 0, 0, fmt.Errorf("too many matches found in statusline: %s", statusline)
		}
	}

	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}

	total, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	active, err = strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}

	return active, total, size, nil
}

func evalRaid0line(statusline string) (size int64, err error) {
	matches := raid0lineRE.FindStringSubmatch(statusline)

	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid raid0 status line: %s", statusline)
	}

	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}

	return size, nil
}

func evalUnknownPersonalitylineRE(statusline string) (size int64, err error) {
	matches := unknownPersonalityLineRE.FindStringSubmatch(statusline)

	if len(matches) != 2+1 {
		return 0, fmt.Errorf("invalid unknown personality status line: %s", statusline)
	}

	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}

	return size, nil
}

// Gets the size that has already been synced out of the sync-line.
func evalBuildline(buildline string) (int64, error) {
	matches := buildlineRE.FindStringSubmatch(buildline)

	// +1 to make it more obvious that the whole string containing the info is also returned as matches[0].
	if len(matches) < 1+1 {
		return 0, fmt.Errorf("too few matches found in buildline: %s", buildline)
	}

	if len(matches) > 1+1 {
		return 0, fmt.Errorf("too many matches found in buildline: %s", buildline)
	}

	syncedSize, err := strconv.ParseInt(matches[1], 10, 64)

	if err != nil {
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	}

	return syncedSize, nil
}

// Parses an mdstat-file and returns a struct with the relevant infos.
func parseMdstat(mdStatusFilePath string) ([]mdStatus, error) {
	content, err := ioutil.ReadFile(mdStatusFilePath)
	if err != nil {
		return []mdStatus{}, fmt.Errorf("error parsing mdstat: %s", err)
	}

	mdStatusFile := string(content)

	lines := strings.Split(mdStatusFile, "\n")
	var (
		currentMD           string
		personality         string
		active, total, size int64
	)

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
		currentMD = mainLine[0]               // The name of the md-device.
		isActive := (mainLine[2] == "active") // The activity status of the md-device.
		personality = ""
		for _, possiblePersonality := range mainLine[3:] {
			if raidPersonalityRE.MatchString(possiblePersonality) {
				personality = possiblePersonality
				break
			}
		}

		if len(lines) <= i+3 {
			return mdStates, fmt.Errorf("error parsing mdstat: entry for %s has fewer lines than expected", currentMD)
		}

		switch {
		case personality == "raid0":
			active = int64(len(mainLine) - 4)     // Get the number of devices from the main line.
			total = active                        // Raid0 active and total is always the same if active.
			size, err = evalRaid0line(lines[i+1]) // Parse statusline, always present.
		case raidPersonalityRE.MatchString(personality):
			active, total, size, err = evalStatusline(lines[i+1]) // Parse statusline, always present.
		default:
			log.Infof("Personality unknown: %s\n", mainLine)
			size, err = evalUnknownPersonalitylineRE(lines[i+1]) // Parse statusline, always present.
		}

		if err != nil {
			return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
		}

		// Now get the number of synced blocks.
		var syncedBlocks int64

		// Get the line number of the syncing-line.
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
			syncedBlocks, err = evalBuildline(lines[j])
			if err != nil {
				return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
			}
		} else {
			syncedBlocks = size
		}

		mdStates = append(mdStates, mdStatus{currentMD, isActive, active, total, size, syncedBlocks})

	}

	return mdStates, nil
}

// Just returns the pointer to an empty struct as we only use throwaway-metrics.
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

	disksActiveDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "disks_active"),
		"Number of active disks of device.",
		[]string{"device"},
		nil,
	)

	disksTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "md", "disks"),
		"Total number of disks of device.",
		[]string{"device"},
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

func (c *mdadmCollector) Update(ch chan<- prometheus.Metric) (err error) {
	statusfile := procFilePath("mdstat")
	// take care we don't crash on non-existent statusfiles
	_, err = os.Stat(statusfile)
	if os.IsNotExist(err) {
		// no such file or directory, nothing to do, just return
		log.Debugf("Not collecting mdstat, file does not exist: %s", statusfile)
		return nil
	}

	if err != nil { // now things get weird, better to return
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

		log.Debugf("collecting metrics for device %s", mds.mdName)

		if mds.isActive {
			isActiveFloat = 1
		} else {
			isActiveFloat = 0
		}

		ch <- prometheus.MustNewConstMetric(
			isActiveDesc,
			prometheus.GaugeValue,
			isActiveFloat,
			mds.mdName,
		)

		ch <- prometheus.MustNewConstMetric(
			disksActiveDesc,
			prometheus.GaugeValue,
			float64(mds.disksActive),
			mds.mdName,
		)

		ch <- prometheus.MustNewConstMetric(
			disksTotalDesc,
			prometheus.GaugeValue,
			float64(mds.disksTotal),
			mds.mdName,
		)

		ch <- prometheus.MustNewConstMetric(
			blocksTotalDesc,
			prometheus.GaugeValue,
			float64(mds.blocksTotal),
			mds.mdName,
		)

		ch <- prometheus.MustNewConstMetric(
			blocksSyncedDesc,
			prometheus.GaugeValue,
			float64(mds.blocksSynced),
			mds.mdName,
		)

	}

	return nil
}
