// +build !nonative

package collector

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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

var (
	diskStatsHeader = []string{
		"reads_completed", "reads_merged",
		"sectors_read", "read_time_ms",
		"writes_completed", "writes_merged",
		"sectors_written", "write_time_ms",
		"io_now", "io_time_ms", "io_time_weighted",
	}
)

type nativeCollector struct {
	loadAvg    prometheus.Gauge
	attributes prometheus.Gauge
	lastSeen   prometheus.Gauge
	memInfo    prometheus.Gauge
	interrupts prometheus.Counter
	netStats   prometheus.Counter
	diskStats  prometheus.Counter
	name       string
	config     Config
}

func init() {
	Factories = append(Factories, NewNativeCollector)
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewNativeCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := nativeCollector{
		name:       "native_collector",
		config:     config,
		loadAvg:    prometheus.NewGauge(),
		attributes: prometheus.NewGauge(),
		lastSeen:   prometheus.NewGauge(),
		memInfo:    prometheus.NewGauge(),
		interrupts: prometheus.NewCounter(),
		netStats:   prometheus.NewCounter(),
		diskStats:  prometheus.NewCounter(),
	}

	registry.Register(
		"node_load",
		"node_exporter: system load.",
		prometheus.NilLabels,
		c.loadAvg,
	)

	registry.Register(
		"node_last_login_seconds",
		"node_exporter: seconds since last login.",
		prometheus.NilLabels,
		c.lastSeen,
	)

	registry.Register(
		"node_attributes",
		"node_exporter: system attributes.",
		prometheus.NilLabels,
		c.attributes,
	)

	registry.Register(
		"node_mem",
		"node_exporter: memory details.",
		prometheus.NilLabels,
		c.memInfo,
	)

	registry.Register(
		"node_interrupts",
		"node_exporter: interrupt details.",
		prometheus.NilLabels,
		c.interrupts,
	)

	registry.Register(
		"node_net",
		"node_exporter: network stats.",
		prometheus.NilLabels,
		c.netStats,
	)

	registry.Register(
		"node_disk",
		"node_exporter: disk stats.",
		prometheus.NilLabels,
		c.diskStats,
	)
	return &c, nil
}

func (c *nativeCollector) Name() string { return c.name }

func (c *nativeCollector) Update() (updates int, err error) {
	last, err := getSecondsSinceLastLogin()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get last seen: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_last_login_seconds: %f", last)
	c.lastSeen.Set(nil, last)

	load, err := getLoad()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get load: %s", err)
	}
	updates++
	debug(c.Name(), "Set node_load: %f", load)
	c.loadAvg.Set(nil, load)

	debug(c.Name(), "Set node_attributes{%v}: 1", c.config.Attributes)
	c.attributes.Set(c.config.Attributes, 1)

	memInfo, err := getMemInfo()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get meminfo: %s", err)
	}
	debug(c.Name(), "Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		updates++
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return updates, fmt.Errorf("Invalid value in meminfo: %s", err)
		}
		c.memInfo.Set(map[string]string{"type": k}, fv)
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
			c.interrupts.Set(labels, fv)
		}
	}

	netStats, err := getNetStats()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netStats {
		for dev, stats := range devStats {
			for t, value := range stats {
				updates++
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return updates, fmt.Errorf("Invalid value %s in interrupts: %s", value, err)
				}
				labels := map[string]string{
					"device":    dev,
					"direction": direction,
					"type":      t,
				}
				c.netStats.Set(labels, v)
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
			labels := map[string]string{"device": dev, "type": k}
			c.diskStats.Set(labels, v)
		}
	}
	return updates, err
}

func getLoad() (float64, error) {
	data, err := ioutil.ReadFile(procLoad)
	if err != nil {
		return 0, err
	}
	parts := strings.Fields(string(data))
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}

func getSecondsSinceLastLogin() (float64, error) {
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

	return float64(time.Now().Sub(last).Seconds()), nil
}

func getMemInfo() (map[string]string, error) {
	memInfo := map[string]string{}
	fh, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		key := ""
		switch len(parts) {
		case 2: // no unit
			key = parts[0][:len(parts[0])-1] // remove trailing : from key
		case 3: // has unit
			key = fmt.Sprintf("%s_%s", parts[0][:len(parts[0])-1], parts[2])
		default:
			return nil, fmt.Errorf("Invalid line in %s: %s", procMemInfo, line)
		}
		memInfo[key] = parts[1]
	}
	return memInfo, nil

}

type interrupt struct {
	info    string
	devices string
	values  []string
}

func getInterrupts() (map[string]interrupt, error) {
	interrupts := map[string]interrupt{}
	fh, err := os.Open(procInterrupts)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
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
	netStats := map[string]map[string]map[string]string{}
	netStats["transmit"] = map[string]map[string]string{}
	netStats["receive"] = map[string]map[string]string{}
	fh, err := os.Open(procNetDev)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
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

func getDiskStats() (map[string]map[string]string, error) {
	diskStats := map[string]map[string]string{}
	fh, err := os.Open(procDiskStats)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		parts := strings.Fields(string(scanner.Text()))
		if len(parts) != len(diskStatsHeader)+3 { // we strip major, minor and dev
			return nil, fmt.Errorf("Invalid line in %s: %s", procDiskStats, scanner.Text())
		}
		dev := parts[2]
		diskStats[dev] = map[string]string{}
		for i, v := range parts[3:] {
			diskStats[dev][diskStatsHeader[i]] = v
		}
	}
	return diskStats, nil
}
