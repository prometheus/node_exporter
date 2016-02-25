package collector

// +build !nofilesystem

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type zfsSysctl string
type zfsComputedMetric string
type zfsMetricValue int

const zfsErrorValue = zfsMetricValue(-1)

type zfsMetricProvider struct {
	fetchedResults  map[zfsSysctl]zfsMetricValue
	computedResults map[zfsComputedMetric]zfsMetricValue
}

type zfsCollector struct {
	metricProvider  *zfsMetricProvider
	fetchedMetrics  map[zfsSysctl]prometheus.Gauge
	computedMetrics map[zfsComputedMetric]prometheus.Gauge
}

//------------------------------------------------------------------------------
//                                 Collector
//------------------------------------------------------------------------------

func init() {
	Factories["zfs"] = NewZFSCollector
}

const zfsArcSubsystemName = "zfs_arc"

func NewZFSCollector() (Collector, error) {
	if err := zfsInitialize(); err != nil {
		return &zfsCollector{}, err
	}

	fetchedMetrics := make(map[zfsSysctl]prometheus.Gauge)

	insertGauge := func(sysctl, name string) {
		fetchedMetrics[zfsSysctl(sysctl)] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: zfsArcSubsystemName,
				Name:      name,
				Help:      sysctl,
			},
		)
	}

	insertGauge("kstat.zfs.misc.arcstats.p", "mru_size")
	insertGauge("kstat.zfs.misc.arcstats.size", "size")
	insertGauge("kstat.zfs.misc.arcstats.c_min", "target_min_size")
	insertGauge("kstat.zfs.misc.arcstats.c", "target_size")
	insertGauge("kstat.zfs.misc.arcstats.c_max", "target_max_size")
	insertGauge("kstat.zfs.misc.arcstats.hits", "hits")
	insertGauge("kstat.zfs.misc.arcstats.misses", "misses")

	return &zfsCollector{
		metricProvider: NewZfsMetricProvider(),
		fetchedMetrics: fetchedMetrics,
		computedMetrics: map[zfsComputedMetric]prometheus.Gauge{
			zfsComputedMetric("computed.mfu_size"): prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: Namespace,
					Subsystem: zfsArcSubsystemName,
					Name:      "mfu_size",
					Help:      "mfu_size",
				},
			),
		},
	}, nil
}

func (c *zfsCollector) Update(ch chan<- prometheus.Metric) (err error) {

	log.Debug("Preparing metrics update")
	if err := c.metricProvider.PrepareUpdate(); err != nil {
		return err
	}
	defer c.metricProvider.InvalidateCache()

	log.Debugf("Fetching %d metrics", len(c.fetchedMetrics))
	for metric := range c.fetchedMetrics {

		value, err := c.metricProvider.GetFetchedMetric(metric)
		if err != nil {
			return err
		}
		c.fetchedMetrics[metric].Set(float64(value))
		c.fetchedMetrics[metric].Collect(ch)

	}

	log.Debugf("Computing %d metrics", len(c.computedMetrics))
	for metric := range c.computedMetrics {

		value, err := c.metricProvider.GetComputedMetric(metric)
		if err != nil {
			return err
		}
		c.computedMetrics[metric].Set(float64(value))
		c.computedMetrics[metric].Collect(ch)

	}

	return err
}

//------------------------------------------------------------------------------
//                                 Metrics Provider
// Platform-dependend parts implemented in zfs_${platform} files.
//------------------------------------------------------------------------------

func NewZfsMetricProvider() *zfsMetricProvider {
	return &zfsMetricProvider{
		fetchedResults:  make(map[zfsSysctl]zfsMetricValue),
		computedResults: make(map[zfsComputedMetric]zfsMetricValue),
	}

}

func (p *zfsMetricProvider) InvalidateCache() {
	p.fetchedResults = make(map[zfsSysctl]zfsMetricValue)
	p.computedResults = make(map[zfsComputedMetric]zfsMetricValue)
}

func (p *zfsMetricProvider) GetFetchedMetric(s zfsSysctl) (value zfsMetricValue, err error) {

	var ok bool
	value = zfsErrorValue

	value, ok = p.fetchedResults[s]
	if !ok {
		value, err = p.handleFetchedMetricCacheMiss(s)
		if err != nil {
			return value, err
		}
	}

	return value, err
}

func (p *zfsMetricProvider) GetComputedMetric(c zfsComputedMetric) (value zfsMetricValue, err error) {

	if c == "computed.mfu_size" {

		value = zfsErrorValue

		arc_size, err := p.GetFetchedMetric(zfsSysctl("kstat.zfs.misc.arcstats.size"))
		if err != nil {
			return value, err
		}
		target_size, err := p.GetFetchedMetric(zfsSysctl("kstat.zfs.misc.arcstats.c"))
		if err != nil {
			return value, err
		}
		mru_size, err := p.GetFetchedMetric(zfsSysctl("kstat.zfs.misc.arcstats.p"))
		if err != nil {
			return value, err
		}

		// See zfs-stats implementation
		if arc_size > target_size {
			return arc_size - mru_size, nil
		} else {
			return target_size - mru_size, nil
		}

	}

	return zfsErrorValue, fmt.Errorf("no implementation for computed metric '%s'", c)
}
