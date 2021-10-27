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

//go:build !notherm
// +build !notherm

package collector

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <stdio.h>
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <IOKit/pwr_mgt/IOPM.h>

struct ref_with_ret {
    CFDictionaryRef ref;
    IOReturn ret;
};

struct ref_with_ret FetchThermal();

struct ref_with_ret FetchThermal() {
    CFDictionaryRef ref;
    IOReturn ret;
    ret = IOPMCopyCPUPowerStatus(&ref);
    struct ref_with_ret result = {
            ref,
            ret,
    };
    return result;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"unsafe"
)

type thermCollector struct {
	cpuSchedulerLimit typedDesc
	cpuAvailableCPU   typedDesc
	cpuSpeedLimit     typedDesc
	logger            log.Logger
}

const thermal = "thermal"

func init() {
	registerCollector(thermal, defaultEnabled, NewThermCollector)
}

// NewThermCollector returns a new Collector exposing current CPU power levels.
func NewThermCollector(logger log.Logger) (Collector, error) {
	return &thermCollector{
		cpuSchedulerLimit: typedDesc{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, thermal, "cpu_scheduler_limit_ratio"),
				"Represents the percentage (0-100) of CPU time available. 100% at normal operation. The OS may limit this time for a percentage less than 100%.",
				nil,
				nil),
			valueType: prometheus.GaugeValue,
		},
		cpuAvailableCPU: typedDesc{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, thermal, "cpu_available_cpu"),
				"Reflects how many, if any, CPUs have been taken offline. Represented as an integer number of CPUs (0 - Max CPUs).",
				nil,
				nil,
			),
			valueType: prometheus.GaugeValue,
		},
		cpuSpeedLimit: typedDesc{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, thermal, "cpu_speed_limit_ratio"),
				"Defines the speed & voltage limits placed on the CPU. Represented as a percentage (0-100) of maximum CPU speed.",
				nil,
				nil,
			),
			valueType: prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

func (c *thermCollector) Update(ch chan<- prometheus.Metric) error {
	cpuPowerStatus, err := fetchCPUPowerStatus()
	if err != nil {
		return err
	}
	if value, ok := cpuPowerStatus[(string(C.kIOPMCPUPowerLimitSchedulerTimeKey))]; ok {
		ch <- c.cpuSchedulerLimit.mustNewConstMetric(float64(value) / 100.0)
	}
	if value, ok := cpuPowerStatus[(string(C.kIOPMCPUPowerLimitProcessorCountKey))]; ok {
		ch <- c.cpuAvailableCPU.mustNewConstMetric(float64(value))
	}
	if value, ok := cpuPowerStatus[(string(C.kIOPMCPUPowerLimitProcessorSpeedKey))]; ok {
		ch <- c.cpuSpeedLimit.mustNewConstMetric(float64(value) / 100.0)
	}
	return nil
}

func fetchCPUPowerStatus() (map[string]int, error) {
	cfDictRef, _ := C.FetchThermal()
	defer func() {
		C.CFRelease(C.CFTypeRef(cfDictRef.ref))
	}()

	if C.kIOReturnNotFound == cfDictRef.ret {
		return nil, errors.New("no CPU power status has been recorded")
	}

	if C.kIOReturnSuccess != cfDictRef.ret {
		return nil, fmt.Errorf("no CPU power status with error code 0x%08x", int(cfDictRef.ret))
	}

	// mapping CFDictionary to map
	cfDict := CFDict(cfDictRef.ref)
	return mappingCFDictToMap(cfDict), nil
}

type CFDict uintptr

func mappingCFDictToMap(dict CFDict) map[string]int {
	if C.CFNullRef(dict) == C.kCFNull {
		return nil
	}
	cfDict := C.CFDictionaryRef(dict)

	var result map[string]int
	count := C.CFDictionaryGetCount(cfDict)
	if count > 0 {
		keys := make([]C.CFTypeRef, count)
		values := make([]C.CFTypeRef, count)
		C.CFDictionaryGetKeysAndValues(cfDict, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])), (*unsafe.Pointer)(unsafe.Pointer(&values[0])))
		result = make(map[string]int, count)
		for i := C.CFIndex(0); i < count; i++ {
			result[mappingCFStringToString(C.CFStringRef(keys[i]))] = mappingCFNumberLongToInt(C.CFNumberRef(values[i]))
		}
	}
	return result
}

// CFStringToString converts a CFStringRef to a string.
func mappingCFStringToString(s C.CFStringRef) string {
	p := C.CFStringGetCStringPtr(s, C.kCFStringEncodingUTF8)
	if p != nil {
		return C.GoString(p)
	}
	length := C.CFStringGetLength(s)
	if length == 0 {
		return ""
	}
	maxBufLen := C.CFStringGetMaximumSizeForEncoding(length, C.kCFStringEncodingUTF8)
	if maxBufLen == 0 {
		return ""
	}
	buf := make([]byte, maxBufLen)
	var usedBufLen C.CFIndex
	_ = C.CFStringGetBytes(s, C.CFRange{0, length}, C.kCFStringEncodingUTF8, C.UInt8(0), C.false, (*C.UInt8)(&buf[0]), maxBufLen, &usedBufLen)
	return string(buf[:usedBufLen])
}

func mappingCFNumberLongToInt(n C.CFNumberRef) int {
	typ := C.CFNumberGetType(n)
	var long C.long
	C.CFNumberGetValue(n, typ, unsafe.Pointer(&long))
	return int(long)
}
