// Copyright 2019 The Prometheus Authors
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

//go:build !nopowersupplyclass
// +build !nopowersupplyclass

package collector

/*
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include <IOKit/ps/IOPowerSources.h>
#include <IOKit/ps/IOPSKeys.h>
#include <CoreFoundation/CFArray.h>

// values collected from IOKit Power Source APIs
// Functions documentation available at
// https://developer.apple.com/documentation/iokit/iopowersources_h
// CFDictionary keys definition
// https://developer.apple.com/documentation/iokit/iopskeys_h/defines
struct macos_powersupply {
   char *Name;
   char *PowerSourceState;
   char *Type;
   char *TransportType;
   char *BatteryHealth;
   char *HardwareSerialNumber;

   int *PowerSourceID;
   int *CurrentCapacity;
   int *MaxCapacity;
   int *DesignCapacity;
   int *NominalCapacity;

   int *TimeToEmpty;
   int *TimeToFullCharge;

   int *Voltage;
   int *Current;

   int *Temperature;

   // boolean values
   int *IsCharged;
   int *IsCharging;
   int *InternalFailure;
   int *IsPresent;
};

int *CFDictionaryGetInt(CFDictionaryRef theDict, const void *key) {
    CFNumberRef tmp;
    int *value;

    tmp = CFDictionaryGetValue(theDict, key);

    if (tmp == NULL)
        return NULL;

    value = (int*)malloc(sizeof(int));
    if (CFNumberGetValue(tmp, kCFNumberIntType, value)) {
        return value;
    }

    free(value);
    return NULL;
}

int *CFDictionaryGetBoolean(CFDictionaryRef theDict, const void *key) {
    CFBooleanRef tmp;
    int *value;

    tmp = CFDictionaryGetValue(theDict, key);

    if (tmp == NULL)
        return NULL;

    value = (int*)malloc(sizeof(int));
    if (CFBooleanGetValue(tmp)) {
        *value = 1;
    } else {
        *value = 0;
    }

    return value;
}

char *CFDictionaryGetSring(CFDictionaryRef theDict, const void *key) {
    CFStringRef tmp;
    CFIndex size;
    char *value;

    tmp = CFDictionaryGetValue(theDict, key);

    if (tmp == NULL)
        return NULL;

    size = CFStringGetLength(tmp) + 1;
    value = (char*)malloc(size);

    if(CFStringGetCString(tmp, value, size, kCFStringEncodingUTF8)) {
         return value;
    }

    free(value);
    return NULL;
}

struct macos_powersupply* getPowerSupplyInfo(CFDictionaryRef powerSourceInformation) {
    struct macos_powersupply *ret;

    if (powerSourceInformation == NULL)
        return NULL;

    ret = (struct macos_powersupply*)malloc(sizeof(struct macos_powersupply));

    ret->PowerSourceID    = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSPowerSourceIDKey));
    ret->CurrentCapacity  = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSCurrentCapacityKey));
    ret->MaxCapacity      = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSMaxCapacityKey));
    ret->DesignCapacity   = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSDesignCapacityKey));
    ret->NominalCapacity  = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSNominalCapacityKey));
    ret->TimeToEmpty      = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSTimeToEmptyKey));
    ret->TimeToFullCharge = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSTimeToFullChargeKey));
    ret->Voltage          = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSVoltageKey));
    ret->Current          = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSCurrentKey));
    ret->Temperature      = CFDictionaryGetInt(powerSourceInformation, CFSTR(kIOPSTemperatureKey));

    ret->Name = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSNameKey));
    ret->PowerSourceState = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSPowerSourceStateKey));
    ret->Type = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSTypeKey));
    ret->TransportType = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSTransportTypeKey));
    ret->BatteryHealth = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSBatteryHealthKey));
    ret->HardwareSerialNumber = CFDictionaryGetSring(powerSourceInformation, CFSTR(kIOPSHardwareSerialNumberKey));

    ret->IsCharged       = CFDictionaryGetBoolean(powerSourceInformation, CFSTR(kIOPSIsChargedKey));
    ret->IsCharging      = CFDictionaryGetBoolean(powerSourceInformation, CFSTR(kIOPSIsChargingKey));
    ret->InternalFailure = CFDictionaryGetBoolean(powerSourceInformation, CFSTR(kIOPSInternalFailureKey));
    ret->IsPresent       = CFDictionaryGetBoolean(powerSourceInformation, CFSTR(kIOPSIsPresentKey));

    return ret;
}



void releasePowerSupply(struct macos_powersupply *ps) {
    free(ps->Name);
    free(ps->PowerSourceState);
    free(ps->Type);
    free(ps->TransportType);
    free(ps->BatteryHealth);
    free(ps->HardwareSerialNumber);

    free(ps->PowerSourceID);
    free(ps->CurrentCapacity);
    free(ps->MaxCapacity);
    free(ps->DesignCapacity);
    free(ps->NominalCapacity);
    free(ps->TimeToEmpty);
    free(ps->TimeToFullCharge);
    free(ps->Voltage);
    free(ps->Current);
    free(ps->Temperature);

    free(ps->IsCharged);
    free(ps->IsCharging);
    free(ps->InternalFailure);
    free(ps->IsPresent);

    free(ps);
}
*/
import "C"

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func (c *powerSupplyClassCollector) Update(ch chan<- prometheus.Metric) error {
	psList, err := getPowerSourceList()
	if err != nil {
		return fmt.Errorf("couldn't get IOPPowerSourcesList: %w", err)
	}

	for _, info := range psList {
		labels := getPowerSourceDescriptorLabels(info)
		powerSupplyName := labels["power_supply"]

		if c.ignoredPattern.MatchString(powerSupplyName) {
			continue
		}

		for name, value := range getPowerSourceDescriptorMap(info) {
			if value == nil {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, c.subsystem, name),
					fmt.Sprintf("IOKit Power Source information field %s for <power_supply>.", name),
					[]string{"power_supply"}, nil,
				),
				prometheus.GaugeValue, *value, powerSupplyName,
			)
		}

		pushEnumMetric(
			ch,
			getPowerSourceDescriptorState(info),
			"power_source_state",
			c.subsystem,
			powerSupplyName,
		)

		pushEnumMetric(
			ch,
			getPowerSourceDescriptorBatteryHealth(info),
			"battery_health",
			c.subsystem,
			powerSupplyName,
		)

		var (
			keys   []string
			values []string
		)
		for name, value := range labels {
			if value != "" {
				keys = append(keys, name)
				values = append(values, value)
			}
		}
		fieldDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"IOKit Power Source information for <power_supply>.",
			keys,
			nil,
		)
		ch <- prometheus.MustNewConstMetric(fieldDesc, prometheus.GaugeValue, 1.0, values...)

		C.releasePowerSupply(info)
	}

	return nil
}

// getPowerSourceList fetches information from IOKit APIs
//
// Data is provided as opaque CoreFoundation references
// C.getPowerSupplyInfo will convert those objects in something
// easily manageable in Go.
// https://developer.apple.com/documentation/iokit/iopowersources_h
func getPowerSourceList() ([]*C.struct_macos_powersupply, error) {
	infos, err := C.IOPSCopyPowerSourcesInfo()
	if err != nil {
		return nil, err
	}
	defer C.CFRelease(infos)

	psList, err := C.IOPSCopyPowerSourcesList(infos)
	if err != nil {
		return nil, err
	}

	if psList == C.CFArrayRef(0) {
		return nil, nil
	}
	defer C.CFRelease(C.CFTypeRef(psList))

	size, err := C.CFArrayGetCount(psList)
	if err != nil {
		return nil, err
	}

	ret := make([]*C.struct_macos_powersupply, size)
	for i := C.CFIndex(0); i < size; i++ {
		ps, err := C.CFArrayGetValueAtIndex(psList, i)
		if err != nil {
			return nil, err
		}

		dict, err := C.IOPSGetPowerSourceDescription(infos, (C.CFTypeRef)(ps))
		if err != nil {
			return nil, err
		}

		info, err := C.getPowerSupplyInfo(dict)
		if err != nil {
			return nil, err
		}

		ret[int(i)] = info
	}

	return ret, nil
}

func getPowerSourceDescriptorMap(info *C.struct_macos_powersupply) map[string]*float64 {
	return map[string]*float64{
		"current_capacity":      convertValue(info.CurrentCapacity),
		"max_capacity":          convertValue(info.MaxCapacity),
		"design_capacity":       convertValue(info.DesignCapacity),
		"nominal_capacity":      convertValue(info.NominalCapacity),
		"time_to_empty_seconds": minutesToSeconds(info.TimeToEmpty),
		"time_to_full_seconds":  minutesToSeconds(info.TimeToFullCharge),
		"voltage_volt":          scaleValue(info.Voltage, 1e3),
		"current_ampere":        scaleValue(info.Current, 1e3),
		"temp_celsius":          convertValue(info.Temperature),
		"present":               convertValue(info.IsPresent),
		"charging":              convertValue(info.IsCharging),
		"charged":               convertValue(info.IsCharged),
		"internal_failure":      convertValue(info.InternalFailure),
	}
}

func getPowerSourceDescriptorLabels(info *C.struct_macos_powersupply) map[string]string {
	return map[string]string{
		"id":             strconv.FormatInt(int64(*info.PowerSourceID), 10),
		"power_supply":   C.GoString(info.Name),
		"type":           C.GoString(info.Type),
		"transport_type": C.GoString(info.TransportType),
		"serial_number":  C.GoString(info.HardwareSerialNumber),
	}
}

func getPowerSourceDescriptorState(info *C.struct_macos_powersupply) map[string]float64 {
	stateMap := map[string]float64{
		"Off Line":      0,
		"AC Power":      0,
		"Battery Power": 0,
	}

	// This field is always present
	// https://developer.apple.com/documentation/iokit/kiopspowersourcestatekey
	stateMap[C.GoString(info.PowerSourceState)] = 1

	return stateMap
}

func getPowerSourceDescriptorBatteryHealth(info *C.struct_macos_powersupply) map[string]float64 {
	// This field is optional
	// https://developer.apple.com/documentation/iokit/kiopsBatteryHealthkey
	if info.BatteryHealth == nil {
		return nil
	}

	stateMap := map[string]float64{
		"Good": 0,
		"Fair": 0,
		"Poor": 0,
	}

	stateMap[C.GoString(info.BatteryHealth)] = 1

	return stateMap
}

func convertValue(value *C.int) *float64 {
	if value == nil {
		return nil
	}

	ret := new(float64)
	*ret = (float64)(*value)
	return ret
}

func scaleValue(value *C.int, scale float64) *float64 {
	ret := convertValue(value)
	if ret == nil {
		return nil
	}

	*ret /= scale

	return ret
}

// minutesToSeconds converts *C.int minutes into *float64 seconds.
//
// Only positive values will be scaled to seconds, because negative ones
// have special meanings. I.e. -1 indicates "Still Calculating the Time"
func minutesToSeconds(minutes *C.int) *float64 {
	ret := convertValue(minutes)
	if ret == nil {
		return nil
	}

	if *ret > 0 {
		*ret *= 60
	}

	return ret
}

func pushEnumMetric(ch chan<- prometheus.Metric, values map[string]float64, name, subsystem, powerSupply string) {
	for state, value := range values {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, name),
				fmt.Sprintf("IOKit Power Source information field %s for <power_supply>.", name),
				[]string{"power_supply", "state"}, nil,
			),
			prometheus.GaugeValue, value, powerSupply, state,
		)
	}
}
