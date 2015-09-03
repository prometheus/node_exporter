package collector

import (
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
