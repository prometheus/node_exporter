package collector

// +build !nozfs

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type zfsMetricValue int

const zfsErrorValue = zfsMetricValue(-1)

var zfsNotAvailableError = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string
type zfsSubsystemName string

const (
	arc = zfsSubsystemName("zfs_arc")
)

// Metrics

type zfsMetric struct {
	subsystem zfsSubsystemName // The Prometheus subsystem name.
	name      string           // The Prometheus name of the metric.
	sysctl    zfsSysctl        // The sysctl of the ZFS metric.
}

func NewZFSMetric(subsystem zfsSubsystemName, sysctl, name string) zfsMetric {
	return zfsMetric{
		sysctl:    zfsSysctl(sysctl),
		subsystem: subsystem,
		name:      name,
	}
}
func (m *zfsMetric) BuildFQName() string {
	return prometheus.BuildFQName(Namespace, string(m.subsystem), m.name)
}

func (m *zfsMetric) HelpString() string {
	return m.name
}

// Collector

func init() {
	Factories["zfs"] = NewZFSCollector
}

type zfsCollector struct {
	zfsMetrics     []zfsMetric
	metricProvider zfsMetricProvider
}

func NewZFSCollector() (Collector, error) {
	err := zfsInitialize()
	switch {
	case err == zfsNotAvailableError:
		log.Debug(err)
		break
		return &zfsCollector{}, err
	}

	return &zfsCollector{
		zfsMetrics: []zfsMetric{
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.p", "mru_size"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.size", "size"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.c_min", "target_min_size"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.c", "target_size"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.c_max", "target_max_size"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.hits", "hits"),
			NewZFSMetric(arc, "kstat.zfs.misc.arcstats.misses", "misses"),
		},
		metricProvider: NewZFSMetricProvider(),
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {

	log.Debug("Preparing metrics update")
	err = c.metricProvider.PrepareUpdate()
	switch {
	case err == zfsNotAvailableError:
		log.Debug(err)
		return nil
		return err
	}
	defer c.metricProvider.InvalidateCache()

	log.Debugf("Fetching %d metrics", len(c.zfsMetrics))
	for _, metric := range c.zfsMetrics {

		value, err := c.metricProvider.Value(metric.sysctl)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				metric.BuildFQName(),
				metric.HelpString(),
				nil,
				nil,
			),
			prometheus.UntypedValue,
			float64(value),
		)

	}

	return err
}

// Metrics Provider
// Platform-dependend parts implemented in zfs_${os} files.

type zfsMetricProvider struct {
	values map[zfsSysctl]zfsMetricValue
}

func NewZFSMetricProvider() zfsMetricProvider {
	return zfsMetricProvider{
		values: make(map[zfsSysctl]zfsMetricValue),
	}

}

func (p *zfsMetricProvider) InvalidateCache() {
	p.values = make(map[zfsSysctl]zfsMetricValue)
}

func (p *zfsMetricProvider) Value(s zfsSysctl) (value zfsMetricValue, err error) {

	var ok bool
	value = zfsErrorValue

	value, ok = p.values[s]
	if !ok {
		value, err = p.handleMiss(s)
		if err != nil {
			return value, err
		}
	}

	return value, err
}
