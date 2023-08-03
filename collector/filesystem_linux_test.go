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
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
)

func Test_parseFilesystemLabelsError(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			name: "too few fields",
			in:   "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseFilesystemLabels(strings.NewReader(tt.in)); err == nil {
				t.Fatal("expected an error, but none occurred")
			}
		})
	}
}

func TestMountPointDetails(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./fixtures/proc"}); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"/":                               "",
		"/sys":                            "",
		"/proc":                           "",
		"/dev":                            "",
		"/dev/pts":                        "",
		"/run":                            "",
		"/sys/kernel/security":            "",
		"/dev/shm":                        "",
		"/run/lock":                       "",
		"/sys/fs/cgroup":                  "",
		"/sys/fs/cgroup/systemd":          "",
		"/sys/fs/pstore":                  "",
		"/sys/fs/cgroup/cpuset":           "",
		"/sys/fs/cgroup/cpu,cpuacct":      "",
		"/sys/fs/cgroup/devices":          "",
		"/sys/fs/cgroup/freezer":          "",
		"/sys/fs/cgroup/net_cls,net_prio": "",
		"/sys/fs/cgroup/blkio":            "",
		"/sys/fs/cgroup/perf_event":       "",
		"/proc/sys/fs/binfmt_misc":        "",
		"/dev/mqueue":                     "",
		"/sys/kernel/debug":               "",
		"/dev/hugepages":                  "",
		"/sys/fs/fuse/connections":        "",
		"/boot":                           "",
		"/run/rpc_pipefs":                 "",
		"/run/user/1000":                  "",
		"/run/user/1000/gvfs":             "",
		"/var/lib/kubelet/plugins/kubernetes.io/vsphere-volume/mounts/[vsanDatastore] bafb9e5a-8856-7e6c-699c-801844e77a4a/kubernetes-dynamic-pvc-3eba5bba-48a3-11e8-89ab-005056b92113.vmdk": "",
		"/var/lib/kubelet/plugins/kubernetes.io/vsphere-volume/mounts/[vsanDatastore]	bafb9e5a-8856-7e6c-699c-801844e77a4a/kubernetes-dynamic-pvc-3eba5bba-48a3-11e8-89ab-005056b92113.vmdk": "",
	}

	filesystems, err := mountPointDetails(log.NewNopLogger())
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}

func TestMountsFallback(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./fixtures_hidepid/proc"}); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"/": "",
	}

	filesystems, err := mountPointDetails(log.NewNopLogger())
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}

func TestPathRootfs(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./fixtures_bindmount/proc", "--path.rootfs", "/host"}); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		// should modify these mountpoints (removes /host, see fixture proc file)
		"/":              "",
		"/media/volume1": "",
		"/media/volume2": "",
		// should not modify these mountpoints
		"/dev/shm":       "",
		"/run/lock":      "",
		"/sys/fs/cgroup": "",
	}

	filesystems, err := mountPointDetails(log.NewNopLogger())
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}
