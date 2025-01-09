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
	"reflect"
	"strings"
	"testing"
)

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

func TestParseOSSupportEnd(t *testing.T) {
	nixOSCurrentSystem = "fixtures/run/current-system"
	nixOSBootedSystem = "fixtures/run/booted-system"
	nixOSstoreDB = "fixtures/nix.sqlite"

	want := &osRelease{
		BuildID:         "23.11.20240328.219951b",
		Name:            "NixOS",
		ID:              "nixos",
		IDLike:          "",
		ImageID:         "511hrn15gag56l7z7lm6zkxy6rz4i9gp",
		ImageVersion:    "",
		PrettyName:      "NixOS 23.11 (Tapir)",
		SupportEnd:      "2024-06-30",
		Variant:         "",
		VariantID:       "",
		Version:         "23.11 (Tapir)",
		VersionID:       "23.11",
		VersionCodename: "tapir",

		BootedImageID:   "511hrn15gag56l7z7lm6zkxy6rz4i9gp",
		BootedImageTime: 1736255100,
		ImageTime:       1736255100,
	}

	got, err := parseOSRelease(strings.NewReader(nixosTapir))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %+v osRelease: got %+v", want, got)
	}
}
