// Copyright 2021 The Prometheus Authors
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

//go:build !noethtool

package collector

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	if a == b {
		return true
	}
	return math.Abs(a-b)/math.Max(math.Abs(a), math.Abs(b)) < 0.1
}

func makeSFF8472(temp int16, voltage, txBias, txPower, rxPower uint16, ddmByte byte) []byte {
	data := make([]byte, 512)
	data[0] = sfpIdentifierSFP
	data[92] = ddmByte
	binary.BigEndian.PutUint16(data[352:], uint16(temp))
	binary.BigEndian.PutUint16(data[354:], voltage)
	binary.BigEndian.PutUint16(data[356:], txBias)
	binary.BigEndian.PutUint16(data[358:], txPower)
	binary.BigEndian.PutUint16(data[360:], rxPower)
	return data
}

func makeSFF8636(temp int16, voltage uint16, rxPower, txBias, txPower [4]uint16) []byte {
	data := make([]byte, 58)
	data[0] = sfpIdentifierQSFP28
	binary.BigEndian.PutUint16(data[22:], uint16(temp))
	binary.BigEndian.PutUint16(data[26:], voltage)
	for i := range 4 {
		binary.BigEndian.PutUint16(data[34+i*2:], rxPower[i])
		binary.BigEndian.PutUint16(data[42+i*2:], txBias[i])
		binary.BigEndian.PutUint16(data[50+i*2:], txPower[i])
	}
	return data
}

var sff8472Cases = []struct {
	name        string
	data        []byte
	wantErr     bool
	wantTemp    float64
	wantVoltage float64
	wantTxBias  float64
	wantTxPower float64
	wantRxPower float64
}{
	{
		name:        "typical values",
		data:        makeSFF8472(25*256, 33000, 10000, 10000, 5000, 0x40),
		wantTemp:    25.0,
		wantVoltage: 33000 * 100e-6,
		wantTxBias:  20.0,
		wantTxPower: 1.0,
		wantRxPower: 0.5,
	},
	{
		name:        "negative temperature",
		data:        makeSFF8472(-10*256, 33000, 5000, 8000, 4000, 0x40),
		wantTemp:    -10.0,
		wantVoltage: 33000 * 100e-6,
		wantTxBias:  10.0,
		wantTxPower: 0.8,
		wantRxPower: 0.4,
	},
	{
		name:        "fractional temperature",
		data:        makeSFF8472(int16(25*256+128), 33000, 0, 0, 0, 0x40),
		wantTemp:    25.5,
		wantVoltage: 33000 * 100e-6,
	},
	{
		name:    "DDM not supported",
		data:    makeSFF8472(0, 0, 0, 0, 0, 0x00),
		wantErr: true,
	},
	{
		name:    "too short for diagnostic type byte",
		data:    make([]byte, 50),
		wantErr: true,
	},
	{
		name: "too short for DDM values",
		data: func() []byte {
			d := make([]byte, 200)
			d[0] = sfpIdentifierSFP
			d[92] = 0x40
			return d
		}(),
		wantErr: true,
	},
	{
		name:    "empty",
		data:    []byte{},
		wantErr: true,
	},
	{
		name: "zero values",
		data: makeSFF8472(0, 0, 0, 0, 0, 0x40),
	},
	{
		name: "alt SFP identifier (0x0B)",
		data: func() []byte {
			d := makeSFF8472(25*256, 33000, 10000, 10000, 5000, 0x40)
			d[0] = sfpIdentifierSFPAlt
			return d
		}(),
		wantTemp:    25.0,
		wantVoltage: 33000 * 100e-6,
		wantTxBias:  20.0,
		wantTxPower: 1.0,
		wantRxPower: 0.5,
	},
}

var sff8636Cases = []struct {
	name        string
	data        []byte
	wantErr     bool
	wantTemp    float64
	wantVoltage float64
	wantLanes   [4]sfpLaneMetrics
}{
	{
		name: "typical values",
		data: makeSFF8636(
			25*256, 33000,
			[4]uint16{5000, 4800, 4900, 5100},
			[4]uint16{10000, 9800, 10200, 10100},
			[4]uint16{9000, 8800, 9100, 9200},
		),
		wantTemp:    25.0,
		wantVoltage: 33000 * 100e-6,
		wantLanes: [4]sfpLaneMetrics{
			{txBias: 20.0, txPower: 0.9, rxPower: 0.5},
			{txBias: 19.6, txPower: 0.88, rxPower: 0.48},
			{txBias: 20.4, txPower: 0.91, rxPower: 0.49},
			{txBias: 20.2, txPower: 0.92, rxPower: 0.51},
		},
	},
	{
		name: "negative temperature",
		data: makeSFF8636(
			-5*256, 33000,
			[4]uint16{}, [4]uint16{}, [4]uint16{},
		),
		wantTemp:    -5.0,
		wantVoltage: 33000 * 100e-6,
	},
	{
		name:    "too short",
		data:    make([]byte, 10),
		wantErr: true,
	},
	{
		name:    "empty",
		data:    []byte{},
		wantErr: true,
	},
	{
		name: "zero values",
		data: makeSFF8636(0, 0, [4]uint16{}, [4]uint16{}, [4]uint16{}),
	},
}

// moduleEepromCases are the table-driven test cases for parseModuleEeprom, also
// used as fuzzing seeds.
var moduleEepromCases = []struct {
	name      string
	data      []byte
	wantErr   bool
	wantLanes int
}{
	{
		name:      "SFP (0x03)",
		data:      makeSFF8472(25*256, 33000, 10000, 10000, 5000, 0x40),
		wantLanes: 1,
	},
	{
		name: "SFP alt (0x0B)",
		data: func() []byte {
			d := makeSFF8472(25*256, 33000, 10000, 10000, 5000, 0x40)
			d[0] = sfpIdentifierSFPAlt
			return d
		}(),
		wantLanes: 1,
	},
	{
		name: "QSFP (0x0C)",
		data: func() []byte {
			d := makeSFF8636(25*256, 33000, [4]uint16{}, [4]uint16{}, [4]uint16{})
			d[0] = sfpIdentifierQSFP
			return d
		}(),
		wantLanes: 4,
	},
	{
		name: "QSFP+ (0x0D)",
		data: func() []byte {
			d := makeSFF8636(25*256, 33000, [4]uint16{}, [4]uint16{}, [4]uint16{})
			d[0] = sfpIdentifierQSFPP
			return d
		}(),
		wantLanes: 4,
	},
	{
		name:      "QSFP28 (0x11)",
		data:      makeSFF8636(25*256, 33000, [4]uint16{}, [4]uint16{}, [4]uint16{}),
		wantLanes: 4,
	},
	{
		name:    "unknown identifier",
		data:    []byte{0x01},
		wantErr: true,
	},
	{
		name:    "empty",
		data:    []byte{},
		wantErr: true,
	},
}

func TestParseSFF8472(t *testing.T) {
	for _, tc := range sff8472Cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSFF8472(tc.data)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got.temperature, tc.wantTemp) {
				t.Errorf("temperature: got %v, want %v", got.temperature, tc.wantTemp)
			}
			if !almostEqual(got.voltage, tc.wantVoltage) {
				t.Errorf("voltage: got %v, want %v", got.voltage, tc.wantVoltage)
			}
			if len(got.lanes) != 1 {
				t.Fatalf("expected 1 lane, got %d", len(got.lanes))
			}
			if !almostEqual(got.lanes[0].txBias, tc.wantTxBias) {
				t.Errorf("txBias: got %v, want %v", got.lanes[0].txBias, tc.wantTxBias)
			}
			if !almostEqual(got.lanes[0].txPower, tc.wantTxPower) {
				t.Errorf("txPower: got %v, want %v", got.lanes[0].txPower, tc.wantTxPower)
			}
			if !almostEqual(got.lanes[0].rxPower, tc.wantRxPower) {
				t.Errorf("rxPower: got %v, want %v", got.lanes[0].rxPower, tc.wantRxPower)
			}
		})
	}
}

func TestParseSFF8636(t *testing.T) {
	for _, tc := range sff8636Cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSFF8636(tc.data)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got.temperature, tc.wantTemp) {
				t.Errorf("temperature: got %v, want %v", got.temperature, tc.wantTemp)
			}
			if !almostEqual(got.voltage, tc.wantVoltage) {
				t.Errorf("voltage: got %v, want %v", got.voltage, tc.wantVoltage)
			}
			if len(got.lanes) != 4 {
				t.Fatalf("expected 4 lanes, got %d", len(got.lanes))
			}
			for i, want := range tc.wantLanes {
				l := got.lanes[i]
				if !almostEqual(l.txBias, want.txBias) {
					t.Errorf("lane %d txBias: got %v, want %v", i, l.txBias, want.txBias)
				}
				if !almostEqual(l.txPower, want.txPower) {
					t.Errorf("lane %d txPower: got %v, want %v", i, l.txPower, want.txPower)
				}
				if !almostEqual(l.rxPower, want.rxPower) {
					t.Errorf("lane %d rxPower: got %v, want %v", i, l.rxPower, want.rxPower)
				}
			}
		})
	}
}

func TestParseModuleEeprom(t *testing.T) {
	for _, tc := range moduleEepromCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseModuleEeprom(tc.data)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got.lanes) != tc.wantLanes {
				t.Errorf("lanes: got %d, want %d", len(got.lanes), tc.wantLanes)
			}
		})
	}
}

func FuzzParseSFF8472(f *testing.F) {
	for _, tc := range sff8472Cases {
		f.Add(tc.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = parseSFF8472(data) //nolint:errcheck
	})
}

func FuzzParseSFF8636(f *testing.F) {
	for _, tc := range sff8636Cases {
		f.Add(tc.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = parseSFF8636(data) //nolint:errcheck
	})
}

func FuzzParseModuleEeprom(f *testing.F) {
	for _, tc := range moduleEepromCases {
		f.Add(tc.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = parseModuleEeprom(data) //nolint:errcheck
	})
}

func Test_parseSFPTemperature(t *testing.T) {
	// Values from Table 9-4 TEC Current Format.
	tests := []struct {
		b    [2]byte
		want float64 // celsius
	}{
		{[2]byte{0x7D, 0}, 125.0},
		{[2]byte{0x19, 0}, 25.0},
		{[2]byte{0x00, 0xFF}, 1},
		{[2]byte{0x00, 0x01}, 0.004},
		{[2]byte{}, 0},
		{[2]byte{0xFF, 0xFF}, -0.004},
		{[2]byte{0xE7, 0x00}, -25.0},
		{[2]byte{0xD8, 0x00}, -40.0},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			got := parseSFPTemperature(tt.b[:])
			if !almostEqual(tt.want, got) {
				t.Fatalf("Expected ~ %v, got %v", tt.want, got)
			}
		})
	}
}

func Test_parseSFPVoltage(t *testing.T) {
	tests := []struct {
		b    [2]byte
		want float64 // volts
	}{
		{[2]byte{0xFF, 0xFF}, 6.55},
		{[2]byte{0x00, 0x01}, 0.0001},
		{[2]byte{}, 0},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			got := parseSFPVoltage(tt.b[:])
			if !almostEqual(tt.want, got) {
				t.Fatalf("Expected ~ %v, got %v", tt.want, got)
			}
		})
	}
}

func Test_parseSFPPower(t *testing.T) {
	tests := []struct {
		b    [2]byte
		want float64 // milliwatts
	}{
		{[2]byte{0xFF, 0xFF}, 6.55},
		{[2]byte{0x00, 0x01}, 0.0001},
		{[2]byte{}, 0},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			got := parseSFPPower(tt.b[:])
			if !almostEqual(tt.want, got) {
				t.Fatalf("Expected ~ %v, got %v", tt.want, got)
			}
		})
	}
}
