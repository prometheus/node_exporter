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

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	procMounts          = "/proc/mounts"
	filesystemSubsystem = "filesystem"
)

var (
	ignoredMountPoints = flag.String("collector.filesystem.ignored-mount-points", "^/(sys|proc|dev)($|/)", "Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	config                    Config
	ignoredMountPointsPattern *regexp.Regexp

	size, free, avail, files, filesFree *prometheus.GaugeVec
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device filesystems.
func NewFilesystemCollector(config Config) (Collector, error) {
	var filesystemLabelNames = []string{"filesystem"}

	return &filesystemCollector{
		config: config,
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
	mps, err := mountPoints()
	if err != nil {
		return err
	}
	for _, mp := range mps {
		if c.ignoredMountPointsPattern.MatchString(mp) {
			glog.V(1).Infof("Ignoring mount point: %s", mp)
			continue
		}
		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(mp, buf)
		if err != nil {
			return fmt.Errorf("Statfs on %s returned %s", mp, err)
		}
		c.size.WithLabelValues(mp).Set(float64(buf.Blocks) * float64(buf.Bsize))
		c.free.WithLabelValues(mp).Set(float64(buf.Bfree) * float64(buf.Bsize))
		c.avail.WithLabelValues(mp).Set(float64(buf.Bavail) * float64(buf.Bsize))
		c.files.WithLabelValues(mp).Set(float64(buf.Files))
		c.filesFree.WithLabelValues(mp).Set(float64(buf.Ffree))
	}
	c.size.Collect(ch)
	c.free.Collect(ch)
	c.avail.Collect(ch)
	c.files.Collect(ch)
	c.filesFree.Collect(ch)
	return err
}

func mountPoints() ([]string, error) {
	file, err := os.Open(procMounts)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mountPoints := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		mountPoints = append(mountPoints, parts[1])
	}
	return mountPoints, nil
}
