// +build !notime

package collector

import (
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	systemTime = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "time",
		Help:      "System time in seconds since epoch (1970).",
	})
)

type timeCollector struct {
	config Config
}

func init() {
	Factories["time"] = NewTimeCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// the current system time in seconds since epoch.
func NewTimeCollector(config Config) (Collector, error) {
	c := timeCollector{
		config: config,
	}

	return &c, nil
}

func (c *timeCollector) Update(ch chan<- prometheus.Metric) (err error) {
	now := time.Now()
	glog.V(1).Infof("Set time: %f", now.Unix())
	systemTime.Set(float64(now.Unix()))
	systemTime.Collect(ch)
	return err
}
