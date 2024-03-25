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
	"errors"
	"time"

	"github.com/go-kit/log/level"
	"golang.org/x/sys/unix"
)

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^devfs$"
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	var mountPointCount int
	nChan := make(chan int, 1)
	errChan := make(chan error, 1)
	go func() {
		var err error
		var n int
		n, err = unix.Getfsstat(nil, unix.MNT_WAIT)
		if err != nil {
			errChan <- err
			return
		}
		nChan <- n
	}()
	select {
	case mountPointCount = <-nChan:
	case err := <-errChan:
		return nil, err
	case <-time.After(*mountTimeout):
		return nil, errors.New("getfsstat timed out")
	}

	buf := make([]unix.Statfs_t, mountPointCount)
	go func(buf []unix.Statfs_t) {
		_, err := unix.Getfsstat(buf, unix.MNT_WAIT)
		errChan <- err
	}(buf)
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case <-time.After(*mountTimeout):
		return nil, errors.New("getfsstat timed out")
	}

	var stats []filesystemStats
	for _, fs := range buf {
		mountpoint := unix.ByteSliceToString(fs.Mntonname[:])
		if c.excludedMountPointsPattern.MatchString(mountpoint) {
			level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := unix.ByteSliceToString(fs.Mntfromname[:])
		fstype := unix.ByteSliceToString(fs.Fstypename[:])
		if c.excludedFSTypesPattern.MatchString(fstype) {
			level.Debug(c.logger).Log("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		if (fs.Flags & unix.MNT_IGNORE) != 0 {
			level.Debug(c.logger).Log("msg", "Ignoring mount flagged as ignore", "mountpoint", mountpoint)
			continue
		}

		var ro float64
		if (fs.Flags & unix.MNT_RDONLY) != 0 {
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
