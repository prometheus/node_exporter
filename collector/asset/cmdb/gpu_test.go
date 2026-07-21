package cmdb

import (
	"reflect"
	"strings"
	"testing"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
)

const npuSMISample = `+------------------------------------------------------------------------------------------------+
| npu-smi 25.5.2                   Version: 25.5.2                                               |
+---------------------------+---------------+----------------------------------------------------+
| NPU   Name                | Health        | Power(W)    Temp(C)           Hugepages-Usage(page)|
| Chip                      | Bus-Id        | AICore(%)   Memory-Usage(MB)  HBM-Usage(MB)        |
+===========================+===============+====================================================+
| 0     910B2C              | OK            | 94.4        42                0    / 0             |
| 0                         | 0000:66:00.0  | 0           0    / 0          3412 / 65536         |
+===========================+===============+====================================================+
| 1     910B2C              | OK            | 96.6        44                0    / 0             |
| 0                         | 0000:19:00.0  | 0           0    / 0          3404 / 65536         |
+===========================+===============+====================================================+
| 15    910B2C              | OK            | 96.4        43                0    / 0             |
| 0                         | 0000:D0:00.0  | 0           0    / 0          3404 / 65536         |
+===========================+===============+====================================================+
`

const nvidiaSMISample = `0, NVIDIA A100-SXM4-40GB, GPU-1a2b3c4d-aaaa-bbbb-cccc-1234567890ab, 161212300123, 90.00.2C.00.01, 40960, 1234, 39726, 0, 35, 70.5, 535.129.03
1, NVIDIA A100-SXM4-40GB, GPU-2b3c4d5d-bbbb-cccc-dddd-234567890abc, 161212300124, 90.00.2C.00.01, 40960, 2048, 38912, 12, 38, 75.2, 535.129.03
`

const npuBoardSample = `        NPU ID                         : 0
        Product Name                   : IT21HMDB02-B2
        Model                          : NA
        Manufacturer                   : Huawei
        Serial Number                  : 102415083974
        Software Version               : 25.5.2
        Firmware Version               : 7.8.0.7.220
        Compatibility                  : OK
        Board ID                       : 0x65
        PCB ID                         : A
        BOM ID                         : 1
        PCIe Bus Info                  : 0000:66:00.0
        Slot ID                        : 0
`

func TestParseHuaweiNPU(t *testing.T) {
	devs := parseHuaweiNPU(npuSMISample)
	if len(devs) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(devs))
	}

	d0 := devs[0]
	if d0.Index != 0 || d0.Vendor != "huawei" || d0.Name != "910B2C" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.Health != "OK" {
		t.Errorf("d0 health = %q, want OK", d0.Health)
	}
	if d0.PowerW != 94.4 {
		t.Errorf("d0 power = %v, want 94.4", d0.PowerW)
	}
	if d0.Temperature != 42 {
		t.Errorf("d0 temp = %v, want 42", d0.Temperature)
	}
	if d0.UUID != "0000:66:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:66:00.0", d0.UUID)
	}
	if d0.Utilization != 0 {
		t.Errorf("d0 util = %v, want 0", d0.Utilization)
	}
	if d0.MemoryUsedMB != 3412 || d0.MemoryTotalMB != 65536 || d0.MemoryFreeMB != 62124 {
		t.Errorf("d0 mem used/total/free = %d/%d/%d, want 3412/65536/62124",
			d0.MemoryUsedMB, d0.MemoryTotalMB, d0.MemoryFreeMB)
	}
	if d0.DriverVersion != "25.5.2" {
		t.Errorf("d0 driver = %q, want 25.5.2", d0.DriverVersion)
	}
	if !d0.RuntimeMetrics {
		t.Error("d0 RuntimeMetrics should be true (npu-smi provided runtime values)")
	}

	if devs[2].Index != 15 || devs[2].UUID != "0000:D0:00.0" {
		t.Errorf("d2 mismatch: %+v", devs[2])
	}
}

func TestParseNVIDIA(t *testing.T) {
	devs := parseNVIDIA(nvidiaSMISample)
	if len(devs) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devs))
	}

	d0 := devs[0]
	if d0.Index != 0 || d0.Vendor != "nvidia" || d0.Name != "NVIDIA A100-SXM4-40GB" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.UUID != "GPU-1a2b3c4d-aaaa-bbbb-cccc-1234567890ab" {
		t.Errorf("d0 uuid = %q", d0.UUID)
	}
	if d0.Serial != "161212300123" {
		t.Errorf("d0 serial = %q, want 161212300123", d0.Serial)
	}
	if d0.MemoryTotalMB != 40960 || d0.MemoryUsedMB != 1234 || d0.MemoryFreeMB != 39726 {
		t.Errorf("d0 mem = %d/%d/%d", d0.MemoryTotalMB, d0.MemoryUsedMB, d0.MemoryFreeMB)
	}
	if d0.Utilization != 0 || d0.Temperature != 35 || d0.PowerW != 70.5 {
		t.Errorf("d0 util/temp/power = %v/%v/%v", d0.Utilization, d0.Temperature, d0.PowerW)
	}
	if d0.DriverVersion != "535.129.03" {
		t.Errorf("d0 driver = %q", d0.DriverVersion)
	}
	if d0.FirmwareVersion != "90.00.2C.00.01" {
		t.Errorf("d0 firmware = %q, want 90.00.2C.00.01", d0.FirmwareVersion)
	}
	if d0.Health != "OK" {
		t.Errorf("d0 health = %q, want OK", d0.Health)
	}
	if !d0.RuntimeMetrics {
		t.Error("d0 RuntimeMetrics should be true (nvidia-smi provided runtime values)")
	}
}

// TestGPUFieldParity verifies both vendors populate the same set of struct fields.
func TestGPUFieldParity(t *testing.T) {
	nv := parseNVIDIA(nvidiaSMISample)
	npu := parseHuaweiNPU(npuSMISample)
	if len(nv) == 0 || len(npu) == 0 {
		t.Fatal("need both vendors' devices for parity check")
	}

	// Simulate serial enrichment for Huawei (normally done by queryHuaweiNPUSerial)
	npu[0].Serial = "1234567890ABCDEF"

	// Both should have the core identifying fields non-empty.
	coreFields := []string{"Vendor", "Name", "DriverVersion"}
	for _, f := range coreFields {
		nvV := reflect.ValueOf(nv[0]).FieldByName(f).String()
		npuV := reflect.ValueOf(npu[0]).FieldByName(f).String()
		if nvV == "" || npuV == "" {
			t.Errorf("field %s empty: nvidia=%q huawei=%q", f, nvV, npuV)
		}
	}

	// Serial must exist on both after enrichment.
	if nv[0].Serial == "" {
		t.Error("nvidia serial should be non-empty")
	}
	if npu[0].Serial == "" {
		t.Error("huawei serial should be non-empty after enrichment")
	}
}

func TestQueryHuaweiNPUBoard(t *testing.T) {
	// Serial
	m := npuSerialRe.FindStringSubmatch(npuBoardSample)
	if m == nil {
		t.Fatal("serial not found in board sample")
	}
	if got := strings.TrimSpace(m[1]); got != "102415083974" {
		t.Errorf("serial = %q, want 102415083974", got)
	}
	// Firmware version
	mf := npuFWRe.FindStringSubmatch(npuBoardSample)
	if mf == nil {
		t.Fatal("firmware not found in board sample")
	}
	if got := strings.TrimSpace(mf[1]); got != "7.8.0.7.220" {
		t.Errorf("firmware = %q, want 7.8.0.7.220", got)
	}
}

func TestParseHuaweiNPUEmpty(t *testing.T) {
	if devs := parseHuaweiNPU(""); len(devs) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devs))
	}
}

func TestParseNVIDIAEmpty(t *testing.T) {
	if devs := parseNVIDIA(""); len(devs) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devs))
	}
}

// Multi-row sample simulating an npu-smi output where each NPU block contains
// several rows whose first column has 2 tokens (NPU idx + chip/alarm id) and
// whose Health column carries numeric alarm codes. With the naive
// "len(col1) >= 2 => new device" rule this produces 3 entries for NPU 0 and 4
// for NPU 2 (matching the duplicate-index bug observed in production). The
// state-machine parser must collapse each `+====+` block into a single device.
const npuMultiRowSample = `+------------------------------------------------------------------------------------------------+
| npu-smi 25.5.2                   Version: 25.5.2                                               |
+---------------------------+---------------+----------------------------------------------------+
| NPU   Name                | Health        | Power(W)    Temp(C)           Hugepages-Usage(page)|
| Chip                      | Bus-Id        | AICore(%)   Memory-Usage(MB)  HBM-Usage(MB)        |
+===========================+===============+====================================================+
| 0     0                   | 4106296       | 94.2        41                0    / 0             |
| 0     0                   | 4107830       | 0                                0    / 0          |
| 0     0                   | 4107829       | 0                                0    / 0          |
+===========================+===============+====================================================+
| 1     0                   | 4107830       | 96.4        44                0    / 0             |
+===========================+===============+====================================================+
| 2     0                   | 4111621       | 88.0        43                0    / 0             |
| 2     0                   | 4111624       | 0                                0    / 0          |
| 2     0                   | 4111632       | 0                                0    / 0          |
| 2     0                   | 4109801       | 0                                0    / 0          |
+===========================+===============+====================================================+
`

func TestParseHuaweiNPUMultiRowPerCard(t *testing.T) {
	devs := parseHuaweiNPU(npuMultiRowSample)
	if len(devs) != 3 {
		t.Fatalf("expected 3 devices (one per NPU), got %d: %+v", len(devs), devs)
	}

	seen := map[int]int{}
	for _, d := range devs {
		seen[d.Index]++
	}
	for idx, c := range seen {
		if c != 1 {
			t.Errorf("index %d appears %d times, want 1", idx, c)
		}
	}

	d0 := devs[0]
	if d0.Index != 0 || d0.Name != "0" || d0.Health != "4106296" {
		t.Errorf("d0 mismatch: idx=%d name=%q health=%q", d0.Index, d0.Name, d0.Health)
	}
	if d0.PowerW != 94.2 || d0.Temperature != 41 {
		t.Errorf("d0 power/temp = %v/%v, want 94.2/41", d0.PowerW, d0.Temperature)
	}
	if d0.DriverVersion != "25.5.2" {
		t.Errorf("d0 driver = %q, want 25.5.2", d0.DriverVersion)
	}

	d2 := devs[2]
	if d2.Index != 2 || d2.Health != "4111621" {
		t.Errorf("d2 mismatch: idx=%d health=%q", d2.Index, d2.Health)
	}
}

// 910B3-style output where each NPU block has a card header row followed by
// per-chip detail rows (all with numeric first column). Verify each NPU
// produces exactly one device and chip rows are merged, not duplicated.
const npuMultiChipSample = `+------------------------------------------------------------------------------------------------+
| npu-smi 25.5.2                   Version: 25.5.2                                               |
+---------------------------+---------------+----------------------------------------------------+
| NPU   Name                | Health        | Power(W)    Temp(C)           Hugepages-Usage(page)|
| Chip                      | Bus-Id        | AICore(%)   Memory-Usage(MB)  HBM-Usage(MB)        |
+===========================+===============+====================================================+
| 0     910B3               | OK            | 94.2        41                0    / 0             |
| 0                         | 0000:66:00.0  | 0           0    / 0          3411 / 65536         |
| 0     0                   | OK            | 0                                0    / 0          |
| 0     1                   | OK            | 0                                0    / 0          |
+===========================+===============+====================================================+
| 1     910B3               | OK            | 96.4        44                0    / 0             |
| 1                         | 0000:19:00.0  | 0           0    / 0          3404 / 65536         |
+===========================+===============+====================================================+
`

func TestParseHuaweiNPUMultiChipPerCard(t *testing.T) {
	devs := parseHuaweiNPU(npuMultiChipSample)
	if len(devs) != 2 {
		t.Fatalf("expected 2 devices, got %d: %+v", len(devs), devs)
	}

	d0 := devs[0]
	if d0.Index != 0 || d0.Name != "910B3" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.UUID != "0000:66:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:66:00.0", d0.UUID)
	}
	if d0.MemoryUsedMB != 3411 || d0.MemoryTotalMB != 65536 {
		t.Errorf("d0 mem = %d/%d, want 3411/65536", d0.MemoryUsedMB, d0.MemoryTotalMB)
	}

	d1 := devs[1]
	if d1.Index != 1 || d1.UUID != "0000:19:00.0" {
		t.Errorf("d1 mismatch: %+v", d1)
	}
}

// 310P3-style output: each chip of a card gets its own block, with the card
// header repeated per chip and "+---+" (dash) separators between chips. Only
// the column-header boundaries use "+===+". Each chip row carries its own
// Memory-Usage (no HBM column). The parser must collapse the two chips of
// each NPU into a single device and must NOT produce duplicate indexes.
const npu310P3Sample = `+--------------------------------------------------------------------------------------------------------+
| npu-smi 25.5.1                                   Version: 25.5.1                                       |
+-------------------------------+-----------------+------------------------------------------------------+
| NPU     Name                  | Health          | Power(W)     Temp(C)           Hugepages-Usage(page) |
| Chip    Device                | Bus-Id          | AICore(%)    Memory-Usage(MB)                        |
+===============================+=================+======================================================+
| 1       310P3                 | OK              | NA           49                0     / 0             |
| 0       0                     | 0000:01:00.0    | 0            1448 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 1       310P3                 | OK              | NA           45                0     / 0             |
| 1       1                     | 0000:01:00.0    | 0            1511 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 2       310P3                 | OK              | NA           46                0     / 0             |
| 0       2                     | 0000:02:00.0    | 0            1384 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 2       310P3                 | OK              | NA           46                0     / 0             |
| 1       3                     | 0000:02:00.0    | 0            1572 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
`

func TestParseHuaweiNPU310P3PerChipBlocks(t *testing.T) {
	devs := parseHuaweiNPU(npu310P3Sample)
	if len(devs) != 2 {
		t.Fatalf("expected 2 devices (one per NPU), got %d: %+v", len(devs), devs)
	}

	seen := map[int]int{}
	for _, d := range devs {
		seen[d.Index]++
	}
	for idx, c := range seen {
		if c != 1 {
			t.Errorf("index %d appears %d times, want 1", idx, c)
		}
	}

	d0 := devs[0]
	if d0.Index != 1 || d0.Name != "310P3" || d0.Health != "OK" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.UUID != "0000:01:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:01:00.0", d0.UUID)
	}
	// Memory is summed across the two chips of NPU 1.
	if d0.MemoryTotalMB != 44278+43693 || d0.MemoryUsedMB != 1448+1511 {
		t.Errorf("d0 mem = %d/%d, want %d/%d", d0.MemoryUsedMB, d0.MemoryTotalMB, 1448+1511, 44278+43693)
	}
	if d0.MemoryFreeMB != d0.MemoryTotalMB-d0.MemoryUsedMB {
		t.Errorf("d0 free = %d, want %d", d0.MemoryFreeMB, d0.MemoryTotalMB-d0.MemoryUsedMB)
	}
	if d0.Temperature != 49 {
		t.Errorf("d0 temp = %v, want 49 (first chip header)", d0.Temperature)
	}
	if d0.DriverVersion != "25.5.1" {
		t.Errorf("d0 driver = %q, want 25.5.1", d0.DriverVersion)
	}

	d1 := devs[1]
	if d1.Index != 2 || d1.UUID != "0000:02:00.0" {
		t.Errorf("d1 mismatch: %+v", d1)
	}
}

// Variant A of the 310P3 output: "+---+" (dash) is used between EVERY data
// block — including between different NPUs (e.g. NPU 1 chip1 -> NPU 2 chip0).
// Only the column-header boundary uses "+===+". This is the exact format
// observed in production that produced 8 duplicated devices (2 per NPU) with
// the old separator-driven flush. The index-driven parser must yield one
// device per NPU (4 total for NPUs 1/2/4/5) with summed per-chip memory.
const npu310P3AllDashSample = `+--------------------------------------------------------------------------------------------------------+
| npu-smi 25.5.1                                   Version: 25.5.1                                       |
+-------------------------------+-----------------+------------------------------------------------------+
| NPU     Name                  | Health          | Power(W)     Temp(C)           Hugepages-Usage(page) |
| Chip    Device                | Bus-Id          | AICore(%)    Memory-Usage(MB)                        |
+===============================+=================+======================================================+
| 1       310P3                 | OK              | NA           49                0     / 0             |
| 0       0                     | 0000:01:00.0    | 0            1448 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 1       310P3                 | OK              | NA           45                0     / 0             |
| 1       1                     | 0000:01:00.0    | 0            1511 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 2       310P3                 | OK              | NA           46                0     / 0             |
| 0       2                     | 0000:02:00.0    | 0            1384 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 2       310P3                 | OK              | NA           46                0     / 0             |
| 1       3                     | 0000:02:00.0    | 0            1572 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 4       310P3                 | OK              | NA           47                0     / 0             |
| 0       4                     | 0000:81:00.0    | 0            1863 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 4       310P3                 | OK              | NA           48                0     / 0             |
| 1       5                     | 0000:81:00.0    | 0            1105 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 5       310P3                 | OK              | NA           48                0     / 0             |
| 0       6                     | 0000:82:00.0    | 0            1694 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
| 5       310P3                 | OK              | NA           46                0     / 0             |
| 1       7                     | 0000:82:00.0    | 0            1269 / 43693                            |
+-------------------------------+-----------------+------------------------------------------------------+
`

func TestParseHuaweiNPU310P3AllDashSeparators(t *testing.T) {
	devs := parseHuaweiNPU(npu310P3AllDashSample)
	if len(devs) != 4 {
		t.Fatalf("expected 4 devices (one per NPU 1/2/4/5), got %d: %+v", len(devs), devs)
	}

	// No duplicate indexes.
	seen := map[int]int{}
	for _, d := range devs {
		seen[d.Index]++
	}
	wantIdx := map[int]int{1: 1, 2: 1, 4: 1, 5: 1}
	if !reflect.DeepEqual(seen, wantIdx) {
		t.Errorf("index counts = %+v, want %+v", seen, wantIdx)
	}

	byIdx := map[int]model.GPUDevice{}
	for _, d := range devs {
		byIdx[d.Index] = d
	}

	d1 := byIdx[1]
	if d1.UUID != "0000:01:00.0" {
		t.Errorf("NPU1 uuid = %q, want 0000:01:00.0", d1.UUID)
	}
	if d1.MemoryUsedMB != 1448+1511 || d1.MemoryTotalMB != 44278+43693 {
		t.Errorf("NPU1 mem = %d/%d, want %d/%d", d1.MemoryUsedMB, d1.MemoryTotalMB, 1448+1511, 44278+43693)
	}
	if d1.MemoryFreeMB != d1.MemoryTotalMB-d1.MemoryUsedMB {
		t.Errorf("NPU1 free = %d, want %d", d1.MemoryFreeMB, d1.MemoryTotalMB-d1.MemoryUsedMB)
	}
	if d1.Temperature != 49 {
		t.Errorf("NPU1 temp = %v, want 49", d1.Temperature)
	}

	d4 := byIdx[4]
	if d4.UUID != "0000:81:00.0" {
		t.Errorf("NPU4 uuid = %q, want 0000:81:00.0", d4.UUID)
	}
	if d4.MemoryUsedMB != 1863+1105 || d4.MemoryTotalMB != 44278+43693 {
		t.Errorf("NPU4 mem = %d/%d, want %d/%d", d4.MemoryUsedMB, d4.MemoryTotalMB, 1863+1105, 44278+43693)
	}
}

// 910B2C full output: one chip row per NPU (data memory is in the HBM-Usage
// column, the last "d/d" match), "+===+" between NPUs, followed by a
// "No running processes" footer section. The footer must NOT create phantom
// devices.
const npu910B2CSample = `+------------------------------------------------------------------------------------------------+
| npu-smi 25.5.2                   Version: 25.5.2                                               |
+---------------------------+---------------+----------------------------------------------------+
| NPU   Name                | Health        | Power(W)    Temp(C)           Hugepages-Usage(page)|
| Chip                      | Bus-Id        | AICore(%)   Memory-Usage(MB)  HBM-Usage(MB)        |
+===========================+===============+====================================================+
| 0     910B2C              | OK            | 93.9        41                0    / 0             |
| 0                         | 0000:66:00.0  | 0           0    / 0          3410 / 65536         |
+===========================+===============+====================================================+
| 1     910B2C              | OK            | 96.0        43                0    / 0             |
| 0                         | 0000:19:00.0  | 0           0    / 0          3404 / 65536         |
+===========================+===============+====================================================+
| 15    910B2C              | OK            | 96.0        41                0    / 0             |
| 0                         | 0000:D0:00.0  | 0           0    / 0          3401 / 65536         |
+===========================+===============+====================================================+
+---------------------------+---------------+----------------------------------------------------+
| NPU     Chip              | Process id    | Process name             | Process memory(MB)      |
+===========================+===============+====================================================+
| No running processes found in NPU 0                                                            |
+===========================+===============+====================================================+
| No running processes found in NPU 15                                                           |
+===========================+===============+====================================================+
`

func TestParseHuaweiNPU910B2CWithFooter(t *testing.T) {
	devs := parseHuaweiNPU(npu910B2CSample)
	if len(devs) != 3 {
		t.Fatalf("expected 3 devices, got %d: %+v", len(devs), devs)
	}

	d0 := devs[0]
	if d0.Index != 0 || d0.Name != "910B2C" || d0.Health != "OK" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.UUID != "0000:66:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:66:00.0", d0.UUID)
	}
	// Single chip: memory comes from HBM-Usage (the last d/d match), not
	// summed (only one chip row), not the placeholder "0/0".
	if d0.MemoryUsedMB != 3410 || d0.MemoryTotalMB != 65536 || d0.MemoryFreeMB != 65536-3410 {
		t.Errorf("d0 mem used/total/free = %d/%d/%d, want 3410/65536/%d",
			d0.MemoryUsedMB, d0.MemoryTotalMB, d0.MemoryFreeMB, 65536-3410)
	}
	if d0.PowerW != 93.9 || d0.Temperature != 41 {
		t.Errorf("d0 power/temp = %v/%v, want 93.9/41", d0.PowerW, d0.Temperature)
	}
	if d0.Utilization != 0 {
		t.Errorf("d0 util = %v, want 0", d0.Utilization)
	}
	if d0.DriverVersion != "25.5.2" {
		t.Errorf("d0 driver = %q, want 25.5.2", d0.DriverVersion)
	}

	d2 := devs[2]
	if d2.Index != 15 || d2.UUID != "0000:D0:00.0" {
		t.Errorf("d2 mismatch: %+v", d2)
	}
	if d2.MemoryUsedMB != 3401 || d2.MemoryTotalMB != 65536 {
		t.Errorf("d2 mem = %d/%d, want 3401/65536", d2.MemoryUsedMB, d2.MemoryTotalMB)
	}
}

// Real-world `lspci | grep -i nvidia` output from a host whose 8 RTX 5090s
// (PCI device 0x2b85) are passed through to guest VMs: nvidia-smi on the
// host reports zero devices, so the collector falls back to lspci. Each GPU
// shows up as a pair — VGA controller + audio subfunction — and the audio
// functions must NOT be counted as separate devices. The first card carries
// "(rev ff)" which is the typical signature of a card bound to vfio-pci
// (its config-space power state is D3cold), but the parser doesn't need to
// care about the revision.
const lspciPlainSample = `16:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev ff)
16:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev ff)
27:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
27:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
38:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
38:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
5a:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
5a:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
98:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
98:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
a8:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
a8:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
b8:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
b8:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
d8:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev a1)
d8:00.1 Audio device: NVIDIA Corporation Device 22e8 (rev a1)
`

// Same setup, but with `lspci -Dnn` output (the flags the collector actually
// uses): the PCI domain is always printed and numeric vendor:device IDs are
// included in brackets so brand-new SKUs can still be identified.
const lspciDnnSample = `0000:16:00.0 VGA compatible controller [0300]: NVIDIA Corporation Device 2b85 [10de:2b85] (rev ff)
0000:16:00.1 Audio device [0403]: NVIDIA Corporation Device 22e8 [10de:22e8] (rev ff)
0000:27:00.0 VGA compatible controller [0300]: NVIDIA Corporation Device 2b85 [10de:2b85] (rev a1)
0000:27:00.1 Audio device [0403]: NVIDIA Corporation Device 22e8 [10de:22e8] (rev a1)
`

// When lspci's pci.ids database knows the SKU, the textual description is the
// real product name (here a GA102 RTX 3090) and should be preserved verbatim.
const lspciNamedSample = `0000:01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GA102 [GeForce RTX 3090] [10de:2204] (rev a1)
`

// A100-SXM4 in pass-through: shows up as "3D controller" (compute-only, no
// VGA), which the parser must still pick up.
const lspci3DControllerSample = `0000:43:00.0 3D controller [0302]: NVIDIA Corporation GA100 [A100 SXM4 40GB] [10de:20b5] (rev a1)
0000:43:00.1 Audio device [0403]: NVIDIA Corporation GA100 High Definition Audio [10de:20b5] (rev a1)
`

// A typical mixed-vendor host: an Intel integrated GPU alongside discrete
// cards from AMD and NVIDIA (here an NVIDIA audio subfunction with no display
// sibling on this host, e.g. a USB-C display output mux). All display-class
// functions of any vendor are captured; audio subfunctions are not.
const lspciMixedVendorSample = `00:02.0 VGA compatible controller [0300]: Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630] [8086:3e98] (rev 02)
00:1f.3 Audio device [0403]: Intel Corporation Cannon Lake PCH cAVS [8086:a348] (rev 10)
01:00.0 VGA compatible controller [0300]: Advanced Micro Devices, Inc. [AMD/ATI] Navi 21 [Radeon RX 6800/6800 XT / 6900 XT] [1002:73bf] (rev c1)
01:00.1 Audio device [0403]: Advanced Micro Devices, Inc. [AMD/ATI] Navi 21 HDMI Audio [1002:ab28] (rev c1)
02:00.0 Audio device [0403]: NVIDIA Corporation GP107GL High Definition Audio Controller [10de:0fb5] (rev a1)
`

// An accelerator from a vendor the collector's pciVendorMap doesn't know —
// the parser must still emit a device, surfacing the (unknown) PCI vendor ID
// in the vendor field so it remains identifiable / mergeable downstream.
const lspciUnknownVendorSample = `0000:41:0c.0 Display controller [0380]: Iluvatar CoreX Triton X100 [1cee:0100] (rev 01)
`

func TestParseLspciGPUPlain(t *testing.T) {
	devs := parseLspciGPU(lspciPlainSample)
	if len(devs) != 8 {
		t.Fatalf("expected 8 GPUs (one per VGA controller, audio skipped), got %d: %+v",
			len(devs), devs)
	}

	// Indices run 0..7 contiguously. Plain lspci has no -nn, so there's no
	// PCI vendor:device ID to recover; Serial stays empty but the textual
	// "Device 2b85" description (which still includes "NVIDIA") lets the
	// text-fallback path recognize the vendor as "nvidia".
	for i, d := range devs {
		if d.Index != i {
			t.Errorf("dev %d index = %d, want %d", i, d.Index, i)
		}
		if d.Vendor != "nvidia" {
			t.Errorf("dev %d vendor = %q, want nvidia", i, d.Vendor)
		}
		if d.Serial != "" {
			t.Errorf("dev %d serial = %q, want empty (no -nn)", i, d.Serial)
		}
		if d.Health != "Unknown" {
			t.Errorf("dev %d health = %q, want Unknown (no driver on host)", i, d.Health)
		}
		if d.DriverVersion != "" || d.FirmwareVersion != "" {
			t.Errorf("dev %d driver/firmware should be empty: %q/%q",
				i, d.DriverVersion, d.FirmwareVersion)
		}
		if d.MemoryTotalMB != 0 || d.Utilization != 0 || d.Temperature != 0 {
			t.Errorf("dev %d runtime metrics should be 0: mem=%d util=%v temp=%v",
				i, d.MemoryTotalMB, d.Utilization, d.Temperature)
		}
		if d.RuntimeMetrics {
			t.Errorf("dev %d RuntimeMetrics should be false for lspci-fallback devices", i)
		}
	}

	d0 := devs[0]
	if d0.Name != "NVIDIA Corporation Device 2b85" {
		t.Errorf("d0 name = %q, want 'NVIDIA Corporation Device 2b85'", d0.Name)
	}
	if d0.UUID != "0000:16:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:16:00.0 (domain added)", d0.UUID)
	}

	// Bus IDs unique and in the order they appeared.
	wantBuses := []string{
		"0000:16:00.0", "0000:27:00.0", "0000:38:00.0", "0000:5a:00.0",
		"0000:98:00.0", "0000:a8:00.0", "0000:b8:00.0", "0000:d8:00.0",
	}
	seen := map[string]bool{}
	for i, d := range devs {
		if d.UUID != wantBuses[i] {
			t.Errorf("dev %d uuid = %q, want %q", i, d.UUID, wantBuses[i])
		}
		if seen[d.UUID] {
			t.Errorf("duplicate uuid %q", d.UUID)
		}
		seen[d.UUID] = true
	}
}

func TestParseLspciGPUDnn(t *testing.T) {
	devs := parseLspciGPU(lspciDnnSample)
	if len(devs) != 2 {
		t.Fatalf("expected 2 GPUs, got %d: %+v", len(devs), devs)
	}

	d0 := devs[0]
	if d0.Name != "NVIDIA Corporation Device 2b85" {
		t.Errorf("d0 name = %q", d0.Name)
	}
	if d0.UUID != "0000:16:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:16:00.0 (domain preserved)", d0.UUID)
	}
	// With -nn the numeric PCI vendor:device ID is recovered into Serial so
	// the SKU is identifiable even though pci.ids doesn't know "0x2b85". The
	// numeric vendor "10de" is also recognized via pciVendorMap as "nvidia".
	if d0.Serial != "10de:2b85" {
		t.Errorf("d0 serial = %q, want 10de:2b85", d0.Serial)
	}
	if d0.Vendor != "nvidia" {
		t.Errorf("d0 vendor = %q, want nvidia (from PCI vendor ID 10de)", d0.Vendor)
	}
	if d1 := devs[1]; d1.UUID != "0000:27:00.0" || d1.Serial != "10de:2b85" {
		t.Errorf("d1 mismatch: %+v", d1)
	}
}

func TestParseLspciGPUNamedModel(t *testing.T) {
	devs := parseLspciGPU(lspciNamedSample)
	if len(devs) != 1 {
		t.Fatalf("expected 1 GPU, got %d: %+v", len(devs), devs)
	}
	d := devs[0]
	// When lspci has the real model name it should be preserved verbatim —
	// don't strip brackets, don't strip the "NVIDIA Corporation" prefix.
	if d.Name != "NVIDIA Corporation GA102 [GeForce RTX 3090]" {
		t.Errorf("name = %q", d.Name)
	}
	if d.UUID != "0000:01:00.0" {
		t.Errorf("uuid = %q", d.UUID)
	}
	if d.Serial != "10de:2204" {
		t.Errorf("serial = %q, want 10de:2204", d.Serial)
	}
}

func TestParseLspciGPU3DController(t *testing.T) {
	devs := parseLspciGPU(lspci3DControllerSample)
	if len(devs) != 1 {
		t.Fatalf("expected 1 GPU (3D controller, audio skipped), got %d: %+v",
			len(devs), devs)
	}
	d := devs[0]
	if d.Name != "NVIDIA Corporation GA100 [A100 SXM4 40GB]" {
		t.Errorf("name = %q", d.Name)
	}
	if d.UUID != "0000:43:00.0" {
		t.Errorf("uuid = %q", d.UUID)
	}
	if d.Serial != "10de:20b5" {
		t.Errorf("serial = %q, want 10de:20b5", d.Serial)
	}
}

// Any display-class device of ANY vendor is captured — not just NVIDIA. This
// is the case that distinguishes the catch-all lspci fallback from the older
// NVIDIA-only parser: an Intel iGPU and a discrete AMD Radeon show up
// alongside an NVIDIA audio subfunction (which is skipped), each with the
// right canonical vendor derived from its PCI vendor ID.
func TestParseLspciGPUMixedVendor(t *testing.T) {
	devs := parseLspciGPU(lspciMixedVendorSample)
	if len(devs) != 2 {
		t.Fatalf("expected 2 display devices (Intel + AMD, audio/NVIDIA-audio skipped), got %d: %+v",
			len(devs), devs)
	}

	d0 := devs[0]
	if d0.Vendor != "intel" {
		t.Errorf("d0 vendor = %q, want intel", d0.Vendor)
	}
	if d0.Name != "Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630]" {
		t.Errorf("d0 name = %q", d0.Name)
	}
	if d0.UUID != "0000:00:02.0" {
		t.Errorf("d0 uuid = %q, want 0000:00:02.0 (domain added)", d0.UUID)
	}
	if d0.Serial != "8086:3e98" {
		t.Errorf("d0 serial = %q, want 8086:3e98", d0.Serial)
	}

	d1 := devs[1]
	if d1.Vendor != "amd" {
		t.Errorf("d1 vendor = %q, want amd", d1.Vendor)
	}
	if d1.Name != "Advanced Micro Devices, Inc. [AMD/ATI] Navi 21 [Radeon RX 6800/6800 XT / 6900 XT]" {
		t.Errorf("d1 name = %q", d1.Name)
	}
	if d1.UUID != "0000:01:00.0" {
		t.Errorf("d1 uuid = %q", d1.UUID)
	}
	if d1.Serial != "1002:73bf" {
		t.Errorf("d1 serial = %q, want 1002:73bf", d1.Serial)
	}
}

// A device whose vendor isn't in the collector's pciVendorMap: it must still
// be emitted (rather than silently dropped) and the PCI vendor ID is surfaced
// as the canonical vendor name so downstream tools can merge it later.
func TestParseLspciGPUUnknownVendor(t *testing.T) {
	devs := parseLspciGPU(lspciUnknownVendorSample)
	if len(devs) != 1 {
		t.Fatalf("expected 1 device, got %d: %+v", len(devs), devs)
	}
	d := devs[0]
	if d.Vendor != "1cee" {
		t.Errorf("vendor = %q, want 1cee (PCI vendor ID hex fallback)", d.Vendor)
	}
	if d.Name != "Iluvatar CoreX Triton X100" {
		t.Errorf("name = %q", d.Name)
	}
	if d.UUID != "0000:41:0c.0" {
		t.Errorf("uuid = %q", d.UUID)
	}
	if d.Serial != "1cee:0100" {
		t.Errorf("serial = %q, want 1cee:0100", d.Serial)
	}
}

func TestParseLspciGPUEmpty(t *testing.T) {
	if devs := parseLspciGPU(""); len(devs) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devs))
	}
}

// Audio devices with no sibling display function (e.g. an NVIDIA HDMI audio
// controller on a headless host) must NOT produce phantom GPUs.
func TestParseLspciGPUAudioOnlyNoDevices(t *testing.T) {
	in := `00:1f.3 Audio device [0403]: Intel Corporation Cannon Lake PCH cAVS [8086:a348] (rev 10)
01:00.0 Audio device [0403]: NVIDIA Corporation GP107GL High Definition Audio Controller [10de:0fb5] (rev a1)
`
	if devs := parseLspciGPU(in); len(devs) != 0 {
		t.Fatalf("expected 0 devices (no display-class functions), got %d: %+v", len(devs), devs)
	}
}

func TestIdentifyLspciVendor(t *testing.T) {
	cases := []struct {
		pciID string
		name  string
		want  string
	}{
		// Known PCI vendor IDs — preferred over text matching.
		{"10de:2b85", "NVIDIA Corporation Device 2b85", "nvidia"},
		{"8086:3e98", "Intel Corporation CoffeeLake-S GT2", "intel"},
		{"1002:73bf", "Advanced Micro Devices, Inc. [AMD/ATI] Navi 21", "amd"},
		{"1002:AB28", "AMD audio (uppercase hex ID, case-insensitive)", "amd"},
		{"19e5:abcd", "Huawei NPU device", "huawei"},
		// Plain lspci (no -nn) — fall back to textual description.
		{"", "NVIDIA Corporation Device 2b85", "nvidia"},
		{"", "Advanced Micro Devices, Inc. [AMD/ATI]", "amd"},
		// Unknown PCI vendor: surface the hex ID so it stays identifiable.
		{"1cee:0100", "Iluvatar CoreX Triton X100", "1cee"},
		// Empty PCI ID + unrecognized text: last-resort "unknown".
		{"", "Some unknown controller", "unknown"},
	}
	for i, c := range cases {
		got := identifyLspciVendor(c.pciID, c.name)
		if got != c.want {
			t.Errorf("case %d identifyLspciVendor(%q, %q) = %q, want %q",
				i, c.pciID, c.name, got, c.want)
		}
	}
}

func TestParseLspciDeviceDesc(t *testing.T) {
	cases := []struct {
		line  string
		name  string
		pciID string
	}{
		{
			line:  "0000:16:00.0 VGA compatible controller [0300]: NVIDIA Corporation Device 2b85 [10de:2b85] (rev ff)",
			name:  "NVIDIA Corporation Device 2b85",
			pciID: "10de:2b85",
		},
		{
			line:  "16:00.0 VGA compatible controller: NVIDIA Corporation Device 2b85 (rev ff)",
			name:  "NVIDIA Corporation Device 2b85",
			pciID: "",
		},
		{
			line:  "0000:01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GA102 [GeForce RTX 3090] [10de:2204] (rev a1)",
			name:  "NVIDIA Corporation GA102 [GeForce RTX 3090]",
			pciID: "10de:2204",
		},
		{
			line:  "0000:43:00.0 3D controller [0302]: NVIDIA Corporation GA100 [A100 SXM4 40GB] [10de:20b5] (rev a1)",
			name:  "NVIDIA Corporation GA100 [A100 SXM4 40GB]",
			pciID: "10de:20b5",
		},
		{
			line:  "00:02.0 VGA compatible controller [0300]: Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630] [8086:3e98] (rev 02)",
			name:  "Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630]",
			pciID: "8086:3e98",
		},
		{
			line:  "0000:41:0c.0 Display controller [0380]: Iluvatar CoreX Triton X100 [1cee:0100] (rev 01)",
			name:  "Iluvatar CoreX Triton X100",
			pciID: "1cee:0100",
		},
	}
	for i, c := range cases {
		name, pciID := parseLspciDeviceDesc(c.line)
		if name != c.name {
			t.Errorf("case %d name = %q, want %q", i, name, c.name)
		}
		if pciID != c.pciID {
			t.Errorf("case %d pciID = %q, want %q", i, pciID, c.pciID)
		}
	}
}

// The lspci catch-all must skip cards of vendors already covered by their
// specialized tool (otherwise every working NVIDIA driver would be double
// counted once by nvidia-smi and again by lspci). It must still add cards of
// any OTHER vendor on the host. Indices continue past whatever the
// specialized tools already populated so they stay globally unique.
func TestAppendLspciGPUDedup(t *testing.T) {
	// Simulate nvidia-smi having already contributed 2 NVIDIA devices.
	g := &model.GPU{Devices: []model.GPUDevice{
		{Index: 0, Vendor: "nvidia", Name: "NVIDIA GeForce RTX 4090", UUID: "GPU-aaa"},
		{Index: 1, Vendor: "nvidia", Name: "NVIDIA GeForce RTX 4090", UUID: "GPU-bbb"},
	}}

	// Same host's lspci: 2 NVIDIA VGA controllers (already covered by
	// nvidia-smi — must be skipped), 1 Intel iGPU (must be added).
	const lspciOut = `00:02.0 VGA compatible controller [0300]: Intel Corporation CoffeeLake-S GT2 [UHD Graphics 630] [8086:3e98] (rev 02)
01:00.0 VGA compatible controller [0300]: NVIDIA Corporation AD102 [GeForce RTX 4090] [10de:2684] (rev a1)
01:00.1 Audio device [0403]: NVIDIA Corporation AD102 High Definition Audio [10de:22ba] (rev a1)
02:00.0 VGA compatible controller [0300]: NVIDIA Corporation AD102 [GeForce RTX 4090] [10de:2684] (rev a1)
`
	appendLspciGPU(g, lspciOut, /*nvidiaOK=*/ true, /*huaweiOK=*/ false)

	if len(g.Devices) != 3 {
		t.Fatalf("expected 3 devices (2 nvidia-smi + 1 Intel), got %d: %+v",
			len(g.Devices), g.Devices)
	}

	// First two devices (from nvidia-smi) untouched.
	if g.Devices[0].Vendor != "nvidia" || g.Devices[0].UUID != "GPU-aaa" {
		t.Errorf("d0 should be the nvidia-smi entry, got %+v", g.Devices[0])
	}
	if g.Devices[1].Vendor != "nvidia" || g.Devices[1].UUID != "GPU-bbb" {
		t.Errorf("d1 should be the nvidia-smi entry, got %+v", g.Devices[1])
	}

	// Intel added via lspci, with index continuing past the existing entries.
	d2 := g.Devices[2]
	if d2.Index != 2 {
		t.Errorf("d2 index = %d, want 2 (continues past nvidia-smi's 0/1)", d2.Index)
	}
	if d2.Vendor != "intel" {
		t.Errorf("d2 vendor = %q, want intel", d2.Vendor)
	}
	if d2.UUID != "0000:00:02.0" {
		t.Errorf("d2 uuid = %q, want 0000:00:02.0", d2.UUID)
	}
}

// When the specialized tool returned nothing for NVIDIA (i.e. cards are
// passthrough'd to guests) lspci must add ALL the NVIDIA cards it sees —
// including the audio-sibling dedup — AND any other-vendor cards.
func TestAppendLspciGPUNvidiaPassthrough(t *testing.T) {
	g := &model.GPU{Devices: []model.GPUDevice{}}
	appendLspciGPU(g, lspciPlainSample, /*nvidiaOK=*/ false, /*huaweiOK=*/ false)

	if len(g.Devices) != 8 {
		t.Fatalf("expected 8 NVIDIA GPUs from lspci (audio skipped), got %d: %+v",
			len(g.Devices), g.Devices)
	}
	for i, d := range g.Devices {
		if d.Index != i {
			t.Errorf("dev %d index = %d, want %d", i, d.Index, i)
		}
		if d.Vendor != "nvidia" {
			t.Errorf("dev %d vendor = %q, want nvidia", i, d.Vendor)
		}
	}
}

// Huawei NPU passthrough mirrors the NVIDIA case: when npu-smi finds nothing
// on the host (cards reserved for guests), lspci enumerates the same cards.
func TestAppendLspciGPUHuaweiPassthrough(t *testing.T) {
	g := &model.GPU{Devices: []model.GPUDevice{}}
	const lspciOut = `0000:c1:00.0 Display controller [0380]: Huawei Technologies Co., Ltd. Ascend NPU [19e5:abcd] (rev 01)
0000:c1:00.1 Processing accelerators [1200]: Huawei Technologies Co., Ltd. Ascend NPU [19e5:abce] (rev 01)
`
	appendLspciGPU(g, lspciOut, /*nvidiaOK=*/ false, /*huaweiOK=*/ false)

	if len(g.Devices) != 1 {
		t.Fatalf("expected 1 NPU (display controller only; Processing accelerators class is not display), got %d: %+v",
			len(g.Devices), g.Devices)
	}
	d := g.Devices[0]
	if d.Vendor != "huawei" {
		t.Errorf("vendor = %q, want huawei", d.Vendor)
	}
	if d.UUID != "0000:c1:00.0" {
		t.Errorf("uuid = %q, want 0000:c1:00.0", d.UUID)
	}
	if d.Serial != "19e5:abcd" {
		t.Errorf("serial = %q, want 19e5:abcd", d.Serial)
	}
}
