package cmdb

import "testing"

// Regression: when `npu-smi info` includes a Process info section that
// actually lists running processes (not just "No running processes found"),
// each process row has the same shape as a card header (numeric NPU idx +
// chip in col1, numeric Process id in the Health column). The parser must
// NOT turn those rows into phantom devices. Previously this sample yielded
// 4 devices (2 real + 2 phantom with name="0", health="<pid>").
const npu910B1WithProcs = `+------------------------------------------------------------------------------------------------+
| npu-smi 25.5.1                   Version: 25.5.1                                               |
+---------------------------+---------------+----------------------------------------------------+
| NPU   Name                | Health        | Power(W)    Temp(C)           Hugepages-Usage(page)|
| Chip                      | Bus-Id        | AICore(%)   Memory-Usage(MB)  HBM-Usage(MB)        |
+===========================+===============+====================================================+
| 0     910B1               | OK            | 105.3       39                0    / 0             |
| 0                         | 0000:C1:00.0  | 0           0    / 0          56017/ 65536         |
+===========================+===============+====================================================+
| 1     910B1               | OK            | 112.3       41                0    / 0             |
| 0                         | 0000:01:00.0  | 0           0    / 0          55358/ 65536         |
+===========================+===============+====================================================+
+---------------------------+---------------+----------------------------------------------------+
| NPU     Chip              | Process id    | Process name             | Process memory(MB)      |
+===========================+===============+====================================================+
| 0       0                 | 3755936       | python                   | 114                     |
| 0       0                 | 3755939       | python                   | 114                     |
| 0       0                 | 3755935       | python                   | 52245                   |
+===========================+===============+====================================================+
| 1       0                 | 3755936       | python                   | 51829                   |
+===========================+===============+====================================================+
`

func TestParseHuaweiNPUWithRunningProcesses(t *testing.T) {
	devs := parseHuaweiNPU(npu910B1WithProcs)
	if len(devs) != 2 {
		t.Fatalf("expected 2 devices, got %d: %+v", len(devs), devs)
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
	if d0.Index != 0 || d0.Name != "910B1" || d0.Health != "OK" {
		t.Errorf("d0 meta mismatch: %+v", d0)
	}
	if d0.UUID != "0000:C1:00.0" {
		t.Errorf("d0 uuid = %q, want 0000:C1:00.0", d0.UUID)
	}
	if d0.MemoryUsedMB != 56017 || d0.MemoryTotalMB != 65536 {
		t.Errorf("d0 mem = %d/%d, want 56017/65536", d0.MemoryUsedMB, d0.MemoryTotalMB)
	}
	if d0.PowerW != 105.3 || d0.Temperature != 39 {
		t.Errorf("d0 power/temp = %v/%v, want 105.3/39", d0.PowerW, d0.Temperature)
	}
	if d0.DriverVersion != "25.5.1" {
		t.Errorf("d0 driver = %q, want 25.5.1", d0.DriverVersion)
	}

	d1 := devs[1]
	if d1.Index != 1 || d1.UUID != "0000:01:00.0" {
		t.Errorf("d1 mismatch: %+v", d1)
	}
	if d1.Health != "OK" {
		t.Errorf("d1 health = %q, want OK (not a process id)", d1.Health)
	}
	if d1.MemoryUsedMB != 55358 || d1.MemoryTotalMB != 65536 {
		t.Errorf("d1 mem = %d/%d, want 55358/65536", d1.MemoryUsedMB, d1.MemoryTotalMB)
	}
}
