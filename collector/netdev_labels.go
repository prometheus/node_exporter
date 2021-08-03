// Copyright 2021 The Prometheus Authors
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
	"io/ioutil"
	"runtime"
	"strings"
)

type label struct {
	key   string
	value string
}

type labels []label

func (l labels) keys() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.key)
	}

	return ret
}

func (l labels) values() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.value)
	}

	return ret
}

func getLabelsFromIfAlias(ifName string) labels {
	if runtime.GOOS != "linux" {
		return nil
	}

	if !labelsFromIfAlias {
		return nil
	}

	ifAliasBytes, err := ioutil.ReadFile("/sys/class/net/" + ifName + "/ifalias")
	if err != nil {
		return nil
	}

	ifAlias := strings.TrimSpace(string(ifAliasBytes))
	keyValueStrings := strings.Split(ifAlias, ",")
	ret := make(labels, 0, len(keyValueStrings))

	for _, kv := range keyValueStrings {
		parts := strings.Split(kv, "=")
		if len(parts) != 2 {
			continue
		}

		ret = append(ret, label{
			key:   parts[0],
			value: parts[1],
		})
	}

	return ret
}
