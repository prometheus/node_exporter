// Copyright 2016 The Prometheus Authors
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
	"path"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

var (
	uptrackUpgradePlanPath = flag.String(
		"collector.uptrack.upgrade_plan_path",
		uptrackDefaultUpgradePlanPath(),
		"Path where Uptrack stores its upgrade_plan YAML file.")
)

type uptrackCollector struct {
	timeDesc      *prometheus.Desc
	planCountDesc *prometheus.Desc
}

func init() {
	Factories["uptrack"] = NewUptrackCollector
}

// The location of the upgrade_plan file is based on the currently running kernel.
func uptrackDefaultUpgradePlanPath() string {
	var uname syscall.Utsname
	syscall.Uname(&uname)
	return path.Join(
		"/var/cache/uptrack",
		unameToString(uname.Sysname), unameToString(uname.Machine),
		unameToString(uname.Release), unameToString(uname.Version),
		"upgrade_plan")
}

// Time type that has a custom unmarshalling function for Uptrack's almost-RFC3339 timestamps.
type uptrackTime struct {
	time.Time
}

func (t *uptrackTime) UnmarshalText(data []byte) error {
	var err error
	t.Time, err = time.Parse("2006-01-02 15:04:05.999999", string(data))
	return err
}

type uptrackUpgradePlan struct {
	Plan []struct{}  `yaml:"Plan"`
	Time uptrackTime `yaml:"Time"`
}

func NewUptrackCollector() (Collector, error) {
	return &uptrackCollector{
		timeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "uptrack", "upgrade_plan_time_seconds"),
			"Time at which Uptrack created the upgrade plan.",
			nil, nil),
		planCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "uptrack", "upgrade_plan_plan_count"),
			"Number of patches that still need to be applied to the running kernel.",
			nil, nil),
	}, nil
}

func (c *uptrackCollector) Update(ch chan<- prometheus.Metric) (err error) {
	contents, err := ioutil.ReadFile(*uptrackUpgradePlanPath)
	if err != nil {
		return err
	}

	var plan uptrackUpgradePlan
	if err := yaml.Unmarshal(contents, &plan); err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.timeDesc,
		prometheus.GaugeValue,
		float64(plan.Time.UnixNano())/1e9)
	ch <- prometheus.MustNewConstMetric(
		c.planCountDesc,
		prometheus.GaugeValue,
		float64(len(plan.Plan)))
	return nil
}
