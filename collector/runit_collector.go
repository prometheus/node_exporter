// +build runit

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/go-runit/runit"
)

const (
	runitSubsystem = "runit"
)

type runitCollector struct {
	config Config

	state, stateDesired, stateNormal *prometheus.GaugeVec
}

func init() {
	Factories["runit"] = NewRunitCollector
}

func NewRunitCollector(config Config) (Collector, error) {
	var labels = []string{"service"}

	return &runitCollector{
		config: config,
		state: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: runitSubsystem,
				Name:      "state",
				Help:      "state of runit service.",
			},
			labels,
		),
		stateDesired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: runitSubsystem,
				Name:      "desired_state",
				Help:      "desired state of runit service.",
			},
			labels,
		),
		stateNormal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: runitSubsystem,
				Name:      "normal_state",
				Help:      "normal state of runit service.",
			},
			labels,
		),
	}, nil
}

func (c *runitCollector) Update(ch chan<- prometheus.Metric) error {
	services, err := runit.GetServices("/etc/service")
	if err != nil {
		return err
	}

	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			glog.V(1).Infof("Couldn't get status for %s: %s, skipping...", service.Name, err)
			continue
		}

		glog.V(1).Infof("%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		c.state.WithLabelValues(service.Name).Set(float64(status.State))
		c.stateDesired.WithLabelValues(service.Name).Set(float64(status.Want))
		if status.NormallyUp {
			c.stateNormal.WithLabelValues(service.Name).Set(1)
		} else {
			c.stateNormal.WithLabelValues(service.Name).Set(0)
		}
	}
	c.state.Collect(ch)
	c.stateDesired.Collect(ch)
	c.stateNormal.Collect(ch)

	return nil
}
