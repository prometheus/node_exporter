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

// SFP/QSFP module EEPROM parsing for Digital Optical Monitoring (DOM) /
// Digital Diagnostic Monitoring (DDM) data.
//
// Standards:
//   - SFF-8472: SFP/SFP+ DDM (A0 + A2 EEPROM pages, 512 bytes total)
//   - SFF-8636: QSFP/QSFP28 DOM (page 0, 256 bytes)

package collector

import (
	"encoding/binary"
	"fmt"
)

// SFP/QSFP module identifier values (EEPROM byte 0).
const (
	sfpIdentifierSFP    = 0x03 // SFP/SFP+/SFP28 (SFF-8472)
	sfpIdentifierSFPAlt = 0x0B // SFP+ alternative identifier
	sfpIdentifierQSFP   = 0x0C // QSFP (SFF-8436)
	sfpIdentifierQSFPP  = 0x0D // QSFP+ (SFF-8436)
	sfpIdentifierQSFP28 = 0x11 // QSFP28 (SFF-8636)
)

// sfpLaneMetrics holds per-lane optical monitoring values.
type sfpLaneMetrics struct {
	txBias  float64 // TX laser bias current in amperes
	txPower float64 // TX optical power in watts
	rxPower float64 // RX optical power in watts
}

// sfpMetrics holds parsed DOM/DDM values from a transceiver module.
type sfpMetrics struct {
	temperature float64          // Module temperature in degrees Celsius
	voltage     float64          // Module supply voltage in volts
	lanes       []sfpLaneMetrics // Per-lane metrics (1 lane for SFP, 4 for QSFP)
}

// parseModuleEeprom parses raw EEPROM bytes returned by ethtool GMODULEEEPROM
// and extracts DOM/DDM values.
//
// Returns an error if the data is too short, the identifier is unrecognised, or DDM is not available.
func parseModuleEeprom(data []byte) (sfpMetrics, error) {
	if len(data) < 1 {
		return sfpMetrics{}, fmt.Errorf("module EEPROM data too short (%d bytes)", len(data))
	}

	switch data[0] {
	case sfpIdentifierSFP, sfpIdentifierSFPAlt:
		return parseSFF8472(data)
	case sfpIdentifierQSFP, sfpIdentifierQSFPP, sfpIdentifierQSFP28:
		return parseSFF8636(data)
	default:
		return sfpMetrics{}, fmt.Errorf("unsupported module identifier 0x%02x", data[0])
	}
}

// parseSFF8472 parses SFP/SFP+ DDM data per SFF-8472.
func parseSFF8472(data []byte) (sfpMetrics, error) {
	const (
		a0DiagnosticType = 92   // A0 page: diagnostic monitoring type byte
		ddmSupportBit    = 0x40 // bit 6: DDM implemented

		// Offsets within the full 512-byte dump (A2 page starts at 256).
		a2PageOffset = 256
		valuesOffset = a2PageOffset + 96

		tempOffset    = valuesOffset
		voltageOffset = tempOffset + 2
		txBiasOffset  = voltageOffset + 2
		txPowerOffset = txBiasOffset + 2
		rxPowerOffset = txPowerOffset + 2
		minLen        = rxPowerOffset + 2
	)

	if len(data) < a0DiagnosticType+1 {
		return sfpMetrics{}, fmt.Errorf("SFF-8472 EEPROM too short for diagnostic type byte (%d bytes)", len(data))
	}
	if data[a0DiagnosticType]&ddmSupportBit == 0 {
		return sfpMetrics{}, fmt.Errorf("SFP module does not support DDM (diagnostic type byte: 0x%02x)", data[a0DiagnosticType])
	}
	if len(data) < minLen {
		return sfpMetrics{}, fmt.Errorf("SFF-8472 EEPROM too short for DDM values (%d bytes, need %d)", len(data), minLen)
	}

	temp := parseSFPTemperature(data[tempOffset:])
	voltage := parseSFPVoltage(data[voltageOffset:])

	txBias := parseSFPBias(data[txBiasOffset:])
	txPower := parseSFPPower(data[txPowerOffset:])
	rxPower := parseSFPPower(data[rxPowerOffset:])

	return sfpMetrics{
		temperature: temp,
		voltage:     voltage,
		lanes: []sfpLaneMetrics{
			{txBias: txBias, txPower: txPower, rxPower: rxPower},
		},
	}, nil
}

// parseSFF8636 parses QSFP/QSFP28 DOM data per SFF-8636.
func parseSFF8636(data []byte) (sfpMetrics, error) {
	// All real-time values are on Page 00h.
	const (
		// Table 6-8 Free Side Monitoring Values
		tempOffset    = 22 // Temperature MSB
		voltageOffset = 26 // Supply voltage MSB

		// Table 6-9 Channel Monitoring Values.
		numLanes      = 4
		rxPowerOffset = 34                         // RX power ch1 MSB
		txBiasOffset  = rxPowerOffset + numLanes*2 // TX bias ch1 MSB
		txPowerOffset = txBiasOffset + numLanes*2  // TX power ch1 MSB

		minLen = txPowerOffset + numLanes*2
	)

	if len(data) < minLen {
		return sfpMetrics{}, fmt.Errorf("SFF-8636 EEPROM too short (%d bytes, need %d)", len(data), minLen)
	}

	temp := parseSFPTemperature(data[tempOffset:])
	voltage := parseSFPVoltage(data[voltageOffset:])

	lanes := make([]sfpLaneMetrics, numLanes)
	for i := range numLanes {
		lanes[i] = sfpLaneMetrics{
			rxPower: parseSFPPower(data[rxPowerOffset+i*2:]),
			txBias:  parseSFPBias(data[txBiasOffset+i*2:]),
			txPower: parseSFPPower(data[txPowerOffset+i*2:]),
		}
	}

	return sfpMetrics{
		temperature: temp,
		voltage:     voltage,
		lanes:       lanes,
	}, nil
}

func parseSFPTemperature(b []byte) float64 {
	// SFF-8472
	//
	// Table 9-1  Bit Weights (°C) for Temperature Reporting Registers
	//
	// +----------------------------------+----------------------------------+-------+-------+
	// | Most Significant Byte (byte 96)  | Least Significant Byte (byte 97) |       |       |
	// +------+----+----+----+---+---+---+---+---+---+----+-----+-----+------+-------+-------+
	// | D7   | D6 | D5 | D4 | D3| D2| D1| D0| D7| D6| D5 | D4  | D3  |  D2  |  D1   |  D0   |
	// +------+----+----+----+---+---+---+---+---+---+----+-----+-----+------+-------+-------+
	// | Sign | 64 | 32 | 16 | 8 | 4 | 2 | 1 |1/2|1/4|1/8 |1/16 |1/32 | 1/64 | 1/128 | 1/256 |
	// +------+----+----+----+---+---+---+---+---+---+----+-----+-----+------+-------+-------+
	//
	rawVal := int16(binary.BigEndian.Uint16(b))
	return float64(rawVal) / 256.0
}

func parseSFPVoltage(b []byte) float64 {
	// SFF-8472
	//
	// 9.2 Internal Calibration
	//
	// ...
	// 2) Internally measured transceiver supply voltage. Represented as a 16-bit unsigned integer with the voltage
	// 	defined as the full 16-bit value (0-65535) with LSB equal to 100 microvolts, yielding a total range of 0 V to +6.55 V.
	rawVal := binary.BigEndian.Uint16(b)
	mV := float64(rawVal) / 10
	V := mV / 1000
	return V
}

func parseSFPBias(b []byte) float64 {
	// SFF-8472
	//
	// 9.2 Internal Calibration
	//
	// ...
	// 3) Measured TX bias current in mA. Represented as a 16-bit unsigned integer with the current defined as the full
	// 	16-bit value (0-65535) with LSB equal to 2 microamps, yielding a total range of 0 to 131 mA.
	rawVal := binary.BigEndian.Uint16(b)
	mA := float64(rawVal) / 500
	return mA
}

func parseSFPPower(b []byte) float64 {
	// SFF-8472
	//
	// 9.2 Internal Calibration
	//
	// ...
	// 4) Measured TX output power in mW. Represented as a 16-bit unsigned integer with the power defined as the
	// 	full 16-bit value (0-65535) with LSB equal to 0.1 microwatts, yielding a total range of 0 to 6.5535 mW (-40 to +8.2 dBm).
	// ...
	// 5) Measured RX received optical power in mW. Value can represent either average received power or OMA
	// 	depending upon how bit 3 of byte 92 (A0h) is set. Represented as a 16-bit unsigned integer with the power
	// 	defined as the full 16-bit value (0-65535) with LSB equal to 0.1 microwatts, yielding a total range of 0 to 6.5535 mW (-40 to +8.2 dBm).
	rawVal := binary.BigEndian.Uint16(b)
	mW := float64(rawVal) / 10000
	return mW
}
