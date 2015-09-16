// +build !nofilesystem

package collector

import (
	"errors"
	"regexp"
	"unsafe"

	"github.com/prometheus/log"
)

/*
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
#include <stdio.h>
*/
import "C"

const (
	defIgnoredMountPoints = "^/(dev)($|/)"
)

var (
	filesystemLabelNames = []string{"filesystem"}
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	var mntbuf *C.struct_statfs
	count := C.getmntinfo(&mntbuf, C.MNT_NOWAIT)
	if count == 0 {
		return nil, errors.New("getmntinfo() failed")
	}

	mnt := (*[1 << 30]C.struct_statfs)(unsafe.Pointer(mntbuf))
	stats = []filesystemStats{}
	for i := 0; i < int(count); i++ {
		name := C.GoString(&mnt[i].f_mntonname[0])
		if c.ignoredMountPointsPattern.MatchString(name) {
			log.Debugf("Ignoring mount point: %s", name)
			continue
		}
		labelValues := []string{name}
		stats = append(stats, filesystemStats{
			labelValues: labelValues,
			size:        float64(mnt[i].f_blocks) * float64(mnt[i].f_bsize),
			free:        float64(mnt[i].f_bfree) * float64(mnt[i].f_bsize),
			avail:       float64(mnt[i].f_bavail) * float64(mnt[i].f_bsize),
			files:       float64(mnt[i].f_files),
			filesFree:   float64(mnt[i].f_ffree),
		})
	}
	return stats, nil
}
