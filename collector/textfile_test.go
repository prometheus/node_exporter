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
	"flag"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestParseTextFiles(t *testing.T) {
	tests := []struct {
		path string
		out  string
	}{
		{
			path: "fixtures/textfile/no_metric_files",
			out:  "fixtures/textfile/no_metric_files.out",
		},
		{
			path: "fixtures/textfile/two_metric_files",
			out:  "fixtures/textfile/two_metric_files.out",
		},
		{
			path: "fixtures/textfile/nonexistent_path",
			out:  "fixtures/textfile/nonexistent_path.out",
		},
	}

	for i, test := range tests {
		c := textFileCollector{
			path: test.path,
		}

		// Suppress a log message about `nonexistent_path` not existing, this is
		// expected and clutters the test output.
		err := flag.Set("log.level", "fatal")
		if err != nil {
			t.Fatal(err)
		}

		mfs := c.parseTextFiles()
		textMFs := make([]string, 0, len(mfs))
		for _, mf := range mfs {
			if mf.GetName() == "node_textfile_mtime" {
				mf.GetMetric()[0].GetGauge().Value = proto.Float64(1)
				mf.GetMetric()[1].GetGauge().Value = proto.Float64(2)
			}
			textMFs = append(textMFs, proto.MarshalTextString(mf))
		}
		sort.Strings(textMFs)
		got := strings.Join(textMFs, "")

		want, err := ioutil.ReadFile(test.out)
		if err != nil {
			t.Fatalf("%d. error reading fixture file %s: %s", i, test.out, err)
		}

		if string(want) != got {
			t.Fatalf("%d. want:\n\n%s\n\ngot:\n\n%s", i, string(want), got)
		}
	}
}
