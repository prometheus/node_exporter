// Copyright 2024 The Prometheus Authors
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
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^(kernfs|procfs|ptyfs|fdesc)$"
	_VFS_NAMELEN           = 32
	_VFS_MNAMELEN          = 1024
)

/*
 * Go uses the NetBSD 9 ABI and thus syscall.SYS_GETVFSSTAT is compat_90_getvfsstat.
 * We have to declare struct statvfs90 because it is not included in the unix package.
 * See NetBSD/src/sys/compat/sys/statvfs.h.
 */
type statvfs90 struct {
	F_flag   uint
	F_bsize  uint
	F_frsize uint
	F_iosize uint

	F_blocks uint64
	F_bfree  uint64
	F_bavail uint64
	F_bresvd uint64

	F_files  uint64
	F_ffree  uint64
	F_favail uint64
	F_fresvd uint64

	F_syncreads  uint64
	F_syncwrites uint64

	F_asyncreads  uint64
	F_asyncwrites uint64

	F_fsidx   [2]uint32
	F_fsid    uint32
	F_namemax uint
	F_owner   uint32
	F_spare   [4]uint32

	F_fstypename  [_VFS_NAMELEN]byte
	F_mntonname   [_VFS_MNAMELEN]byte
	F_mntfromname [_VFS_MNAMELEN]byte

	cgo_pad [4]byte
}

func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	var mnt []statvfs90
	if syscall.SYS_GETVFSSTAT != 356 /* compat_90_getvfsstat */ {
		/*
		 * Catch if golang ever updates to newer ABI and bail.
		 */
		return nil, fmt.Errorf("getvfsstat: ABI mismatch")
	}
	for {
		r1, _, errno := syscall.Syscall(syscall.SYS_GETVFSSTAT, uintptr(0), 0, unix.ST_NOWAIT)
		if errno != 0 {
			return nil, fmt.Errorf("getvfsstat: %s", string(errno))
		}
		mnt = make([]statvfs90, r1, r1)
		r2, _, errno := syscall.Syscall(syscall.SYS_GETVFSSTAT, uintptr(unsafe.Pointer(&mnt[0])), unsafe.Sizeof(mnt[0])*r1, unix.ST_NOWAIT /* ST_NOWAIT */)
		if errno != 0 {
			return nil, fmt.Errorf("getvfsstat: %s", string(errno))
		}
		if r1 == r2 {
			break
		}
	}

	stats = []filesystemStats{}
	for _, v := range mnt {
		mountpoint := unix.ByteSliceToString(v.F_mntonname[:])
		if c.mountPointFilter.ignored(mountpoint) {
			c.logger.Debug("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := unix.ByteSliceToString(v.F_mntfromname[:])
		fstype := unix.ByteSliceToString(v.F_fstypename[:])
		if c.fsTypeFilter.ignored(fstype) {
			c.logger.Debug("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (v.F_flag & unix.MNT_RDONLY) != 0 {
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
