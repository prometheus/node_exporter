// +build !nofilesystem

package collector

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	procMounts  = "/proc/mounts"
	subsystem   = "filesystem"
)

var (
	ignoredMountPoints = flag.String(
		"collector.filesystem.ignored-mount-points",
		"^/(sys|proc|dev)($|/)",
		"Regexp of mount points to ignore for filesystem collector.")
)

type filesystemDetails struct {
	device     string
	mountPoint string
	fsType     string
}

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
	mpds, err := mountPointDetails()
	if err != nil {
		return err
	}
	for _, mpd := range mpds {
		if c.ignoredMountPointsPattern.MatchString(mpd.mountPoint) {
			log.Debugf("Ignoring mount point: %s", mpd.mountPoint)
			continue
		}
		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(mpd.mountPoint, buf)
		if err != nil {
			return fmt.Errorf("Statfs on %s returned %s",
				mpd.mountPoint, err)
		}

		ch <- prometheus.MustNewConstMetric(
			sizeDesc,
			prometheus.GaugeValue,
			float64(buf.Blocks) * float64(buf.Bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			freeDesc,
			prometheus.GaugeValue,
			float64(buf.Bfree) * float64(buf.Bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			availDesc,
			prometheus.GaugeValue,
			float64(buf.Bavail) * float64(buf.Bsize),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			filesDesc,
			prometheus.GaugeValue,
			float64(buf.Files),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)

		ch <- prometheus.MustNewConstMetric(
			filesFreeDesc,
			prometheus.GaugeValue,
			float64(buf.Ffree),
			mpd.device, mpd.mountPoint, mpd.fsType,
		)
	}
	return nil
}

func mountPointDetails() ([]filesystemDetails, error) {
	file, err := os.Open(procMounts)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	filesystems := []filesystemDetails{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		filesystems = append(filesystems, filesystemDetails{parts[0], parts[1], parts[2]})
	}
	return filesystems, nil
}
