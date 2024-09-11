// Copyright 2021 The Prometheus Authors
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
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"testing"
)

const debianBullseye string = `PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
NAME="Debian GNU/Linux"
VERSION_ID="11"
VERSION="11 (bullseye)"
VERSION_CODENAME=bullseye
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"
`

const nixosTapir string = `BUG_REPORT_URL="https://github.com/NixOS/nixpkgs/issues"
BUILD_ID="23.11.20240328.219951b"
DOCUMENTATION_URL="https://nixos.org/learn.html"
HOME_URL="https://nixos.org/"
ID=nixos
LOGO="nix-snowflake"
NAME=NixOS
PRETTY_NAME="NixOS 23.11 (Tapir)"
SUPPORT_END="2024-06-30"
SUPPORT_URL="https://nixos.org/community.html"
VERSION="23.11 (Tapir)"
VERSION_CODENAME=tapir
VERSION_ID="23.11"
`

func TestParseOSRelease(t *testing.T) {
	want := &osRelease{
		Name:            "Ubuntu",
		ID:              "ubuntu",
		IDLike:          "debian",
		PrettyName:      "Ubuntu 20.04.2 LTS",
		SupportEnd:      "",
		Version:         "20.04.2 LTS (Focal Fossa)",
		VersionID:       "20.04",
		VersionCodename: "focal",
	}

	osReleaseFile, err := os.Open("fixtures" + usrLibOSRelease)
	if err != nil {
		t.Fatal(err)
	}
	defer osReleaseFile.Close()

	got, err := parseOSRelease(osReleaseFile)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %+v osRelease: got %+v", want, got)
	}

	want = &osRelease{
		Name:            "Debian GNU/Linux",
		ID:              "debian",
		PrettyName:      "Debian GNU/Linux 11 (bullseye)",
		Version:         "11 (bullseye)",
		VersionID:       "11",
		VersionCodename: "bullseye",
	}
	got, err = parseOSRelease(strings.NewReader(debianBullseye))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %+v osRelease: got %+v", want, got)
	}
}

func TestParseOSSupportEnd(t *testing.T) {
	want := &osRelease{
		BuildID:         "23.11.20240328.219951b",
		Name:            "NixOS",
		ID:              "nixos",
		IDLike:          "",
		ImageID:         "",
		ImageVersion:    "",
		PrettyName:      "NixOS 23.11 (Tapir)",
		SupportEnd:      "2024-06-30",
		Variant:         "",
		VariantID:       "",
		Version:         "23.11 (Tapir)",
		VersionID:       "23.11",
		VersionCodename: "tapir",
	}

	got, err := parseOSRelease(strings.NewReader(nixosTapir))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %+v osRelease: got %+v", want, got)
	}
}

func TestUpdateStruct(t *testing.T) {
	wantedOS := &osRelease{
		Name:            "Ubuntu",
		ID:              "ubuntu",
		IDLike:          "debian",
		PrettyName:      "Ubuntu 20.04.2 LTS",
		Version:         "20.04.2 LTS (Focal Fossa)",
		VersionID:       "20.04",
		VersionCodename: "focal",
	}
	wantedVersion := 20.04

	collector, err := NewOSCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	c := collector.(*osReleaseCollector)

	err = c.UpdateStruct("fixtures" + usrLibOSRelease)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantedOS, c.os) {
		t.Fatalf("should have %+v osRelease: got %+v", wantedOS, c.os)
	}
	if wantedVersion != c.version {
		t.Errorf("Expected '%v' but got '%v'", wantedVersion, c.version)
	}
}
