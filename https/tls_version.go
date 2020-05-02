// Copyright 2020 The Prometheus Authors
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

// +build go1.12

package https

import (
	"crypto/tls"
	"strings"

	"github.com/pkg/errors"
)

func pickMinVersion(s string) (uint16, error) {
	switch strings.ToUpper(s) {
	case "TLS1.3":
		return tls.VersionTLS13, nil
	case "TLS1.2", "":
		return tls.VersionTLS12, nil
	case "TLS1.1":
		return tls.VersionTLS11, nil
	case "TLS1.0":
		return tls.VersionTLS10, nil
	default:
		return 0, errors.New("unknown min_version: " + s)
	}
}

func pickMaxVersion(s string) (uint16, error) {
	switch strings.ToUpper(s) {
	case "TLS1.3", "":
		return tls.VersionTLS13, nil
	case "TLS1.2":
		return tls.VersionTLS12, nil
	case "TLS1.1":
		return tls.VersionTLS11, nil
	case "TLS1.0":
		return tls.VersionTLS10, nil
	default:
		return 0, errors.New("unknown max_version: " + s)
	}
}
