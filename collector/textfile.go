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

// +build !notextfile

package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	textFileDirectory = kingpin.Flag("collector.textfile.directory", "Directory to read text files with metrics from.").Default("").String()
	mtimeDesc         = prometheus.NewDesc(
		"node_textfile_mtime_seconds",
		"Unixtime mtime of textfiles successfully read.",
		[]string{"file"},
		nil,
	)
)

type textFileCollector struct {
	path string
	// Only set for testing to get predictable output.
	mtime  *float64
	logger log.Logger
}

func init() {
	registerCollector("textfile", defaultEnabled, NewTextFileCollector)
}

// NewTextFileCollector returns a new Collector exposing metrics read from files
// in the given textfile directory.
func NewTextFileCollector(logger log.Logger) (Collector, error) {
	c := &textFileCollector{
		path:   *textFileDirectory,
		logger: logger,
	}
	return c, nil
}

func convertMetricFamily(metricFamily *dto.MetricFamily, ch chan<- prometheus.Metric, logger log.Logger) {
	var valType prometheus.ValueType
	var val float64

	allLabelNames := map[string]struct{}{}
	for _, metric := range metricFamily.Metric {
		labels := metric.GetLabel()
		for _, label := range labels {
			if _, ok := allLabelNames[label.GetName()]; !ok {
				allLabelNames[label.GetName()] = struct{}{}
			}
		}
	}

	for _, metric := range metricFamily.Metric {
		if metric.TimestampMs != nil {
			level.Warn(logger).Log("msg", "Ignoring unsupported custom timestamp on textfile collector metric", "metric", metric)
		}

		labels := metric.GetLabel()
		var names []string
		var values []string
		for _, label := range labels {
			names = append(names, label.GetName())
			values = append(values, label.GetValue())
		}

		for k := range allLabelNames {
			present := false
			for _, name := range names {
				if k == name {
					present = true
					break
				}
			}
			if !present {
				names = append(names, k)
				values = append(values, "")
			}
		}

		metricType := metricFamily.GetType()
		switch metricType {
		case dto.MetricType_COUNTER:
			valType = prometheus.CounterValue
			val = metric.Counter.GetValue()

		case dto.MetricType_GAUGE:
			valType = prometheus.GaugeValue
			val = metric.Gauge.GetValue()

		case dto.MetricType_UNTYPED:
			valType = prometheus.UntypedValue
			val = metric.Untyped.GetValue()

		case dto.MetricType_SUMMARY:
			quantiles := map[float64]float64{}
			for _, q := range metric.Summary.Quantile {
				quantiles[q.GetQuantile()] = q.GetValue()
			}
			ch <- prometheus.MustNewConstSummary(
				prometheus.NewDesc(
					*metricFamily.Name,
					metricFamily.GetHelp(),
					names, nil,
				),
				metric.Summary.GetSampleCount(),
				metric.Summary.GetSampleSum(),
				quantiles, values...,
			)
		case dto.MetricType_HISTOGRAM:
			buckets := map[float64]uint64{}
			for _, b := range metric.Histogram.Bucket {
				buckets[b.GetUpperBound()] = b.GetCumulativeCount()
			}
			ch <- prometheus.MustNewConstHistogram(
				prometheus.NewDesc(
					*metricFamily.Name,
					metricFamily.GetHelp(),
					names, nil,
				),
				metric.Histogram.GetSampleCount(),
				metric.Histogram.GetSampleSum(),
				buckets, values...,
			)
		default:
			panic("unknown metric type")
		}
		if metricType == dto.MetricType_GAUGE || metricType == dto.MetricType_COUNTER || metricType == dto.MetricType_UNTYPED {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					*metricFamily.Name,
					metricFamily.GetHelp(),
					names, nil,
				),
				valType, val, values...,
			)
		}
	}
}

func (c *textFileCollector) exportMTimes(mtimes map[string]time.Time, ch chan<- prometheus.Metric) {
	if len(mtimes) == 0 {
		return
	}

	// Export the mtimes of the successful files.
	// Sorting is needed for predictable output comparison in tests.
	filenames := make([]string, 0, len(mtimes))
	for filename := range mtimes {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		mtime := float64(mtimes[filename].UnixNano() / 1e9)
		if c.mtime != nil {
			mtime = *c.mtime
		}
		ch <- prometheus.MustNewConstMetric(mtimeDesc, prometheus.GaugeValue, mtime, filename)
	}
}

// Update implements the Collector interface.
func (c *textFileCollector) Update(ch chan<- prometheus.Metric) error {
	// Iterate over files and accumulate their metrics, but also track any
	// parsing errors so an error metric can be reported.
	var errored bool
	files, err := ioutil.ReadDir(c.path)
	if err != nil && c.path != "" {
		errored = true
		level.Error(c.logger).Log("msg", "failed to read textfile collector directory", "path", c.path, "err", err)
	}

	mtimes := make(map[string]time.Time, len(files))
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".prom") {
			continue
		}

		mtime, err := c.processFile(f.Name(), ch)
		if err != nil {
			errored = true
			level.Error(c.logger).Log("msg", "failed to collect textfile data", "file", f.Name(), "err", err)
			continue
		}

		mtimes[f.Name()] = *mtime
	}

	c.exportMTimes(mtimes, ch)

	// Export if there were errors.
	var errVal float64
	if errored {
		errVal = 1.0
	}

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"node_textfile_scrape_error",
			"1 if there was an error opening or reading a file, 0 otherwise",
			nil, nil,
		),
		prometheus.GaugeValue, errVal,
	)

	return nil
}

// processFile processes a single file, returning its modification time on success.
func (c *textFileCollector) processFile(name string, ch chan<- prometheus.Metric) (*time.Time, error) {
	path := filepath.Join(c.path, name)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open textfile data file %q: %w", path, err)
	}
	defer f.Close()

	var parser expfmt.TextParser
	families, err := parser.TextToMetricFamilies(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse textfile data from %q: %w", path, err)
	}

	if hasTimestamps(families) {
		return nil, fmt.Errorf("textfile %q contains unsupported client-side timestamps, skipping entire file", path)
	}

	for _, mf := range families {
		if mf.Help == nil {
			help := fmt.Sprintf("Metric read from %s", path)
			mf.Help = &help
		}
	}

	for _, mf := range families {
		convertMetricFamily(mf, ch, c.logger)
	}

	// Only stat the file once it has been parsed and validated, so that
	// a failure does not appear fresh.
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat %q: %w", path, err)
	}

	t := stat.ModTime()
	return &t, nil
}

// hasTimestamps returns true when metrics contain unsupported timestamps.
func hasTimestamps(parsedFamilies map[string]*dto.MetricFamily) bool {
	for _, mf := range parsedFamilies {
		for _, m := range mf.Metric {
			if m.TimestampMs != nil {
				return true
			}
		}
	}
	return false
}
