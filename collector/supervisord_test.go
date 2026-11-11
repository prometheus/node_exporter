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

//go:build !nosupervisord

package collector

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

// supervisordServer answers getAllProcessInfo with one running process.
func supervisordServer(t *testing.T, process string, connState func(net.Conn, http.ConnState)) string {
	t.Helper()
	// The 32-bit CI job runs without sockets.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot listen: %v", err)
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, `<?xml version="1.0"?>
<methodResponse><params><param><value><array><data><value><struct>
<member><name>name</name><value><string>%[1]s</string></value></member>
<member><name>group</name><value><string>%[1]s</string></value></member>
<member><name>state</name><value><int>20</int></value></member>
<member><name>start</name><value><int>100</int></value></member>
</struct></value></data></array></value></param></params></methodResponse>`, process)
	}))
	srv.Config.ConnState = connState
	srv.Listener.Close()
	srv.Listener = l
	srv.Start()
	t.Cleanup(srv.Close)
	return srv.URL + "/RPC2"
}

func supervisordRuntime(t *testing.T, url string) *Runtime {
	t.Helper()
	cfg := config.NewConfigWithDefaults()
	cfg.CollectorDisableDefaults = true
	cfg.Collectors = map[string]bool{"supervisord": true}
	cfg.Supervisord.URL = url
	r, err := NewRuntime(cfg, testLogger())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestSupervisordRuntimesKeepTheirOwnURL(t *testing.T) {
	runtimes := map[string]*Runtime{
		"a": supervisordRuntime(t, supervisordServer(t, "a", nil)),
		"b": supervisordRuntime(t, supervisordServer(t, "b", nil)),
	}

	for name, r := range runtimes {
		want := fmt.Sprintf(`# HELP node_supervisord_up Process Up
# TYPE node_supervisord_up gauge
node_supervisord_up{group=%[1]q,name=%[1]q} 1
`, name)
		registry, err := r.Registry()
		if err != nil {
			t.Fatal(err)
		}
		if err := testutil.GatherAndCompare(registry, strings.NewReader(want), "node_supervisord_up"); err != nil {
			t.Errorf("runtime %s: %v", name, err)
		}
	}
}

func TestSupervisordShutdownClosesConnections(t *testing.T) {
	closed := make(chan struct{}, 1)
	url := supervisordServer(t, "a", func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			select {
			case closed <- struct{}{}:
			default:
			}
		}
	})
	r := supervisordRuntime(t, url)
	registry, err := r.Registry()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Gather(); err != nil {
		t.Fatal(err)
	}

	if err := r.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("the connection to supervisord is still open after Shutdown")
	}
}
