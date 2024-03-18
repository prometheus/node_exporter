// Copyright 2023 The Prometheus Authors
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

//go:build linux && !nokey_users
// +build linux,!nokey_users

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	keyUsersSubsystem = "keys"
)

var (
	// See https://man7.org/linux/man-pages/man7/keyrings.7.html
	keyUsersLineRegex = regexp.MustCompile(`^\s*(?P<uid>\d+):\s+(?P<usage>\d+)\s+(?P<nkeys>\d+)/(?P<nikeys>\d+)\s+(?P<qnkeys>\d+)/(?P<maxkeys>\d+)\s+(?P<qnbytes>\d+)/(?P<maxbytes>\d+)$`)
)

type keyUsersCollector struct {
	logger log.Logger
}

type keyUsersEntry struct {
	uid  int64
	name string
}

func init() {
	registerCollector("key_users", defaultDisabled, NewKeyUsersCollector)
}

func NewKeyUsersCollector(logger log.Logger) (Collector, error) {
	return &keyUsersCollector{logger}, nil
}

func (c *keyUsersCollector) Update(ch chan<- prometheus.Metric) error {
	keyUsers, err := c.getKeyUsers()
	if err != nil {
		return fmt.Errorf("couldn't get key-users: %w", err)
	}
	level.Debug(c.logger).Log("msg", "Set node_keys", "keyUsers", keyUsers)
	for k, v := range keyUsers {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, keyUsersSubsystem, k.name),
				fmt.Sprintf("Key utilization field %s.", k.name),
				[]string{"uid"}, nil,
			),
			prometheus.GaugeValue, v, strconv.FormatInt(int64(k.uid), 10),
		)
	}
	return nil
}

func (c *keyUsersCollector) getKeyUsers() (map[keyUsersEntry]float64, error) {
	file, err := os.Open(procFilePath("key-users"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseKeyUsers(file)
}

func parseKeyUsers(r io.Reader) (map[keyUsersEntry]float64, error) {
	var scanner = bufio.NewScanner(r)
	var keyUsersInfo = map[keyUsersEntry]float64{}

	for scanner.Scan() {
		line := scanner.Text()
		match := keyUsersLineRegex.FindStringSubmatch(line)
		names := keyUsersLineRegex.SubexpNames()
		if len(match) != len(names) {
			return nil, fmt.Errorf("invalid line in key-users: %s", line)
		}
		uid, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid UID in key-users: %s", match[0])
		}
		for i := 2; i < len(names); i++ {
			v, err := strconv.ParseInt(match[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid metric '%s' in key-users: %s", names[i], match[i])
			}
			keyUsersInfo[keyUsersEntry{uid, names[i]}] = float64(v)
		}
	}
	return keyUsersInfo, scanner.Err()
}
