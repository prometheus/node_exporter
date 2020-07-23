// Copyright 2017 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporter_shared

import (
	"crypto/subtle"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	authFileF = kingpin.Flag("web.auth-file", "Path to YAML file with server_user, server_password keys for HTTP Basic authentication "+
		"(overrides HTTP_AUTH environment variable).").String()
)

// basicAuth combines username and password.
type basicAuth struct {
	Username string `yaml:"server_user,omitempty"`
	Password string `yaml:"server_password,omitempty"`
}

// readBasicAuth returns basicAuth from -web.auth-file file, or HTTP_AUTH environment variable, or empty one.
func readBasicAuth() *basicAuth {
	var auth basicAuth
	httpAuth := os.Getenv("HTTP_AUTH")
	switch {
	case *authFileF != "":
		bytes, err := ioutil.ReadFile(*authFileF)
		if err != nil {
			log.Fatalf("cannot read auth file %q: %s", *authFileF, err)
		}
		if err = yaml.Unmarshal(bytes, &auth); err != nil {
			log.Fatalf("cannot parse auth file %q: %s", *authFileF, err)
		}
	case httpAuth != "":
		data := strings.SplitN(httpAuth, ":", 2)
		if len(data) != 2 || data[0] == "" || data[1] == "" {
			log.Fatalf("HTTP_AUTH should be formatted as user:password")
		}
		auth.Username = data[0]
		auth.Password = data[1]
	default:
		// that's fine, return empty one below
	}

	return &auth
}

// basicAuthHandler checks username and password before invoking provided next handler.
type basicAuthHandler struct {
	basicAuth
	nextHandler http.Handler
}

// ServeHTTP implements http.Handler.
func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	usernameOk := subtle.ConstantTimeCompare([]byte(h.Username), []byte(username)) == 1
	passwordOk := subtle.ConstantTimeCompare([]byte(h.Password), []byte(password)) == 1
	if !usernameOk || !passwordOk {
		w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	h.nextHandler.ServeHTTP(w, r)
}

// authHandler wraps provided handler with basic authentication if it is configured.
func authHandler(handler http.Handler) http.Handler {
	auth := readBasicAuth()
	if auth.Username != "" && auth.Password != "" {
		log.Infoln("HTTP Basic authentication is enabled.")
		return &basicAuthHandler{basicAuth: *auth, nextHandler: handler}
	}

	return handler
}

// check interfaces
var (
	_ http.Handler = (*basicAuthHandler)(nil)
)
