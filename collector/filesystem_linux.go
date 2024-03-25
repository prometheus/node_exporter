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

var statWorkerCount = kingpin.Flag("collector.filesystem.stat-workers",
	"how many stat calls to process simultaneously").
	Hidden().Default("4").Int()
var stuckMountsMap = make(map[string]struct{})
var stuckMountsMutex = &sync.Mutex{}

// GetStats returns filesystem stats.
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	fsLabels, err := mountPointDetails(c.logger)
	if err != nil {
		return nil, err
	}
	var fsStats []filesystemStats
	fsLabelChan := make(chan filesystemLabels)
	fsStatChan := make(chan filesystemStats)
	wg := sync.WaitGroup{}
	workerCount := *statWorkerCount
	if workerCount < 1 {
		workerCount = 1
	}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fsLabel := range fsLabelChan {
				fsStatChan <- c.processStat(fsLabel)
			}
		}()
	}

	go func() {
		for _, fsLabel := range fsLabels {
			if c.excludedMountPointsPattern.MatchString(fsLabel.mountPoint) {
				level.Debug(c.logger).Log("msg", "Ignoring mount point", "mountpoint", fsLabel.mountPoint)
				continue
			}
			if c.excludedFSTypesPattern.MatchString(fsLabel.fsType) {
				level.Debug(c.logger).Log("msg", "Ignoring fs", "type", fsLabel.fsType)
				continue
			}
			fsLabelChan <- fsLabel
		}
		close(fsLabelChan)
		wg.Wait()
		close(fsStatChan)
	}()

	for fsStat := range fsStatChan {
		fsStats = append(fsStats, fsStat)
	}
	return fsStats, nil
}

func (c *filesystemCollector) processStat(fsLabel filesystemLabels) filesystemStats {
	var ro float64
	for _, option := range strings.Split(fsLabel.options, ",") {
		if option == "ro" {
			ro = 1
			break
		}
	}

	// If the mount point is stuck, mark it as such and return early.
	// This is done to avoid blocking the stat call indefinitely.
	// NOTE: For instance, this can happen when an NFS mount is unreachable.
	buf := new(unix.Statfs_t)
	statFsErrChan := make(chan error, 1)
	go func(buf *unix.Statfs_t) {
		statFsErrChan <- unix.Statfs(rootfsFilePath(fsLabel.mountPoint), buf)
		close(statFsErrChan)
	}(buf)

	select {
	case err := <-statFsErrChan:
		if err != nil {
			level.Debug(c.logger).Log("msg", "Error on statfs() system call", "rootfs", rootfsFilePath(fsLabel.mountPoint), "err", err)
			fsLabel.deviceError = err.Error()
		}
	case <-time.After(*mountTimeout):
		stuckMountsMutex.Lock()
		if _, ok := stuckMountsMap[fsLabel.mountPoint]; ok {
			level.Debug(c.logger).Log("msg", "Mount point timed out, it is being labeled as stuck and will not be monitored", "mountpoint", fsLabel.mountPoint)
			stuckMountsMap[fsLabel.mountPoint] = struct{}{}
			fsLabel.deviceError = "mountpoint timeout"
		}
		stuckMountsMutex.Unlock()
	}

	// Check if the mount point has recovered and remove it from the stuck map.
	if _, isOpen := <-statFsErrChan; !isOpen {
		stuckMountsMutex.Lock()
		if _, ok := stuckMountsMap[fsLabel.mountPoint]; ok {
			level.Debug(c.logger).Log("msg", "Mount point has recovered, monitoring will resume", "mountpoint", fsLabel.mountPoint)
			delete(stuckMountsMap, fsLabel.mountPoint)
		}
		stuckMountsMutex.Unlock()
	}

	// If the mount point is stuck or statfs errored, mark it as such and return.
	if fsLabel.deviceError != "" {
		return filesystemStats{
			labels:      fsLabel,
			deviceError: 1,
			ro:          ro,
		}
	}

	return filesystemStats{
		labels:    fsLabel,
		size:      float64(buf.Blocks) * float64(buf.Bsize),
		free:      float64(buf.Bfree) * float64(buf.Bsize),
		avail:     float64(buf.Bavail) * float64(buf.Bsize),
		files:     float64(buf.Files),
		filesFree: float64(buf.Ffree),
		ro:        ro,
	}
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
			device:      parts[0],
			mountPoint:  rootfsStripPrefix(parts[1]),
			fsType:      parts[2],
			options:     parts[3],
			deviceError: "",
		})
	}

	return filesystems, scanner.Err()
}
