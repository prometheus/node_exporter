// Copyright 2016 The Prometheus Authors
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

package collector

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// zpool metrics

func (c *zfsCollector) parseZpoolOutput(reader io.Reader, handler func(string, string, float64)) (err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())
		if len(fields) != 4 {
			return fmt.Errorf("unexpected output of zpool command")
		}

		valueString := fields[2]
		switch {
		case strings.HasSuffix(fields[2], "%"):
			percentage := strings.TrimSuffix(fields[2], "%")
			valueString = "0." + percentage
		case strings.HasSuffix(fields[2], "x"):
			valueString = strings.TrimSuffix(fields[2], "x")
		}

		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return err
		}
		handler(fields[0], fields[1], value)

	}
	return scanner.Err()
}
