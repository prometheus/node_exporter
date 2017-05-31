package procfs

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"os/exec"
	"bytes"
)

// IPVSStats holds IPVS statistics, as exposed by the kernel in `/proc/net/ip_vs_stats`.
type IPVSStats struct {
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

// IPVSBackendStatus holds current metrics of one virtual / real address pair.
type IPVSBackendStatus struct {
	// The local (virtual) IP address.
	LocalAddress net.IP
	// The local (virtual) port.
	LocalPort uint16
	// The local firewall mark
	LocalMark string
	// The transport protocol (TCP, UDP).
	Proto string
	// The remote (real) IP address.
	RemoteAddress net.IP
	// The remote (real) port.
	RemotePort uint16
	// The current number of active connections for this virtual/real address pair.
	ActiveConn uint64
	// The current number of inactive connections for this virtual/real address pair.
	InactConn uint64
	// The current weight of this virtual/real address pair.
	Weight uint64
	// Total incoming connections per second
	IncomingConnectionsPerSecond uint64
	// Total incoming packages per second
	IncomingPackgesPerSecond uint64
	// Total outgoing packges per second
	OutgoingPackgesPerSecond uint64
	// Total incomingBytes per second
	IncomingBytesPerSecond uint64
	// Total outgoingBytes per second
	OutgoingBytesPerSecond uint64


}

// NewIPVSStats reads the IPVS statistics.
func NewIPVSStats() (IPVSStats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return IPVSStats{}, err
	}

	return fs.NewIPVSStats()
}

// NewIPVSStats reads the IPVS statistics from the specified `proc` filesystem.
func (fs FS) NewIPVSStats() (IPVSStats, error) {
	file, err := os.Open(fs.Path("net/ip_vs_stats"))
	if err != nil {
		return IPVSStats{}, err
	}
	defer file.Close()

	return parseIPVSStats(file)
}

// parseIPVSStats performs the actual parsing of `ip_vs_stats`.
func parseIPVSStats(file io.Reader) (IPVSStats, error) {
	var (
		statContent []byte
		statLines   []string
		statFields  []string
		stats       IPVSStats
	)

	statContent, err := ioutil.ReadAll(file)
	if err != nil {
		return IPVSStats{}, err
	}

	statLines = strings.SplitN(string(statContent), "\n", 4)
	if len(statLines) != 4 {
		return IPVSStats{}, errors.New("ip_vs_stats corrupt: too short")
	}

	statFields = strings.Fields(statLines[2])
	if len(statFields) != 5 {
		return IPVSStats{}, errors.New("ip_vs_stats corrupt: unexpected number of fields")
	}

	stats.Connections, err = strconv.ParseUint(statFields[0], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.IncomingPackets, err = strconv.ParseUint(statFields[1], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.OutgoingPackets, err = strconv.ParseUint(statFields[2], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.IncomingBytes, err = strconv.ParseUint(statFields[3], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}
	stats.OutgoingBytes, err = strconv.ParseUint(statFields[4], 16, 64)
	if err != nil {
		return IPVSStats{}, err
	}

	return stats, nil
}

// NewIPVSBackendStatus reads and returns the status of all (virtual,real) server pairs.
func NewIPVSBackendStatus() ([]IPVSBackendStatus, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return []IPVSBackendStatus{}, err
	}

	return fs.NewIPVSBackendStatus()
}

// NewIPVSBackendStatus reads and returns the status of all (virtual,real) server pairs from the specified `proc` filesystem.
func (fs FS) NewIPVSBackendStatus() ([]IPVSBackendStatus, error) {
	file, err := os.Open(fs.Path("net/ip_vs"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// For CPS, PPS, BPS
	rate, err := exec.Command("/sbin/ipvsadm", "-ln", "--rate").Output()
	if err != nil {
		return nil, err
	}
	rates := bytes.NewReader(rate)

	return parseIPVSBackendStatus(file, rates)
}

func parseIPVSBackendStatus(file, rates io.Reader) ([]IPVSBackendStatus, error) {
	var (
		status       []IPVSBackendStatus
		scanner      = bufio.NewScanner(file)
		rateScanner  = bufio.NewScanner(rates)
		rateMap      = make(map[string]uint64)
		proto        string
		localMark    string
		localAddress net.IP
		localPort    uint16
		localInfo    string
		err          error
	)

	for rateScanner.Scan() {
		fields := strings.Fields(string(rateScanner.Text()))
		if len(fields) == 0 {
			continue
		}
		switch {
		case fields[0] == "IP" || fields[0] == "Prot" || fields[1] == "RemoteAddress:Port":
			continue
		case fields[0] == "TCP" || fields[0] == "UDP":
			localInfo = fields[0] + fields[1]
			if err != nil {
				return nil, err
			}
		case fields[0] == "->":
			cps, err := strconv.ParseUint(fields[2], 10, 64)
			if err != nil {
				return nil, err
			}
			inpps, err := strconv.ParseUint(fields[3], 10, 64)
			if err != nil {
				return nil, err
			}
			outpps, err := strconv.ParseUint(fields[4], 10, 64)
			if err != nil {
				return nil, err
			}
			inbps, err := strconv.ParseUint(fields[5], 10, 64)
			if err != nil {
				return nil, err
			}
			outbps, err := strconv.ParseUint(fields[6], 10, 64)
			if err != nil {
				return nil, err
			}

			rateMap[localInfo + fields[1] + "cps"] = cps
			rateMap[localInfo + fields[1] + "inpps"] = inpps
			rateMap[localInfo + fields[1] + "outpps"] = outpps
			rateMap[localInfo + fields[1] + "inbps"] = inbps
			rateMap[localInfo + fields[1] + "outbps"] = outbps
		}
	}

	for scanner.Scan() {
		fields := strings.Fields(string(scanner.Text()))
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

			cpsKey := proto + localAddress.String() + ":" + fmt.Sprint(localPort) + remoteAddress.String() + ":" + fmt.Sprint(remotePort) + "cps"
			inppsKey := proto + localAddress.String() + ":" + fmt.Sprint(localPort) + remoteAddress.String() + ":" + fmt.Sprint(remotePort) + "inpps"
			outppsKey := proto + localAddress.String() + ":" + fmt.Sprint(localPort) + remoteAddress.String() + ":" + fmt.Sprint(remotePort) + "outpps"
			inbpsKey := proto + localAddress.String() + ":" + fmt.Sprint(localPort) + remoteAddress.String() + ":" + fmt.Sprint(remotePort) + "inbps"
			outbpsKey := proto + localAddress.String() + ":" + fmt.Sprint(localPort) + remoteAddress.String() + ":" + fmt.Sprint(remotePort) + "outpbs"

			status = append(status, IPVSBackendStatus{
				LocalAddress:  localAddress,
				LocalPort:     localPort,
				LocalMark:     localMark,
				RemoteAddress: remoteAddress,
				RemotePort:    remotePort,
				Proto:         proto,
				Weight:        weight,
				ActiveConn:    activeConn,
				InactConn:     inactConn,
				IncomingConnectionsPerSecond: rateMap[cpsKey],
				IncomingPackgesPerSecond: rateMap[inppsKey],
				OutgoingPackgesPerSecond: rateMap[outppsKey],
				IncomingBytesPerSecond: rateMap[inbpsKey],
				OutgoingBytesPerSecond: rateMap[outbpsKey],
			})
		}
	}
	return status, nil
}

func parseIPPort(s string) (net.IP, uint16, error) {
	var (
		ip  net.IP
		err error
	)

	switch len(s) {
	case 13:
		ip, err = hex.DecodeString(s[0:8])
		if err != nil {
			return nil, 0, err
		}
	case 46:
		ip = net.ParseIP(s[1:40])
		if ip == nil {
			return nil, 0, fmt.Errorf("invalid IPv6 address: %s", s[1:40])
		}
	default:
		return nil, 0, fmt.Errorf("unexpected IP:Port: %s", s)
	}

	portString := s[len(s)-4:]
	if len(portString) != 4 {
		return nil, 0, fmt.Errorf("unexpected port string format: %s", portString)
	}
	port, err := strconv.ParseUint(portString, 16, 16)
	if err != nil {
		return nil, 0, err
	}

	return ip, uint16(port), nil
}
