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

package config

import "time"

const (
	DefaultCPUEnableGuest = true
)

// CPUConfig contains the options for the cpu collector.
type CPUConfig struct {
	// EnableGuest exposes node_cpu_guest_seconds_total, like
	// --collector.cpu.guest.
	EnableGuest bool
	// EnableInfo exposes node_cpu_info, like --collector.cpu.info.
	EnableInfo bool
	// FlagsInclude and BugsInclude are regular expressions that filter the
	// flags and bugs fields of node_cpu_info, like
	// --collector.cpu.info.flags-include and --collector.cpu.info.bugs-include.
	// Setting either one also enables node_cpu_info.
	FlagsInclude string
	BugsInclude  string
}

// TextfileConfig contains the options for the textfile collector.
type TextfileConfig struct {
	// Directories to read text files with metrics from, like
	// --collector.textfile.directory. Glob patterns are supported.
	Directories []string
}

const (
	DefaultNTPServer               = "127.0.0.1"
	DefaultNTPServerPort           = 123
	DefaultNTPProtocolVersion      = 4
	DefaultNTPIPTTL                = 1
	DefaultNTPMaxDistance          = 3466080 * time.Microsecond
	DefaultNTPLocalOffsetTolerance = time.Millisecond
)

// NTPConfig contains the options for the ntp collector.
type NTPConfig struct {
	// Server is --collector.ntp.server: NTP server to use for ntp collector.
	Server string
	// ServerPort is --collector.ntp.server-port: UDP port number to connect to on NTP server.
	ServerPort int
	// ProtocolVersion is --collector.ntp.protocol-version: NTP protocol version.
	ProtocolVersion int
	// ServerIsLocal is --collector.ntp.server-is-local: Certify that collector.ntp.server address is not a public ntp server.
	ServerIsLocal bool
	// IPTTL is --collector.ntp.ip-ttl: IP TTL to use while sending NTP query.
	IPTTL int
	// MaxDistance is --collector.ntp.max-distance: Max accumulated distance to the root.
	MaxDistance time.Duration
	// LocalOffsetTolerance is --collector.ntp.local-offset-tolerance: Offset between local clock and local ntpd time to tolerate.
	LocalOffsetTolerance time.Duration
}

// StatConfig contains the options for the stat collector.
type StatConfig struct {
	// Softirq is --collector.stat.softirq: Export softirq calls per vector.
	Softirq bool
}

const (
	DefaultVMStatFields = "^(oom_kill|pgpg|pswp|pg.*fault).*"
)

// VMStatConfig contains the options for the vmstat collector.
type VMStatConfig struct {
	// Fields is --collector.vmstat.fields: Regexp of fields to return for vmstat collector.
	Fields string
}

const (
	DefaultNetstatFields = "^(.*_(InErrors|InErrs)|Ip_Forwarding|Ip(6|Ext)_(InOctets|OutOctets)|Icmp6?_(InMsgs|OutMsgs)|TcpExt_(Listen.*|Syncookies.*|TCPSynRetrans|TCPTimeouts|TCPOFOQueue|TCPRcvQDrop)|Tcp_(ActiveOpens|InSegs|OutSegs|OutRsts|PassiveOpens|RetransSegs|CurrEstab)|Udp6?_(InDatagrams|OutDatagrams|NoPorts|RcvbufErrors|SndbufErrors))$"
)

// NetstatConfig contains the options for the netstat collector.
type NetstatConfig struct {
	// Fields is --collector.netstat.fields: Regexp of fields to return for netstat collector.
	Fields string
}

// BcacheConfig contains the options for the bcache collector.
type BcacheConfig struct {
	// PriorityStats is --collector.bcache.priorityStats: Expose expensive priority stats.
	PriorityStats bool
}

const (
	DefaultNVMeNamespaceUsage = true
)

// NVMeConfig contains the options for the nvme collector.
type NVMeConfig struct {
	// NamespaceUsage is --collector.nvme.namespace-usage: Expose namespace used bytes (reads nuse, which makes the kernel query the controller).
	NamespaceUsage bool
}

// RAPLConfig contains the options for the rapl collector.
type RAPLConfig struct {
	// EnableZoneLabel is --collector.rapl.enable-zone-label: Enables service unit metric unit_start_time_seconds.
	EnableZoneLabel bool
}
