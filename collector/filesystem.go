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
	procMounts = "/proc/mounts"
)

var (
	fsSizeMetric      = prometheus.NewGauge()
	fsFreeMetric      = prometheus.NewGauge()
	fsAvailMetric     = prometheus.NewGauge()
	fsFilesMetric     = prometheus.NewGauge()
	fsFilesFreeMetric = prometheus.NewGauge()

	ignoredMountPoints = flag.String("filesystemIgnoredMountPoints", "^/(sys|proc|dev)($|/)", "Regexp of mount points to ignore for filesystem collector.")
)

type filesystemCollector struct {
	registry                  prometheus.Registry
	config                    Config
	ignoredMountPointsPattern *regexp.Regexp
}

func init() {
	Factories["filesystem"] = NewFilesystemCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device filesystems.
func NewFilesystemCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := filesystemCollector{
		config:                    config,
		registry:                  registry,
		ignoredMountPointsPattern: regexp.MustCompile(*ignoredMountPoints),
	}
	registry.Register(
		"node_filesystem_size",
		"Filesystem size in bytes.",
		prometheus.NilLabels,
		fsSizeMetric,
	)
	registry.Register(
		"node_filesystem_free",
		"Filesystem free space in bytes.",
		prometheus.NilLabels,
		fsFreeMetric,
	)
	registry.Register(
		"node_filesystem_avail",
		"Filesystem space available to non-root users in bytes.",
		prometheus.NilLabels,
		fsAvailMetric,
	)
	registry.Register(
		"node_filesystem_files",
		"Filesystem total file nodes.",
		prometheus.NilLabels,
		fsFilesMetric,
	)
	registry.Register(
		"node_filesystem_files_free",
		"Filesystem total free file nodes.",
		prometheus.NilLabels,
		fsFilesFreeMetric,
	)
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
		fsSizeMetric.Set(map[string]string{"filesystem": mp}, float64(buf.Blocks)*float64(buf.Bsize))
		fsFreeMetric.Set(map[string]string{"filesystem": mp}, float64(buf.Bfree)*float64(buf.Bsize))
		fsAvailMetric.Set(map[string]string{"filesystem": mp}, float64(buf.Bavail)*float64(buf.Bsize))
		fsFilesMetric.Set(map[string]string{"filesystem": mp}, float64(buf.Files))
		fsFilesFreeMetric.Set(map[string]string{"filesystem": mp}, float64(buf.Ffree))
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
