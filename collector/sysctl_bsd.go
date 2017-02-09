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

// +build freebsd dragonfly
// +build !nomeminfo

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

type bsdSysctlType uint8

// BSD-specific sysctl value types.  There is an impedience mismatch between
// native C types, e.g. int vs long, and the golang unix.Sysctl variables

const (
	// Default to uint32.
	bsdSysctlTypeUint32 bsdSysctlType = iota
	bsdSysctlTypeUint64
)

// Contains all the info needed to map a single bsd-sysctl to a prometheus value.
type bsdSysctl struct {
	// Prometheus name
	name string

	// Simple prometheus description
	description string

	// Cached prometheus description
	cached_descriptor *prometheus.Desc

	// Prometheus type
	valueType prometheus.ValueType

	// Convert sysctl to desired units
	conversion func(uint64) float64

	// Sysctl name
	mib string

	// Sysctl data-type
	dataType bsdSysctlType
}

func (b bsdSysctl) GetValue() (float64, error) {
	var tmp32 uint32
	var tmp64 uint64
	var err error

	switch b.dataType {
	case bsdSysctlTypeUint32:
		tmp32, err = unix.SysctlUint32(b.mib)
		tmp64 = uint64(tmp32)
	case bsdSysctlTypeUint64:
		tmp64, err = unix.SysctlUint64(b.mib)
	}
	if err != nil {
		return 0, err
	}

	if b.conversion != nil {
		return float64(b.conversion(tmp64)), nil
	}
	return float64(tmp64), nil
}

func (b bsdSysctl) GetDesc(subsystem string) *prometheus.Desc {
	if b.cached_descriptor == nil {
		b.cached_descriptor = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, b.name),
			b.description,
			nil, nil,
		)
	}

	return b.cached_descriptor
}
