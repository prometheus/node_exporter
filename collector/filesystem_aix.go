// Copyright 2024 The Prometheus Authors
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
	"github.com/power-devops/perfstat"
)

const (
	defMountPointsExcluded = "^/(dev|aha)($|/)"
	defFSTypesExcluded     = "^procfs$"
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	fsStat, err := perfstat.FileSystemStat()
	if err != nil {
		return nil, err
	}
	for _, stat := range fsStat {
		if c.mountPointFilter.ignored(stat.MountPoint) {
			c.logger.Debug("Ignoring mount point", "mountpoint", stat.MountPoint)
			continue
		}
		fstype := stat.TypeString()
		if c.fsTypeFilter.ignored(fstype) {
			c.logger.Debug("Ignoring fs type", "type", fstype)
			continue
		}

		ro := 0.0
		if stat.Flags&perfstat.VFS_READONLY != 0 {
			ro = 1.0
		}

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     stat.Device,
				mountPoint: stat.MountPoint,
				fsType:     fstype,
			},
			size:      float64(stat.TotalBlocks / 512.0),
			free:      float64(stat.FreeBlocks / 512.0),
			avail:     float64(stat.FreeBlocks / 512.0), // AIX doesn't distinguish between free and available blocks.
			files:     float64(stat.TotalInodes),
			filesFree: float64(stat.FreeInodes),
			ro:        ro,
		})
	}
	return stats, nil
}
