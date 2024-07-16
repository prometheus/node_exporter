// Copyright 2022 The Prometheus Authors
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

//go:build linux && !noselinux
// +build linux,!noselinux

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getAVCStats(path string) (avcStats map[string]uint64, err error) {
	file, err := os.Open(sysFilePath(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseAVCStats(file)
}

func getAVCHashStats(path string) (avcStats map[string]uint64, err error) {
	file, err := os.Open(sysFilePath(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseAVCHashStats(file)
}

func parseAVCStats(r io.Reader) (stats map[string]uint64, err error) {
	avcValuesRE := regexp.MustCompile(`\d+`)
	stats = make(map[string]uint64)
	scanner := bufio.NewScanner(r)
	scanner.Scan() // Skip header

	for scanner.Scan() {
		avcValues := avcValuesRE.FindAllString(scanner.Text(), -1)

		if len(avcValues) != 6 { //
			return nil, fmt.Errorf("invalid AVC stat line: %s",
				scanner.Text())
		}

		lookups, err := strconv.ParseUint(avcValues[0], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for lookups")
		}

		hits, err := strconv.ParseUint(avcValues[1], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for hits")
		}

		misses, err := strconv.ParseUint(avcValues[2], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for misses")
		}

		allocations, err := strconv.ParseUint(avcValues[3], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for allocations")
		}

		reclaims, err := strconv.ParseUint(avcValues[4], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for reclaims")
		}

		frees, err := strconv.ParseUint(avcValues[5], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse expected integer value for frees")
		}

		stats["lookups"] += lookups
		stats["hits"] += hits
		stats["misses"] += misses
		stats["allocations"] += allocations
		stats["reclaims"] += reclaims
		stats["frees"] += frees
	}

	return stats, err
}

func parseAVCHashStats(r io.Reader) (stats map[string]uint64, err error) {
	stats = make(map[string]uint64)
	scanner := bufio.NewScanner(r)

	scanner.Scan()
	entriesValue := strings.TrimPrefix(scanner.Text(), "entries: ")

	scanner.Scan()
	bucketsValues := strings.Split(scanner.Text(), "buckets used: ")
	bucketsValuesTuple := strings.Split(bucketsValues[1], "/")

	scanner.Scan()
	longestChainValue := strings.TrimPrefix(scanner.Text(), "longest chain: ")

	stats["entries"], err = strconv.ParseUint(entriesValue, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse expected integer value for hash entries")
	}

	stats["buckets_used"], err = strconv.ParseUint(bucketsValuesTuple[0], 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse expected integer value for hash buckets used")
	}

	stats["buckets_available"], err = strconv.ParseUint(bucketsValuesTuple[1], 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse expected integer value for hash buckets available")
	}

	stats["longest_chain"], err = strconv.ParseUint(longestChainValue, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse expected integer value for hash longest chain")
	}

	return stats, err
}
