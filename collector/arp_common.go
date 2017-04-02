// +build !noarp

package collector

import "github.com/prometheus/client_golang/prometheus"

type arpCollector struct {
	count *prometheus.Desc
}

func init() {
	Factories["arp"] = NewArpCollector
}

// NewArpCollector returns a new Collector exposing ARP stats.
func NewArpCollector() (Collector, error) {
	return &arpCollector{
		count: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "arp",
				"count"), "ARP entries by device",
			[]string{"device"}, nil,
		),
	}, nil
}
