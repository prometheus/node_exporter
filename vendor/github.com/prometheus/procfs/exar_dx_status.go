/*
 * @Author: David
 * @version:
 * @Description:
 * @Date: 2020-03-12 17:11:43
 * @LastEditors: David
 * @LastEditTime: 2020-03-12 18:13:57
 */
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
	"fmt"
	"io"
	"os"
	"strings"
)

// For the proc file format details,
// see https://elixir.bootlin.com/linux/v4.17/source/net/unix/af_unix.c#L2815
// and https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/net.h#L48.

// Constants for the various /proc/exar/dx_dev_status enumerations.

var (
	MaxStatusFields = 2
)

// ExarDx represents a line of /proc/exar/dx_dev_status.
type ExarDxStatus struct {
	Status map[string]string
}

// ExarDx returns data read from /proc/exar/dx_dev_status.
func (fs FS) ExarDxStatus() (*ExarDxStatus, error) {
	return readExarDxStatus(fs.proc.Path("exar/dx_dev_status"))
}

// readExarDx reads data in /proc/exar/dx_dev_status format from the specified file.
func readExarDxStatus(file string) (*ExarDxStatus, error) {
	// This file could be quite large and a streaming read is desirable versus
	// reading the entire contents at once.
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseExarDxStatus(f)
}

// parseExarDx creates a ExarDx structure from the incoming stream.
func parseExarDxStatus(r io.Reader) (*ExarDxStatus, error) {
	// Begin scanning by checking for the existence of Inode.
	s := bufio.NewScanner(r)

	edStatus := ExarDxStatus{
		Status: make(map[string]string),
	}
	for s.Scan() {
		line := s.Text()
		item, err := edStatus.parseLine(line)

		if err != nil {
			// return nil, fmt.Errorf("failed to parse /proc/exar/dx_dev_status data %q: %v", line, err)
			continue
		}
		edStatus.Status[item[0]] = item[1]

	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan /proc/exar/dx_dev_status data: %v", err)
	}

	return &edStatus, nil
}

func (e *ExarDxStatus) parseLine(line string) ([]string, error) {
	fields := strings.Split(strings.TrimSpace(line), "=")

	l := len(fields)
	if l != MaxStatusFields {
		return nil, fmt.Errorf("expected at least %d fields but got %d,=> %v", MaxStatusFields, l, fields)
	}

	// Field offsets are as follows:
	// device	ring	type	max_water_mark	water_mark	failedCmds	totalFinCmds

	return fields, nil
}
