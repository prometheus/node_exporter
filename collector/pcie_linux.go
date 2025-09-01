package collector

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	pciIdsPaths = []string{
		"/usr/share/misc/pci.ids",
		"/usr/share/hwdata/pci.ids",
	}
	pciVendors    = make(map[string]string)
	pciDevices    = make(map[string]map[string]string)
	pciSubsystems = make(map[string]map[string]string)
)

type pcieCollector struct {
	info   *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("pcie", defaultDisabled, NewPCIeCollector)
	loadPCIIds()
}

func NewPCIeCollector(logger *slog.Logger) (Collector, error) {
	return &pcieCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "pcie_device", "info"),
			"Detailed PCIe device information from /sys/bus/pci/devices/. "+
				"Power state values: D0 (fully powered), D1/D2 (intermediate states), D3hot/D3cold (lowest power). "+
				"Link speeds are in GT/s (e.g., 2.5, 5.0, 8.0, 16.0). "+
				"Link widths are lanes (x1, x2, x4, x8, x16). "+
				"D3cold_allowed indicates if deepest power saving state is supported (0/1).",
			[]string{
				"slot", // Changed from "device"
				"vendor_id",
				"vendor_name", // Changed from "vendor"
				"device_id",
				"device_name", // Changed from "device"
				"subsystem_vendor_id",
				"subsystem_vendor_name", // Changed from "subsystem_vendor"
				"subsystem_device_id",
				"subsystem_device_name", // Changed from "subsystem_device"
				"class",
				"revision",
				"current_speed",
				"current_width",
				"max_speed",
				"max_width",
				"power_state",
				"d3cold_allowed",
			},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *pcieCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := filepath.Glob("/sys/bus/pci/devices/*")
	if err != nil {
		return fmt.Errorf("failed to list PCI devices: %w", err)
	}

	for _, devicePath := range devices {
		deviceID := filepath.Base(devicePath)
		if err := c.collectDeviceMetrics(ch, devicePath, deviceID); err != nil {
			c.logger.Debug("failed collecting metrics for device", "device", deviceID, "err", err)
			continue
		}
	}

	return nil
}

func (c *pcieCollector) collectDeviceMetrics(ch chan<- prometheus.Metric, devicePath, deviceID string) error {
	// Read IDs first
	vendorID := readFileContent(filepath.Join(devicePath, "vendor"))
	devID := readFileContent(filepath.Join(devicePath, "device"))
	subsysVendorID := readFileContent(filepath.Join(devicePath, "subsystem_vendor"))
	subsysDeviceID := readFileContent(filepath.Join(devicePath, "subsystem_device"))

	// Get human-readable names from pci.ids database
	vendor := getPCIVendorName(vendorID)
	device := getPCIDeviceName(vendorID, devID)
	subsysVendor := getPCIVendorName(subsysVendorID)
	subsysDevice := getPCISubsystemName(vendorID, devID, subsysVendorID, subsysDeviceID)

	labels := []string{
		deviceID,
		vendorID,
		vendor,
		devID,
		device,
		subsysVendorID,
		subsysVendor,
		subsysDeviceID,
		subsysDevice,
		readFileContent(filepath.Join(devicePath, "class")),
		readFileContent(filepath.Join(devicePath, "revision")),
		readFileContent(filepath.Join(devicePath, "current_link_speed")),
		readFileContent(filepath.Join(devicePath, "current_link_width")),
		readFileContent(filepath.Join(devicePath, "max_link_speed")),
		readFileContent(filepath.Join(devicePath, "max_link_width")),
		readFileContent(filepath.Join(devicePath, "power_state")),
		readFileContent(filepath.Join(devicePath, "d3cold_allowed")),
	}

	ch <- prometheus.MustNewConstMetric(
		c.info,
		prometheus.GaugeValue,
		1,
		labels...,
	)

	return nil
}

func readFileContent(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(content))
}

func loadPCIIds() {
	var file *os.File
	var err error

	// Try each possible path
	for _, path := range pciIdsPaths {
		file, err = os.Open(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentVendor, currentDevice string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle vendor lines (no leading whitespace)
		if !strings.HasPrefix(line, "\t") {
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) >= 2 {
				currentVendor = strings.TrimSpace(parts[0])
				pciVendors[currentVendor] = strings.TrimSpace(parts[1])
				currentDevice = ""
			}
			continue
		}

		// Handle device lines (single tab)
		if strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "\t\t") {
			line = strings.TrimPrefix(line, "\t")
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) >= 2 && currentVendor != "" {
				currentDevice = strings.TrimSpace(parts[0])
				if pciDevices[currentVendor] == nil {
					pciDevices[currentVendor] = make(map[string]string)
				}
				pciDevices[currentVendor][currentDevice] = strings.TrimSpace(parts[1])
			}
			continue
		}

		// Handle subsystem lines (double tab)
		if strings.HasPrefix(line, "\t\t") {
			line = strings.TrimPrefix(line, "\t\t")
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) >= 2 && currentVendor != "" && currentDevice != "" {
				subsysID := strings.TrimSpace(parts[0])
				subsysName := strings.TrimSpace(parts[1])
				key := fmt.Sprintf("%s:%s", currentVendor, currentDevice)
				if pciSubsystems[key] == nil {
					pciSubsystems[key] = make(map[string]string)
				}
				pciSubsystems[key][subsysID] = subsysName
			}
		}
	}
}

func getPCIVendorName(vendorID string) string {
	// Remove "0x" prefix if present
	vendorID = strings.TrimPrefix(vendorID, "0x")
	vendorID = strings.ToLower(vendorID)

	if name, ok := pciVendors[vendorID]; ok {
		return name
	}
	return vendorID // Return ID if name not found
}

func getPCIDeviceName(vendorID, deviceID string) string {
	// Remove "0x" prefix if present
	vendorID = strings.TrimPrefix(vendorID, "0x")
	deviceID = strings.TrimPrefix(deviceID, "0x")
	vendorID = strings.ToLower(vendorID)
	deviceID = strings.ToLower(deviceID)

	if devices, ok := pciDevices[vendorID]; ok {
		if name, ok := devices[deviceID]; ok {
			return name
		}
	}
	return deviceID // Return ID if name not found
}

func getPCISubsystemName(vendorID, deviceID, subsysVendorID, subsysDeviceID string) string {
	// Normalize all IDs
	vendorID = strings.TrimPrefix(vendorID, "0x")
	deviceID = strings.TrimPrefix(deviceID, "0x")
	subsysVendorID = strings.TrimPrefix(subsysVendorID, "0x")
	subsysDeviceID = strings.TrimPrefix(subsysDeviceID, "0x")

	key := fmt.Sprintf("%s:%s", vendorID, deviceID)
	subsysKey := fmt.Sprintf("%s:%s", subsysVendorID, subsysDeviceID)

	if subsystems, ok := pciSubsystems[key]; ok {
		if name, ok := subsystems[subsysKey]; ok {
			return name
		}
	}
	return subsysDeviceID
}
