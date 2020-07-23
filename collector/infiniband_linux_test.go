// Copyright 2017 The Prometheus Authors
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
	"testing"
)

func TestInfiniBandDevices(t *testing.T) {
	devices, err := infinibandDevices("fixtures/sys/class/infiniband")
	if err != nil {
		t.Fatal(err)
	}

	if l := len(devices); l != 2 {
		t.Fatalf("Retrieved an unexpected number of InfiniBand devices: %d", l)
	}
}

func TestInfiniBandPorts(t *testing.T) {
	ports, err := infinibandPorts("fixtures/sys/class/infiniband", "mlx4_0")
	if err != nil {
		t.Fatal(err)
	}

	if l := len(ports); l != 2 {
		t.Fatalf("Retrieved an unexpected number of InfiniBand ports: %d", l)
	}
}
