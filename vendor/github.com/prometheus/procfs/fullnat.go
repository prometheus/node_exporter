// Copyright 2018 The Prometheus Authors
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

package procfs

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// FNATStats holds FNAT statistics, as exposed by the kernel in `/proc/net/ip_vs_stats`.

type FNATStats struct {
	Stat []FNATStatsPerCpu
}
type FNATStatsPerCpu struct {
	Cpu string
	// Total count of connections.
	Connections uint64
	// Total incoming packages processed.
	IncomingPackets uint64
	// Total outgoing packages processed.
	OutgoingPackets uint64
	// Total incoming traffic.
	IncomingBytes uint64
	// Total outgoing traffic.
	OutgoingBytes uint64
}

type ExtStatsPerCpu = map[string]uint64
type FNATExtStats struct {
	FullnatAddToaOk             ExtStatsPerCpu
	FullnatAddToaFailLen        ExtStatsPerCpu
	FullnatAddToaHeadFull       ExtStatsPerCpu
	FullnatAddToaFailMem        ExtStatsPerCpu
	FullnatAddToaFailProto      ExtStatsPerCpu
	FullnatConnReused           ExtStatsPerCpu
	FullnatConnReusedClose      ExtStatsPerCpu
	FullnatConnReusedTimewait   ExtStatsPerCpu
	FullnatConnReusedFinwait    ExtStatsPerCpu
	FullnatConnReusedClosewait  ExtStatsPerCpu
	FullnatConnReusedLastack    ExtStatsPerCpu
	FullnatConnReusedEstab      ExtStatsPerCpu
	SynproxyRsError             ExtStatsPerCpu
	SynproxyNullAck             ExtStatsPerCpu
	SynproxyBadAck              ExtStatsPerCpu
	SynproxyOkAck               ExtStatsPerCpu
	SynproxySynCnt              ExtStatsPerCpu
	SynproxyAckstorm            ExtStatsPerCpu
	SynproxySynsendQlen         ExtStatsPerCpu
	SynproxyConnReused          ExtStatsPerCpu
	SynproxyConnReusedClose     ExtStatsPerCpu
	SynproxyConnReusedTimewait  ExtStatsPerCpu
	SynproxyConnReusedFinwait   ExtStatsPerCpu
	SynproxyConnReusedClosewait ExtStatsPerCpu
	SynproxyConnReusedLastack   ExtStatsPerCpu
	DefenceIpFragDrop           ExtStatsPerCpu
	DefenceIpFragGather         ExtStatsPerCpu
	DefenceTcpDrop              ExtStatsPerCpu
	DefenceUdpDrop              ExtStatsPerCpu
	FastXmitReject              ExtStatsPerCpu
	FastXmitPass                ExtStatsPerCpu
	FastXmitSkbCopy             ExtStatsPerCpu
	FastXmitNoMac               ExtStatsPerCpu
	FastXmitSynproxySave        ExtStatsPerCpu
	FastXmitDevLost             ExtStatsPerCpu
	RstInSynSent                ExtStatsPerCpu
	RstOutSynSent               ExtStatsPerCpu
	RstInEstablished            ExtStatsPerCpu
	RstOutEstablished           ExtStatsPerCpu
	GroPass                     ExtStatsPerCpu
	LroReject                   ExtStatsPerCpu
	XmitUnexpectedMtu           ExtStatsPerCpu
	ConnSchedUnreach            ExtStatsPerCpu
}

// FNATBackendStatus holds current metrics of one virtual / real address pair.
type FNATBackendStatus struct {
	// The local (virtual) IP address.
	LocalAddress net.IP
	// The remote (real) IP address.
	RemoteAddress net.IP
	// The local (virtual) port.
	LocalPort uint16
	// The remote (real) port.
	RemotePort uint16
	// The local firewall mark
	LocalMark string
	// The transport protocol (TCP, UDP).
	Proto string
	// The current number of active connections for this virtual/real address pair.
	ActiveConn uint64
	// The current number of inactive connections for this virtual/real address pair.
	InactConn uint64
	// The current weight of this virtual/real address pair.
	Weight uint64
}

// NewFNATStats reads the FNAT statistics.
func NewFNATStats() (FNATStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return FNATStats{}, err
	}

	return fs.NewFNATStats()
}

// NewFNATStats reads the FNAT statistics from the specified `proc` filesystem.
func (fs FS) NewFNATStats() (FNATStats, error) {
	file, err := os.Open(fs.Path("net/ip_vs_stats"))
	if err != nil {
		return FNATStats{}, err
	}
	defer file.Close()

	return parseFNATStats(file)
}

// parseFNATStats performs the actual parsing of `ip_vs_stats`.
func parseFNATStats(file io.Reader) (FNATStats, error) {
	var (
		statContent []byte
		statLines   []string
		statFields  []string
		stats       FNATStats = FNATStats{
			Stat: make([]FNATStatsPerCpu, 0),
		}
	)

	statContent, err := ioutil.ReadAll(file)
	if err != nil {
		return FNATStats{}, err
	}

	statLines = strings.Split(string(statContent), "\n")
	if len(statLines) < 4 {
		return FNATStats{}, errors.New("ip_vs_stats corrupt: too short")
	}
	for i := 2; i < len(statLines); i++ {
		statsPerCpu := FNATStatsPerCpu{}

		statFieldsAll := strings.Split(statLines[i], ":")
		if len(statFieldsAll) < 2 {
			continue
		}
		statsPerCpu.Cpu = strings.Replace(statFieldsAll[0], " ", "", -1)
		statFields = strings.Fields(statFieldsAll[1])
		if len(statFields) != 5 {
			return FNATStats{}, errors.New("ip_vs_stats corrupt: unexpected number of fields")
		}

		statsPerCpu.Connections, err = strconv.ParseUint(statFields[0], 16, 64)
		if err != nil {
			return FNATStats{}, err
		}
		statsPerCpu.IncomingPackets, err = strconv.ParseUint(statFields[1], 16, 64)
		if err != nil {
			return FNATStats{}, err
		}
		statsPerCpu.OutgoingPackets, err = strconv.ParseUint(statFields[2], 16, 64)
		if err != nil {
			return FNATStats{}, err
		}
		statsPerCpu.IncomingBytes, err = strconv.ParseUint(statFields[3], 16, 64)
		if err != nil {
			return FNATStats{}, err
		}
		statsPerCpu.OutgoingBytes, err = strconv.ParseUint(statFields[4], 16, 64)
		if err != nil {
			return FNATStats{}, err
		}
		stats.Stat = append(stats.Stat, statsPerCpu)
	}

	return stats, nil
}

// NewFNATExtStats reads the FNAT ext statistics.
func NewFNATExtStats() (FNATExtStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return FNATExtStats{}, err
	}

	return fs.NewFNATExtStats()
}

// NewFNATStats reads the FNAT  ext statistics from the specified `proc` filesystem.
func (fs FS) NewFNATExtStats() (FNATExtStats, error) {
	file, err := os.Open(fs.Path("net/ip_vs_ext_stats"))
	if err != nil {
		return FNATExtStats{}, err
	}
	defer file.Close()

	return parseFNATExtStats(file)
}

func parseFNATExtStats(file io.Reader) (FNATExtStats, error) {
	var (
		status  FNATExtStats = FNATExtStats{}
		title   []string     = make([]string, 0)
		scanner              = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) == 1 {
			title = strings.Fields(line[0])
			continue
		}
		fields := strings.Fields(line[1])

		if len(fields) == 0 {
			continue
		}

		switch {
		case strings.TrimSpace(line[0]) == "fullnat_add_toa_ok":
			status.FullnatAddToaOk = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_add_toa_fail_len":
			status.FullnatAddToaFailLen = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_add_toa_head_full":
			status.FullnatAddToaHeadFull = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_add_toa_fail_mem":
			status.FullnatAddToaFailMem = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_add_toa_fail_proto":
			status.FullnatAddToaFailProto = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused":
			status.FullnatConnReused = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_close":
			status.FullnatConnReusedClose = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_timewait":
			status.FullnatConnReusedTimewait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_finwait":
			status.FullnatConnReusedFinwait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_closewait":
			status.FullnatConnReusedClosewait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_lastack":
			status.FullnatConnReusedLastack = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fullnat_conn_reused_estab":
			status.FullnatConnReusedEstab = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_rs_error":
			status.SynproxyRsError = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_null_ack":
			status.SynproxyNullAck = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_bad_ack":
			status.SynproxyBadAck = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_ok_ack":
			status.SynproxyOkAck = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_syn_cnt":
			status.SynproxySynCnt = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_ackstorm":
			status.SynproxyAckstorm = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_synsend_qlen":
			status.SynproxySynsendQlen = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused":
			status.SynproxyConnReused = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused_close":
			status.SynproxyConnReusedClose = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused_timewait":
			status.SynproxyConnReusedTimewait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused_finwait":
			status.SynproxyConnReusedFinwait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused_closewait":
			status.SynproxyConnReusedClosewait = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "synproxy_conn_reused_lastack":
			status.SynproxyConnReusedLastack = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "defence_ip_frag_drop":
			status.DefenceIpFragDrop = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "defence_ip_frag_gather":
			status.DefenceIpFragGather = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "defence_tcp_drop":
			status.DefenceTcpDrop = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "defence_udp_drop":
			status.DefenceUdpDrop = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_reject":
			status.FastXmitReject = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_pass":
			status.FastXmitPass = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_skb_copy":
			status.FastXmitSkbCopy = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_no_mac":
			status.FastXmitNoMac = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_synproxy_save":
			status.FastXmitSynproxySave = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "fast_xmit_dev_lost":
			status.FastXmitDevLost = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "rst_in_syn_sent":
			status.RstInSynSent = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "rst_out_syn_sent":
			status.RstOutSynSent = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "rst_in_established":
			status.RstInEstablished = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "rst_out_established":
			status.RstOutEstablished = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "gro_pass":
			status.GroPass = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "lro_reject":
			status.LroReject = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "xmit_unexpected_mtu":
			status.XmitUnexpectedMtu = initExtStatsPerCpu(title, fields)
		case strings.TrimSpace(line[0]) == "conn_sched_unreach":
			status.ConnSchedUnreach = initExtStatsPerCpu(title, fields)

		}

	}
	return status, nil
}

func initExtStatsPerCpu(title, f []string) ExtStatsPerCpu {
	ret := make(ExtStatsPerCpu)
	var err error
	for k, v := range f {
		ret[title[k]], err = strconv.ParseUint(v, 16, 64)
		if err != nil {
			continue
		}
	}
	return ret
}

// NewFNATBackendStatus reads and returns the status of all (virtual,real) server pairs.
func NewFNATBackendStatus() ([]FNATBackendStatus, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return []FNATBackendStatus{}, err
	}

	return fs.NewFNATBackendStatus()
}

// NewFNATBackendStatus reads and returns the status of all (virtual,real) server pairs from the specified `proc` filesystem.
func (fs FS) NewFNATBackendStatus() ([]FNATBackendStatus, error) {
	file, err := os.Open(fs.Path("net/ip_vs"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseFNATBackendStatus(file)
}

func parseFNATBackendStatus(file io.Reader) ([]FNATBackendStatus, error) {
	var (
		status       []FNATBackendStatus
		scanner      = bufio.NewScanner(file)
		proto        string
		localMark    string
		localAddress net.IP
		localPort    uint16
		err          error
	)

	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		switch {
		case fields[0] == "IP" || fields[0] == "Prot" || fields[1] == "RemoteAddress:Port":
			continue
		case fields[0] == "TCP" || fields[0] == "UDP":
			if len(fields) < 2 {
				continue
			}
			proto = fields[0]
			localMark = ""
			localAddress, localPort, err = parseIPPort(fields[1])
			if err != nil {
				return nil, err
			}
		case fields[0] == "FWM":
			if len(fields) < 2 {
				continue
			}
			proto = fields[0]
			localMark = fields[1]
			localAddress = nil
			localPort = 0
		case fields[0] == "->":
			if len(fields) < 6 {
				continue
			}
			remoteAddress, remotePort, err := parseIPPort(fields[1])
			if err != nil {
				return nil, err
			}
			weight, err := strconv.ParseUint(fields[3], 10, 64)
			if err != nil {
				return nil, err
			}
			activeConn, err := strconv.ParseUint(fields[4], 10, 64)
			if err != nil {
				return nil, err
			}
			inactConn, err := strconv.ParseUint(fields[5], 10, 64)
			if err != nil {
				return nil, err
			}
			status = append(status, FNATBackendStatus{
				LocalAddress:  localAddress,
				LocalPort:     localPort,
				LocalMark:     localMark,
				RemoteAddress: remoteAddress,
				RemotePort:    remotePort,
				Proto:         proto,
				Weight:        weight,
				ActiveConn:    activeConn,
				InactConn:     inactConn,
			})
		}
	}
	return status, nil
}

//func parseIPPort(s string) (net.IP, uint16, error) {
//	var (
//		ip  net.IP
//		err error
//	)
//
//	switch len(s) {
//	case 13:
//		ip, err = hex.DecodeString(s[0:8])
//		if err != nil {
//			return nil, 0, err
//		}
//	case 46:
//		ip = net.ParseIP(s[1:40])
//		if ip == nil {
//			return nil, 0, fmt.Errorf("invalid IPv6 address: %s", s[1:40])
//		}
//	default:
//		return nil, 0, fmt.Errorf("unexpected IP:Port: %s", s)
//	}
//
//	portString := s[len(s)-4:]
//	if len(portString) != 4 {
//		return nil, 0, fmt.Errorf("unexpected port string format: %s", portString)
//	}
//	port, err := strconv.ParseUint(portString, 16, 16)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	return ip, uint16(port), nil
//}
//
