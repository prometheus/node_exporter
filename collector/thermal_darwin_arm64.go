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

//go:build !notherm && darwin && arm64

package collector

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>
#include <stdlib.h>

typedef struct __IOHIDEventSystemClient * IOHIDEventSystemClientRef;
typedef struct __IOHIDServiceClient * IOHIDServiceClientRef;
typedef struct __IOHIDEvent * IOHIDEventRef;

#define kIOHIDEventTypeTemperature 15
#define IOHIDEventFieldBase(type) (type << 16)

int32_t GetIOHIDEventFieldBase(int32_t type) {
    return IOHIDEventFieldBase(type);
}

// External functions
IOHIDEventSystemClientRef IOHIDEventSystemClientCreate(CFAllocatorRef allocator);
void IOHIDEventSystemClientSetMatching(IOHIDEventSystemClientRef client, CFDictionaryRef match);
CFArrayRef IOHIDEventSystemClientCopyServices(IOHIDEventSystemClientRef client);
IOHIDEventRef IOHIDServiceClientCopyEvent(IOHIDServiceClientRef service, int64_t type, int32_t options, int64_t timestamp);
double IOHIDEventGetFloatValue(IOHIDEventRef event, int32_t field);
CFTypeRef IOHIDServiceClientCopyProperty(IOHIDServiceClientRef service, CFStringRef key);
uint64_t IOHIDServiceClientGetRegistryID(IOHIDServiceClientRef service);
*/
import "C"

import (
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

const absoluteZeroCelsius = -273.15

// A thermal sensor is labelled with its IOHID "Product" property, which is not
// unique. Apple Silicon reports several services under one product name: some
// are distinct sensors that only share the name and can be told apart by their
// location, while others share the location too and differ only by the registry
// ID of the service.
type thermalSensor struct {
	name       string
	location   string
	registryID string
	temp       float64
}

func (c *thermCollector) updateTemperatures(ch chan<- prometheus.Metric) error {
	client := C.IOHIDEventSystemClientCreate(C.kCFAllocatorDefault)
	if client == nil {
		return nil
	}
	defer C.CFRelease(C.CFTypeRef(unsafe.Pointer(client)))

	page := 0xff00
	usage := 5

	pageNum := C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberIntType, unsafe.Pointer(&page))
	defer C.CFRelease(C.CFTypeRef(pageNum))
	usageNum := C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberIntType, unsafe.Pointer(&usage))
	defer C.CFRelease(C.CFTypeRef(usageNum))

	keyPage := C.CString("PrimaryUsagePage")
	defer C.free(unsafe.Pointer(keyPage))
	keyUsage := C.CString("PrimaryUsage")
	defer C.free(unsafe.Pointer(keyUsage))

	cfKeyPage := C.CFStringCreateWithCString(C.kCFAllocatorDefault, keyPage, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfKeyPage))
	cfKeyUsage := C.CFStringCreateWithCString(C.kCFAllocatorDefault, keyUsage, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfKeyUsage))

	keys := []C.CFTypeRef{C.CFTypeRef(cfKeyPage), C.CFTypeRef(cfKeyUsage)}
	values := []C.CFTypeRef{C.CFTypeRef(pageNum), C.CFTypeRef(usageNum)}

	matching := C.CFDictionaryCreate(C.kCFAllocatorDefault,
		(*unsafe.Pointer)(unsafe.Pointer(&keys[0])),
		(*unsafe.Pointer)(unsafe.Pointer(&values[0])),
		2,
		&C.kCFTypeDictionaryKeyCallBacks,
		&C.kCFTypeDictionaryValueCallBacks)
	defer C.CFRelease(C.CFTypeRef(matching))

	C.IOHIDEventSystemClientSetMatching(client, matching)

	services := C.IOHIDEventSystemClientCopyServices(client)
	if services == 0 {
		return nil
	}
	defer C.CFRelease(C.CFTypeRef(services))

	count := C.CFArrayGetCount(services)

	prodKey := C.CString("Product")
	defer C.free(unsafe.Pointer(prodKey))
	cfProdKey := C.CFStringCreateWithCString(C.kCFAllocatorDefault, prodKey, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfProdKey))

	locKey := C.CString("LocationID")
	defer C.free(unsafe.Pointer(locKey))
	cfLocKey := C.CFStringCreateWithCString(C.kCFAllocatorDefault, locKey, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfLocKey))

	// Read every sensor first, so that colliding product names can be detected
	// before any metric is emitted.
	sensors := make([]thermalSensor, 0, int(count))
	for i := 0; i < int(count); i++ {
		service := (C.IOHIDServiceClientRef)(C.CFArrayGetValueAtIndex(services, C.CFIndex(i)))

		event := C.IOHIDServiceClientCopyEvent(service, C.kIOHIDEventTypeTemperature, 0, 0)
		if event == nil {
			continue
		}

		temp := C.IOHIDEventGetFloatValue(event, C.GetIOHIDEventFieldBase(C.kIOHIDEventTypeTemperature))
		C.CFRelease(C.CFTypeRef(unsafe.Pointer(event)))

		// Observed invalid values on some Apple Silicon devices are around -9200.
		// Filter out physically impossible temperatures.
		if temp < absoluteZeroCelsius {
			continue
		}

		nameRef := C.IOHIDServiceClientCopyProperty(service, cfProdKey)
		name := "Unknown"
		if nameRef != 0 {
			name = cfStringToString((C.CFStringRef)(nameRef))
			C.CFRelease(C.CFTypeRef(nameRef))
		}

		sensors = append(sensors, thermalSensor{
			name:       name,
			location:   serviceNumberProperty(service, cfLocKey),
			registryID: strconv.FormatUint(uint64(C.IOHIDServiceClientGetRegistryID(service)), 10),
			temp:       float64(temp),
		})
	}

	for i, label := range resolveSensorNames(sensors) {
		ch <- c.temperature.mustNewConstMetric(sensors[i].temp, label)
	}

	return nil
}

// resolveSensorNames returns the sensor label to report for each sensor, in the
// order the sensors were given.
//
// Emitting the same label set twice makes the registry reject the samples and
// fail the whole scrape. A product name shared by several services is therefore
// qualified with the service location, falling back to the registry ID, which is
// unique per service, when the location does not tell them apart either.
func resolveSensorNames(sensors []thermalSensor) []string {
	nameCount := make(map[string]int, len(sensors))
	nameLocationCount := make(map[string]int, len(sensors))
	for _, s := range sensors {
		nameCount[s.name]++
		nameLocationCount[s.name+"\x00"+s.location]++
	}

	labels := make([]string, 0, len(sensors))
	for _, s := range sensors {
		label := s.name
		if nameCount[s.name] > 1 {
			suffix := s.location
			if suffix == "" || nameLocationCount[s.name+"\x00"+s.location] > 1 {
				suffix = s.registryID
			}
			label = s.name + "_" + suffix
		}
		labels = append(labels, label)
	}

	return labels
}

// serviceNumberProperty reads a numeric IOHID service property, returning an
// empty string when it is absent or not a number.
func serviceNumberProperty(service C.IOHIDServiceClientRef, key C.CFStringRef) string {
	ref := C.IOHIDServiceClientCopyProperty(service, key)
	if ref == 0 {
		return ""
	}
	defer C.CFRelease(ref)

	if C.CFGetTypeID(ref) != C.CFNumberGetTypeID() {
		return ""
	}

	var value C.longlong
	if C.CFNumberGetValue(C.CFNumberRef(ref), C.kCFNumberLongLongType, unsafe.Pointer(&value)) == 0 {
		return ""
	}

	return strconv.FormatInt(int64(value), 10)
}

func cfStringToString(s C.CFStringRef) string {
	p := C.CFStringGetCStringPtr(s, C.kCFStringEncodingUTF8)
	if p != nil {
		return C.GoString(p)
	}
	length := C.CFStringGetLength(s)
	if length <= 0 {
		return ""
	}
	maxBufLen := C.CFStringGetMaximumSizeForEncoding(length, C.kCFStringEncodingUTF8)
	if maxBufLen <= 0 {
		return ""
	}
	if maxBufLen > 4096 {
		maxBufLen = 4096
	}
	buf := make([]byte, maxBufLen)
	var usedBufLen C.CFIndex
	_ = C.CFStringGetBytes(s, C.CFRange{0, length}, C.kCFStringEncodingUTF8, C.UInt8(0), C.false, (*C.UInt8)(&buf[0]), maxBufLen, &usedBufLen)
	return SafeBytesToString(buf[:usedBufLen])
}
