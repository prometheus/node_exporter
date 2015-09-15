// +build !nofilesystem

package collector

import (
	"errors"
	"flag"
	"regexp"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
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
	subsystem = "filesystem"
)

var (
	ignoredMountPoints = flag.String(
		"collector.filesystem.ignored-mount-points",
		"^/(dev)($|/)",
		"Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	ignoredMountPointsPattern *regexp.Regexp
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

func NewFilesystemCollector() (Collector, error) {
	pattern := regexp.MustCompile(*ignoredMountPoints)
	return &filesystemCollector{
		ignoredMountPointsPattern: pattern,
	}, nil
}

var (
	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}

	sizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "size"),
		"Filesystem size in bytes.",
		filesystemLabelNames, nil,
	)

	freeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "free"),
		"Filesystem free space in bytes.",
		filesystemLabelNames, nil,
	)

	availDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "avail"),
		"Filesystem space available to non-root users in bytes.",
		filesystemLabelNames, nil,
	)

	filesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files"),
		"Filesystem total file nodes.",
		filesystemLabelNames, nil,
	)

	filesFreeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "files_free"),
		"Filesystem total free file nodes.",
		filesystemLabelNames, nil,
	)
)

// Expose filesystem fullness.
func (c *filesystemCollector) Update(ch chan<- prometheus.Metric) (err error) {
	var mntbuf *C.struct_statfs
	count := C.getmntinfo(&mntbuf, C.MNT_NOWAIT)
	if count == 0 {
		return errors.New("getmntinfo() failed")
	}

	mnt := (*[1 << 30]C.struct_statfs)(unsafe.Pointer(mntbuf))
	for i := 0; i < int(count); i++ {
		name := C.GoString(&mnt[i].f_mntonname[0])
		if c.ignoredMountPointsPattern.MatchString(name) {
			log.Debugf("Ignoring mount point: %s", name)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			sizeDesc,
			prometheus.GaugeValue,
			float64(mnt[i].f_blocks) * float64(mnt[i].f_bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			freeDesc,
			prometheus.GaugeValue,
			float64(mnt[i].f_bfree) * float64(mnt[i].f_bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			availDesc,
			prometheus.GaugeValue,
			float64(mnt[i].f_bavail) * float64(mnt[i].f_bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			filesDesc,
			prometheus.GaugeValue,
			float64(mnt[i].f_files),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			filesFreeDesc,
			prometheus.GaugeValue,
			float64(mnt[i].f_ffree),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)
	}
	return nil
}
