// Copyright 2020 The Prometheus Authors
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

// +build go1.14

package https

import (
	"crypto/tls"
	"strings"

	"github.com/pkg/errors"
)

func (c *cipher) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal((*string)(&s))
	if err != nil {
		return err
	}
	s = strings.ToUpper(s)
	for _, cs := range tls.CipherSuites() {
		if cs.Name == s {
			*c = (cipher)(cs.ID)
			return nil
		}
	}
	return errors.New("unknown cipher: " + s)
}

func (c cipher) MarshalYAML() (interface{}, error) {
	return tls.CipherSuiteName((uint16)(c)), nil
}
