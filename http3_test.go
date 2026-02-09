// Copyright 2024 The Prometheus Authors
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

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/quic-go/quic-go/http3"
)

func TestLoadWebConfig(t *testing.T) {
	t.Run("valid TLS config", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		c, err := loadWebConfig(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Paths should be resolved relative to config dir.
		if c.TLSConfig.TLSCertPath != filepath.Join(dir, "server.crt") {
			t.Errorf("expected cert path %q, got %q", filepath.Join(dir, "server.crt"), c.TLSConfig.TLSCertPath)
		}
		if c.TLSConfig.TLSKeyPath != filepath.Join(dir, "server.key") {
			t.Errorf("expected key path %q, got %q", filepath.Join(dir, "server.key"), c.TLSConfig.TLSKeyPath)
		}
		if !c.hasTLS() {
			t.Error("expected hasTLS() to be true")
		}
	})

	t.Run("with basic auth users", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		c, err := loadWebConfig(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(c.Users) != 1 {
			t.Errorf("expected 1 user, got %d", len(c.Users))
		}
		if _, ok := c.Users["alice"]; !ok {
			t.Error("expected user 'alice' to be present")
		}
	})

	t.Run("no TLS config", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		c, err := loadWebConfig(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.hasTLS() {
			t.Error("expected hasTLS() to be false")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := loadWebConfig("/nonexistent/web.yml")
		if err == nil {
			t.Error("expected error for missing file")
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		if err := os.WriteFile(configPath, []byte("invalid: [yaml: content"), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := loadWebConfig(configPath)
		if err == nil {
			t.Error("expected error for invalid YAML")
		}
	})

	t.Run("default TLS versions", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		c, err := loadWebConfig(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if uint16(c.TLSConfig.MinVersion) != tls.VersionTLS12 {
			t.Errorf("expected default min TLS 1.2, got %d", c.TLSConfig.MinVersion)
		}
		if uint16(c.TLSConfig.MaxVersion) != tls.VersionTLS13 {
			t.Errorf("expected default max TLS 1.3, got %d", c.TLSConfig.MaxVersion)
		}
	})
}

// generateTestCert creates a self-signed TLS certificate for testing.
func generateTestCert(t *testing.T, dir string) (certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")

	certFile, err := os.Create(certPath)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	certFile.Close()

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	keyFile, err := os.Create(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	keyFile.Close()

	return certPath, keyPath
}

func TestAltSvcMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	dir := t.TempDir()
	certPath, keyPath := generateTestCert(t, dir)

	// Build a real TLS config so we can start a QUIC listener.
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		t.Fatal(err)
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{http3.NextProtoH3},
	}

	// Open a real UDP listener so SetQUICHeaders has port info.
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer udpConn.Close()

	h3Server := &http3.Server{
		TLSConfig: tlsConf,
		Handler:   inner,
	}
	// Start serving in background so the listener is registered.
	go h3Server.Serve(udpConn)
	defer h3Server.Close()

	// Give the server a moment to register the listener.
	time.Sleep(50 * time.Millisecond)

	middleware := &altSvcMiddleware{
		inner:    inner,
		h3Server: h3Server,
	}

	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// The Alt-Svc header should be set.
	altSvc := rec.Header().Get("Alt-Svc")
	if altSvc == "" {
		t.Error("expected Alt-Svc header to be set")
	}

	if rec.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rec.Body.String())
	}
}

func TestHTTP3AuthHandler(t *testing.T) {
	logger := slog.Default()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	t.Run("no users configured passes through", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("valid credentials", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		// bcrypt hash of "fakepassword"
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		req.SetBasicAuth("alice", "fakepassword")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		req.SetBasicAuth("alice", "wrongpassword")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("no credentials when required", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}

		if rec.Header().Get("WWW-Authenticate") != "Basic" {
			t.Errorf("expected WWW-Authenticate header 'Basic', got %q", rec.Header().Get("WWW-Authenticate"))
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		req.SetBasicAuth("bob", "fakepassword")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("broken config file", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		if err := os.WriteFile(configPath, []byte("invalid: [yaml"), 0644); err != nil {
			t.Fatal(err)
		}

		handler := &http3AuthHandler{
			inner:      inner,
			configPath: configPath,
			logger:     logger,
		}

		req := httptest.NewRequest("GET", "/metrics", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}
	})
}

func TestHTTP3RequiresTLS(t *testing.T) {
	logger := slog.Default()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("no TLS config returns error", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "web.yml")
		content := `
basic_auth_users:
  alice: "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		_, _, err := startHTTP3Server([]string{":0"}, configPath, handler, logger)
		if err == nil {
			t.Error("expected error when TLS is not configured")
		}
	})

	t.Run("missing config file returns error", func(t *testing.T) {
		_, _, err := startHTTP3Server([]string{":0"}, "/nonexistent/web.yml", handler, logger)
		if err == nil {
			t.Error("expected error when config file is missing")
		}
	})
}
