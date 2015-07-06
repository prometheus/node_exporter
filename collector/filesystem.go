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
	procMounts          = "/proc/mounts"
	filesystemSubsystem = "filesystem"
)

var (
	ignoredMountPoints = flag.String("collector.filesystem.ignored-mount-points", "^/(sys|proc|dev)($|/)", "Regexp of mount points to ignore for filesystem collector.")
)

type filesystemDetails struct {
	device     string
	mountPoint string
	fsType     string
}

type filesystemCollector struct {
	ignoredMountPointsPattern *regexp.Regexp

	size, free, avail, files, filesFree *prometheus.GaugeVec
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// network device filesystems.
func NewFilesystemCollector() (Collector, error) {
	var filesystemLabelNames = []string{"device", "mountpoint", "fstype"}

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
			return fmt.Errorf("Statfs on %s returned %s", mpd.mountPoint, err)
		}

		c.size.WithLabelValues(mpd.device, mpd.mountPoint, mpd.fsType).Set(float64(buf.Blocks) * float64(buf.Bsize))
		c.free.WithLabelValues(mpd.device, mpd.mountPoint, mpd.fsType).Set(float64(buf.Bfree) * float64(buf.Bsize))
		c.avail.WithLabelValues(mpd.device, mpd.mountPoint, mpd.fsType).Set(float64(buf.Bavail) * float64(buf.Bsize))
		c.files.WithLabelValues(mpd.device, mpd.mountPoint, mpd.fsType).Set(float64(buf.Files))
		c.filesFree.WithLabelValues(mpd.device, mpd.mountPoint, mpd.fsType).Set(float64(buf.Ffree))
	}
	c.size.Collect(ch)
	c.free.Collect(ch)
	c.avail.Collect(ch)
	c.files.Collect(ch)
	c.filesFree.Collect(ch)
	return err
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
