// +build !nouname

package collector

import (
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
)

var unameDesc = prometheus.NewDesc(
	prometheus.BuildFQName(Namespace, "uname", "info"),
	"Labeled system information as provided by the uname system call.",
	[]string{
		"sysname",
		"release",
		"version",
		"machine",
		"nodename",
		"domainname",
	},
	nil,
)

type unameCollector struct{}

func init() {
	Factories["uname"] = newUnameCollector
}

// NewUnameCollector returns new unameCollector.
func newUnameCollector() (Collector, error) {
	return &unameCollector{}, nil
}

func intArrayToString(array [65]int8) string {
	var str string
	for _, a := range array {
		if a == 0 {
			break
		}
		str += string(a)
	}
	return str
}

func (c unameCollector) Update(ch chan<- prometheus.Metric) error {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return err
	}

	labelValues := []string{
		intArrayToString(uname.Sysname),
		intArrayToString(uname.Release),
		intArrayToString(uname.Version),
		intArrayToString(uname.Machine),
		intArrayToString(uname.Nodename),
		intArrayToString(uname.Domainname),
	}
	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1, labelValues...)
	return nil
}
