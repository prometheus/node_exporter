// +build !nofilesystem

package collector

import (
	"errors"
	"flag"
	"regexp"
	"unsafe"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <sys/param.h>
#include <sys/ucred.h>
#include <sys/mount.h>
#include <stdio.h>
*/
import "C"

const (
	procMounts          = "/proc/mounts"
	filesystemSubsystem = "filesystem"
)

var (
	ignoredMountPoints = flag.String("collector.filesystem.ignored-mount-points", "^/(dev)($|/)", "Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	ignoredMountPointsPattern *regexp.Regexp

	size, free, avail, files, filesFree *prometheus.GaugeVec
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// filesystems stats.
func NewFilesystemCollector() (Collector, error) {
	var filesystemLabelNames = []string{"filesystem"}

	return &filesystemCollector{
		ignoredMountPointsPattern: regexp.MustCompile(*ignoredMountPoints),
		size: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: filesystemSubsystem,
				Name:      "size",
				Help:      "Filesystem size in bytes.",
			},
			filesystemLabelNames,
		),
		free: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: filesystemSubsystem,
				Name:      "free",
				Help:      "Filesystem free space in bytes.",
			},
			filesystemLabelNames,
		),
		avail: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: filesystemSubsystem,
				Name:      "avail",
				Help:      "Filesystem space available to non-root users in bytes.",
			},
			filesystemLabelNames,
		),
		files: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: filesystemSubsystem,
				Name:      "files",
				Help:      "Filesystem total file nodes.",
			},
			filesystemLabelNames,
		),
		filesFree: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: filesystemSubsystem,
				Name:      "files_free",
				Help:      "Filesystem total free file nodes.",
			},
			filesystemLabelNames,
		),
	}, nil
}

// Expose filesystem fullness.
func (c *filesystemCollector) Update(ch chan<- prometheus.Metric) (err error) {
	var mntbuf *C.struct_statfs
	count := C.getmntinfo(&mntbuf, C.MNT_NOWAIT)
	if count == 0 {
		return errors.New("getmntinfo() failed!")
	}

	mnt := (*[1 << 30]C.struct_statfs)(unsafe.Pointer(mntbuf))
	for i := 0; i < int(count); i++ {
		//printf("path: %s\t%lu\n", mntbuf[i].f_mntonname, mntbuf[i].f_bfree)
		name := C.GoString(&mnt[i].f_mntonname[0])
		if c.ignoredMountPointsPattern.MatchString(name) {
			glog.V(1).Infof("Ignoring mount point: %s", name)
			continue
		}
		c.size.WithLabelValues(name).Set(float64(mnt[i].f_blocks) * float64(mnt[i].f_bsize))
		c.free.WithLabelValues(name).Set(float64(mnt[i].f_bfree) * float64(mnt[i].f_bsize))
		c.avail.WithLabelValues(name).Set(float64(mnt[i].f_bavail) * float64(mnt[i].f_bsize))
		c.files.WithLabelValues(name).Set(float64(mnt[i].f_files))
		c.filesFree.WithLabelValues(name).Set(float64(mnt[i].f_ffree))
	}

	c.size.Collect(ch)
	c.free.Collect(ch)
	c.avail.Collect(ch)
	c.files.Collect(ch)
	c.filesFree.Collect(ch)
	return err
}
