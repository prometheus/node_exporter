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

// +build !nofilesystem

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/sys/unix"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	defMountPointsExcluded = "^/(dev|proc|sys|var/lib/docker/.+)($|/)"
	defFSTypesExcluded     = "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
)

var mountTimeout = kingpin.Flag("collector.filesystem.mount-timeout",
	"how long to wait for a mount to respond before marking it as stale").
	Hidden().Default("5s").Duration()
var statWorkerCount = kingpin.Flag("collector.filesystem.stat-workers",
	"how many stat calls to process simultaneously. Providing 0 will attempt to process all stat calls simultaneously").
	Hidden().Default("4").Int()
var stuckMounts = make(map[string]struct{})
var stuckMountsMtx = &sync.Mutex{}

// GetStats returns filesystem stats.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	mpsAll, err := mountPointDetails(c.logger)
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}

	labelChan := make(chan filesystemLabels)
	statChan := make(chan filesystemStats)
	wg := sync.WaitGroup{}

	mps := []filesystemLabels{}

	// Remove all ignored mount points.
	for _, labels := range mpsAll {
		if c.excludedMountPointsPattern.MatchString(labels.mountPoint) {
			level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", labels.mountPoint)
			continue
		}
		if c.excludedFSTypesPattern.MatchString(labels.fsType) {
			level.Debug(c.logger).Log("msg", "Ignoring fs", "type", labels.fsType)
			continue
		}
		mps = append(mps, labels)
	}

	workerCount := *statWorkerCount
	if workerCount == 0 {
		workerCount = len(mps)
	} else if workerCount < 0 {
		workerCount = 1
	}

	// Effectively creates a threadpool to distribute stat calls amongst. Allows multiple stat calls to be run in parallel.
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for labels := range labelChan {
				statChan <- c.processStat(labels)
			}
		}()
	}

	go func() {
		for _, labels := range mps {
			stuckMountsMtx.Lock()
			if _, ok := stuckMounts[labels.mountPoint]; ok {
				statChan <- filesystemStats{
					labels:      labels,
					deviceError: 1,
				}
				level.Debug(c.logger).Log("msg", "Mount point is in an unresponsive state", "mountpoint", labels.mountPoint)
				stuckMountsMtx.Unlock()
				continue
			}
			stuckMountsMtx.Unlock()
			labelChan <- labels
		}
		close(labelChan)
		wg.Wait()
		close(statChan)
	}()

	for stat := range statChan {
		stats = append(stats, stat)
	}
	return stats, nil
}

func (c *filesystemCollector) processStat(labels filesystemLabels) filesystemStats {
	// The done channel is used do tell the "watcher" that the stat call completed.
	done := make(chan error)

	// Run Statfs call which may hang
	buf := new(unix.Statfs_t)
	go runStatfs(buf, labels, done, c.logger)

	// Waits for Statfs call to complete or timeout. If Statfs call does not timeout
	// return whatever Statfs returns. Otherwise, return error if call times out.
	err := waitStatfs(labels.mountPoint, done, c.logger)
	close(done)

	if err != nil {
		level.Debug(c.logger).Log("msg", "Error on statfs() system call", "rootfs", rootfsFilePath(labels.mountPoint), "err", err)
		return filesystemStats{
			labels:      labels,
			deviceError: 1,
		}
	}

	var ro float64
	for _, option := range strings.Split(labels.options, ",") {
		if option == "ro" {
			ro = 1
			break
		}
	}

	return filesystemStats{
		labels:    labels,
		size:      float64(buf.Blocks) * float64(buf.Bsize),
		free:      float64(buf.Bfree) * float64(buf.Bsize),
		avail:     float64(buf.Bavail) * float64(buf.Bsize),
		files:     float64(buf.Files),
		filesFree: float64(buf.Ffree),
		ro:        ro,
	}
}

// runStatfs runs the unix.Statfs command which may hang. On completion of the Statfs call
// stuckMounts is used to determine whether the call resulted in a timeout, if so unset the
// stuckMounts flag.
func runStatfs(buf *unix.Statfs_t, labels filesystemLabels, done chan error, logger log.Logger) {
	err := unix.Statfs(rootfsFilePath(labels.mountPoint), buf)
	stuckMountsMtx.Lock()
	// If the mount has been marked as stuck, unmark it and log it's recovery.
	if _, ok := stuckMounts[labels.mountPoint]; ok {
		level.Debug(logger).Log("msg", "Mount point has recovered, monitoring will resume", "mountpoint", labels.mountPoint)
		delete(stuckMounts, labels.mountPoint)
	} else {
		done <- err
	}
	stuckMountsMtx.Unlock()
}

// waitStatfs listens on the given done channel and returns whatever error is
// returned by the Statfs call. If instead the timeout is reached, the mount point
// is marked as stuck and returns with error.
func waitStatfs(mountPoint string, done chan error, logger log.Logger) error {
	var err error
	select {
	case err = <-done:
		// Success.
	case <-time.After(*mountTimeout):
		// Timed out, mark mount as stuck.
		stuckMountsMtx.Lock()
		select {
		case err = <-done:
			// Success came in just after the timeout was reached, don't label the mount as stuck.
		default:
			err = errors.New("Statfs timed out")
			level.Debug(logger).Log("msg", "Mount point timed out, it is being labeled as stuck and will not be monitored", "mountpoint", mountPoint)
			stuckMounts[mountPoint] = struct{}{}
		}
		stuckMountsMtx.Unlock()
	}
	return err
}

func mountPointDetails(logger log.Logger) ([]filesystemLabels, error) {
	file, err := os.Open(procFilePath("1/mounts"))
	if errors.Is(err, os.ErrNotExist) {
		// Fallback to `/proc/mounts` if `/proc/1/mounts` is missing due hidepid.
		level.Debug(logger).Log("msg", "Reading root mounts failed, falling back to system mounts", "err", err)
		file, err = os.Open(procFilePath("mounts"))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseFilesystemLabels(file)
}

func parseFilesystemLabels(r io.Reader) ([]filesystemLabels, error) {
	var filesystems []filesystemLabels

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		if len(parts) < 4 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}

		// Ensure we handle the translation of \040 and \011
		// as per fstab(5).
		parts[1] = strings.Replace(parts[1], "\\040", " ", -1)
		parts[1] = strings.Replace(parts[1], "\\011", "\t", -1)

		filesystems = append(filesystems, filesystemLabels{
			device:     parts[0],
			mountPoint: rootfsStripPrefix(parts[1]),
			fsType:     parts[2],
			options:    parts[3],
		})
	}

	return filesystems, scanner.Err()
}
