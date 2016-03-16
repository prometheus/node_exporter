package collector

// +build linux freebsd
// +build !nozfs

import (
	"errors"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type zfsMetricValue int

const zfsErrorValue = zfsMetricValue(-1)

var zfsNotAvailableError = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string
type zfsSubsystemName string

const (
	arc            = zfsSubsystemName("zfsArc")
	zpoolSubsystem = zfsSubsystemName("zfsPool")
)

// Metrics

type zfsMetric struct {
	subsystem zfsSubsystemName // The Prometheus subsystem name.
	name      string           // The Prometheus name of the metric.
	sysctl    zfsSysctl        // The sysctl of the ZFS metric.
}

type datasetMetric struct {
	subsystem zfsSubsystemName
	name      string
}

// Collector

func init() {
	Factories["zfs"] = NewZFSCollector
}

type zfsCollector struct {
	zfsMetrics []zfsMetric
}

func NewZFSCollector() (Collector, error) {
	return &zfsCollector{}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {

	err = c.zfsAvailable()
	switch {
	case err == zfsNotAvailableError:
		log.Debug(err)
		return nil
	case err != nil:
		return err
	}

	// Arcstats
	err = c.updateArcstats(ch)
	if err != nil {
		return err
	}

	// Pool stats
	err = c.updatePoolStats(ch)
	if err != nil {
		return err
	}

	return err
}

func (s zfsSysctl) metricName() string {
	parts := strings.Split(string(s), ".")
	return parts[len(parts)-1]
}

func (c *zfsCollector) ConstSysctlMetric(subsystem zfsSubsystemName, sysctl zfsSysctl, value zfsMetricValue) prometheus.Metric {

	metricName := sysctl.metricName()

	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, string(subsystem), metricName),
			string(sysctl),
			nil,
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
	)
}

func (c *zfsCollector) ConstZpoolMetric(pool, name string, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, string(zpoolSubsystem), name),
			name,
			[]string{"pool"},
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
		pool,
	)
}
