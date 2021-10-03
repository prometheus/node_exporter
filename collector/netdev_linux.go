// Copyright 2015 The Prometheus Authors
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

//go:build !nonetdev
// +build !nonetdev

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	procNetDevInterfaceRE = regexp.MustCompile(`^(.+): *(.+)$`)
	procNetDevFieldSep    = regexp.MustCompile(` +`)
)

func getNetDevStats(filter *netDevFilter, logger log.Logger) (netDevStats, error) {
	file, err := os.Open(procFilePath("net/dev"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetDevStats(file, filter, logger)
}

func parseNetDevStats(r io.Reader, filter *netDevFilter, logger log.Logger) (netDevStats, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan() // skip first header
	scanner.Scan()
	parts := strings.Split(scanner.Text(), "|")
	if len(parts) != 3 { // interface + receive + transmit
		return nil, fmt.Errorf("invalid header line in net/dev: %s",
			scanner.Text())
	}

	receiveHeader := strings.Fields(parts[1])
	transmitHeader := strings.Fields(parts[2])
	headerLength := len(receiveHeader) + len(transmitHeader)

	netDev := netDevStats{}
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		parts := procNetDevInterfaceRE.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("couldn't get interface name, invalid line in net/dev: %q", line)
		}

		dev := parts[1]
		if filter.ignored(dev) {
			level.Debug(logger).Log("msg", "Ignoring device", "device", dev)
			continue
		}

		values := procNetDevFieldSep.Split(strings.TrimLeft(parts[2], " "), -1)
		if len(values) != headerLength {
			return nil, fmt.Errorf("couldn't get values, invalid line in net/dev: %q", parts[2])
		}

		devStats := map[string]uint64{}
		addStats := func(key, value string) {
			v, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				level.Debug(logger).Log("msg", "invalid value in netstats", "key", key, "value", value, "err", err)
				return
			}

			devStats[key] = v
		}

		for i := 0; i < len(receiveHeader); i++ {
			addStats("receive_"+receiveHeader[i], values[i])
		}

		for i := 0; i < len(transmitHeader); i++ {
			addStats("transmit_"+transmitHeader[i], values[i+len(receiveHeader)])
		}

		netDev[dev] = devStats
	}
	return netDev, scanner.Err()
}
