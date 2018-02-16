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
	defIgnoredMountPoints = "^/(dev|proc|sys|var/lib/docker)($|/)"
	defIgnoredFSTypes     = "^(autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$"
	readOnly              = 0x1 // ST_RDONLY
)

// GetStats returns filesystem stats.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	mps, err := mountPointDetails()
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}
	for _, labels := range mps {
		if c.ignoredMountPointsPattern.MatchString(labels.mountPoint) {
			log.Debugf("Ignoring mount point: %s", labels.mountPoint)
			continue
		}
		if c.ignoredFSTypesPattern.MatchString(labels.fsType) {
			log.Debugf("Ignoring fs type: %s", labels.fsType)
			continue
		}

		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(labels.mountPoint, buf)
		if err != nil {
			stats = append(stats, filesystemStats{
				labels:      labels,
				deviceError: 1,
			})
			log.Debugf("Error on statfs() system call for %q: %s", labels.mountPoint, err)
			continue
		}

		var ro float64
		if (buf.Flags & readOnly) != 0 {
			ro = 1
		}

		stats = append(stats, filesystemStats{
			labels:    labels,
			size:      float64(buf.Blocks) * float64(buf.Bsize),
			free:      float64(buf.Bfree) * float64(buf.Bsize),
			avail:     float64(buf.Bavail) * float64(buf.Bsize),
			files:     float64(buf.Files),
			filesFree: float64(buf.Ffree),
			ro:        ro,
		})
	}
	return stats, nil
}

func mountPointDetails() ([]filesystemLabels, error) {
	file, err := os.Open(procFilePath("mounts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filesystems := []filesystemLabels{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		filesystems = append(filesystems, filesystemLabels{
			device:     parts[0],
			mountPoint: parts[1],
			fsType:     parts[2],
		})
	}
	return filesystems, scanner.Err()
}
