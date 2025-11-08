// Copyright 2016 The Prometheus Authors
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

//go:build !nohwmon
// +build !nohwmon

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

var (
	collectorHWmonChipInclude   = kingpin.Flag("collector.hwmon.chip-include", "Regexp of hwmon chip to include (mutually exclusive to device-exclude).").String()
	collectorHWmonChipExclude   = kingpin.Flag("collector.hwmon.chip-exclude", "Regexp of hwmon chip to exclude (mutually exclusive to device-include).").String()
	collectorHWmonSensorInclude = kingpin.Flag("collector.hwmon.sensor-include", "Regexp of hwmon sensor to include (mutually exclusive to sensor-exclude).").String()
	collectorHWmonSensorExclude = kingpin.Flag("collector.hwmon.sensor-exclude", "Regexp of hwmon sensor to exclude (mutually exclusive to sensor-include).").String()

	hwmonInvalidMetricChars = regexp.MustCompile("[^a-z0-9:_]")
	hwmonFilenameFormat     = regexp.MustCompile(`^(?P<type>[^0-9]+)(?P<id>[0-9]*)?(_(?P<property>.+))?$`)
	hwmonLabelDesc          = []string{"chip", "sensor"}
	hwmonChipNameLabelDesc  = []string{"chip", "chip_name"}
	hwmonSensorTypes        = []string{
		"vrm", "beep_enable", "update_interval", "in", "cpu", "fan",
		"pwm", "temp", "curr", "power", "energy", "humidity",
		"intrusion", "freq",
	}
)

func init() {
	registerCollector("hwmon", defaultEnabled, NewHwMonCollector)
}

type hwMonCollector struct {
	deviceFilter deviceFilter
	sensorFilter deviceFilter
	logger       *slog.Logger
}

// NewHwMonCollector returns a new Collector exposing /sys/class/hwmon stats
// (similar to lm-sensors).
func NewHwMonCollector(logger *slog.Logger) (Collector, error) {

	return &hwMonCollector{
		logger:       logger,
		deviceFilter: newDeviceFilter(*collectorHWmonChipExclude, *collectorHWmonChipInclude),
		sensorFilter: newDeviceFilter(*collectorHWmonSensorExclude, *collectorHWmonSensorInclude),
	}, nil
}

func cleanMetricName(name string) string {
	lower := strings.ToLower(name)
	replaced := hwmonInvalidMetricChars.ReplaceAllLiteralString(lower, "_")
	cleaned := strings.Trim(replaced, "_")
	return cleaned
}

func addValueFile(data map[string]map[string]string, sensor string, prop string, file string) {
	raw, err := sysReadFile(file)
	if err != nil {
		return
	}
	value := strings.Trim(string(raw), "\n")

	if _, ok := data[sensor]; !ok {
		data[sensor] = make(map[string]string)
	}

	data[sensor][prop] = value
}

// sysReadFile is a simplified os.ReadFile that invokes syscall.Read directly.
func sysReadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// On some machines, hwmon drivers are broken and return EAGAIN.  This causes
	// Go's os.ReadFile implementation to poll forever.
	//
	// Since we either want to read data or bail immediately, do the simplest
	// possible read using system call directly.
	b := make([]byte, 128)
	n, err := unix.Read(int(f.Fd()), b)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, fmt.Errorf("failed to read file: %q, read returned negative bytes value: %d", file, n)
	}

	return b[:n], nil
}

// explodeSensorFilename splits a sensor name into <type><num>_<property>.
func explodeSensorFilename(filename string) (ok bool, sensorType string, sensorNum int, sensorProperty string) {
	matches := hwmonFilenameFormat.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return false, sensorType, sensorNum, sensorProperty
	}
	for i, match := range hwmonFilenameFormat.SubexpNames() {
		if i >= len(matches) {
			return true, sensorType, sensorNum, sensorProperty
		}
		if match == "type" {
			sensorType = matches[i]
		}
		if match == "property" {
			sensorProperty = matches[i]
		}
		if match == "id" && len(matches[i]) > 0 {
			if num, err := strconv.Atoi(matches[i]); err == nil {
				sensorNum = num
			} else {
				return false, sensorType, sensorNum, sensorProperty
			}
		}
	}
	return true, sensorType, sensorNum, sensorProperty
}

func collectSensorData(dir string, data map[string]map[string]string) error {
	sensorFiles, dirError := os.ReadDir(dir)
	if dirError != nil {
		return dirError
	}
	for _, file := range sensorFiles {
		filename := file.Name()
		ok, sensorType, sensorNum, sensorProperty := explodeSensorFilename(filename)
		if !ok {
			continue
		}

		if slices.Contains(hwmonSensorTypes, sensorType) {
			addValueFile(data, sensorType+strconv.Itoa(sensorNum), sensorProperty, filepath.Join(dir, file.Name()))
		}
	}
	return nil
}

func (c *hwMonCollector) updateHwmon(ch chan<- prometheus.Metric, dir string) error {
	hwmonName, err := c.hwmonName(dir)
	if err != nil {
		return err
	}

	if c.deviceFilter.ignored(hwmonName) {
		c.logger.Debug("ignoring hwmon chip", "chip", hwmonName)
		return nil
	}

	data := make(map[string]map[string]string)
	err = collectSensorData(dir, data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(dir, "device")); err == nil {
		err := collectSensorData(filepath.Join(dir, "device"), data)
		if err != nil {
			return err
		}
	}

	hwmonChipName, err := c.hwmonHumanReadableChipName(dir)
	if err == nil {
		// sensor chip metadata
		desc := prometheus.NewDesc(
			"node_hwmon_chip_names",
			"Annotation metric for human-readable chip names",
			hwmonChipNameLabelDesc,
			nil,
		)

		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			1.0,
			hwmonName,
			hwmonChipName,
		)
	}

	// Format all sensors.
	for sensor, sensorData := range data {

		// Filtering for sensors is done on concatenated device name and sensor name
		// separated by a semicolon. This allows for excluding or including of specific
		// sensors on specific devices. For example, to exclude the sensor "temp3" on
		// the device "platform_coretemp_0", use "platform_coretemp_0;temp3"
		if c.sensorFilter.ignored(hwmonName + ";" + sensor) {
			c.logger.Debug("ignoring sensor", "sensor", sensor)
			continue
		}

		_, sensorType, _, _ := explodeSensorFilename(sensor)

		labels := []string{hwmonName, sensor}
		if labelText, ok := sensorData["label"]; ok {
			label := strings.ToValidUTF8(labelText, "�")
			desc := prometheus.NewDesc("node_hwmon_sensor_label", "Label for given chip and sensor",
				[]string{"chip", "sensor", "label"}, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1.0, hwmonName, sensor, label)
		}

		if sensorType == "beep_enable" {
			value := 0.0
			if sensorData[""] == "1" {
				value = 1.0
			}
			metricName := "node_hwmon_beep_enabled"
			desc := prometheus.NewDesc(metricName, "Hardware beep enabled", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(
				desc, prometheus.GaugeValue, value, labels...)
			continue
		}
		if sensorType == "vrm" {
			parsedValue, err := strconv.ParseFloat(sensorData[""], 64)
			if err != nil {
				continue
			}
			metricName := "node_hwmon_voltage_regulator_version"
			desc := prometheus.NewDesc(metricName, "Hardware voltage regulator", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(
				desc, prometheus.GaugeValue, parsedValue, labels...)
			continue
		}
		if sensorType == "update_interval" {
			parsedValue, err := strconv.ParseFloat(sensorData[""], 64)
			if err != nil {
				continue
			}
			metricName := "node_hwmon_update_interval_seconds"
			desc := prometheus.NewDesc(metricName, "Hardware monitor update interval", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(
				desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
			continue
		}

		prefix := "node_hwmon_" + sensorType

		for element, value := range sensorData {

			if element == "label" {
				continue
			}

			name := prefix
			if element == "input" {
				// input is actually the value
				if _, ok := sensorData[""]; ok {
					name = name + "_input"
				}
			} else if element != "" {
				name = name + "_" + cleanMetricName(element)
			}
			parsedValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				continue
			}

			// special elements, fault, alarm & beep should be handed out without units
			if element == "fault" || element == "alarm" {
				desc := prometheus.NewDesc(name, "Hardware sensor "+element+" status ("+sensorType+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}
			if element == "beep" {
				desc := prometheus.NewDesc(name+"_enabled", "Hardware monitor sensor has beeping enabled", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}

			// everything else should get a unit
			if sensorType == "in" || sensorType == "cpu" {
				desc := prometheus.NewDesc(name+"_volts", "Hardware monitor for voltage ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "temp" && element != "type" {
				if element == "" {
					element = "input"
				}
				desc := prometheus.NewDesc(name+"_celsius", "Hardware monitor for temperature ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "curr" {
				desc := prometheus.NewDesc(name+"_amps", "Hardware monitor for current ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "energy" {
				desc := prometheus.NewDesc(name+"_joule_total", "Hardware monitor for joules used so far ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.CounterValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "power" && element == "accuracy" {
				desc := prometheus.NewDesc(name, "Hardware monitor power meter accuracy, as a ratio", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "power" && (element == "average_interval" || element == "average_interval_min" || element == "average_interval_max") {
				desc := prometheus.NewDesc(name+"_seconds", "Hardware monitor power usage update interval ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "power" {
				desc := prometheus.NewDesc(name+"_watt", "Hardware monitor for power usage in watts ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}

			if sensorType == "humidity" {
				desc := prometheus.NewDesc(name, "Hardware monitor for humidity, as a ratio (multiply with 100.0 to get the humidity as a percentage) ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}

			if sensorType == "fan" && (element == "input" || element == "min" || element == "max" || element == "target") {
				desc := prometheus.NewDesc(name+"_rpm", "Hardware monitor for fan revolutions per minute ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(
					desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}

			if sensorType == "freq" && element == "input" {
				if label, ok := sensorData["label"]; ok {
					sensorLabel := cleanMetricName(label)
					desc := prometheus.NewDesc(name+"_freq_mhz", "Hardware monitor for GPU frequency in MHz", hwmonLabelDesc, nil)
					ch <- prometheus.MustNewConstMetric(
						desc, prometheus.GaugeValue, parsedValue/1000000.0, append(labels[:len(labels)-1], sensorLabel)...)
				}
				continue
			}
			// fallback, just dump the metric as is

			desc := prometheus.NewDesc(name, "Hardware monitor "+sensorType+" element "+element, hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(
				desc, prometheus.GaugeValue, parsedValue, labels...)
		}
	}

	return nil
}

func (c *hwMonCollector) hwmonName(dir string) (string, error) {
	// generate a name for a sensor path

	// sensor numbering depends on the order of linux module loading and
	// is thus unstable.
	// However the path of the device has to be stable:
	// - /sys/devices/<bus>/<device>
	// Some hardware monitors have a "name" file that exports a human
	// readable name that can be used.

	// human readable names would be bat0 or coretemp, while a path string
	// could be platform_applesmc.768

	// preference 1: construct a name based on device name, always unique

	devicePath, devErr := filepath.EvalSymlinks(filepath.Join(dir, "device"))
	if devErr == nil {
		devPathPrefix, devName := filepath.Split(devicePath)
		_, devType := filepath.Split(strings.TrimRight(devPathPrefix, "/"))

		cleanDevName := cleanMetricName(devName)
		cleanDevType := cleanMetricName(devType)

		if cleanDevType != "" && cleanDevName != "" {
			return cleanDevType + "_" + cleanDevName, nil
		}

		if cleanDevName != "" {
			return cleanDevName, nil
		}
	}

	// preference 2: is there a name file
	sysnameRaw, nameErr := os.ReadFile(filepath.Join(dir, "name"))
	if nameErr == nil && string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}

	// it looks bad, name and device don't provide enough information
	// return a hwmon[0-9]* name

	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}

	// take the last path element, this will be hwmonX
	_, name := filepath.Split(realDir)
	cleanName := cleanMetricName(name)
	if cleanName != "" {
		return cleanName, nil
	}
	return "", errors.New("Could not derive a monitoring name for " + dir)
}

// hwmonHumanReadableChipName is similar to the methods in hwmonName, but with
// different precedences -- we can allow duplicates here.
func (c *hwMonCollector) hwmonHumanReadableChipName(dir string) (string, error) {
	sysnameRaw, nameErr := os.ReadFile(filepath.Join(dir, "name"))
	if nameErr != nil {
		return "", nameErr
	}

	if string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}

	return "", errors.New("Could not derive a human-readable chip type for " + dir)
}

func (c *hwMonCollector) Update(ch chan<- prometheus.Metric) error {
	// Step 1: scan /sys/class/hwmon, resolve all symlinks and call
	//         updatesHwmon for each folder

	hwmonPathName := filepath.Join(sysFilePath("class"), "hwmon")

	hwmonFiles, err := os.ReadDir(hwmonPathName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("hwmon collector metrics are not available for this system")
			return ErrNoData
		}

		return err
	}

	var lastErr error
	for _, hwDir := range hwmonFiles {
		hwmonXPathName := filepath.Join(hwmonPathName, hwDir.Name())
		fileInfo, err := os.Lstat(hwmonXPathName)
		if err != nil {
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink > 0 {
			fileInfo, err = os.Stat(hwmonXPathName)
			if err != nil {
				continue
			}
		}

		if !fileInfo.IsDir() {
			continue
		}

		if err = c.updateHwmon(ch, hwmonXPathName); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
