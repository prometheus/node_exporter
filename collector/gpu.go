// +build darwin linux openbsd
// +build !nogpu

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type gpuCollector struct {
	info        gpuCache
	total       *prometheus.Desc
	used        *prometheus.Desc
	free        *prometheus.Desc
	utilization *prometheus.Desc
	temp        *prometheus.Desc
}

func init() {
	registerCollector("gpu", defaultEnabled, NewGpuCollector)
}

// NewGpuCollector data come from nvidia-smi -q
func NewGpuCollector() (Collector, error) {
	info := gpuCache{}

	return &gpuCollector{
		info: info,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "total"),
			"Framebuffer memory total (in MiB).",
			gpuLabelNames, nil,
		),
		used: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "used"),
			"Framebuffer memory used (in MiB).",
			gpuLabelNames, nil,
		),
		free: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "free"),
			"Framebuffer memory free (in MiB).",
			gpuLabelNames, nil,
		),
		utilization: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "utilization"),
			"GPU utilization (in %).",
			gpuLabelNames, nil,
		),
		temp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, gpuCollectorSubsystem, "temp"),
			"GPU temperature (in C).",
			gpuLabelNames, nil,
		),
	}, nil
}

func (this *gpuCollector) Update(ch chan<- prometheus.Metric) error {
	if err := this.updateStat(ch); err != nil {
		return err
	}
	return nil
}

func (this *gpuCollector) updateStat(ch chan<- prometheus.Metric) error {
	stats, err := this.info.Stat()
	if err != nil {
		return err
	}

	for _, gpuStat := range stats {
		ch <- prometheus.MustNewConstMetric(this.total, prometheus.CounterValue, gpuStat.TotalMem, gpuStat.Host, gpuStat.Count, gpuStat.UUID)
		ch <- prometheus.MustNewConstMetric(this.used, prometheus.CounterValue, gpuStat.UsedMem, gpuStat.Host, gpuStat.Count, gpuStat.UUID)
		ch <- prometheus.MustNewConstMetric(this.free, prometheus.CounterValue, gpuStat.FreeMem, gpuStat.Host, gpuStat.Count, gpuStat.UUID)
		ch <- prometheus.MustNewConstMetric(this.utilization, prometheus.CounterValue, gpuStat.Utilization, gpuStat.Host, gpuStat.Count, gpuStat.UUID)
		ch <- prometheus.MustNewConstMetric(this.temp, prometheus.CounterValue, gpuStat.Temp, gpuStat.Host, gpuStat.Count, gpuStat.UUID)
	}
	return nil
}
