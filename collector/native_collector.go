// +build !nonative

package collector

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	procLoad       = "/proc/loadavg"
	procMemInfo    = "/proc/meminfo"
	procInterrupts = "/proc/interrupts"
	procNetDev     = "/proc/net/dev"
	procDiskStats  = "/proc/diskstats"
)

type diskStat struct {
	name          string
	metric        prometheus.Metric
	documentation string
}

var (
	// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
	diskStatsMetrics = []diskStat{
		{"reads_completed", prometheus.NewCounter(), "The total number of reads completed successfully."},
		{"reads_merged", prometheus.NewCounter(), "The number of reads merged. Reads and writes which are adjacent to each other may be merged for efficiency. Thus two 4K reads may become one 8K read before it is ultimately handed to the disk, and so it will be counted (and queued) as only one I/O. This metric lets you know how often this was done."},
		{"sectors_read", prometheus.NewCounter(), "The total number of sectors read successfully."},
		{"read_time_ms", prometheus.NewCounter(), "the total number of milliseconds spent by all reads (as measured from __make_request() to end_that_request_last())."},
		{"writes_completed", prometheus.NewCounter(), "The total number of writes completed successfully."},
		{"writes_merged", prometheus.NewCounter(), "The number of writes merged. Reads and writes which are adjacent to each other may be merged for efficiency. Thus two 4K reads may become one 8K read before it is ultimately handed to the disk, and so it will be counted (and queued) as only one I/O. This metric lets you know how often this was done."},
		{"sectors_written", prometheus.NewCounter(), "The total number of sectors written successfully."},
		{"write_time_ms", prometheus.NewCounter(), "This is the total number of milliseconds spent by all writes (as measured from __make_request() to end_that_request_last())."},
		{"io_now", prometheus.NewGauge(), "The number of I/Os currently in progress. Incremented as requests are given to appropriate struct request_queue and decremented as they finish."},
		{"io_time_ms", prometheus.NewCounter(), "Milliseconds spent doing I/Os. This metric increases so long as node_disk_io_now is nonzero."},
		{"io_time_weighted", prometheus.NewCounter(), "The weighted # of milliseconds spent doing I/Os. This metric is incremented at each I/O start, I/O completion, I/O merge, or read of these stats by the number of I/Os in progress (node_disk_io_now) times the number of milliseconds spent doing I/O since the last update of this field.  This can provide an easy measure of both I/O completion time and the backlog that may be accumulating."},
	}
	lastSeen         = prometheus.NewGauge()
	load1            = prometheus.NewGauge()
	attributes       = prometheus.NewGauge()
	memInfoMetrics   = map[string]prometheus.Gauge{}
	netStatsMetrics  = map[string]prometheus.Gauge{}
	interruptsMetric = prometheus.NewCounter()
)

type nativeCollector struct {
	registry prometheus.Registry
	name     string
	config   Config
}

func init() {
	Factories = append(Factories, NewNativeCollector)
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewNativeCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := nativeCollector{
		name:     "native_collector",
		config:   config,
		registry: registry,
	}

	registry.Register(
		"node_load1",
		"1m load average",
		prometheus.NilLabels,
		load1,
	)

	registry.Register(
		"node_last_login_time",
		"The time of the last login.",
		prometheus.NilLabels,
		lastSeen,
	)

	registry.Register(
		"node_attributes",
		"node_exporter attributes",
		prometheus.NilLabels,
		attributes,
	)

	registry.Register(
		"node_interrupts",
		"Interrupt details from /proc/interrupts",
		prometheus.NilLabels,
		interruptsMetric,
	)

	for _, v := range diskStatsMetrics {
		registry.Register(
			"node_disk_"+v.name,
			v.documentation,
			prometheus.NilLabels,
			v.metric,
		)
	}
	return &c, nil
}

func (c *nativeCollector) Name() string { return c.name }

func (c *nativeCollector) Update() (updates int, err error) {
	last, err := getLastLoginTime()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get last seen: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_last_login_time: %f", last)
	lastSeen.Set(nil, last)

	load, err := getLoad1()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get load: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_load: %f", load)
	load1.Set(nil, load)

	debug(c.Name(), "Set node_attributes{%v}: 1", c.config.Attributes)
	attributes.Set(c.config.Attributes, 1)

	memInfo, err := getMemInfo()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get meminfo: %s", err)
	}
	debug(c.Name(), "Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		if _, ok := memInfoMetrics[k]; !ok {
			memInfoMetrics[k] = prometheus.NewGauge()
			c.registry.Register(
				"node_memory_"+k,
				k+" from /proc/meminfo",
				prometheus.NilLabels,
				memInfoMetrics[k],
			)
		}
		updates++
		memInfoMetrics[k].Set(nil, v)
	}

	interrupts, err := getInterrupts()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get interrupts: %s", err)
	}
	for name, interrupt := range interrupts {
		for cpuNo, value := range interrupt.values {
			updates++
			fv, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return updates, fmt.Errorf("Invalid value %s in interrupts: %s", value, err)
			}
			labels := map[string]string{
				"CPU":     strconv.Itoa(cpuNo),
				"type":    name,
				"info":    interrupt.info,
				"devices": interrupt.devices,
			}
			interruptsMetric.Set(labels, fv)
		}
	}

	netStats, err := getNetStats()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netStats {
		for dev, stats := range devStats {
			for t, value := range stats {
				key := direction + "_" + t
				if _, ok := netStatsMetrics[key]; !ok {
					netStatsMetrics[key] = prometheus.NewGauge()
					c.registry.Register(
						"node_network_"+key,
						t+" "+direction+" from /proc/net/dev",
						prometheus.NilLabels,
						netStatsMetrics[key],
					)
				}
				updates++
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return updates, fmt.Errorf("Invalid value %s in netstats: %s", value, err)
				}
				netStatsMetrics[key].Set(map[string]string{"device": dev}, v)
			}
		}
	}

	diskStats, err := getDiskStats()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get diskstats: %s", err)
	}
	for dev, stats := range diskStats {
		for k, value := range stats {
			updates++
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return updates, fmt.Errorf("Invalid value %s in diskstats: %s", value, err)
			}
			labels := map[string]string{"device": dev}
			counter, ok := diskStatsMetrics[k].metric.(prometheus.Counter)
			if ok {
				counter.Set(labels, v)
			} else {
				var gauge = diskStatsMetrics[k].metric.(prometheus.Gauge)
				gauge.Set(labels, v)
			}
		}
	}
	return updates, err
}

func getLoad1() (float64, error) {
	data, err := ioutil.ReadFile(procLoad)
	if err != nil {
		return 0, err
	}
	return parseLoad(string(data))
}

func parseLoad(data string) (float64, error) {
	parts := strings.Fields(data)
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}

func getLastLoginTime() (float64, error) {
	who := exec.Command("who", "/var/log/wtmp", "-l", "-u", "-s")

	output, err := who.StdoutPipe()
	if err != nil {
		return 0, err
	}

	err = who.Start()
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(output)

	var last time.Time
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return 0, fmt.Errorf("line to long: %s(...)", line)
		}

		fields := strings.Fields(string(line))
		lastDate := fields[2]
		lastTime := fields[3]

		dateParts, err := splitToInts(lastDate, "-") // 2013-04-16
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse date in line '%s': %s", fields, err)
		}

		timeParts, err := splitToInts(lastTime, ":") // 11:33
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse time in line '%s': %s", fields, err)
		}

		last_t := time.Date(dateParts[0], time.Month(dateParts[1]), dateParts[2], timeParts[0], timeParts[1], 0, 0, time.UTC)
		last = last_t
	}
	err = who.Wait()
	if err != nil {
		return 0, err
	}

	return float64(last.Unix()), nil
}

func getMemInfo() (map[string]float64, error) {
	file, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	return parseMemInfo(file)
}

func parseMemInfo(r io.ReadCloser) (map[string]float64, error) {
	defer r.Close()
	memInfo := map[string]float64{}
	scanner := bufio.NewScanner(r)
	re := regexp.MustCompile("\\((.*)\\)")
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid value in meminfo: %s", err)
		}
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
		default:
			return nil, fmt.Errorf("Invalid line in %s: %s", procMemInfo, line)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = re.ReplaceAllString(key, "_${1}")
		memInfo[key] = fv
	}
	return memInfo, nil

}

type interrupt struct {
	info    string
	devices string
	values  []string
}

func getInterrupts() (map[string]interrupt, error) {
	file, err := os.Open(procInterrupts)
	if err != nil {
		return nil, err
	}
	return parseInterrupts(file)
}

func parseInterrupts(r io.ReadCloser) (map[string]interrupt, error) {
	defer r.Close()
	interrupts := map[string]interrupt{}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return nil, fmt.Errorf("%s empty", procInterrupts)
	}
	cpuNum := len(strings.Fields(string(scanner.Text()))) // one header per cpu

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		if len(parts) < cpuNum+2 { // irq + one column per cpu + details,
			continue // we ignore ERR and MIS for now
		}
		intName := parts[0][:len(parts[0])-1] // remove trailing :
		intr := interrupt{
			values: parts[1:cpuNum],
		}

		if _, err := strconv.Atoi(intName); err == nil { // numeral interrupt
			intr.info = parts[cpuNum+1]
			intr.devices = strings.Join(parts[cpuNum+2:], " ")
		} else {
			intr.info = strings.Join(parts[cpuNum+1:], " ")
		}
		interrupts[intName] = intr
	}
	return interrupts, nil
}

func getNetStats() (map[string]map[string]map[string]string, error) {
	file, err := os.Open(procNetDev)
	if err != nil {
		return nil, err
	}
	return parseNetStats(file)
}

func parseNetStats(r io.ReadCloser) (map[string]map[string]map[string]string, error) {
	defer r.Close()
	netStats := map[string]map[string]map[string]string{}
	netStats["transmit"] = map[string]map[string]string{}
	netStats["receive"] = map[string]map[string]string{}

	scanner := bufio.NewScanner(r)
	scanner.Scan() // skip first header
	scanner.Scan()
	parts := strings.Split(string(scanner.Text()), "|")
	if len(parts) != 3 { // interface + receive + transmit
		return nil, fmt.Errorf("Invalid header line in %s: %s",
			procNetDev, scanner.Text())
	}
	header := strings.Fields(parts[1])
	for scanner.Scan() {
		parts := strings.Fields(string(scanner.Text()))
		if len(parts) != 2*len(header)+1 {
			return nil, fmt.Errorf("Invalid line in %s: %s",
				procNetDev, scanner.Text())
		}

		dev := parts[0][:len(parts[0])-1]
		receive, err := parseNetDevLine(parts[1:len(header)+1], header)
		if err != nil {
			return nil, err
		}

		transmit, err := parseNetDevLine(parts[len(header)+1:], header)
		if err != nil {
			return nil, err
		}
		netStats["transmit"][dev] = transmit
		netStats["receive"][dev] = receive
	}
	return netStats, nil
}

func parseNetDevLine(parts []string, header []string) (map[string]string, error) {
	devStats := map[string]string{}
	for i, v := range parts {
		devStats[header[i]] = v
	}
	return devStats, nil
}

func getDiskStats() (map[string]map[int]string, error) {
	file, err := os.Open(procDiskStats)
	if err != nil {
		return nil, err
	}
	return parseDiskStats(file)
}

func parseDiskStats(r io.ReadCloser) (map[string]map[int]string, error) {
	defer r.Close()
	diskStats := map[string]map[int]string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(string(scanner.Text()))
		if len(parts) != len(diskStatsMetrics)+3 { // we strip major, minor and dev
			return nil, fmt.Errorf("Invalid line in %s: %s", procDiskStats, scanner.Text())
		}
		dev := parts[2]
		diskStats[dev] = map[int]string{}
		for i, v := range parts[3:] {
			diskStats[dev][i] = v
		}
	}
	return diskStats, nil
}
