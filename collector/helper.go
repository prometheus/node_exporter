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

package collector

import (
	"fmt"
	"strconv"
	"strings"
)

func splitToInts(str string, sep string) (ints []int, err error) {
	for _, part := range strings.Split(str, sep) {
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Could not split '%s' because %s is no int: %s", str, part, err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}
