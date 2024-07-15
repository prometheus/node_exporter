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

//go:build (darwin || dragonfly) && !nofilesystem
// +build darwin dragonfly
// +build !nofilesystem

package collector

import (
	"errors"
	"time"
	"unsafe"

	"github.com/go-kit/log/level"
)

/*
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
#include <stdio.h>
*/
import "C"

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^devfs$"
	readOnly               = 0x1 // MNT_RDONLY
)

// GetStats exposes filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	// `getmntinfo` relies on `getfsstat` in some variants, and is blocking in general.
	count := 0
	countCh := make(chan int, 1)
	var mountBuf *C.struct_statfs
	go func(mountBuf **C.struct_statfs) {
		countCh <- int(C.getmntinfo(mountBuf, C.MNT_WAIT))
		close(countCh)
	}(&mountBuf)
	select {
	case count = <-countCh:
		if count <= 0 {
			return nil, errors.New("getmntinfo failed")
		}
	case <-time.After(*mountTimeout):
		return nil, errors.New("getmntinfo timed out")
	}

	mnt := (*[1 << 20]C.struct_statfs)(unsafe.Pointer(mountBuf))
	stats = []filesystemStats{}
	for i := 0; i < count; i++ {
		mountpoint := C.GoString(&mnt[i].f_mntonname[0])
		if c.excludedMountPointsPattern.MatchString(mountpoint) {
			level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := C.GoString(&mnt[i].f_mntfromname[0])
		fstype := C.GoString(&mnt[i].f_fstypename[0])
		if c.excludedFSTypesPattern.MatchString(fstype) {
			level.Debug(c.logger).Log("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (mnt[i].f_flags & readOnly) != 0 {
			ro = 1
		}

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     device,
				mountPoint: rootfsStripPrefix(mountpoint),
				fsType:     fstype,
			},
			size:      float64(mnt[i].f_blocks) * float64(mnt[i].f_bsize),
			free:      float64(mnt[i].f_bfree) * float64(mnt[i].f_bsize),
			avail:     float64(mnt[i].f_bavail) * float64(mnt[i].f_bsize),
			files:     float64(mnt[i].f_files),
			filesFree: float64(mnt[i].f_ffree),
			ro:        ro,
		})
	}
	return stats, nil
}
