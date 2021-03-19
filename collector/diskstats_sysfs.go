// +build sysfsdiskstats,!nodiskstats

package collector

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	DefaultSectorSize  = 512
	DefaultSysPath     = "/sys"
	sysBlockStatFormat = "%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"
)

// Conversion is enumeration of unit Conversions for prometheus metrics
type Conversion int

// Possible unit Conversions
const (
	_ Conversion = iota
	Ms
	Sectors
	None
)

type metricDesc struct {
	Name        string
	Description string
	ValueType   prometheus.ValueType
	Conv        Conversion
}

type metric struct {
	Desc        *metricDesc
	Value       float64
	LabelNames  []string
	LabelValues []string
}

var (
	ignoredSysDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$").String()
	blockClassPath    = filepath.Join("class", "block")
	// https://www.kernel.org/doc/html/latest/block/stat.html
	// https://www.kernel.org/doc/Documentation/admin-guide/iostats.rst
	// Name            units         Description
	// ----            -----         -----------
	// read I/Os       requests      number of read I/Os processed
	// read merges     requests      number of read I/Os merged with in-queue I/O
	// read sectors    sectors       number of sectors read
	// read ticks      milliseconds  total wait time for read requests
	// write I/Os      requests      number of write I/Os processed
	// write merges    requests      number of write I/Os merged with in-queue I/O
	// write sectors   sectors       number of sectors written
	// write ticks     milliseconds  total wait time for write requests
	// in_flight       requests      number of I/Os currently in flight
	// io_ticks        milliseconds  total time this block device has been active
	// time_in_queue   milliseconds  total wait time for all requests
	// discard I/Os    requests      number of discard I/Os processed
	// discard merges  requests      number of discard I/Os merged with in-queue I/O
	// discard sectors sectors       number of sectors discarded
	// discard ticks   milliseconds  total wait time for discard requests
	// flush requests  requests      number of flush requests completed successfully.
	// flush time      milliseconds  total time spent by all flush requests.
	metricDescs = []metricDesc{
		{
			// Field 1
			Name:        "reads_completed_total",
			Description: "Total number of reads completed successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			//Field 2
			Name: "reads_merged_total",
			Description: "The number of adjacent reads merged with in-queue I/O. " +
				"Thus two 4K reads may become one 8K read. " +
				"This field lets you know how often this was done.",
			ValueType: prometheus.CounterValue,
			Conv:      None,
		},
		{
			// Field 3,
			Name:        "read_bytes_total",
			Description: "Total number of bytes read successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        Sectors,
		},
		{
			// Field 4
			Name:        "read_time_seconds_total",
			Description: "Total number of seconds spent by all reads",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},

		{
			// Field 5
			Name:        "writes_completed_total",
			Description: "Total number of writes completed successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			// Field 6
			Name:        "writes_merged_total",
			Description: "Total number of writes merged with in-queue I/O.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			// Field 7
			Name:        "written_bytes_total",
			Description: "Total number of bytes written successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        Sectors,
		},
		{
			// Field 8
			Name:        "write_time_seconds_total",
			Description: "Total number of seconds spent by all writes.",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},

		{
			// Field 9
			Name:        "io_now",
			Description: "The number of I/Os currently in progress.",
			ValueType:   prometheus.GaugeValue,
			Conv:        None,
		},
		{
			// Field 10
			Name:        "io_time_seconds_total",
			Description: "Total number of seconds when io_now was > 0.",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},
		{
			// Field 11
			Name:        "io_time_weighted_seconds_total",
			Description: "Total wait time in seconds for all completed and ongoing I/O requests.",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},
		// Kernel >5.5
		{
			// Field 12
			Name:        "discards_completed_total",
			Description: "Total number of discard requests completed successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			// Field 13
			Name:        "discards_merged_total",
			Description: "Total number of discards merges with in-queue I/O.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			// Field 14
			Name:        "discarded_sectors_total",
			Description: "Total number of sectors discarded successfully.",
			ValueType:   prometheus.CounterValue,
			Conv:        None,
		},
		{
			// Field 15
			Name:        "discard_time_seconds_total",
			Description: "Total number of seconds spent by all discards.",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},
		// Not tracked for partitions
		{
			// Field 16
			Name: "flush_requests_total",
			Description: "Total number of flush requests completed successfully. " +
				"Block layer combines flush requests and executes at most one at a time. " +
				"This counts flush requests executed by disk.",
			ValueType: prometheus.CounterValue,
			Conv:      None,
		},
		{
			// Field 17
			Name:        "flush_requests_time_seconds_total",
			Description: "Total number of seconds spent by all flush requests.",
			ValueType:   prometheus.CounterValue,
			Conv:        Ms,
		},
	}
)

type dmDevice struct {
	Name       string  // /sys/class/block/Name
	Label      string  // /sys/class/block/<Name>/dm/name
	UUID       string  // /sys/class/block/<Name>/dm/uuid
	SectorSize float64 // /sys/block/<Name>/queue/hw_sector_size
}

type SysFS struct {
	path                  *string
	ignoredDevicesPattern *regexp.Regexp
}

type sysfsDiskCollector struct {
	descs  []metricDesc
	logger log.Logger
	fs     SysFS
}

func (m *metric) mustNewConstMetric() prometheus.Metric {
	promMetric := prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, m.Desc.Name),
			m.Desc.Description,
			m.LabelNames,
			nil,
		),
		m.Desc.ValueType,
		m.Value,
		m.LabelValues...,
	)

	return promMetric
}

// NewFS returns a new sysfs using sysPath from collector/paths.go
// It will error if sysfs can't be read
func NewSysFS() (*SysFS, error) {
	info, err := os.Stat(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %q: %w", *sysPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("Sysfs path %q is not a directory", *sysPath)
	}

	return &SysFS{
		path:                  sysPath,
		ignoredDevicesPattern: regexp.MustCompile(*ignoredSysDevices),
	}, nil
}

// Path appends the given path elements to the filesystem path, adding separators
// as necessary.
func (fs SysFS) Path(p ...string) string {
	return filepath.Join(append([]string{*fs.path}, p...)...)
}

// ReadFile is using syscall.Read directly to avoid sysfs limitations
// Maximum read size is 256B
func (fs SysFS) ReadFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// https://github.com/prometheus/node_exporter/pull/728/files
	const sysFileBufferSize = 256
	b := make([]byte, sysFileBufferSize)
	n, err := syscall.Read(int(f.Fd()), b)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", file, err)
	}

	return string(bytes.TrimSpace(b[:n])), nil
}

// ListSysBlockDevices lists the device names from /sys/block/<dev>
func (c *sysfsDiskCollector) listSysBlockDevices() ([]string, error) {
	deviceDirs, err := ioutil.ReadDir(c.fs.Path(blockClassPath))
	if err != nil {
		return nil, err
	}
	devices := []string{}

	for _, deviceDir := range deviceDirs {
		if c.fs.ignoredDevicesPattern.MatchString(deviceDir.Name()) {
			level.Debug(c.logger).Log("msg", "Ignoring device", "device", deviceDir.Name())
			continue
		}

		devices = append(devices, deviceDir.Name())

	}
	return devices, nil
}

// newSysFsDiskCollector returns new Linux disk stats collector with sysfs source
// Doc from https://www.kernel.org/doc/html/latest/block/stat.html
// Metrics name matches diskstatsCollector
// Label sets include extra block device info such as fs Label and UUID
func newSysFsDiskCollector(logger log.Logger) (Collector, error) {
	level.Debug(logger).Log("msg", "Loading sysfs disk collector")
	fs, err := NewSysFS()

	if err != nil {
		return nil, err
	}

	collector := sysfsDiskCollector{
		descs:  metricDescs,
		logger: logger,
		fs:     *fs,
	}

	return &collector, nil
}

func init() {
	registerCollector("diskstats", defaultDisabled, newSysFsDiskCollector)
}

// getDeviceInfo returns static sysfs info for device name
func (c *sysfsDiskCollector) getDeviceInfo(name string) (*dmDevice, error) {
	var (
		sectorSize float64
		label      string
		uuid       string
	)

	hwSectorSize, err := c.fs.ReadFile(c.fs.Path(blockClassPath, name, "queue", "hw_sector_size"))
	if err != nil {
		sectorSize = DefaultSectorSize
	} else {
		sz, err := strconv.ParseFloat(hwSectorSize, 64)
		if err != nil {
			sectorSize = DefaultSectorSize
		} else {
			sectorSize = sz
		}
	}

	// Extra info for dm devices
	info, err := os.Stat(c.fs.Path(blockClassPath, name, "dm"))
	if err == nil && info.IsDir() {
		label, err = c.fs.ReadFile(c.fs.Path(blockClassPath, name, "dm", "name"))
		if err != nil {
			level.Warn(c.logger).Log("msg", "Failed to get fs label for device", "device", name, "err", err.Error)
			label = ""
		}

		uuid, err = c.fs.ReadFile(c.fs.Path(blockClassPath, name, "dm", "uuid"))
		if err != nil {
			level.Warn(c.logger).Log("msg", "Failed to get fs uuid for device", "device", name, "err", err.Error)
			uuid = ""
		}
	}

	device := dmDevice{
		Name:       name,
		Label:      label,
		UUID:       uuid,
		SectorSize: sectorSize,
	}

	return &device, nil
}

func (c *sysfsDiskCollector) getDeviceStats(d *dmDevice) ([]metric, error) {
	var m []metric
	stats, err := c.fs.ReadFile(c.fs.Path(blockClassPath, d.Name, "stat"))
	if err != nil {
		return nil, err
	}

	labelNames := []string{"device"}
	labelValues := []string{d.Name}

	if d.Label != "" {
		labelNames = append(labelNames, "label")
		labelValues = append(labelValues, d.Label)
	}

	if d.UUID != "" {
		labelNames = append(labelNames, "uuid")
		labelValues = append(labelValues, d.UUID)

	}

	parts := strings.Fields(stats)
	for i, val := range parts {
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value %s in diskstats: %w", val, err)
		}

		if i > len(metricDescs) {
			// Ignore any unknown stats
			break
		}

		switch conv := metricDescs[i].Conv; conv {
		case Ms:
			v = v / 1000.0
		case Sectors:
			v = v * d.SectorSize
		}

		m = append(m, metric{
			Desc:        &metricDescs[i],
			Value:       v,
			LabelNames:  labelNames,
			LabelValues: labelValues,
		})
	}
	return m, nil
}

func (c *sysfsDiskCollector) Update(ch chan<- prometheus.Metric) error {
	disks, err := c.listSysBlockDevices()
	if err != nil {
		level.Warn(c.logger).Log("msg", "Failed to get list of devices", "err", err.Error)
		return err
	}

	for _, name := range disks {
		dev, err := c.getDeviceInfo(name)
		if err != nil {
			level.Warn(c.logger).Log("msg", "Failed to get info", "device", name, "err", err.Error)
			continue
		}
		stats, err := c.getDeviceStats(dev)
		if err != nil {
			level.Warn(c.logger).Log("msg", "Failed to get device stats", "device", name, "err", err.Error)
		}
		for _, stat := range stats {
			ch <- stat.mustNewConstMetric()
		}
	}
	return nil
}
