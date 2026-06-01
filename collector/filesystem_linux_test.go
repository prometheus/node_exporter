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

package collector

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"

	"github.com/prometheus/procfs"
)

func Test_parseFilesystemLabelsError(t *testing.T) {
	tests := []struct {
		name string
		in   []*procfs.MountInfo
	}{
		{
			name: "malformed Major:Minor",
			in: []*procfs.MountInfo{
				{
					MajorMinorVer: "nope",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseFilesystemLabels(tt.in); err == nil {
				t.Fatal("expected an error, but none occurred")
			}
		})
	}
}

func Test_isFilesystemReadOnly(t *testing.T) {
	tests := map[string]struct {
		labels   filesystemLabels
		expected bool
	}{
		"/media/volume1": {
			labels: filesystemLabels{
				mountOptions: "rw,nosuid,nodev,noexec,relatime",
				superOptions: "rw,devices",
			},
			expected: false,
		},
		"/media/volume2": {
			labels: filesystemLabels{
				mountOptions: "ro,relatime",
				superOptions: "rw,fd=22,pgrp=1,timeout=300,minproto=5,maxproto=5,direct",
			}, expected: true,
		},
		"/media/volume3": {
			labels: filesystemLabels{
				mountOptions: "rw,user_id=1000,group_id=1000",
				superOptions: "ro",
			}, expected: true,
		},
		"/media/volume4": {
			labels: filesystemLabels{
				mountOptions: "ro,nosuid,noexec",
				superOptions: "ro,nodev",
			}, expected: true,
		},
	}

	for _, tt := range tests {
		if got := isFilesystemReadOnly(tt.labels); got != tt.expected {
			t.Errorf("Expected %t, got %t", tt.expected, got)
		}
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
		"/var/lib/containers/storage/overlay": "",
	}

	filesystems, err := mountPointDetails(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Log(err)
	}

	foundSet := map[string]bool{}
	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
		foundSet[fs.mountPoint] = true
	}

	for mountPoint := range expected {
		if _, ok := foundSet[mountPoint]; !ok {
			t.Errorf("Expected %s, got nothing", mountPoint)
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

	filesystems, err := mountPointDetails(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}

func TestGetMountInfoWithRetry(t *testing.T) {
	t.Run("retries parse errors", func(t *testing.T) {
		attempts := 0
		expected := []*procfs.MountInfo{{MountPoint: "/"}}

		got, err := getMountInfoWithRetry(slog.New(slog.NewTextHandler(io.Discard, nil)), func() ([]*procfs.MountInfo, error) {
			attempts++
			if attempts < mountInfoParseRetries {
				return nil, fmt.Errorf("%w: Too few fields in mount string: truncated", procfs.ErrFileParse)
			}
			return expected, nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempts != mountInfoParseRetries {
			t.Fatalf("expected %d attempts, got %d", mountInfoParseRetries, attempts)
		}
		if len(got) != len(expected) || got[0].MountPoint != expected[0].MountPoint {
			t.Fatalf("unexpected mountinfo: %+v", got)
		}
	})

	t.Run("does not retry non parse errors", func(t *testing.T) {
		attempts := 0

		_, err := getMountInfoWithRetry(slog.New(slog.NewTextHandler(io.Discard, nil)), func() ([]*procfs.MountInfo, error) {
			attempts++
			return nil, os.ErrNotExist
		})
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
		if attempts != 1 {
			t.Fatalf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("returns parse error after max retries", func(t *testing.T) {
		attempts := 0

		_, err := getMountInfoWithRetry(slog.New(slog.NewTextHandler(io.Discard, nil)), func() ([]*procfs.MountInfo, error) {
			attempts++
			return nil, fmt.Errorf("%w: Too few fields in mount string: truncated", procfs.ErrFileParse)
		})
		if !errors.Is(err, procfs.ErrFileParse) {
			t.Fatalf("expected procfs.ErrFileParse, got %v", err)
		}
		if attempts != mountInfoParseRetries {
			t.Fatalf("expected %d attempts, got %d", mountInfoParseRetries, attempts)
		}
	})
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

	filesystems, err := mountPointDetails(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Log(err)
	}

	for _, fs := range filesystems {
		if _, ok := expected[fs.mountPoint]; !ok {
			t.Errorf("Got unexpected %s", fs.mountPoint)
		}
	}
}
