// Copyright The Prometheus Authors
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
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/prometheus/procfs"
)

// TestNetStatDescsInSync guards against drift between the committed descriptor
// table (netstat_descs_linux.go) and the procfs structs it is generated from.
// If a procfs upgrade adds or removes statistics fields, this test fails until
// the table is regenerated with `GOOS=linux go generate ./collector`.
func TestNetStatDescsInSync(t *testing.T) {
	src, err := os.ReadFile("netstat_descs_linux.go")
	if err != nil {
		t.Fatalf("read netstat_descs_linux.go: %v", err)
	}

	re := regexp.MustCompile(`"([A-Za-z0-9]+_[A-Za-z0-9]+)": prometheus\.NewDesc`)
	committed := map[string]bool{}
	for _, m := range re.FindAllSubmatch(src, -1) {
		committed[string(m[1])] = true
	}

	var n procfs.ProcNetstat
	var s procfs.ProcSnmp
	var s6 procfs.ProcSnmp6
	expected := map[string]bool{}
	for _, v := range []any{
		n.TcpExt, n.IpExt,
		s.Ip, s.Icmp, s.IcmpMsg, s.Tcp, s.Udp, s.UdpLite,
		s6.Ip6, s6.Icmp6, s6.Udp6, s6.UdpLite6,
	} {
		tt := reflect.TypeOf(v)
		for i := 0; i < tt.NumField(); i++ {
			f := tt.Field(i)
			if f.Type != reflect.TypeFor[*float64]() {
				continue
			}
			expected[tt.Name()+"_"+f.Name] = true
		}
	}

	var missing, extra []string
	for k := range expected {
		if !committed[k] {
			missing = append(missing, k)
		}
	}
	for k := range committed {
		if !expected[k] {
			extra = append(extra, k)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)

	if len(missing) != 0 || len(extra) != 0 {
		t.Fatalf("netstat_descs_linux.go is out of sync with procfs (missing=%v extra=%v); run `GOOS=linux go generate ./collector`", missing, extra)
	}
}
