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
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"
)

var (
	// The path of the proc filesystem.
	procPath     = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	sysPath      = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	rootfsPath   = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
	udevDataPath = kingpin.Flag("path.udev.data", "udev data path.").Default("/run/udev/data").String()
)

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
