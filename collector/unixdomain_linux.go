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

// +build !nounixdomain

package collector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// For the proc file format details,
// see https://elixir.bootlin.com/linux/v4.17/source/net/unix/af_unix.c#L2815
// and https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/net.h#L48
const (
	unixDomainRefCountIdx = 1

	unixDomainTypeIdx  = 4
	unixDomainFlagsIdx = 3
	unixDomainStateIdx = 5
)

const (
	unixDomainTypeStream    = 1
	unixDomainTypeDgram     = 2
	unixDomainTypeSeqpacket = 5

	unixDomainFlagListen = 1 << 16

	unixDomainStateUnconnected  = 1
	unixDomainStateConnecting   = 2
	unixDomainStateConnected    = 3
	unixDomainStateDisconnected = 4
)

var (
	errUnixDomainUnknownType  = errors.New("Unknown type")
	errUnixDomainUnknownState = errors.New("Unknown state")
)

type unixDomainCollector struct {
	connections *typedDesc
	users       *typedDesc
	// used to generate a zero value for labels combination we saw before but absent now
	knownLabelsCombs map[unixDomainLabelsComb]struct{}
}

type unixDomainLabelsComb struct {
	typ   string
	flags string
	state string
}

type unixDomainEntry struct {
	labelsComb  *unixDomainLabelsComb
	connections int64
	users       int64
}

func init() {
	registerCollector("unixdomain", defaultDisabled, NewUnixDomainCollector)
}

// NewUnixDomainCollector returns a new Collector exposing Unix domain stats.
func NewUnixDomainCollector() (Collector, error) {
	return &unixDomainCollector{
		connections: &typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "unixdomain", "connections"),
			"Number of Unix domain connections",
			[]string{"type", "flags", "state"}, nil,
		), prometheus.GaugeValue},
		users: &typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "unixdomain", "users"),
			"Number of Unix domain users(RefCount)",
			[]string{"type", "flags", "state"}, nil,
		), prometheus.GaugeValue},
		knownLabelsCombs: make(map[unixDomainLabelsComb]struct{}),
	}, nil
}

func (c *unixDomainCollector) Update(ch chan<- prometheus.Metric) error {
	f, err := os.Open(procFilePath("net/unix"))
	if err != nil {
		return fmt.Errorf("unable to open Unix domain proc file: %s", err)
	}

	defer f.Close()

	return c.process(f, ch)
}

func (c *unixDomainCollector) parseUsers(hexStr string) (int64, error) {
	return strconv.ParseInt(hexStr, 16, 32)
}

func (c *unixDomainCollector) parseType(hexStr string) (string, error) {
	typ, err := strconv.ParseInt(hexStr, 16, 16)
	if err != nil {
		return "", err
	}
	switch typ {
	case unixDomainTypeStream:
		return "stream", nil
	case unixDomainTypeDgram:
		return "dgram", nil
	case unixDomainTypeSeqpacket:
		return "seqpacket", nil
	}
	return "", errUnixDomainUnknownType
}

func (c *unixDomainCollector) parseFlags(hexStr string) (string, error) {
	flag, err := strconv.ParseInt(hexStr, 16, 32)
	if err != nil {
		return "", err
	}
	switch flag {
	case unixDomainFlagListen:
		return "accepton", nil
	default:
		return "default", nil
	}
}

func (c *unixDomainCollector) parseState(hexStr string) (string, error) {
	st, err := strconv.ParseInt(hexStr, 16, 8)
	if err != nil {
		return "", err
	}
	switch st {
	case unixDomainStateUnconnected:
		return "unconnected", nil
	case unixDomainStateConnecting:
		return "connecting", nil
	case unixDomainStateConnected:
		return "connected", nil
	case unixDomainStateDisconnected:
		return "disconnected", nil
	}
	return "", errUnixDomainUnknownState
}

func (c *unixDomainCollector) parseItem(line string) (*unixDomainEntry, error) {
	fields := strings.Fields(line)
	typ, err := c.parseType(fields[unixDomainTypeIdx])
	if err != nil {
		return nil, fmt.Errorf("Parse Unix domain type(%s) failed: %s", fields[unixDomainTypeIdx], err)
	}
	flags, err := c.parseFlags(fields[unixDomainFlagsIdx])
	if err != nil {
		return nil, fmt.Errorf("Parse Unix domain flags(%s) failed: %s", fields[unixDomainFlagsIdx], err)
	}
	state, err := c.parseState(fields[unixDomainStateIdx])
	if err != nil {
		return nil, fmt.Errorf("Parse Unix domain state(%s) failed: %s", fields[unixDomainStateIdx], err)
	}

	users, err := c.parseUsers(fields[unixDomainRefCountIdx])
	if err != nil {
		return nil, fmt.Errorf("Parse Unix domain ref count(%s) failed: %s", fields[unixDomainRefCountIdx], err)
	}
	return &unixDomainEntry{
		labelsComb: &unixDomainLabelsComb{
			typ:   typ,
			flags: flags,
			state: state,
		},
		users: users,
	}, nil
}

func (c *unixDomainCollector) process(reader io.Reader, ch chan<- prometheus.Metric) error {
	labelsCombConns := make(map[unixDomainLabelsComb]int64)
	labelsCombUsers := make(map[unixDomainLabelsComb]int64)
	scanner := bufio.NewScanner(reader)
	//omit the first line
	scanner.Scan()
	for scanner.Scan() {
		entry, err := c.parseItem(scanner.Text())
		if err != nil {
			return err
		}
		labelsCombConns[*entry.labelsComb]++
		labelsCombUsers[*entry.labelsComb] += entry.users
		c.knownLabelsCombs[*entry.labelsComb] = struct{}{}
	}

	for labelsComb := range c.knownLabelsCombs {
		conns := labelsCombConns[labelsComb]
		users := labelsCombUsers[labelsComb]

		ch <- c.connections.mustNewConstMetric(float64(conns), labelsComb.typ, labelsComb.flags, labelsComb.state)
		ch <- c.users.mustNewConstMetric(float64(users), labelsComb.typ, labelsComb.flags, labelsComb.state)
	}
	return nil
}
