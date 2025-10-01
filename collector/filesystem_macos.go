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

//go:build darwin && !nofilesystem
// +build darwin,!nofilesystem

package collector

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
#include <stdio.h>

double purgeable(char *path) {
  double value = -1.0f;

  @autoreleasepool {
    NSError *error = nil;
    NSString *str = [NSString stringWithUTF8String:path];
    NSURL *fileURL = [[NSURL alloc] initFileURLWithPath:str];

    NSDictionary *results = [fileURL resourceValuesForKeys:@[NSURLVolumeAvailableCapacityForImportantUsageKey] error:&error];
    if (results) {
      CFNumberRef tmp = CFDictionaryGetValue((CFDictionaryRef)results, NSURLVolumeAvailableCapacityForImportantUsageKey);
      if (tmp != NULL) {
        CFNumberGetValue(tmp, kCFNumberFloat64Type, &value);
      }
    }

    [fileURL release];
  }

  return value;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^devfs$"
	readOnly               = 0x1 // MNT_RDONLY
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	var mntbuf *C.struct_statfs
	count := C.getmntinfo(&mntbuf, C.MNT_NOWAIT)
	if count == 0 {
		return nil, errors.New("getmntinfo() failed")
	}

	mnt := (*[1 << 20]C.struct_statfs)(unsafe.Pointer(mntbuf))
	stats = []filesystemStats{}
	for i := 0; i < int(count); i++ {
		mountpoint := C.GoString(&mnt[i].f_mntonname[0])
		if c.mountPointFilter.ignored(mountpoint) {
			c.logger.Debug("Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := C.GoString(&mnt[i].f_mntfromname[0])
		fstype := C.GoString(&mnt[i].f_fstypename[0])
		if c.fsTypeFilter.ignored(fstype) {
			c.logger.Debug("Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (mnt[i].f_flags & readOnly) != 0 {
			ro = 1
		}

		mountpointCString := C.CString(mountpoint)
		defer C.free(unsafe.Pointer(mountpointCString))

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
			purgeable: float64(C.purgeable(mountpointCString)),
			ro:        ro,
		})
	}
	return stats, nil
}
