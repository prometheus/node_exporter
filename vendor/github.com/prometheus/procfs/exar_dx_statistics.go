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
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// For the proc file format details,
// see https://elixir.bootlin.com/linux/v4.17/source/net/unix/af_unix.c#L2815
// and https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/net.h#L48.

// Constants for the various /proc/exar/dx_cmd_statistics enumerations.

var (
	ErrorTitleLine = errors.New("ExarDxLine: Title Line")
	MaxFields      = 7
)

// ExarDxLine represents a line of /proc/net/unix.
type ExarDxLine struct {
	Device       string
	Ring         string
	Type         string
	MaxWaterMark string
	WaterMark    string
	FailedCmds   uint64
	TotalFinCmds uint64
}

// ExarDx holds the data read from /proc/exar/dx_cmd_statistics.
type ExarDx struct {
	Rows []*ExarDxLine
}

// ExarDx returns data read from /proc/exar/dx_cmd_statistics.
func (fs FS) ExarDx() (*ExarDx, error) {
	return readExarDx(fs.proc.Path("exar/dx_cmd_statistics"))
}

// readExarDx reads data in /proc/exar/dx_cmd_statistics format from the specified file.
func readExarDx(file string) (*ExarDx, error) {
	// This file could be quite large and a streaming read is desirable versus
	// reading the entire contents at once.
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseExarDx(f)
}

// parseExarDx creates a ExarDx structure from the incoming stream.
func parseExarDx(r io.Reader) (*ExarDx, error) {
	// Begin scanning by checking for the existence of Inode.
	s := bufio.NewScanner(r)

	var ed ExarDx
	for s.Scan() {
		line := s.Text()
		item, err := ed.parseLine(line)

		if err == ErrorTitleLine {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse /proc/exar/dx_cmd_statistics data %q: %v", line, err)
		}

		ed.Rows = append(ed.Rows, item)
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan /proc/exar/dx_cmd_statistics data: %v", err)
	}

	return &ed, nil
}

func (e *ExarDx) parseLine(line string) (*ExarDxLine, error) {
	fields := strings.Fields(line)

	l := len(fields)
	if l != MaxFields {
		return nil, fmt.Errorf("expected at least %d fields but got %d", MaxFields, l)
	}

	// Field offsets are as follows:
	// device	ring	type	max_water_mark	water_mark	failedCmds	totalFinCmds

	Device := strings.TrimSpace(fields[0])
	if "device" == Device {
		return nil, ErrorTitleLine
	}
	Ring := strings.TrimSpace(fields[1])

	Type := strings.TrimSpace(fields[2])

	MaxWaterMark := strings.TrimSpace(fields[3])

	WaterMark := strings.TrimSpace(fields[4])

	FailedCmds, err := e.parseUint(strings.TrimSpace(fields[5]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse FailedCmds count(%s): %v", fields[5], err)
	}
	TotalFinCmds, err := e.parseUint(strings.TrimSpace(fields[6]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse TotalFinCmds count(%s): %v", fields[6], err)
	}
	n := &ExarDxLine{
		Device:       Device,
		Ring:         Ring,
		Type:         Type,
		MaxWaterMark: MaxWaterMark,
		WaterMark:    WaterMark,
		FailedCmds:   FailedCmds,
		TotalFinCmds: TotalFinCmds,
	}

	return n, nil
}
func (e ExarDx) parseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 16, 32)
}
