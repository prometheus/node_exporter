package cmdb

import (
	"regexp"
	"strings"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectMemory(machineType string) (*model.Memory, error) {
	mm := &model.Memory{}

	if vm, err := mem.VirtualMemory(); err == nil {
		mm.TotalBytes = vm.Total
	}

	if machineType != "virtual" {
		mm.Modules = parseDMIDecodeMemory()
	}

	return mm, nil
}

func parseDMIDecodeMemory() []model.MemoryModule {
	if !commandExists("dmidecode") {
		return nil
	}

	out, err := runCmd("dmidecode", "-t", "memory")
	if err != nil {
		return nil
	}

	var modules []model.MemoryModule
	var cur *model.MemoryModule
	inDevice := false

	kvRe := regexp.MustCompile(`^\s+([A-Z][\w ]+):\s*(.*)$`)

	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(line, "Memory Device") || strings.HasPrefix(trimmed, "Memory Device") {
			if inDevice && cur != nil {
				modules = append(modules, *cur)
			}
			cur = &model.MemoryModule{}
			inDevice = true
			continue
		}

		if line == "" || (!strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")) {
			if inDevice && cur != nil {
				modules = append(modules, *cur)
				cur = nil
				inDevice = false
			}
			continue
		}

		if !inDevice || cur == nil {
			continue
		}

		if m := kvRe.FindStringSubmatch(line); m != nil {
			key := strings.TrimSpace(m[1])
			val := strings.TrimSpace(m[2])
			switch key {
			case "Locator":
				cur.Locator = val
			case "Bank Locator":
				cur.BankLocator = val
			case "Size":
				cur.Size = val
			case "Type":
				cur.Type = val
			case "Speed":
				cur.Speed = val
			case "Manufacturer":
				cur.Manufacturer = val
			case "Serial Number":
				cur.Serial = val
			case "Part Number":
				cur.PartNumber = val
			}
		}
	}

	if inDevice && cur != nil {
		modules = append(modules, *cur)
	}

	filtered := modules[:0]
	for _, m := range modules {
		if strings.EqualFold(m.Size, "No Module Installed") || (m.Size == "" && m.Type == "") {
			continue
		}
		filtered = append(filtered, m)
	}

	return filtered
}
