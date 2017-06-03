package iostat

import (
	"testing"
)

func TestReadDriveStats(t *testing.T) {
	a, err := ReadDriveStats()
	if err != nil {
		t.Fatal(err)
	}
	if len(a) == 0 {
		t.Errorf("ReadDriveStats() = 0; want >1")
		return
	}
	for _, stat := range a {
		t.Run(stat.Name, func(t *testing.T) {
			if stat.Size <= 0 {
				t.Errorf("Size = %d; want >0", stat.Size)
			}
			if stat.BlockSize <= 0 {
				t.Errorf("BlockSize = %d; want >0", stat.BlockSize)
			}
			if stat.BytesRead <= 0 {
				t.Errorf("BytesRead = %d; want >0", stat.BytesRead)
			}
			t.Logf("Size = %d\n", stat.Size)
			t.Logf("BlockSize = %d\n", stat.BlockSize)
			t.Logf("BytesRead = %d\n", stat.BytesRead)
			t.Logf("BytesWritten = %d\n", stat.BytesWritten)
			t.Logf("NumRead = %d\n", stat.NumRead)
			t.Logf("NumWrite = %d\n", stat.NumWrite)
			t.Logf("TotalReadTime = %v\n", stat.TotalReadTime)
			t.Logf("TotalWriteTime = %v\n", stat.TotalWriteTime)
			t.Logf("ReadLatency = %v\n", stat.ReadLatency)
			t.Logf("WriteLatency = %v\n", stat.WriteLatency)
		})
	}
}
