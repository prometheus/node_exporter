// Copyright 2015 The Prometheus Authors
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

// +build !noudpqueues

package collector

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_parseUDPqueues(t *testing.T) {
	noFile, _ := os.Open("follow the white rabbit")
	defer noFile.Close()

	if _, err := parseUDPqueues(noFile); err == nil {
		t.Fatal("expected an error, but none occurred")
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]float64
		wantErr bool
	}{
		{
			name: "reading valid lines, no issue should happened",
			args: args{
				strings.NewReader(
					"sl  local_address rem_address   st tx_queue rx_queue \n" +
						"1: 00000000:0000 00000000:0000 07 00000000:00000001 \n" +
						"2: 00000000:0000 00000000:0000 07 00000002:00000001 \n"),
			},
			want:    map[string]float64{"tx_queue": 2, "rx_queue": 2},
			wantErr: false,
		},
		{
			name: "error case - invalid line - number of fields < 5",
			args: args{
				strings.NewReader(
					"sl  local_address rem_address   st tx_queue rx_queue \n" +
						"1: 00000000:0000 00000000:0000 07 00000000:00000001 \n" +
						"2: 00000000:0000 00000000:0000 07 \n"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - cannot parse line - missing colon",
			args: args{
				strings.NewReader(
					"sl  local_address rem_address   st tx_queue rx_queue \n" +
						"1: 00000000:0000 00000000:0000 07 00000000:00000001 \n" +
						"2: 00000000:0000 00000000:0000 07 0000000200000001 \n"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - parse tx_queue - not an valid hex",
			args: args{
				strings.NewReader(
					"sl  local_address rem_address   st tx_queue rx_queue \n" +
						"1: 00000000:0000 00000000:0000 07 0000000G:00000001 \n" +
						"2: 00000000:0000 00000000:0000 07 00000002:00000001 \n"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error case - parse rx_queue - not an valid hex",
			args: args{
				strings.NewReader(
					"sl  local_address rem_address   st tx_queue rx_queue \n" +
						"1: 00000000:0000 00000000:0000 07 00000000:00000001 \n" +
						"2: 00000000:0000 00000000:0000 07 00000002:0000000G \n"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUDPqueues(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUDPqueues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tt.want) != len(got) {
				t.Errorf("parseUDPqueues() = %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if _, ok := got[k]; !ok {
					t.Errorf("parseUDPqueues() = %v, want %v", got, tt.want)
				}
				if got[k] != v {
					t.Errorf("parseUDPqueues() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_getUDPqueues(t *testing.T) {
	type args struct {
		statsFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "file found",
			args:    args{statsFile: "fixtures/proc/net/udp"},
			wantErr: false,
		},
		{
			name:    "error case - file not found",
			args:    args{statsFile: "somewhere over the rainbow"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getUDPqueues(tt.args.statsFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUDPqueues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// other cases are covered by Test_getUDPqueues()
		})
	}
}
