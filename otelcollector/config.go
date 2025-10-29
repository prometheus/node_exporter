// Copyright The Prometheus Authors
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

package otelcollector

import "fmt"

type Config struct {
	DisableDefaults   bool
	EnableCollectors  []string
	ExcludeCollectors []string
}

func (c Config) Validate() error {
	if len(c.EnableCollectors) > 0 && len(c.ExcludeCollectors) > 0 {
		return fmt.Errorf("%q and %q can't be used at the same time", "EnableCollectors", "ExcludeCollectors")
	}
	return nil
}
