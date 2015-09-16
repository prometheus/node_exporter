// +build !norunit

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
	"github.com/soundcloud/go-runit/runit"
)

type runitCollector struct {
	state, stateDesired, stateNormal, stateTimestamp *prometheus.GaugeVec
}

func init() {
	Factories["runit"] = NewRunitCollector
}

func NewRunitCollector() (Collector, error) {
	var (
		subsystem   = "service"
		constLabels = prometheus.Labels{"supervisor": "runit"}
		labelNames  = []string{"service"}
	)

	return &runitCollector{
		state: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "state",
				Help:        "State of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateDesired: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "desired_state",
				Help:        "Desired state of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateNormal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "normal_state",
				Help:        "Normal state of runit service.",
				ConstLabels: constLabels,
			},
			labelNames,
		),
		stateTimestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   Namespace,
				Subsystem:   subsystem,
				Name:        "state_last_change_timestamp_seconds",
				Help:        "Unix timestamp of the last runit service state change.",
				ConstLabels: constLabels,
			},
			labelNames,
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
			log.Debugf("Couldn't get status for %s: %s, skipping...", service.Name, err)
			continue
		}

		log.Debugf("%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		c.state.WithLabelValues(service.Name).Set(float64(status.State))
		c.stateDesired.WithLabelValues(service.Name).Set(float64(status.Want))
		c.stateTimestamp.WithLabelValues(service.Name).Set(float64(status.Timestamp.Unix()))
		if status.NormallyUp {
			c.stateNormal.WithLabelValues(service.Name).Set(1)
		} else {
			c.stateNormal.WithLabelValues(service.Name).Set(0)
		}
	}
	c.state.Collect(ch)
	c.stateDesired.Collect(ch)
	c.stateNormal.Collect(ch)
	c.stateTimestamp.Collect(ch)

	return nil
}
