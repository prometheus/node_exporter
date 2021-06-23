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
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

func readUintFromFile(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Take a []byte{} and return a string based on null termination.
// This is useful for situations where the OS has returned a null terminated
// string to use.
// If this function happens to receive a byteArray that contains no nulls, we
// simply convert the array to a string with no bounding.
func bytesToString(byteArray []byte) string {
	n := bytes.IndexByte(byteArray, 0)
	if n < 0 {
		return string(byteArray)
	}
	return string(byteArray[:n])
}
