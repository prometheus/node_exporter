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

// Package collector includes all individual collectors to gather and export system metrics.
package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

func RegisterCollectorPublic(collector string, isDefaultEnabled bool, factory func(logger log.Logger) (Collector, error)) {
	registerCollector(collector, isDefaultEnabled, factory)
}

func ReplaceCollector(collector string, isDefaultEnabled bool, factory func(logger log.Logger) (Collector, error)) {
	delete(factories, collector)
	/*var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}*/

	flagName := fmt.Sprintf("collector.%s", collector)
	//flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)
	flagModel := kingpin.CommandLine.GetFlag(flagName)

	flag := flagModel.Default(defaultValue).Action(collectorFlagAction(collector)).Bool()
	collectorState[collector] = flag

	factories[collector] = factory
}
