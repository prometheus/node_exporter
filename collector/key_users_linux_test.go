// Copyright 2023 The Prometheus Authors
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

//go:build linux && !nokey_users
// +build linux,!nokey_users

package collector

import (
	"os"
	"testing"
)

func TestKeyUsers(t *testing.T) {
	file, err := os.Open("fixtures/proc/key-users")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	keyUsers, err := parseKeyUsers(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 51.0, keyUsers[keyUsersEntry{0, "usage"}]; want != got {
		t.Errorf("want uid 0 usage %f, got %f", want, got)
	}
	if want, got := 50.0, keyUsers[keyUsersEntry{0, "nkeys"}]; want != got {
		t.Errorf("want uid 0 nkeys %f, got %f", want, got)
	}
	if want, got := 50.0, keyUsers[keyUsersEntry{0, "nikeys"}]; want != got {
		t.Errorf("want uid 0 nikeys %f, got %f", want, got)
	}
	if want, got := 44.0, keyUsers[keyUsersEntry{0, "qnkeys"}]; want != got {
		t.Errorf("want uid 0 qnkeys %f, got %f", want, got)
	}
	if want, got := 1000000.0, keyUsers[keyUsersEntry{0, "maxkeys"}]; want != got {
		t.Errorf("want uid 0 maxkeys %f, got %f", want, got)
	}
	if want, got := 883.0, keyUsers[keyUsersEntry{0, "qnbytes"}]; want != got {
		t.Errorf("want uid 0 qnbytes %f, got %f", want, got)
	}
	if want, got := 25000000.0, keyUsers[keyUsersEntry{0, "maxbytes"}]; want != got {
		t.Errorf("want uid 0 maxbytes %f, got %f", want, got)
	}

	if want, got := 11.0, keyUsers[keyUsersEntry{1000, "usage"}]; want != got {
		t.Errorf("want uid 0 usage %f, got %f", want, got)
	}
	if want, got := 11.0, keyUsers[keyUsersEntry{1000, "nkeys"}]; want != got {
		t.Errorf("want uid 0 nkeys %f, got %f", want, got)
	}
	if want, got := 12.0, keyUsers[keyUsersEntry{1000, "nikeys"}]; want != got {
		t.Errorf("want uid 0 nikeys %f, got %f", want, got)
	}
	if want, got := 13.0, keyUsers[keyUsersEntry{1000, "qnkeys"}]; want != got {
		t.Errorf("want uid 0 qnkeys %f, got %f", want, got)
	}
	if want, got := 200.0, keyUsers[keyUsersEntry{1000, "maxkeys"}]; want != got {
		t.Errorf("want uid 0 maxkeys %f, got %f", want, got)
	}
	if want, got := 250.0, keyUsers[keyUsersEntry{1000, "qnbytes"}]; want != got {
		t.Errorf("want uid 0 qnbytes %f, got %f", want, got)
	}
	if want, got := 20000.0, keyUsers[keyUsersEntry{1000, "maxbytes"}]; want != got {
		t.Errorf("want uid 0 maxbytes %f, got %f", want, got)
	}
}
