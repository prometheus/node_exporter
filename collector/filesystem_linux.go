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

// +build !nofilesystem

package collector

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"github.com/prometheus/common/log"
)

const (
	defIgnoredMountPoints = "^/(sys|proc|dev)($|/)"
	defIgnoredFSTypes     = "^(sys|proc)fs$"
	ST_RDONLY             = 0x1
)

type filesystemDetails struct {
	device     string
	mountPoint string
	fsType     string
}

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	mpds, err := mountPointDetails()
	if err != nil {
		return nil, err
	}
	stats = []filesystemStats{}
	for _, mpd := range mpds {
		if c.ignoredMountPointsPattern.MatchString(mpd.mountPoint) {
			log.Debugf("Ignoring mount point: %s", mpd.mountPoint)
			continue
		}
		if c.ignoredFSTypesPattern.MatchString(mpd.fsType) {
			log.Debugf("Ignoring fs type: %s", mpd.fsType)
			continue
		}
		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(mpd.mountPoint, buf)
		if err != nil {
			log.Debugf("Statfs on %s returned %s",
				mpd.mountPoint, err)
			continue
		}

		var ro float64
		if buf.Flags&ST_RDONLY != 0 {
			ro = 1
		}

		labelValues := []string{mpd.device, mpd.mountPoint, mpd.fsType}
		stats = append(stats, filesystemStats{
			labelValues: labelValues,
			size:        float64(buf.Blocks) * float64(buf.Bsize),
			free:        float64(buf.Bfree) * float64(buf.Bsize),
			avail:       float64(buf.Bavail) * float64(buf.Bsize),
			files:       float64(buf.Files),
			filesFree:   float64(buf.Ffree),
			ro:          ro,
		})
	}
	return stats, nil
}

func mountPointDetails() ([]filesystemDetails, error) {
	file, err := os.Open(procFilePath("mounts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filesystems := []filesystemDetails{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		filesystems = append(filesystems, filesystemDetails{parts[0], parts[1], parts[2]})
	}
	return filesystems, nil
}
