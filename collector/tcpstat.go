// +build !notcpstat

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	procTCPStat  = "/proc/net/tcp"
	procTCP6Stat = "/proc/net/tcp6"
)

type TCPConnectionState int

const (
	TCP_ESTABLISHED TCPConnectionState = iota + 1
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
)

type tcpStatCollector struct {
	metric *prometheus.GaugeVec
}

func init() {
	Factories["tcpstat"] = NewTCPStatCollector
}

// NewTCPStatCollector takes a returns
// a new Collector exposing network stats.
func NewTCPStatCollector() (Collector, error) {
	return &tcpStatCollector{
		metric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "tcp_connection_states",
				Help:      "Number of connection states.",
			},
			[]string{"state"},
		),
	}, nil
}

func (c *tcpStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	tcpStats, err := getTCPStats(procTCPStat)
	if err != nil {
		return fmt.Errorf("couldn't get tcpstats: %s", err)
	}

	// if enabled ipv6 system
	if _, hasIPv6 := os.Stat(procTCP6Stat); hasIPv6 == nil {
		tcp6Stats, err := getTCPStats(procTCP6Stat)
		if err != nil {
			return fmt.Errorf("couldn't get tcp6stats: %s", err)
		}

		for st, value := range tcp6Stats {
			tcpStats[st] += value
		}
	}

	for st, value := range tcpStats {
		c.metric.WithLabelValues(st.String()).Set(value)
	}

	c.metric.Collect(ch)
	return err
}

func getTCPStats(statsFile string) (map[TCPConnectionState]float64, error) {
	file, err := os.Open(statsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseTCPStats(file)
}

func parseTCPStats(r io.Reader) (map[TCPConnectionState]float64, error) {
	var (
		tcpStats = map[TCPConnectionState]float64{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		if strings.HasPrefix(parts[0], "sl") {
			continue
		}
		st, err := strconv.ParseInt(parts[3], 16, 8)
		if err != nil {
			return nil, err
		}

		tcpStats[TCPConnectionState(st)]++
	}

	return tcpStats, nil
}

func (st TCPConnectionState) String() string {
	switch st {
	case TCP_ESTABLISHED:
		return "established"
	case TCP_SYN_SENT:
		return "syn_sent"
	case TCP_SYN_RECV:
		return "syn_recv"
	case TCP_FIN_WAIT1:
		return "fin_wait1"
	case TCP_FIN_WAIT2:
		return "fin_wait2"
	case TCP_TIME_WAIT:
		return "time_wait"
	case TCP_CLOSE:
		return "close"
	case TCP_CLOSE_WAIT:
		return "close_wait"
	case TCP_LAST_ACK:
		return "last_ack"
	case TCP_LISTEN:
		return "listen"
	case TCP_CLOSING:
		return "closing"
	default:
		return "unknown"
	}
}
