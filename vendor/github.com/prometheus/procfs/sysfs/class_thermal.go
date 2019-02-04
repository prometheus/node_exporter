// Copyright 2018 The Prometheus Authors
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

// +build !windows

package sysfs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// ClassThermalZoneStats contains info from files in /sys/class/thermal/thermal_zone<zone>
// for a single <zone>.
// https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt
type ClassThermalZoneStats struct {
	Name    string  // The name of the zone from the directory structure.
	Type    string  // The type of thermal zone.
	Temp    uint64  // Temperature in millidegree Celsius.
	Policy  string  // One of the various thermal governors used for a particular zone.
	Mode    *bool   // Optional: One of the predefined values in [enabled, disabled].
	Passive *uint64 // Optional: millidegrees Celsius. (0 for disabled, > 1000 for enabled+value)
}

// NewClassThermalZoneStats returns Thermal Zone metrics for all zones.
func (fs FS) NewClassThermalZoneStats() ([]ClassThermalZoneStats, error) {
	zones, err := filepath.Glob(fs.Path("class/thermal/thermal_zone[0-9]*"))
	if err != nil {
		return []ClassThermalZoneStats{}, err
	}

	var zoneStats = ClassThermalZoneStats{}
	stats := make([]ClassThermalZoneStats, len(zones))
	for i, zone := range zones {
		zoneName := strings.TrimPrefix(filepath.Base(zone), "thermal_zone")

		zoneStats, err = parseClassThermalZone(zone)
		if err != nil {
			return []ClassThermalZoneStats{}, err
		}
		zoneStats.Name = zoneName
		stats[i] = zoneStats
	}
	return stats, nil
}

func parseClassThermalZone(zone string) (ClassThermalZoneStats, error) {
	// Required attributes.
	zoneType, err := util.SysReadFile(filepath.Join(zone, "type"))
	if err != nil {
		return ClassThermalZoneStats{}, err
	}
	zonePolicy, err := util.SysReadFile(filepath.Join(zone, "policy"))
	if err != nil {
		return ClassThermalZoneStats{}, err
	}
	zoneTemp, err := util.ReadUintFromFile(filepath.Join(zone, "temp"))
	if err != nil {
		return ClassThermalZoneStats{}, err
	}

	// Optional attributes.
	mode, err := util.SysReadFile(filepath.Join(zone, "mode"))
	if err != nil && !os.IsNotExist(err) && !os.IsPermission(err) {
		return ClassThermalZoneStats{}, err
	}
	zoneMode := util.ParseBool(mode)

	var zonePassive *uint64
	passive, err := util.ReadUintFromFile(filepath.Join(zone, "passive"))
	if os.IsNotExist(err) || os.IsPermission(err) {
		zonePassive = nil
	} else if err != nil {
		return ClassThermalZoneStats{}, err
	} else {
		zonePassive = &passive
	}

	return ClassThermalZoneStats{
		Type:    zoneType,
		Policy:  zonePolicy,
		Temp:    zoneTemp,
		Mode:    zoneMode,
		Passive: zonePassive,
	}, nil
}
