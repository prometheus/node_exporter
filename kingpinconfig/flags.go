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

package kingpinconfig

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/node_exporter/collector"
)

func AddFlags(a *kingpin.Application) collector.NodeCollectorConfig {
	config := collector.NodeCollectorConfig{}

	config.Arp.DeviceInclude = a.Flag("collector.arp.device-include", "Regexp of arp devices to include (mutually exclusive to device-exclude).").String()
	config.Arp.DeviceExclude = a.Flag("collector.arp.device-exclude", "Regexp of arp devices to exclude (mutually exclusive to device-include).").String()

	return config
}
