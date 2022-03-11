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
        "bytes"
        "fmt"
        "io/ioutil"
        "os"
        "strconv"
        "github.com/go-kit/log"
        "github.com/prometheus/client_golang/prometheus"
)

const (
        fipsStatSubsystem = "fips"
)

type fipsStatCollector struct {
        logger log.Logger
}

func init() {
        registerCollector(fipsStatSubsystem, defaultDisabled, NewFipsStatCollector)
}

// NewFipsStatCollector returns a new Collector exposing fips status.
func NewFipsStatCollector(logger log.Logger) (Collector, error) {
        return &fipsStatCollector{logger}, nil
}

func (c *fipsStatCollector) Update(ch chan<- prometheus.Metric) error {
        fipsStat, err := parseFipsStats(procFilePath("sys/crypto/fips_enabled"))
        if err != nil {
                return fmt.Errorf("couldn't get fips_enabled: %w", err)
        }
        for name, value := range fipsStat {
                v, err := strconv.ParseFloat(value, 64)
                if err != nil {
                        return fmt.Errorf("invalid value %s in fips_enabled: %w", value, err)
                }
                ch <- prometheus.MustNewConstMetric(
                        prometheus.NewDesc(
                                prometheus.BuildFQName(namespace, fipsStatSubsystem, name),
                                fmt.Sprintf("FIPS status (0-disabled/1-enabled) from /proc/sys/crypto/fips_enabled."),
                                nil, nil,
                        ),
                        prometheus.GaugeValue, v,
                )
        }
        return nil
}

func parseFipsStats(filename string) (map[string]string, error) {
        file, err := os.Open(filename)
        if err != nil {
                return nil, err
        }
        defer file.Close()

        content, err := ioutil.ReadAll(file)
        if err != nil {
                return nil, err
        }
        parts := bytes.Split(bytes.TrimSpace(content), []byte("\u0009"))
        if len(parts) < 1 {
                return nil, fmt.Errorf("unexpected number of file stats in %q", filename)
        }

        var fipsStat = map[string]string{}
        // The fips_enabled proc is only 1 line with 1 value.
        fipsStat["status"] = string(parts[0])

        return fipsStat, nil
}
