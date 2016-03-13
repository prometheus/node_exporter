package collector

// +build linux freebsd
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

func (m *zfsMetric) ConstMetric(value zfsMetricValue) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, string(m.subsystem), m.name),
			m.name,
			nil,
			nil,
		),
		prometheus.UntypedValue,
		float64(value),
	)
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
	return &zfsCollector{
		zfsMetrics: []zfsMetric{
			{arc, "mru_size", zfsSysctl("kstat.zfs.misc.arcstats.p")},
			{arc, "size", zfsSysctl("kstat.zfs.misc.arcstats.size")},
			{arc, "target_min_size", zfsSysctl("kstat.zfs.misc.arcstats.c_min")},
			{arc, "target_size", zfsSysctl("kstat.zfs.misc.arcstats.c")},
			{arc, "target_max_size", zfsSysctl("kstat.zfs.misc.arcstats.c_max")},
			{arc, "hits", zfsSysctl("kstat.zfs.misc.arcstats.hits")},
			{arc, "misses", zfsSysctl("kstat.zfs.misc.arcstats.misses")},
		},
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {

	err = c.PrepareUpdate()
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

	return err
}
