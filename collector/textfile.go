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

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	textFileDirectory = kingpin.Flag("collector.textfile.directory", "Directory to read text files with metrics from.").Default("").String()

	mtimeDesc = prometheus.NewDesc(
		"node_textfile_mtime",
		"Unixtime mtime of textfiles successfully read.",
		[]string{"file"},
		nil,
	)
	errorDesc = prometheus.NewDesc(
		"node_textfile_scrape_error",
		"1 if there was an error opening or reading a file, 0 otherwise",
		nil,
		nil,
	)
)

type textFileCollector struct {
	path string
	// Only set for testing to get predictable output.
	mtime *float64
}

func init() {
	registerCollector("textfile", defaultEnabled, NewTextFileCollector)
}

// NewTextFileCollector returns a new Collector exposing metrics read from files
// in the given textfile directory.
func NewTextFileCollector() (Collector, error) {
	c := &textFileCollector{
		path: *textFileDirectory,
	}
	return c, nil
}

func convertMetricFamilies(metricFamilies []*dto.MetricFamily, ch chan<- prometheus.Metric) {
	var valType prometheus.ValueType
	var val float64

	for _, mf := range metricFamilies {
		labelNames := map[string]struct{}{}
		for _, metric := range mf.Metric {
			labelPairs := metric.GetLabel()
			for _, label := range labelPairs {
				if _, ok := labelNames[label.GetName()]; !ok {
					labelNames[label.GetName()] = struct{}{}
				}
			}
		}

		for _, metric := range mf.Metric {
			labelPairs := metric.GetLabel()
			var labels []string
			var labelVals []string
			for _, label := range labelPairs {
				labels = append(labels, label.GetName())
				labelVals = append(labelVals, label.GetValue())
			}

			for k := range labelNames {
				present := false
				for _, label := range labels {
					if k == label {
						present = true
						break
					}
				}
				if present == false {
					labels = append(labels, k)
					labelVals = append(labelVals, "")
				}
			}

			metricType := mf.GetType()
			switch metricType {
			case dto.MetricType_COUNTER:
				if metric.Counter != nil {
					valType = prometheus.CounterValue
					val = metric.Counter.GetValue()
				}
			case dto.MetricType_GAUGE:
				if metric.Gauge != nil {
					valType = prometheus.GaugeValue
					val = metric.Gauge.GetValue()
				}
			case dto.MetricType_UNTYPED:
				if metric.Untyped != nil {
					valType = prometheus.UntypedValue
					val = metric.Untyped.GetValue()
				}
			case dto.MetricType_SUMMARY:
				if metric.Summary != nil {
					quantiles := map[float64]float64{}
					for _, q := range metric.Summary.Quantile {
						quantiles[q.GetQuantile()] = q.GetValue()
					}
					ch <- prometheus.MustNewConstSummary(
						prometheus.NewDesc(
							*mf.Name,
							mf.GetHelp(),
							labels, nil,
						),
						metric.Summary.GetSampleCount(),
						metric.Summary.GetSampleSum(),
						quantiles, labelVals...,
					)
				}
			case dto.MetricType_HISTOGRAM:
				if metric.Histogram != nil {
					buckets := map[float64]uint64{}
					for _, b := range metric.Histogram.Bucket {
						buckets[b.GetUpperBound()] = b.GetCumulativeCount()
					}
					ch <- prometheus.MustNewConstHistogram(
						prometheus.NewDesc(
							*mf.Name,
							mf.GetHelp(),
							labels, nil,
						),
						metric.Histogram.GetSampleCount(),
						metric.Histogram.GetSampleSum(),
						buckets, labelVals...,
					)
				}

			}
			if metricType == dto.MetricType_GAUGE || metricType == dto.MetricType_COUNTER || metricType == dto.MetricType_UNTYPED {
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						*mf.Name,
						mf.GetHelp(),
						labels, nil,
					),
					valType, val, labelVals...,
				)
			}
		}
	}
}

func exportMTimes(mtimes map[string]time.Time, ch chan<- prometheus.Metric) {
	// Export the mtimes of the successful files.
	if len(mtimes) > 0 {
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
						{
							Name:  proto.String("file"),
							Value: proto.String(filename),
						},
					},
					Gauge: &dto.Gauge{Value: proto.Float64(float64(mtimes[filename].UnixNano()) / 1e9)},
				},
			)
			labels = append(labels, "file")
			labelVals = append(labelVals, filename)
		}
		for _, metric := range mtimeMetricFamily.Metric {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					*mtimeMetricFamily.Name,
					mtimeMetricFamily.GetHelp(),
					labels, nil,
				),
				prometheus.GaugeValue, metric.Gauge.GetValue(), labelVals...,
			)
		}
	}
}

// Update implements the Collector interface.
func (c *textFileCollector) Update(ch chan<- prometheus.Metric) error {
	var metricFamilies []*dto.MetricFamily
	error := 0.0
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
		var parser expfmt.TextParser
		parsedFamilies, err := parser.TextToMetricFamilies(file)
		file.Close()
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

	convertMetricFamilies(metricFamilies, ch)

	exportMTimes(mtimes, ch)

	// Export if there were errors.
	var labels []string
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"node_textfile_scrape_error",
			"1 if there was an error opening or reading a file, 0 otherwise",
			labels, nil,
		),
		prometheus.GaugeValue, error,
	)
	return nil
}

func (c *textFileCollector) parseTextFiles() []*dto.MetricFamily {
	error := 0.0
	var metricFamilies []*dto.MetricFamily
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
		var parser expfmt.TextParser
		parsedFamilies, err := parser.TextToMetricFamilies(file)
		file.Close()
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
			ch <- prometheus.MustNewConstMetric(mtimeDesc, prometheus.GaugeValue, mtime, filename)
		}
	}

	// Export if there were errors.
	ch <- prometheus.MustNewConstMetric(errorDesc, prometheus.GaugeValue, error)
	return nil
}
