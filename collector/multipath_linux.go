// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nomultipath

package collector

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"path/filepath"
	"strconv"
	"time"
	"unsafe"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// sizeOfSizeT matches C's sizeof(size_t) for the multipathd wire protocol.
const sizeOfSizeT = unsafe.Sizeof(uintptr(0))

const maxReplyLen = 32 * 1024 * 1024

var (
	multipathSocketPath = kingpin.Flag(
		"collector.multipath.socket-path",
		"Path to the multipathd unix socket.",
	).Default("/run/multipathd.socket").String()

	multipathTimeout = kingpin.Flag(
		"collector.multipath.timeout",
		"Timeout for multipathd socket communication.",
	).Default("4s").Duration()
)

// multipathTopology mirrors the JSON from "show maps json".
type multipathTopology struct {
	MajorVersion int            `json:"major_version"`
	MinorVersion int            `json:"minor_version"`
	Maps         []multipathMap `json:"maps"`
}

type multipathMap struct {
	Name       string        `json:"name"`
	UUID       string        `json:"uuid"`
	Sysfs      string        `json:"sysfs"`
	Failback   string        `json:"failback"`
	Queueing   string        `json:"queueing"`
	Paths      int           `json:"paths"`
	WriteProt  string        `json:"write_prot"`
	DmSt       string        `json:"dm_st"`
	Features   string        `json:"features"`
	HwHandler  string        `json:"hwhandler"`
	Action     string        `json:"action"`
	PathFaults int           `json:"path_faults"`
	Vendor     string        `json:"vend"`
	Product    string        `json:"prod"`
	Revision   string        `json:"rev"`
	SwitchGrp  int           `json:"switch_grp"`
	MapLoads   int           `json:"map_loads"`
	TotalQTime int           `json:"total_q_time"`
	QTimeouts  int           `json:"q_timeouts"`
	PathGroups []multipathPG `json:"path_groups"`
}

type multipathPG struct {
	Selector   string          `json:"selector"`
	Priority   int             `json:"pri"`
	DmSt       string          `json:"dm_st"`
	MarginalSt string          `json:"marginal_st"`
	Group      int             `json:"group"`
	Paths      []multipathPath `json:"paths"`
}

type multipathPath struct {
	Dev         string `json:"dev"`
	DevT        string `json:"dev_t"`
	DmSt        string `json:"dm_st"`
	DevSt       string `json:"dev_st"`
	ChkSt       string `json:"chk_st"`
	Checker     string `json:"checker"`
	Priority    int    `json:"pri"`
	HostWWNN    string `json:"host_wwnn"`
	TargetWWNN  string `json:"target_wwnn"`
	HostWWPN    string `json:"host_wwpn"`
	TargetWWPN  string `json:"target_wwpn"`
	HostAdapter string `json:"host_adapter"`
	LunHex      string `json:"lun_hex"`
	MarginalSt  string `json:"marginal_st"`
}

var checkerStates = []string{
	"ready", "faulty", "shaky", "ghost",
	"pending", "timeout", "delayed", "disconnected", "unknown",
}

func normalizeCheckerState(raw string) string {
	switch raw {
	case "ready", "faulty", "shaky", "ghost", "delayed", "disconnected":
		return raw
	case "i/o pending":
		return "pending"
	case "i/o timeout":
		return "timeout"
	default:
		return "unknown"
	}
}

type multipathCollector struct {
	logger             *slog.Logger
	queryTopology      func() (*multipathTopology, error)
	readDeviceSize     func(sysfsName string) (uint64, error)
	scanNVMeSubsystems func() ([]nvmeSubsystem, error)
	daemonUp           *prometheus.Desc
	deviceInfo         *prometheus.Desc
	deviceActive       *prometheus.Desc
	deviceSizeBytes    *prometheus.Desc
	devicePathsTotal   *prometheus.Desc
	devicePathsActive  *prometheus.Desc
	devicePathsFailed  *prometheus.Desc
	devicePathFaults   *prometheus.Desc
	deviceSwitchGroup  *prometheus.Desc
	deviceMapLoads     *prometheus.Desc
	deviceQTimeouts    *prometheus.Desc
	pathActive         *prometheus.Desc
	pathCheckerState   *prometheus.Desc

	// NVMe subsystem metrics
	nvmeSubsystemInfo       *prometheus.Desc
	nvmeSubsystemPathsTotal *prometheus.Desc
	nvmeSubsystemPathsLive  *prometheus.Desc
	nvmePathState           *prometheus.Desc
}

func init() {
	registerCollector("multipath", defaultDisabled, NewMultipathCollector)
}

// NewMultipathCollector returns a new Collector exposing multipath device metrics.
func NewMultipathCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "multipath"

	deviceLabels := []string{"device"}
	pathLabels := []string{"device", "path", "path_group"}

	c := &multipathCollector{
		logger: logger,
		daemonUp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "daemon_up"),
			"Whether the multipathd daemon is reachable (1 = up, 0 = down).",
			nil, nil,
		),
		deviceInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_info"),
			"Non-numeric information about a multipath device.",
			[]string{"device", "uuid", "sysfs", "vendor", "product", "revision", "write_protect"},
			nil,
		),
		deviceActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_active"),
			"Whether the multipath device-mapper device is active (1) or suspended (0).",
			deviceLabels, nil,
		),
		deviceSizeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_size_bytes"),
			"Size of the multipath device in bytes, read from /sys/block/<dm>/size.",
			deviceLabels, nil,
		),
		devicePathsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_total"),
			"Total number of paths for a multipath device.",
			deviceLabels, nil,
		),
		devicePathsActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_active"),
			"Number of active and healthy paths for a multipath device.",
			deviceLabels, nil,
		),
		devicePathsFailed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_paths_failed"),
			"Number of failed paths for a multipath device.",
			deviceLabels, nil,
		),
		devicePathFaults: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_path_faults_total"),
			"Cumulative number of path faults for a multipath device.",
			deviceLabels, nil,
		),
		deviceSwitchGroup: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_switch_group_total"),
			"Cumulative number of path group switches for a multipath device.",
			deviceLabels, nil,
		),
		deviceMapLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_map_loads_total"),
			"Cumulative number of map reloads for a multipath device.",
			deviceLabels, nil,
		),
		deviceQTimeouts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "device_queueing_timeouts_total"),
			"Cumulative number of queueing timeouts for a multipath device.",
			deviceLabels, nil,
		),
		pathActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "path_active"),
			"Whether the path's device-mapper state is active (1) or failed (0).",
			pathLabels, nil,
		),
		pathCheckerState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "path_checker_state"),
			"Current path checker state (1 for the current state, 0 for all others).",
			append(pathLabels, "state"), nil,
		),

		nvmeSubsystemInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "nvme_subsystem_info"),
			"Non-numeric information about an NVMe subsystem.",
			[]string{"subsystem", "nqn", "model", "serial", "iopolicy"}, nil,
		),
		nvmeSubsystemPathsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "nvme_subsystem_paths_total"),
			"Total number of controller paths for an NVMe subsystem.",
			[]string{"subsystem"}, nil,
		),
		nvmeSubsystemPathsLive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "nvme_subsystem_paths_live"),
			"Number of controller paths in live state for an NVMe subsystem.",
			[]string{"subsystem"}, nil,
		),
		nvmePathState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "nvme_path_state"),
			"Current NVMe controller path state (1 for the current state, 0 for all others).",
			[]string{"subsystem", "controller", "transport", "state"}, nil,
		),
	}

	c.queryTopology = func() (*multipathTopology, error) {
		return queryMultipathd(*multipathSocketPath, *multipathTimeout)
	}
	c.readDeviceSize = func(sysfsName string) (uint64, error) {
		sectors, err := readUintFromFile(sysFilePath(filepath.Join("block", sysfsName, "size")))
		if err != nil {
			return 0, err
		}
		return sectors * uint64(unixSectorSize), nil
	}
	c.scanNVMeSubsystems = func() ([]nvmeSubsystem, error) {
		return scanNVMeSubsystems(*sysPath)
	}

	return c, nil
}

func (c *multipathCollector) Update(ch chan<- prometheus.Metric) error {
	var dmErr error

	topology, err := c.queryTopology()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.daemonUp, prometheus.GaugeValue, 0)
		c.logger.Debug("multipathd not reachable", "err", err)
		dmErr = err
	} else {
		ch <- prometheus.MustNewConstMetric(c.daemonUp, prometheus.GaugeValue, 1)
		for _, m := range topology.Maps {
			c.emitDeviceMetrics(ch, m)
		}
	}

	subsystems, err := c.scanNVMeSubsystems()
	if err != nil {
		c.logger.Debug("NVMe subsystem sysfs not available", "err", err)
	} else if len(subsystems) > 0 {
		c.emitNVMeSubsystemMetrics(ch, subsystems)
	}

	if dmErr != nil && (err != nil || len(subsystems) == 0) {
		return fmt.Errorf("failed to query multipathd: %w", dmErr)
	}

	return nil
}

func (c *multipathCollector) emitDeviceMetrics(ch chan<- prometheus.Metric, m multipathMap) {
	device := m.Name

	ch <- prometheus.MustNewConstMetric(c.deviceInfo, prometheus.GaugeValue, 1,
		device, m.UUID, m.Sysfs, m.Vendor, m.Product, m.Revision, m.WriteProt)

	active := 0.0
	if m.DmSt == "active" {
		active = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.deviceActive, prometheus.GaugeValue, active, device)

	if sizeBytes, err := c.readDeviceSize(m.Sysfs); err == nil {
		ch <- prometheus.MustNewConstMetric(c.deviceSizeBytes, prometheus.GaugeValue, float64(sizeBytes), device)
	} else {
		c.logger.Debug("failed to read device size", "device", device, "sysfs", m.Sysfs, "err", err)
	}

	ch <- prometheus.MustNewConstMetric(c.devicePathsTotal, prometheus.GaugeValue, float64(m.Paths), device)
	ch <- prometheus.MustNewConstMetric(c.devicePathFaults, prometheus.CounterValue, float64(m.PathFaults), device)
	ch <- prometheus.MustNewConstMetric(c.deviceSwitchGroup, prometheus.CounterValue, float64(m.SwitchGrp), device)
	ch <- prometheus.MustNewConstMetric(c.deviceMapLoads, prometheus.CounterValue, float64(m.MapLoads), device)
	ch <- prometheus.MustNewConstMetric(c.deviceQTimeouts, prometheus.CounterValue, float64(m.QTimeouts), device)

	var activePaths, failedPaths float64
	for _, pg := range m.PathGroups {
		pgStr := strconv.Itoa(pg.Group)
		for _, p := range pg.Paths {
			c.emitPathMetrics(ch, device, pgStr, p)
			if p.DmSt == "active" && p.ChkSt == "ready" {
				activePaths++
			}
			if p.DmSt == "failed" || p.ChkSt == "faulty" {
				failedPaths++
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(c.devicePathsActive, prometheus.GaugeValue, activePaths, device)
	ch <- prometheus.MustNewConstMetric(c.devicePathsFailed, prometheus.GaugeValue, failedPaths, device)
}

func (c *multipathCollector) emitPathMetrics(ch chan<- prometheus.Metric, device, pgStr string, p multipathPath) {
	pathActive := 0.0
	if p.DmSt == "active" {
		pathActive = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.pathActive, prometheus.GaugeValue, pathActive,
		device, p.Dev, pgStr)

	currentState := normalizeCheckerState(p.ChkSt)
	for _, state := range checkerStates {
		val := 0.0
		if state == currentState {
			val = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.pathCheckerState, prometheus.GaugeValue, val,
			device, p.Dev, pgStr, state)
	}
}

// queryMultipathd connects to the multipathd socket and runs "show maps json".
func queryMultipathd(socketPath string, timeout time.Duration) (*multipathTopology, error) {
	conn, err := net.DialTimeout("unix", socketPath, timeout)
	if err != nil {
		return nil, fmt.Errorf("connecting to multipathd socket %q: %w", socketPath, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("setting socket deadline: %w", err)
	}

	if err := sendMultipathCmd(conn, "show maps json"); err != nil {
		return nil, fmt.Errorf("sending command: %w", err)
	}

	reply, err := recvMultipathReply(conn)
	if err != nil {
		return nil, fmt.Errorf("receiving reply: %w", err)
	}

	var topo multipathTopology
	if err := json.Unmarshal(reply, &topo); err != nil {
		return nil, fmt.Errorf("parsing JSON reply: %w", err)
	}

	return &topo, nil
}

// sendMultipathCmd writes a command using the multipathd wire protocol:
// native-endian size_t length prefix, then the command bytes with a null terminator.
func sendMultipathCmd(w io.Writer, cmd string) error {
	cmdBytes := append([]byte(cmd), 0)
	length := uint64(len(cmdBytes))

	buf := make([]byte, sizeOfSizeT)
	switch sizeOfSizeT {
	case 8:
		binary.NativeEndian.PutUint64(buf, length)
	case 4:
		binary.NativeEndian.PutUint32(buf, uint32(length))
	default:
		return fmt.Errorf("unsupported size_t: %d bytes", sizeOfSizeT)
	}

	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("writing length: %w", err)
	}
	if _, err := w.Write(cmdBytes); err != nil {
		return fmt.Errorf("writing command: %w", err)
	}
	return nil
}

// recvMultipathReply reads a reply using the multipathd wire protocol:
// native-endian size_t length prefix, then the reply data.
func recvMultipathReply(r io.Reader) ([]byte, error) {
	buf := make([]byte, sizeOfSizeT)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("reading reply length: %w", err)
	}

	var length uint64
	switch sizeOfSizeT {
	case 8:
		length = binary.NativeEndian.Uint64(buf)
	case 4:
		length = uint64(binary.NativeEndian.Uint32(buf))
	default:
		return nil, fmt.Errorf("unsupported size_t: %d bytes", sizeOfSizeT)
	}

	if length == 0 || length >= maxReplyLen {
		return nil, fmt.Errorf("reply length out of range: %d", length)
	}

	reply := make([]byte, length)
	if _, err := io.ReadFull(r, reply); err != nil {
		return nil, fmt.Errorf("reading reply data: %w", err)
	}

	if len(reply) > 0 && reply[len(reply)-1] == 0 {
		reply = reply[:len(reply)-1]
	}

	return reply, nil
}
