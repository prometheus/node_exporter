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

//go:build !nodrbd

package collector

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestDRBDNetworkAndDiskVolumes(t *testing.T) {
	oldProcPath := *procPath
	t.Cleanup(func() { *procPath = oldProcPath })
	*procPath = "fixtures/proc"

	c, err := newDRBDCollector(slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 100)
	if err := c.Update(ch); err != nil {
		t.Fatal(err)
	}
	close(ch)

	got := map[string]float64{}
	for m := range ch {
		var pb dto.Metric
		if err := m.Write(&pb); err != nil {
			t.Fatal(err)
		}
		desc := m.Desc().String()
		for _, name := range []string{
			"node_drbd_network_sent_bytes_total",
			"node_drbd_network_received_bytes_total",
			"node_drbd_disk_written_bytes_total",
			"node_drbd_disk_read_bytes_total",
		} {
			if strings.Contains(desc, `"`+name+`"`) {
				got[name] = pb.GetCounter().GetValue()
			}
		}
	}

	// The /proc/drbd counters ns, nr, dw and dr are all in KiB.
	want := map[string]float64{
		"node_drbd_network_sent_bytes_total":     17324442 * 1024,
		"node_drbd_network_received_bytes_total": 10961011 * 1024,
		"node_drbd_disk_written_bytes_total":     28263521 * 1024,
		"node_drbd_disk_read_bytes_total":        118696670 * 1024,
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("%s = %v, want %v", name, got[name], w)
		}
	}
}
