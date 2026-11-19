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

package collector

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"
)

var (
	// The path of the proc filesystem.
	procPathSetByUser     bool
	procPath              = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).IsSetByUser(&procPathSetByUser).String()
	sysPathSetByUser      bool
	sysPath               = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").IsSetByUser(&sysPathSetByUser).String()
	rootfsPathSetByUser   bool
	rootfsPath            = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").IsSetByUser(&rootfsPathSetByUser).String()
	udevDataPathSetByUser bool
	udevDataPath          = kingpin.Flag("path.udev.data", "udev data path.").Default("/run/udev/data").IsSetByUser(&udevDataPathSetByUser).String()
)

// configuredPath is a filesystem path the user can override with a flag.
type configuredPath struct {
	name      string
	path      *string
	setByUser *bool
	// markers, if non-empty, name entries that should exist for this path to
	// look like the expected filesystem. Any one match is enough.
	markers []string
}

func configuredPaths() []configuredPath {
	return []configuredPath{
		{name: "procfs", path: procPath, setByUser: &procPathSetByUser, markers: []string{"self", "stat"}},
		{name: "sysfs", path: sysPath, setByUser: &sysPathSetByUser, markers: []string{"devices", "class"}},
		{name: "rootfs", path: rootfsPath, setByUser: &rootfsPathSetByUser},
		{name: "udev data", path: udevDataPath, setByUser: &udevDataPathSetByUser},
	}
}

// WarnUnusablePaths logs a warning for filesystem paths the user explicitly
// configured when those paths are missing, not directories, unreadable, or
// (for procfs and sysfs) clearly not the expected filesystem. Default paths
// are not checked. Warnings never stop startup.
func WarnUnusablePaths(logger *slog.Logger) {
	if logger == nil {
		return
	}
	for _, cp := range configuredPaths() {
		if cp.setByUser == nil || !*cp.setByUser || cp.path == nil {
			continue
		}
		if reason := unusablePathReason(*cp.path, cp.markers); reason != "" {
			logger.Warn(fmt.Sprintf("Configured %s path is not usable", cp.name), "path", *cp.path, "reason", reason)
		}
	}
}

// unusablePathReason reports why path cannot be used, or "" when it is usable.
// markers are optional names; if set, at least one must exist in the directory.
func unusablePathReason(path string, markers []string) string {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "path does not exist"
		}
		return fmt.Sprintf("path is not accessible: %s", err)
	}
	if !info.IsDir() {
		return "path is not a directory"
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("path is not accessible: %s", err)
	}
	if len(markers) == 0 {
		return ""
	}
	present := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		present[entry.Name()] = struct{}{}
	}
	for _, marker := range markers {
		if _, ok := present[marker]; ok {
			return ""
		}
	}
	return fmt.Sprintf("path is missing expected entries (%s)", strings.Join(markers, " or "))
}

func procFilePath(name string) string {
	return filepath.Join(*procPath, name)
}

func sysFilePath(name string) string {
	return filepath.Join(*sysPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(*rootfsPath, name)
}

func udevDataFilePath(name string) string {
	return filepath.Join(*udevDataPath, name)
}

func rootfsStripPrefix(path string) string {
	if *rootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, *rootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
