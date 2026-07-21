package cmdb

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
)

func CollectGPU() (*model.GPU, error) {
	g := &model.GPU{Devices: []model.GPUDevice{}}
	nvidiaOK := collectNVIDIA(g)
	huaweiOK := collectHuaweiNPU(g)
	collectLspciGPU(g, nvidiaOK, huaweiOK)
	return g, nil
}

func collectNVIDIA(g *model.GPU) bool {
	if !commandExists("nvidia-smi") {
		return false
	}
	out, err := runCmd("nvidia-smi",
		"--query-gpu=index,name,uuid,serial,vbios_version,memory.total,memory.used,memory.free,utilization.gpu,temperature.gpu,power.draw,driver_version",
		"--format=csv,noheader,nounits")
	if err != nil {
		return false
	}
	devs := parseNVIDIA(out)
	if len(devs) == 0 {
		return false
	}
	g.Devices = append(g.Devices, devs...)
	return true
}

func parseNVIDIA(out string) []model.GPUDevice {
	var devices []model.GPUDevice
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) < 12 {
			continue
		}
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		devices = append(devices, model.GPUDevice{
			Index:           atoiSafe(fields[0]),
			Vendor:          "nvidia",
			Name:            fields[1],
			UUID:            fields[2],
			Serial:          fields[3],
			FirmwareVersion: fields[4],
			MemoryTotalMB:   atouSafe(fields[5]),
			MemoryUsedMB:    atouSafe(fields[6]),
			MemoryFreeMB:    atouSafe(fields[7]),
			Utilization:     atofSafe(fields[8]),
			Temperature:     atofSafe(fields[9]),
			PowerW:          atofSafe(fields[10]),
			DriverVersion:   fields[11],
			Health:          "OK",
			RuntimeMetrics:  true,
		})
	}
	return devices
}

var (
	npuVerRe    = regexp.MustCompile(`Version:\s*(\S+)`)
	npuMemRe    = regexp.MustCompile(`(\d+)\s*/\s*(\d+)`)
	npuBusIDRe  = regexp.MustCompile(`^[0-9A-Fa-f]{4}:[0-9A-Fa-f]{2}:\d{2}\.\d`)
	npuSerialRe = regexp.MustCompile(`(?im)^\s*Serial\s*Number\s*:\s*(\S+)`)
	npuFWRe     = regexp.MustCompile(`(?im)^\s*Firmware\s*Version\s*:\s*(\S+)`)
	pciIDRe     = regexp.MustCompile(`\[([0-9A-Fa-f]{4}:[0-9A-Fa-f]{4})\]`)
)

func collectHuaweiNPU(g *model.GPU) bool {
	if !commandExists("npu-smi") {
		return false
	}
	out, err := runCmd("npu-smi", "info")
	if err != nil {
		return false
	}
	devices := parseHuaweiNPU(out)
	for i := range devices {
		boardOut, err := runCmd("npu-smi", "info", "-t", "board", "-i", strconv.Itoa(devices[i].Index))
		if err != nil {
			continue
		}
		if m := npuSerialRe.FindStringSubmatch(boardOut); m != nil {
			devices[i].Serial = strings.TrimSpace(m[1])
		}
		if m := npuFWRe.FindStringSubmatch(boardOut); m != nil {
			devices[i].FirmwareVersion = strings.TrimSpace(m[1])
		}
	}
	if len(devices) == 0 {
		return false
	}
	g.Devices = append(g.Devices, devices...)
	return true
}

func parseHuaweiNPU(out string) []model.GPUDevice {
	var driverVer string
	if m := npuVerRe.FindStringSubmatch(out); m != nil {
		driverVer = m[1]
	}

	var devices []model.GPUDevice
	var cur *model.GPUDevice
	// afterSep is true after a "+---+ / +===+" separator: only the first
	// data line following one is a candidate card header. A single NPU card
	// may span several separator-bounded blocks — e.g. 310P3 repeats the
	// card header once per chip, with "+---+" between chips. We therefore do
	// NOT flush on every separator; we keep merging into the current card as
	// long as the NPU index (first column) is unchanged, and only start a
	// new card when the NPU index changes.
	afterSep := false
	flush := func() {
		if cur != nil {
			cur.RuntimeMetrics = true
			devices = append(devices, *cur)
			cur = nil
		}
	}

	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// The trailing "Process info" section lists one row per running
		// process: "| NPU  Chip | Process id | Process name | Process
		// memory |". Those rows have the same shape as card headers
		// (numeric NPU idx + chip in col1, numeric Process id in the
		// Health column), so without an explicit guard the state machine
		// mints one phantom device per running process whose Name is the
		// chip number and Health is the process id. Stop as soon as we
		// cross into that section, detected by its column header.
		if strings.Contains(line, "Process id") {
			break
		}
		if strings.HasPrefix(trimmed, "+") {
			afterSep = true
			continue
		}
		if !strings.Contains(line, "|") {
			continue
		}
		cols := splitPipe(line)
		if len(cols) < 3 {
			continue
		}
		col1 := strings.Fields(cols[0])
		if len(col1) == 0 {
			continue
		}
		npuIdx, err := strconv.Atoi(col1[0])
		if err != nil {
			continue
		}

		isHeader := afterSep && len(col1) >= 2
		afterSep = false

		if isHeader {
			if cur != nil && npuIdx == cur.Index {
				// Another block of the same NPU (e.g. a per-chip header on
				// 310P3). Keep merging into the existing card; name/health/
				// power/temp come from the first header and are not reset.
				continue
			}
			flush()
			cur = &model.GPUDevice{
				Index:         npuIdx,
				Vendor:        "huawei",
				Name:          col1[1],
				Health:        strings.TrimSpace(cols[1]),
				DriverVersion: driverVer,
			}
			col3 := strings.Fields(cols[2])
			if len(col3) >= 1 {
				cur.PowerW = atofSafe(col3[0])
			}
			if len(col3) >= 2 {
				cur.Temperature = atofSafe(col3[1])
			}
			continue
		}

		// Non-header line within the current card: only chip detail rows
		// whose middle column is a Bus-Id carry useful metrics. Extra rows
		// (alarm-event rows, chip-id rows, etc.) are ignored so their
		// placeholder "0 / 0" values won't clobber real metrics. Memory is
		// summed across chips so a multi-chip card reports its aggregate.
		if cur == nil {
			continue
		}
		busField := strings.TrimSpace(cols[1])
		if !npuBusIDRe.MatchString(busField) {
			continue
		}
		cur.UUID = busField
		col3Fields := strings.Fields(cols[2])
		if len(col3Fields) > 0 {
			cur.Utilization = atofSafe(col3Fields[0])
		}
		matches := npuMemRe.FindAllStringSubmatch(cols[2], -1)
		if len(matches) > 0 {
			last := matches[len(matches)-1]
			cur.MemoryUsedMB += atouSafe(last[1])
			cur.MemoryTotalMB += atouSafe(last[2])
		}
	}
	flush()
	for i := range devices {
		if devices[i].MemoryTotalMB > devices[i].MemoryUsedMB {
			devices[i].MemoryFreeMB = devices[i].MemoryTotalMB - devices[i].MemoryUsedMB
		}
	}
	return devices
}

// collectLspciGPU is the catch-all fallback for any display-class PCI device
// that its vendor's specialized tool didn't capture. Triggers when:
//   - A card is passed through to a virtual machine (bound to vfio-pci): the
//     host-side nvidia-smi / npu-smi can't see it, but lspci still scans the
//     PCIe bus, so the card can at least be enumerated.
//   - No specialized tool exists for the vendor (Intel iGPU, AMD, Iluvatar,
//     Biren, Moore Threads, Cambricon, etc.).
//
// nvidiaOK / huaweiOK flag whether nvidia-smi / npu-smi already reported at
// least one card from that vendor on this host. When true the corresponding
// vendor's lspci entries are skipped to avoid duplicating cards the
// driver-level tool already covered (a working NVIDIA driver enumerates every
// NVIDIA GPU bound to it, so one match implies full coverage). Mixed hosts
// where SOME cards of a vendor are passthrough'd and others aren't are a
// known trade-off: the vendor tool lists only the non-passthrough'd subset
// and lspci is then skipped for that vendor for the passthrough'd ones.
func collectLspciGPU(g *model.GPU, nvidiaOK, huaweiOK bool) {
	if !commandExists("lspci") {
		return
	}
	// -D: always print the PCI domain (so bus IDs are unique across hosts).
	// -nn: print both textual names and numeric vendor:device IDs — needed
	//      to identify brand-new SKUs whose PCI ID isn't in pci.ids yet
	//      (e.g. the 0x2b85 RTX 5090 shows up as "Device 2b85" without it).
	out, err := runCmd("lspci", "-Dnn")
	if err != nil {
		return
	}
	appendLspciGPU(g, out, nvidiaOK, huaweiOK)
}

// appendLspciGPU is the testable core of collectLspciGPU (no shell-out): it
// appends one GPUDevice per display-class PCI function from the given lspci
// output, skipping cards of vendors already covered by their specialized tool
// and continuing the index sequence past whatever g already holds.
func appendLspciGPU(g *model.GPU, out string, nvidiaOK, huaweiOK bool) {
	idx := len(g.Devices)
	for _, d := range parseLspciGPU(out) {
		switch d.Vendor {
		case "nvidia":
			if nvidiaOK {
				continue
			}
		case "huawei":
			if huaweiOK {
				continue
			}
		}
		d.Index = idx
		idx++
		g.Devices = append(g.Devices, d)
	}
}

// parseLspciGPU parses `lspci -Dnn` output and returns one GPUDevice per PCI
// display-class function (VGA / 3D / Display / XGA controller), for ANY
// vendor — not just NVIDIA. The audio subfunction paired with most consumer
// GPUs is skipped so each card is counted once via its display function.
//
// Example input lines:
//
//	0000:16:00.0 VGA compatible controller [0300]: NVIDIA Corporation Device 2b85 [10de:2b85] (rev ff)
//	00:02.0 VGA compatible controller [0300]: Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630] [8086:3e98] (rev 02)
//	0000:43:00.0 3D controller [0302]: NVIDIA Corporation GA100 [A100 SXM4 40GB] [10de:20b5] (rev a1)
//
// Only PCI-level fields are populated: memory/utilization/temperature/power/
// driver/firmware require a host-bound driver and are left empty. Vendor is
// derived from the numeric PCI vendor:device ID when available (preferred,
// requires `lspci -nn`), with a text-based fallback for plain `lspci` output.
func parseLspciGPU(out string) []model.GPUDevice {
	var devices []model.GPUDevice
	idx := 0
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !isDisplayController(line) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		bus := fields[0]
		// lspci without -D omits the domain; normalize to canonical PCI BDF
		// (dom:bus:dev.fn) to align with nvidia-smi / npu-smi bus IDs.
		if strings.Count(bus, ":") == 1 {
			bus = "0000:" + bus
		}

		name, pciID := parseLspciDeviceDesc(line)
		devices = append(devices, model.GPUDevice{
			Index:  idx,
			Vendor: identifyLspciVendor(pciID, name),
			Name:   name,
			UUID:   bus,
			Health: "Unknown",
			Serial: pciID,
		})
		idx++
	}
	return devices
}

// isDisplayController reports whether an lspci line describes a PCI display
// controller (base class 0x03). lspci prints these subclass names: "VGA
// compatible controller" (0x0300), "XGA compatible controller" (0x0301),
// "3D controller" (0x0302, used by compute-only cards like A100/H100),
// "Display controller" (0x0380).
func isDisplayController(line string) bool {
	return strings.Contains(line, "VGA compatible controller") ||
		strings.Contains(line, "XGA compatible controller") ||
		strings.Contains(line, "3D controller") ||
		strings.Contains(line, "Display controller")
}

// pciVendorMap maps well-known PCI vendor IDs (lowercase 4-digit hex) to the
// canonical lowercase vendor name used by the collector. Source:
// https://pci-ids.ucw.cz/ — extend as new vendors appear in the fleet.
var pciVendorMap = map[string]string{
	"10de": "nvidia", // NVIDIA Corporation
	"1002": "amd",    // Advanced Micro Devices, Inc.
	"8086": "intel",  // Intel Corporation
	"19e5": "huawei", // Huawei Technologies Co., Ltd.
}

// identifyLspciVendor resolves the canonical vendor name from the numeric PCI
// vendor:device ID (preferred, available with `lspci -nn`) and falls back to
// substring matching on the textual description (plain `lspci` output). When
// neither yields a known vendor the lowercase 4-digit PCI vendor ID is
// returned (e.g. "1cee") so the device stays identifiable downstream; if even
// that is unavailable it returns "unknown".
func identifyLspciVendor(pciID, name string) string {
	if len(pciID) >= 4 {
		if canonical, ok := pciVendorMap[strings.ToLower(pciID[:4])]; ok {
			return canonical
		}
	}
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "nvidia"):
		return "nvidia"
	case strings.Contains(lower, "huawei"), strings.Contains(lower, "ascend"):
		return "huawei"
	case strings.Contains(lower, "advanced micro devices"), strings.Contains(lower, "amd"):
		return "amd"
	case strings.Contains(lower, "intel"):
		return "intel"
	}
	if len(pciID) >= 4 {
		return strings.ToLower(pciID[:4])
	}
	return "unknown"
}

// parseLspciDeviceDesc extracts the textual device description and the numeric
// PCI vendor:device ID from a single lspci line. The line may carry numeric
// IDs from -nn (preferred) or be plain lspci output.
func parseLspciDeviceDesc(line string) (name, pciID string) {
	// Strip the leading "<bus> <class> [classcode]: " (with -nn) or
	// "<bus> <class>: " (plain lspci) prefix to get the vendor + device text.
	rest := line
	if i := strings.Index(rest, "]: "); i >= 0 {
		rest = rest[i+3:]
	} else if i := strings.Index(rest, ": "); i >= 0 {
		rest = rest[i+2:]
	} else {
		return "", ""
	}
	// The vendor:device bracket (e.g. "[10de:2b85]") is the LAST "[hhhh:hhhh]"
	// on the line: the class-code bracket was stripped above together with
	// the class name. Cut it (and anything after, such as "(rev X)") off the
	// name and capture the PCI ID.
	if i := strings.LastIndex(rest, " ["); i >= 0 {
		tail := rest[i+1:]
		rest = rest[:i]
		if m := pciIDRe.FindStringSubmatch(tail); m != nil {
			pciID = m[1]
		}
	}
	// A trailing "(rev X)" can survive in plain lspci output (no -nn) or when
	// the vendor:device bracket is absent because lspci doesn't know the IDs.
	if i := strings.LastIndex(rest, "(rev"); i >= 0 {
		rest = rest[:i]
	}
	name = strings.TrimSpace(rest)
	return name, pciID
}

func splitPipe(line string) []string {
	parts := strings.Split(line, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func atoiSafe(s string) int {
	v, _ := strconv.Atoi(strings.TrimSpace(s))
	return v
}

func atouSafe(s string) uint64 {
	v, _ := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	return v
}

func atofSafe(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}
