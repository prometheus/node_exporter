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

//go:build !novmstat

package collector

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func TestVMStatMalformedLine(t *testing.T) {
	for name, content := range map[string]string{
		"empty line":    "pgfault 10\n\npgmajfault 2\n",
		"missing value": "pgfault 10\npgmajfault\n",
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "vmstat"), []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", dir}); err != nil {
				t.Fatal(err)
			}
			c, err := NewvmStatCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatal(err)
			}
			ch := make(chan prometheus.Metric, 16)
			if err := c.Update(ch); err == nil {
				t.Error("expected an error for a malformed vmstat line")
			}
		})
	}
}
