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

//go:build !nofilesystem
// +build !nofilesystem

package collector

import (
	"github.com/go-kit/log/level"
	"golang.org/x/sys/unix"
)

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^devfs$"
	readOnly               = 0x1 // MNT_RDONLY
	noWait                 = 0x2 // MNT_NOWAIT
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	n, err := unix.Getfsstat(nil, noWait)
	if err != nil {
		return nil, err
	}
	buf := make([]unix.Statfs_t, n)
	_, err = unix.Getfsstat(buf, noWait)
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}
	for _, fs := range buf {
		mountpoint := bytesToString(fs.Mntonname[:])
		if c.excludedMountPointsPattern.MatchString(mountpoint) {
			level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := bytesToString(fs.Mntfromname[:])
		fstype := bytesToString(fs.Fstypename[:])
		if c.excludedFSTypesPattern.MatchString(fstype) {
			level.Debug(c.logger).Log("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (fs.Flags & readOnly) != 0 {
			ro = 1
		}

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     device,
				mountPoint: rootfsStripPrefix(mountpoint),
				fsType:     fstype,
			},
			size:      float64(fs.Blocks) * float64(fs.Bsize),
			free:      float64(fs.Bfree) * float64(fs.Bsize),
			avail:     float64(fs.Bavail) * float64(fs.Bsize),
			files:     float64(fs.Files),
			filesFree: float64(fs.Ffree),
			ro:        ro,
		})
	}
	return stats, nil
}
