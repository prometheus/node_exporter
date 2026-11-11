// Copyright The Prometheus Authors
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

//go:build !nonetworkroute

package collector

import (
	"net"
	"testing"

	"github.com/jsimonetti/rtnetlink/v2"
	"golang.org/x/sys/unix"
)

func TestNetworkRouteGateway(t *testing.T) {
	tests := []struct {
		name    string
		gateway net.IP
		via     *rtnetlink.RouteVia
		want    string
	}{
		{
			name:    "same-family IPv4 gateway",
			gateway: net.IPv4(10, 0, 0, 1),
			via:     nil,
			want:    "10.0.0.1",
		},
		{
			name:    "same-family IPv6 gateway",
			gateway: net.ParseIP("fe80::1"),
			via:     nil,
			want:    "fe80::1",
		},
		{
			name:    "RFC 5549 IPv4 route with IPv6 next-hop via RTA_VIA",
			gateway: nil,
			via: &rtnetlink.RouteVia{
				Family: unix.AF_INET6,
				Addr:   net.ParseIP("fe80::7683:efff:fe20:31ed"),
			},
			want: "fe80::7683:efff:fe20:31ed",
		},
		{
			name:    "gateway takes precedence over via",
			gateway: net.IPv4(10, 0, 0, 1),
			via: &rtnetlink.RouteVia{
				Family: unix.AF_INET6,
				Addr:   net.ParseIP("fe80::1"),
			},
			want: "10.0.0.1",
		},
		{
			name:    "empty gateway and nil via",
			gateway: nil,
			via:     nil,
			want:    "",
		},
		{
			name:    "empty gateway and via with empty addr",
			gateway: nil,
			via: &rtnetlink.RouteVia{
				Family: unix.AF_INET6,
				Addr:   net.IP{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := networkRouteGateway(tt.gateway, tt.via)
			if got != tt.want {
				t.Errorf("networkRouteGateway(%v, %v) = %q, want %q", tt.gateway, tt.via, got, tt.want)
			}
		})
	}
}
