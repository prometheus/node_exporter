// +build runit

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soundcloud/go-runit/runit"
)

var (
	runitLabelNames = []string{"service"}

	runitState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "service_state",
			Help:      "node_exporter: state of runit service.",
		},
		runitLabelNames,
	)
	runitStateDesired = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "service_desired_state",
			Help:      "node_exporter: desired state of runit service.",
		},
		runitLabelNames,
	)
	runitStateNormal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "service_normal_state",
			Help:      "node_exporter: normal state of runit service.",
		},
		runitLabelNames,
	)
)

type runitCollector struct {
	config Config
}

func init() {
	Factories["runit"] = NewRunitCollector
}

func NewRunitCollector(config Config) (Collector, error) {
	c := runitCollector{
		config: config,
	}

	return &c, nil
}

func (c *runitCollector) Update(ch chan<- prometheus.Metric) (err error) {
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
		runitState.WithLabelValues(service.Name).Set(float64(status.State))
		runitStateDesired.WithLabelValues(service.Name).Set(float64(status.Want))
		if status.NormallyUp {
			runitStateNormal.WithLabelValues(service.Name).Set(1)
		} else {
			runitStateNormal.WithLabelValues(service.Name).Set(1)
		}
	}
	runitState.Collect(ch)
	runitStateDesired.Collect(ch)
	runitStateNormal.Collect(ch)
	return err
}
