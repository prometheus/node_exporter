// Copyright 2017 The Prometheus Authors
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

package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/procfs"

	"github.com/prometheus/node_exporter/config"
)

var binary = filepath.Join(os.Getenv("GOPATH"), "bin/node_exporter")

const (
	address = "localhost:19100"
)

// configFromArgs parses args with a new kingpin app that has only the
// collector flags, and returns the Config they describe.
func configFromArgs(t *testing.T, args ...string) (config.Config, error) {
	t.Helper()
	app := kingpin.New("node_exporter", "")
	configFromFlags := addFlags(app)
	if _, err := app.Parse(args); err != nil {
		t.Fatal(err)
	}
	return configFromFlags(slog.New(slog.DiscardHandler))
}

func parseConfigFlags(t *testing.T, args ...string) config.Config {
	t.Helper()
	cfg, err := configFromArgs(t, args...)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestFlagDefaultsMatchConfigDefaults(t *testing.T) {
	if got, want := parseConfigFlags(t), config.NewConfigWithDefaults(); !reflect.DeepEqual(got, want) {
		t.Fatalf("flag defaults differ from config defaults:\nflags:  %+v\nconfig: %+v", got, want)
	}
}

// TestEveryFlagChangesConfig sets each collector and path flag to a value
// other than its default, and checks that the Config changes, so a new flag
// cannot be left out of configFromFlags.
func TestEveryFlagChangesConfig(t *testing.T) {
	app := kingpin.New("node_exporter", "")
	addFlags(app)
	defaults := config.NewConfigWithDefaults()
	for _, flag := range app.Model().Flags {
		if !strings.HasPrefix(flag.Name, "collector.") && !strings.HasPrefix(flag.Name, "path.") {
			continue
		}
		var def string
		if len(flag.Default) > 0 {
			def = flag.Default[0]
		}
		var arg string
		switch {
		case flag.IsBoolFlag() && def == "true":
			arg = "--no-" + flag.Name
		case flag.IsBoolFlag():
			arg = "--" + flag.Name
		default:
			value := "other"
			if _, err := time.ParseDuration(def); err == nil {
				value = "7s"
			} else if _, err := strconv.Atoi(def); err == nil {
				value = "7"
			}
			arg = "--" + flag.Name + "=" + value
		}
		cfg, err := configFromArgs(t, arg)
		if err != nil {
			t.Errorf("%s: %v", arg, err)
			continue
		}
		if reflect.DeepEqual(cfg, defaults) {
			t.Errorf("%s does not change the Config", arg)
		}
	}
}

func TestFlagsSetConfig(t *testing.T) {
	got := parseConfigFlags(t,
		"--path.procfs=/host/proc",
		"--path.sysfs=/host/sys",
		"--path.rootfs=/host",
		"--path.udev.data=/host/run/udev/data",
		"--collector.disable-defaults",
		"--collector.ntp",
		"--no-collector.textfile",
		"--collector.cpu.info",
		"--no-collector.cpu.guest",
		"--collector.cpu.info.flags-include=^fpu$",
		"--collector.cpu.info.bugs-include=^spectre",
		"--collector.textfile.directory=a",
		"--collector.textfile.directory=b",
		"--collector.ntp.server=192.0.2.1",
		"--collector.ntp.server-port=1123",
		"--collector.ntp.protocol-version=3",
		"--collector.ntp.server-is-local",
		"--collector.ntp.ip-ttl=2",
		"--collector.ntp.max-distance=2s",
		"--collector.ntp.local-offset-tolerance=5ms",
		"--collector.stat.softirq",
		"--collector.vmstat.fields=^pg",
		"--collector.netstat.fields=^Tcp_",
		"--collector.bcache.priorityStats",
		"--no-collector.nvme.namespace-usage",
		"--collector.rapl.enable-zone-label",
		"--collector.arp.device-include=^eth",
		"--collector.arp.device-exclude=^lo$",
		"--no-collector.arp.netlink",
		"--collector.ethtool.device-include=^eth0$",
		"--collector.ethtool.device-exclude=^eth1$",
		"--collector.ethtool.metrics-include=^rx_",
		"--collector.infiniband.device-include=^mlx",
		"--collector.infiniband.device-exclude=^ib1$",
		"--collector.netclass.ignored-devices=^veth",
		"--collector.netclass.ignore-invalid-speed",
		"--collector.netclass.netlink",
		"--collector.netclass_rtnl.with-stats",
		"--collector.netdev.device-include=^en",
		"--collector.netdev.address-info",
		"--collector.netdev.enable-detailed-metrics",
		"--no-collector.netdev.netlink",
		"--collector.netdev.label-ifalias",
		"--collector.qdisc.fixtures=qdisc/",
		"--collector.qdisc.device-exclude=^lo$",
		"--collector.ipvs.backend-labels=local_port",
		"--collector.wifi.fixtures=wifi/",
		"--collector.diskstats.device-exclude=^loop",
		"--collector.filesystem.mount-points-exclude=^/boot",
		"--collector.filesystem.fs-types-exclude=^tmpfs$",
		"--collector.filesystem.mount-timeout=7s",
		"--collector.filesystem.stat-workers=2",
		"--collector.systemd.unit-include=^ssh",
		"--collector.systemd.unit-exclude=^user",
		"--collector.systemd.private",
		"--collector.systemd.enable-task-metrics",
		"--collector.systemd.enable-restarts-metrics",
		"--collector.systemd.enable-start-time-metrics",
		"--collector.perf.cpus=1-2",
		"--collector.perf.tracepoint=sched:sched_switch",
		"--collector.perf.disable-hardware-profilers",
		"--collector.perf.hardware-profilers=CpuCycles",
		"--collector.perf.disable-software-profilers",
		"--collector.perf.software-profilers=PageFault",
		"--collector.perf.disable-cache-profilers",
		"--collector.perf.cache-profilers=L1DataReadHit",
		"--collector.sysctl.include=fs.file-nr",
		"--collector.sysctl.include-info=kernel.seccomp.actions_avail",
		"--collector.hwmon.chip-include=^coretemp",
		"--collector.hwmon.chip-exclude=^acpi",
		"--collector.hwmon.sensor-include=^temp1",
		"--collector.hwmon.sensor-exclude=^fan",
		"--collector.cpufreq.enable-cpufreq-prefix",
		"--collector.interrupts.name-include=^eth",
		"--collector.interrupts.name-exclude=^LOC$",
		"--no-collector.interrupts.include-zeros",
		"--collector.pcidevice.idsfile=pci.ids",
		"--collector.pcidevice.names",
		"--collector.powersupply.ignored-supplies=^AC$",
		"--collector.resctrl.path=resctrl/",
		"--collector.runit.servicedir=/run/service",
		"--collector.supervisord.url=unix:///run/supervisord.sock",
		"--collector.slabinfo.slabs-include=^dentry$",
		"--collector.slabinfo.slabs-exclude=^kmalloc",
		"--collector.tapestats.ignored-devices=^st1$",
	)

	want := config.NewConfigWithDefaults()
	want.Path = config.PathConfig{ProcFS: "/host/proc", SysFS: "/host/sys", RootFS: "/host", UdevData: "/host/run/udev/data"}
	want.CollectorDisableDefaults = true
	// Collectors has only the collector flags that the command line sets.
	want.Collectors = map[string]bool{"ntp": true, "textfile": false}
	want.CPU = config.CPUConfig{EnableInfo: true, FlagsInclude: "^fpu$", BugsInclude: "^spectre"}
	want.Textfile = config.TextfileConfig{Directories: []string{"a", "b"}}
	want.NTP = config.NTPConfig{
		Server:               "192.0.2.1",
		ServerPort:           1123,
		ProtocolVersion:      3,
		ServerIsLocal:        true,
		IPTTL:                2,
		MaxDistance:          2 * time.Second,
		LocalOffsetTolerance: 5 * time.Millisecond,
	}
	want.Stat = config.StatConfig{Softirq: true}
	want.VMStat = config.VMStatConfig{Fields: "^pg"}
	want.Netstat = config.NetstatConfig{Fields: "^Tcp_"}
	want.Bcache = config.BcacheConfig{PriorityStats: true}
	want.NVMe = config.NVMeConfig{NamespaceUsage: false}
	want.RAPL = config.RAPLConfig{EnableZoneLabel: true}
	want.ARP = config.ARPConfig{DeviceInclude: "^eth", DeviceExclude: "^lo$", Netlink: false}
	want.Ethtool = config.EthtoolConfig{DeviceInclude: "^eth0$", DeviceExclude: "^eth1$", MetricsInclude: "^rx_"}
	want.Infiniband = config.InfinibandConfig{DeviceInclude: "^mlx", DeviceExclude: "^ib1$"}
	want.NetClass = config.NetClassConfig{IgnoredDevices: "^veth", IgnoreInvalidSpeed: true, Netlink: true, RTNLWithStats: true}
	want.NetDev = config.NetDevConfig{DeviceInclude: "^en", AddressInfo: true, DetailedMetrics: true, Netlink: false, LabelIfAlias: true}
	want.Qdisc = config.QdiscConfig{Fixtures: "qdisc/", DeviceExclude: "^lo$"}
	want.IPVS = config.IPVSConfig{BackendLabels: "local_port"}
	want.Wifi = config.WifiConfig{Fixtures: "wifi/"}
	want.Systemd = config.SystemdConfig{
		UnitInclude:            "^ssh",
		UnitExclude:            "^user",
		Private:                true,
		EnableTaskMetrics:      true,
		EnableRestartsMetrics:  true,
		EnableStartTimeMetrics: true,
	}
	want.Perf = config.PerfConfig{
		CPUs:                     "1-2",
		Tracepoints:              []string{"sched:sched_switch"},
		DisableHardwareProfilers: true,
		HardwareProfilers:        []string{"CpuCycles"},
		DisableSoftwareProfilers: true,
		SoftwareProfilers:        []string{"PageFault"},
		DisableCacheProfilers:    true,
		CacheProfilers:           []string{"L1DataReadHit"},
	}
	want.Sysctl = config.SysctlConfig{Include: []string{"fs.file-nr"}, IncludeInfo: []string{"kernel.seccomp.actions_avail"}}
	want.Hwmon = config.HwmonConfig{ChipInclude: "^coretemp", ChipExclude: "^acpi", SensorInclude: "^temp1", SensorExclude: "^fan"}
	want.CPUFreq = config.CPUFreqConfig{EnablePrefix: true}
	want.Interrupts = config.InterruptsConfig{NameInclude: "^eth", NameExclude: "^LOC$", IncludeZeros: false}
	want.PCIDevice = config.PCIDeviceConfig{IDsFile: "pci.ids", Names: true}
	want.PowerSupply = config.PowerSupplyConfig{IgnoredSupplies: "^AC$"}
	want.Resctrl = config.ResctrlConfig{Path: "resctrl/"}
	want.Runit = config.RunitConfig{ServiceDir: "/run/service"}
	want.Supervisord = config.SupervisordConfig{URL: "unix:///run/supervisord.sock"}
	want.Slabinfo = config.SlabinfoConfig{SlabsInclude: "^dentry$", SlabsExclude: "^kmalloc"}
	want.Tapestats = config.TapestatsConfig{IgnoredDevices: "^st1$"}
	want.Diskstats = config.DiskstatsConfig{DeviceExclude: "^loop"}
	want.Filesystem = config.FilesystemConfig{
		MountPointsExclude: "^/boot",
		FSTypesExclude:     "^tmpfs$",
		MountTimeout:       7 * time.Second,
		StatWorkerCount:    2,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got:  %+v\nwant: %+v", got, want)
	}
}

func TestDeprecatedFlags(t *testing.T) {
	for _, test := range []struct {
		args    []string
		want    func(*config.Config)
		wantErr string
	}{
		{
			args: []string{"--collector.netdev.device-whitelist=^en"},
			want: func(c *config.Config) { c.NetDev.DeviceInclude = "^en" },
		},
		{
			args: []string{"--collector.netdev.device-blacklist=^lo$"},
			want: func(c *config.Config) { c.NetDev.DeviceExclude = "^lo$" },
		},
		{
			args: []string{"--collector.qdisk.device-include=^eth"},
			want: func(c *config.Config) { c.Qdisc.DeviceInclude = "^eth" },
		},
		{
			args: []string{"--collector.qdisk.device-exclude=^lo$"},
			want: func(c *config.Config) { c.Qdisc.DeviceExclude = "^lo$" },
		},
		{
			args: []string{"--collector.diskstats.ignored-devices=^loop"},
			want: func(c *config.Config) { c.Diskstats.DeviceExclude = "^loop" },
		},
		{
			args: []string{"--collector.filesystem.ignored-mount-points=^/boot"},
			want: func(c *config.Config) { c.Filesystem.MountPointsExclude = "^/boot" },
		},
		{
			args: []string{"--collector.filesystem.ignored-fs-types=^tmpfs$"},
			want: func(c *config.Config) { c.Filesystem.FSTypesExclude = "^tmpfs$" },
		},
		{
			args: []string{"--collector.systemd.unit-whitelist=^ssh"},
			want: func(c *config.Config) { c.Systemd.UnitInclude = "^ssh" },
		},
		{
			args: []string{"--collector.systemd.unit-blacklist=^user"},
			want: func(c *config.Config) { c.Systemd.UnitExclude = "^user" },
		},
		{
			args:    []string{"--collector.systemd.unit-blacklist=^user", "--collector.systemd.unit-exclude=^dev"},
			wantErr: "--collector.systemd.unit-blacklist and --collector.systemd.unit-exclude are mutually exclusive",
		},
		{
			args:    []string{"--collector.diskstats.ignored-devices=^loop", "--collector.diskstats.device-exclude=^ram"},
			wantErr: "--collector.diskstats.ignored-devices and --collector.diskstats.device-exclude are mutually exclusive",
		},
		{
			args:    []string{"--collector.netdev.device-whitelist=^en", "--collector.netdev.device-include=^eth"},
			wantErr: "--collector.netdev.device-whitelist and --collector.netdev.device-include are mutually exclusive",
		},
		{
			args:    []string{"--collector.qdisk.device-exclude=^lo$", "--collector.qdisc.device-exclude=^eth"},
			wantErr: "--collector.qdisk.device-exclude and --collector.qdisc.device-exclude are mutually exclusive",
		},
	} {
		t.Run(test.args[0], func(t *testing.T) {
			got, err := configFromArgs(t, test.args...)
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Fatalf("want error %q, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := config.NewConfigWithDefaults()
			test.want(&want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got:  %+v\nwant: %+v", got, want)
			}
		})
	}
}

func TestIncludeFlagReplacesDefaultExclude(t *testing.T) {
	for _, test := range []struct {
		args []string
		want func(*config.Config)
	}{
		{
			args: []string{"--collector.diskstats.device-include=^sd"},
			want: func(c *config.Config) { c.Diskstats = config.DiskstatsConfig{DeviceInclude: "^sd"} },
		},
		{
			args: []string{"--collector.filesystem.mount-points-include=^/data"},
			want: func(c *config.Config) {
				c.Filesystem.MountPointsExclude, c.Filesystem.MountPointsInclude = "", "^/data"
			},
		},
		{
			args: []string{"--collector.filesystem.fs-types-include=^ext4$"},
			want: func(c *config.Config) { c.Filesystem.FSTypesExclude, c.Filesystem.FSTypesInclude = "", "^ext4$" },
		},
		{
			// Both flags on the command line: the collector rejects the pair.
			args: []string{"--collector.diskstats.device-include=^sd", "--collector.diskstats.device-exclude=^sdb"},
			want: func(c *config.Config) {
				c.Diskstats = config.DiskstatsConfig{DeviceInclude: "^sd", DeviceExclude: "^sdb"}
			},
		},
	} {
		t.Run(test.args[0], func(t *testing.T) {
			want := config.NewConfigWithDefaults()
			test.want(&want)
			if got := parseConfigFlags(t, test.args...); !reflect.DeepEqual(got, want) {
				t.Fatalf("got:  %+v\nwant: %+v", got, want)
			}
		})
	}
}

func TestFileDescriptorLeak(t *testing.T) {
	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		t.Skipf("proc filesystem is not available, but currently required to read number of open file descriptors: %s", err)
	}
	if _, err := fs.Stat(); err != nil {
		t.Errorf("unable to read process stats: %s", err)
	}

	exporter := exec.Command(binary, "--web.listen-address", address)
	test := func(pid int) error {
		if err := queryExporter(address); err != nil {
			return err
		}
		proc, err := procfs.NewProc(pid)
		if err != nil {
			return err
		}
		fdsBefore, err := proc.FileDescriptors()
		if err != nil {
			return err
		}
		for range 5 {
			if err := queryExporter(address); err != nil {
				return err
			}
		}
		fdsAfter, err := proc.FileDescriptors()
		if err != nil {
			return err
		}
		if want, have := len(fdsBefore), len(fdsAfter); want != have {
			return fmt.Errorf("want %d open file descriptors after metrics scrape, have %d", want, have)
		}
		return nil
	}

	if err := runCommandAndTests(exporter, address, test); err != nil {
		t.Error(err)
	}
}

func TestHandlingOfDuplicatedMetrics(t *testing.T) {
	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}

	exporter := exec.Command(binary, "--web.listen-address", address, "--collector.textfile.directory", "collector/fixtures/textfile/duplicate_metrics")
	test := func(_ int) error {
		return queryExporter(address)
	}

	if err := runCommandAndTests(exporter, address, test); err != nil {
		t.Error(err)
	}
}

func queryExporter(address string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", address))
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		return fmt.Errorf("want /metrics status code %d, have %d. Body:\n%s", want, have, b)
	}
	return nil
}

func runCommandAndTests(cmd *exec.Cmd, address string, fn func(pid int) error) error {
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %s", err)
	}
	time.Sleep(50 * time.Millisecond)
	for i := range 10 {
		if err := queryExporter(address); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
		if cmd.Process == nil || i == 9 {
			return fmt.Errorf("can't start command")
		}
	}

	errc := make(chan error)
	go func(pid int) {
		errc <- fn(pid)
	}(cmd.Process.Pid)

	err := <-errc
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	return err
}
