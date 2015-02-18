// +build !nontp

package collector

import (
	"flag"
	"fmt"
	"time"

	"github.com/beevik/ntp"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ntpServer = flag.String("collector.ntp.server", "", "NTP server to use for ntp collector.")
)

type ntpCollector struct {
	drift prometheus.Gauge
}

func init() {
	Factories["ntp"] = NewNtpCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// the offset between ntp and the current system time.
func NewNtpCollector(config Config) (Collector, error) {
	if *ntpServer == "" {
		return nil, fmt.Errorf("No NTP server specifies, see --ntpServer")
	}

	return &ntpCollector{
		drift: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "ntp_drift_seconds",
			Help:      "Time between system time and ntp time.",
		}),
	}, nil
}

func (c *ntpCollector) Update(ch chan<- prometheus.Metric) (err error) {
	t, err := ntp.Time(*ntpServer)
	if err != nil {
		return fmt.Errorf("Couldn't get ntp drift: %s", err)
	}
	drift := t.Sub(time.Now())
	glog.V(1).Infof("Set ntp_drift_seconds: %f", drift.Seconds())
	c.drift.Set(drift.Seconds())
	c.drift.Collect(ch)
	return err
}
