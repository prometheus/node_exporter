// +build !notextfile

package collector

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

var (
	textFileDirectory = flag.String("collector.textfile.directory", "", "Directory to read text files with metrics from.")
)

type textFileCollector struct {
	path string
}

func init() {
	Factories["textfile"] = NewTextFileCollector
}

// Takes a registers a
// SetMetricFamilyInjectionHook.
func NewTextFileCollector() (Collector, error) {
	c := &textFileCollector{
		path: *textFileDirectory,
	}

	if c.path == "" {
		// This collector is enabled by default, so do not fail if
		// the flag is not passed.
		log.Infof("No directory specified, see --textfile.directory")
	} else {
		prometheus.SetMetricFamilyInjectionHook(c.parseTextFiles)
	}

	return c, nil
}

// textFile collector works via SetMetricFamilyInjectionHook in parseTextFiles.
func (c *textFileCollector) Update(ch chan<- prometheus.Metric) (err error) {
	return nil
}

func (c *textFileCollector) parseTextFiles() []*dto.MetricFamily {
	error := 0.0
	metricFamilies := make([]*dto.MetricFamily, 0)
	mtimes := map[string]time.Time{}

	// Iterate over files and accumulate their metrics.
	files, err := ioutil.ReadDir(c.path)
	if err != nil && c.path != "" {
		log.Errorf("Error reading textfile collector directory %s: %s", c.path, err)
		error = 1.0
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".prom") {
			continue
		}
		path := filepath.Join(c.path, f.Name())
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("Error opening %s: %v", path, err)
			error = 1.0
			continue
		}
		parsedFamilies, err := (&expfmt.TextParser{}).TextToMetricFamilies(file)
		if err != nil {
			log.Errorf("Error parsing %s: %v", path, err)
			error = 1.0
			continue
		}
		// Only set this once it has been parsed, so that
		// a failure does not appear fresh.
		mtimes[f.Name()] = f.ModTime()
		for _, mf := range parsedFamilies {
			if mf.Help == nil {
				help := fmt.Sprintf("Metric read from %s", path)
				mf.Help = &help
			}
			metricFamilies = append(metricFamilies, mf)
		}
	}

	// Export the mtimes of the successful files.
	if len(mtimes) > 0 {
		mtimeMetricFamily := dto.MetricFamily{
			Name:   proto.String("node_textfile_mtime"),
			Help:   proto.String("Unixtime mtime of textfiles successfully read."),
			Type:   dto.MetricType_GAUGE.Enum(),
			Metric: []*dto.Metric{},
		}

		// Sorting is needed for predictable output comparison in tests.
		filenames := make([]string, 0, len(mtimes))
		for filename := range mtimes {
			filenames = append(filenames, filename)
		}
		sort.Strings(filenames)

		for _, filename := range filenames {
			mtimeMetricFamily.Metric = append(mtimeMetricFamily.Metric,
				&dto.Metric{
					Label: []*dto.LabelPair{
						&dto.LabelPair{
							Name:  proto.String("file"),
							Value: proto.String(filename),
						},
					},
					Gauge: &dto.Gauge{Value: proto.Float64(float64(mtimes[filename].UnixNano()) / 1e9)},
				},
			)
		}
		metricFamilies = append(metricFamilies, &mtimeMetricFamily)
	}
	// Export if there were errors.
	metricFamilies = append(metricFamilies, &dto.MetricFamily{
		Name: proto.String("node_textfile_scrape_error"),
		Help: proto.String("1 if there was an error opening or reading a file, 0 otherwise"),
		Type: dto.MetricType_GAUGE.Enum(),
		Metric: []*dto.Metric{
			&dto.Metric{
				Gauge: &dto.Gauge{Value: &error},
			},
		},
	})

	return metricFamilies
}
