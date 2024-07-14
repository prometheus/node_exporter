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

//go:build !nofilesystem
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

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/sys/unix"
)

const (
	defMountPointsExcluded = "^/(dev|proc|run/credentials/.+|sys|var/lib/docker/.+|var/lib/containers/storage/.+)($|/)"
	defFSTypesExcluded     = "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
)

var mountTimeout = kingpin.Flag("collector.filesystem.mount-timeout",
	"how long to wait for a mount to respond before marking it as stale").
	Hidden().Default("5s").Duration()
var statWorkerCount = kingpin.Flag("collector.filesystem.stat-workers",
	"how many stat calls to process simultaneously").
	Hidden().Default("4").Int()
var stuckMounts = make(map[string]struct{})
var stuckMountsMtx = &sync.Mutex{}

// GetStats returns filesystem stats.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	mps, err := mountPointDetails(c.logger)
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}
	labelChan := make(chan filesystemLabels)
	statChan := make(chan filesystemStats)
	wg := sync.WaitGroup{}

	workerCount := *statWorkerCount
	if workerCount < 1 {
		workerCount = 1
	}

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
			if c.excludedMountPointsPattern.MatchString(labels.mountPoint) {
				level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", labels.mountPoint)
				continue
			}
			if c.excludedFSTypesPattern.MatchString(labels.fsType) {
				level.Debug(c.logger).Log("msg", "Ignoring fs", "type", labels.fsType)
				continue
			}

			stuckMountsMtx.Lock()
			if _, ok := stuckMounts[labels.mountPoint]; ok {
				labels.deviceError = "mountpoint timeout"
				stats = append(stats, filesystemStats{
					labels:      labels,
					deviceError: 1,
				})
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

func (c *filesystemCollector) mountWatcher(mountPoint string, buf *unix.Statfs_t, successCh chan struct{}, errCh chan error) {
	err := unix.Statfs(mountPoint, buf)
	defer func() {
		close(successCh)
		close(errCh)
	}()
	if err != nil {
		level.Debug(c.logger).Log("msg", "Error on statfs() system call", "rootfs", rootfsFilePath(mountPoint), "err", err)
		errCh <- err
		return
	}
	stuckMountsMtx.Lock()
	successCh <- struct{}{}
	// If the mount has been marked as stuck, unmark it and log it's recovery.
	if _, ok := stuckMounts[mountPoint]; ok {
		level.Debug(c.logger).Log("msg", "Mount point has recovered, monitoring will resume", "mountpoint", mountPoint)
		delete(stuckMounts, mountPoint)
	}
	stuckMountsMtx.Unlock()
}

func (c *filesystemCollector) processStat(labels filesystemLabels) filesystemStats {
	var ro float64
	for _, option := range strings.Split(labels.options, ",") {
		if option == "ro" {
			ro = 1
			break
		}
	}

	buf := new(unix.Statfs_t)
	success := make(chan struct{}, 1)
	errCh := make(chan error, 1)

	mountCheckTimer := time.NewTimer(*mountTimeout)
	defer mountCheckTimer.Stop()

	go c.mountWatcher(labels.mountPoint, buf, success, errCh)

	res := filesystemStats{
		labels: labels,
		ro:     ro,
	}

	select {
	case <-success:
		res.size = float64(buf.Blocks) * float64(buf.Bsize)
		res.free = float64(buf.Bfree) * float64(buf.Bsize)
		res.avail = float64(buf.Bavail) * float64(buf.Bsize)
		res.files = float64(buf.Files)
		res.filesFree = float64(buf.Ffree)
	case err := <-errCh:
		labels.deviceError = err.Error()
		res.deviceError = 1
	case <-mountCheckTimer.C:
		// Timed out, mark mount as stuck
		stuckMountsMtx.Lock()
		level.Debug(c.logger).Log("msg", "Mount point timed out, it is being labeled as stuck and will not be monitored", "mountpoint", labels.mountPoint)
		stuckMounts[labels.mountPoint] = struct{}{}
		stuckMountsMtx.Unlock()
		labels.deviceError = "mountpoint timeout"
		res.deviceError = 1
	}

	return res
}

func mountPointDetails(logger log.Logger) ([]filesystemLabels, error) {
	file, err := os.Open(procFilePath("1/mountinfo"))
	if errors.Is(err, os.ErrNotExist) {
		// Fallback to `/proc/self/mountinfo` if `/proc/1/mountinfo` is missing due hidepid.
		level.Debug(logger).Log("msg", "Reading root mounts failed, falling back to self mounts", "err", err)
		file, err = os.Open(procFilePath("self/mountinfo"))
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

		if len(parts) < 10 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}

		major, minor := 0, 0
		_, err := fmt.Sscanf(parts[2], "%d:%d", &major, &minor)
		if err != nil {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}

		m := 5
		for parts[m+1] != "-" {
			m++
		}

		// Ensure we handle the translation of \040 and \011
		// as per fstab(5).
		parts[4] = strings.Replace(parts[4], "\\040", " ", -1)
		parts[4] = strings.Replace(parts[4], "\\011", "\t", -1)

		filesystems = append(filesystems, filesystemLabels{
			device:      parts[m+3],
			mountPoint:  rootfsStripPrefix(parts[4]),
			fsType:      parts[m+2],
			options:     parts[5],
			major:       fmt.Sprint(major),
			minor:       fmt.Sprint(minor),
			deviceError: "",
		})
	}

	return filesystems, scanner.Err()
}
