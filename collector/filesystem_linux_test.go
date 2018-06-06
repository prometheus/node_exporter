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
	"testing"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

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

	filesystems, err := mountPointDetails()
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}
