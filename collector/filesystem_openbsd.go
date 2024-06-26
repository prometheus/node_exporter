// Copyright 2020 The Prometheus Authors
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

// GetStats exposes filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, fsstatErr error) {
	var mountPointCount int
	nChan := make(chan int, 1)
	errChan := make(chan error, 1)
	go func() {
		var statErr error
		var n int
		n, statErr = unix.Getfsstat(nil, unix.MNT_WAIT)
		if statErr != nil {
			errChan <- statErr
			return
		}
		nChan <- n
	}()
	select {
	case mountPointCount = <-nChan:
	case statErr := <-errChan:
		return nil, statErr
	case <-time.After(*mountTimeout):
		return nil, errors.New("getfsstat timed out")
	}

	buf := make([]unix.Statfs_t, mountPointCount)
	go func(buf []unix.Statfs_t) {
		_, fsstatErr = unix.Getfsstat(buf, unix.MNT_WAIT)
		errChan <- fsstatErr
	}(buf)
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case <-time.After(*mountTimeout):
		return nil, errors.New("getfsstat timed out")
	}

	stats = []filesystemStats{}
	for _, v := range buf {
		mountpoint := unix.ByteSliceToString(v.F_mntonname[:])
		if c.excludedMountPointsPattern.MatchString(mountpoint) {
			level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := unix.ByteSliceToString(v.F_mntfromname[:])
		fstype := unix.ByteSliceToString(v.F_fstypename[:])
		if c.excludedFSTypesPattern.MatchString(fstype) {
			level.Debug(c.logger).Log("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (v.F_flags & unix.MNT_RDONLY) != 0 {
			ro = 1
		}

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     device,
				mountPoint: mountpoint,
				fsType:     fstype,
			},
			size:      float64(v.F_blocks) * float64(v.F_bsize),
			free:      float64(v.F_bfree) * float64(v.F_bsize),
			avail:     float64(v.F_bavail) * float64(v.F_bsize),
			files:     float64(v.F_files),
			filesFree: float64(v.F_ffree),
			ro:        ro,
		})
	}
	return stats, nil
}
