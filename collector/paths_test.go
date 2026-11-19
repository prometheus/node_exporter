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
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"
)

func TestDefaultProcPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", procfs.DefaultMountPoint}); err != nil {
		t.Fatal(err)
	}

	if got, want := procFilePath("somefile"), "/proc/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := procFilePath("some/file"), "/proc/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomProcPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}

	if got, want := procFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := procFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestDefaultSysPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "/sys"}); err != nil {
		t.Fatal(err)
	}

	if got, want := sysFilePath("somefile"), "/sys/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := sysFilePath("some/file"), "/sys/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestCustomSysPath(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}

	if got, want := sysFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}

	if got, want := sysFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}

func TestWarnUnusablePaths(t *testing.T) {
	t.Run("unset paths produce no warning", func(t *testing.T) {
		procPathSetByUser = false
		sysPathSetByUser = false
		rootfsPathSetByUser = false
		udevDataPathSetByUser = false

		logs := capturePathWarnings(t)
		if logs != "" {
			t.Fatalf("expected no warnings for unset paths, got: %s", logs)
		}
	})

	t.Run("default flag values produce no warning", func(t *testing.T) {
		// Parsing only unrelated flags must not mark path flags as user-set,
		// even when their default values are applied.
		if _, err := kingpin.CommandLine.Parse([]string{}); err != nil {
			t.Fatal(err)
		}
		if procPathSetByUser || sysPathSetByUser || rootfsPathSetByUser || udevDataPathSetByUser {
			t.Fatalf("default paths were treated as explicitly set: proc=%v sys=%v rootfs=%v udev=%v",
				procPathSetByUser, sysPathSetByUser, rootfsPathSetByUser, udevDataPathSetByUser)
		}

		logs := capturePathWarnings(t)
		if logs != "" {
			t.Fatalf("expected no warnings for default paths, got: %s", logs)
		}
	})

	t.Run("missing path", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does-not-exist")
		logs := withExplicitPath(t, &procPathSetByUser, procPath, missing)
		assertPathWarning(t, logs, "procfs", missing, "path does not exist")
	})

	t.Run("file instead of directory", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "not-a-dir")
		if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		logs := withExplicitPath(t, &procPathSetByUser, procPath, file)
		assertPathWarning(t, logs, "procfs", file, "path is not a directory")
	})

	t.Run("unreadable directory", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("directory permissions are not enforced the same way on Windows")
		}
		if os.Geteuid() == 0 {
			t.Skip("root can read directories regardless of mode")
		}
		dir := filepath.Join(t.TempDir(), "unreadable")
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(dir, 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = os.Chmod(dir, 0o755)
		})

		logs := withExplicitPath(t, &rootfsPathSetByUser, rootfsPath, dir)
		assertPathWarning(t, logs, "rootfs", dir, "path is not accessible")
	})

	t.Run("procfs missing expected entries", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "uptime"), []byte("1"), 0o644); err != nil {
			t.Fatal(err)
		}
		logs := withExplicitPath(t, &procPathSetByUser, procPath, dir)
		assertPathWarning(t, logs, "procfs", dir, "path is missing expected entries (self or stat)")
	})

	t.Run("procfs with one expected entry is usable", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "self"), 0o755); err != nil {
			t.Fatal(err)
		}
		logs := withExplicitPath(t, &procPathSetByUser, procPath, dir)
		if logs != "" {
			t.Fatalf("expected no warning when procfs contains self, got: %s", logs)
		}
	})

	t.Run("sysfs missing expected entries", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "uevent"), []byte("1"), 0o644); err != nil {
			t.Fatal(err)
		}
		logs := withExplicitPath(t, &sysPathSetByUser, sysPath, dir)
		assertPathWarning(t, logs, "sysfs", dir, "path is missing expected entries (devices or class)")
	})

	t.Run("sysfs with one expected entry is usable", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "class"), 0o755); err != nil {
			t.Fatal(err)
		}
		logs := withExplicitPath(t, &sysPathSetByUser, sysPath, dir)
		if logs != "" {
			t.Fatalf("expected no warning when sysfs contains class, got: %s", logs)
		}
	})
}

// withExplicitPath marks one path flag as user-set and clears the others so
// host defaults are not inspected.
func withExplicitPath(t *testing.T, setByUser *bool, path *string, value string) string {
	t.Helper()

	oldProc, oldProcSet := *procPath, procPathSetByUser
	oldSys, oldSysSet := *sysPath, sysPathSetByUser
	oldRootfs, oldRootfsSet := *rootfsPath, rootfsPathSetByUser
	oldUdev, oldUdevSet := *udevDataPath, udevDataPathSetByUser
	t.Cleanup(func() {
		*procPath, procPathSetByUser = oldProc, oldProcSet
		*sysPath, sysPathSetByUser = oldSys, oldSysSet
		*rootfsPath, rootfsPathSetByUser = oldRootfs, oldRootfsSet
		*udevDataPath, udevDataPathSetByUser = oldUdev, oldUdevSet
	})

	procPathSetByUser = false
	sysPathSetByUser = false
	rootfsPathSetByUser = false
	udevDataPathSetByUser = false
	*setByUser = true
	*path = value

	return capturePathWarnings(t)
}

func capturePathWarnings(t *testing.T) string {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
	WarnUnusablePaths(logger)
	return buf.String()
}

func assertPathWarning(t *testing.T, logs, name, path, reason string) {
	t.Helper()
	wantMsg := "Configured " + name + " path is not usable"
	if !strings.Contains(logs, wantMsg) {
		t.Fatalf("logs %q do not contain %q", logs, wantMsg)
	}
	if !strings.Contains(logs, "path="+path) && !strings.Contains(logs, "path=\""+path+"\"") {
		t.Fatalf("logs %q do not contain path %q", logs, path)
	}
	if !strings.Contains(logs, reason) {
		t.Fatalf("logs %q do not contain reason %q", logs, reason)
	}
}
