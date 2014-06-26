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
	filesystemLabelNames = []string{"filesystem"}

	fsSizeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: filesystemSubsystem,
			Name:      "size",
			Help:      "Filesystem size in bytes.",
		},
		filesystemLabelNames,
	)
	fsFreeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: filesystemSubsystem,
			Name:      "free",
			Help:      "Filesystem free space in bytes.",
		},
		filesystemLabelNames,
	)
	fsAvailMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: filesystemSubsystem,
			Name:      "avail",
			Help:      "Filesystem space available to non-root users in bytes.",
		},
		filesystemLabelNames,
	)
	fsFilesMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: filesystemSubsystem,
			Name:      "files",
			Help:      "Filesystem total file nodes.",
		},
		filesystemLabelNames,
	)
	fsFilesFreeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: filesystemSubsystem,
			Name:      "files_free",
			Help:      "Filesystem total free file nodes.",
		},
		filesystemLabelNames,
	)

	ignoredMountPoints = flag.String("filesystemIgnoredMountPoints", "^/(sys|proc|dev)($|/)", "Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	config                    Config
	ignoredMountPointsPattern *regexp.Regexp
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device filesystems.
func NewFilesystemCollector(config Config) (Collector, error) {
	c := filesystemCollector{
		config: config,
		ignoredMountPointsPattern: regexp.MustCompile(*ignoredMountPoints),
	}
	if _, err := prometheus.RegisterOrGet(fsSizeMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(fsFreeMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(fsAvailMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(fsFilesMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(fsFilesFreeMetric); err != nil {
		return nil, err
	}
	return &c, nil
}

// Expose filesystem fullness.
func (c *filesystemCollector) Update() (updates int, err error) {
	mps, err := mountPoints()
	if err != nil {
		return updates, err
	}
	for _, mp := range mps {
		if c.ignoredMountPointsPattern.MatchString(mp) {
			glog.V(1).Infof("Ignoring mount point: %s", mp)
			continue
		}
		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(mp, buf)
		if err != nil {
			return updates, fmt.Errorf("Statfs on %s returned %s", mp, err)
		}
		fsSizeMetric.WithLabelValues(mp).Set(float64(buf.Blocks) * float64(buf.Bsize))
		fsFreeMetric.WithLabelValues(mp).Set(float64(buf.Bfree) * float64(buf.Bsize))
		fsAvailMetric.WithLabelValues(mp).Set(float64(buf.Bavail) * float64(buf.Bsize))
		fsFilesMetric.WithLabelValues(mp).Set(float64(buf.Files))
		fsFilesFreeMetric.WithLabelValues(mp).Set(float64(buf.Ffree))
		updates++
	}
	return updates, err
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
