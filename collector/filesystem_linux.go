// +build !nofilesystem

package collector

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/prometheus/log"
)

const (
	defIgnoredMountPoints = "^/(sys|proc|dev)($|/)"
)

var (
	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}
)

type filesystemDetails struct {
	device     string
	mountPoint string
	fsType     string
}

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	mpds, err := mountPointDetails()
	if err != nil {
		return nil, err
	}
	stats = []filesystemStats{}
	for _, mpd := range mpds {
		if c.ignoredMountPointsPattern.MatchString(mpd.mountPoint) {
			log.Debugf("Ignoring mount point: %s", mpd.mountPoint)
			continue
		}
		buf := new(syscall.Statfs_t)
		err := syscall.Statfs(mpd.mountPoint, buf)
		if err != nil {
			return nil, fmt.Errorf("Statfs on %s returned %s",
				mpd.mountPoint, err)
		}

		labelValues := []string{mpd.device, mpd.mountPoint, mpd.fsType}
		stats = append(stats, filesystemStats{
			labelValues: labelValues,
			size:        float64(buf.Blocks) * float64(buf.Bsize),
			free:        float64(buf.Bfree) * float64(buf.Bsize),
			avail:       float64(buf.Bavail) * float64(buf.Bsize),
			files:       float64(buf.Files),
			filesFree:   float64(buf.Ffree),
		})
	}
	return stats, nil
}

func mountPointDetails() ([]filesystemDetails, error) {
	file, err := os.Open(procFilePath("mounts"))
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
