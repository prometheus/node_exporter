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
	"flag"
	"path"

	"github.com/prometheus/procfs"
)

var (
	// The path of the proc filesystem.
	procPath = flag.String("collector.procfs", procfs.DefaultMountPoint, "procfs mountpoint.")
	sysPath  = flag.String("collector.sysfs", "/sys", "sysfs mountpoint.")
)

func procFilePath(name string) string {
	return path.Join(*procPath, name)
}

func sysFilePath(name string) string {
	return path.Join(*sysPath, name)
}
