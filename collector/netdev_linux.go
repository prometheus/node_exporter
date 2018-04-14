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

// +build !nonetdev

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/prometheus/common/log"
)

var (
	procNetDevFieldSep = regexp.MustCompile("[ :] *")
)

func getNetDevStats(ignore *regexp.Regexp) (map[string]map[string]string, error) {
	file, err := os.Open(procFilePath("net/dev"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseNetDevStats(file, ignore)
}

func parseNetDevStats(r io.Reader, ignore *regexp.Regexp) (map[string]map[string]string, error) {
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
	headerLength := len(receiveHeader)+len(transmitHeader)+1

	netDev := map[string]map[string]string{}
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		parts := procNetDevFieldSep.Split(line, -1)
		if len(parts) != headerLength {
			return nil, fmt.Errorf("invalid line in net/dev: %s", scanner.Text())
		}

		dev := parts[0][:len(parts[0])]
		if ignore.MatchString(dev) {
			log.Debugf("Ignoring device: %s", dev)
			continue
		}
		netDev[dev] = map[string]string{}
                for i := 0; i < len(receiveHeader); i++ {
                        netDev[dev]["receive_"+receiveHeader[i]] = parts[i+1]
                }

                for i := 0; i < len(transmitHeader); i++ {
                        netDev[dev]["transmit_"+transmitHeader[i]] = parts[i+1+len(receiveHeader)]
                }
	}
	return netDev, scanner.Err()
}
