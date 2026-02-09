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
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/crypto/bcrypt"
	"go.yaml.in/yaml/v2"
)

// webConfig mirrors the exporter-toolkit web.Config struct for YAML parsing.
// We need our own copy because getConfig() in exporter-toolkit is unexported.
type webConfig struct {
	TLSConfig tlsConfig                     `yaml:"tls_server_config"`
	Users     map[string]config_util.Secret `yaml:"basic_auth_users"`
}

// tlsConfig mirrors the fields we need from web.TLSConfig.
type tlsConfig struct {
	TLSCert    string             `yaml:"cert"`
	TLSKey     config_util.Secret `yaml:"key"`
	TLSCertPath string            `yaml:"cert_file"`
	TLSKeyPath  string            `yaml:"key_file"`
	ClientCAs   string            `yaml:"client_ca_file"`
	ClientAuth  string            `yaml:"client_auth_type"`
	MinVersion  web.TLSVersion    `yaml:"min_version"`
	MaxVersion  web.TLSVersion    `yaml:"max_version"`
}

// loadWebConfig reads and parses the web config YAML file, resolving relative
// cert paths against the config file's directory. This replicates the behavior
// of the unexported getConfig() in exporter-toolkit.
func loadWebConfig(configPath string) (*webConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	c := &webConfig{
		TLSConfig: tlsConfig{
			MinVersion: web.TLSVersion(tls.VersionTLS12),
			MaxVersion: web.TLSVersion(tls.VersionTLS13),
		},
	}
	if err := yaml.UnmarshalStrict(content, c); err != nil {
		return nil, fmt.Errorf("parsing web config: %w", err)
	}

	// Resolve relative cert paths against config file directory.
	dir := filepath.Dir(configPath)
	c.TLSConfig.TLSCertPath = config_util.JoinDir(dir, c.TLSConfig.TLSCertPath)
	c.TLSConfig.TLSKeyPath = config_util.JoinDir(dir, c.TLSConfig.TLSKeyPath)
	c.TLSConfig.ClientCAs = config_util.JoinDir(dir, c.TLSConfig.ClientCAs)
	return c, nil
}

// hasTLS returns true if the web config has TLS certificate configuration.
func (c *webConfig) hasTLS() bool {
	return c.TLSConfig.TLSCertPath != "" || c.TLSConfig.TLSCert != ""
}

// altSvcMiddleware wraps an http.Handler to inject Alt-Svc headers advertising
// HTTP/3 availability via the given http3.Server.
type altSvcMiddleware struct {
	inner    http.Handler
	h3Server *http3.Server
}

func (m *altSvcMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := m.h3Server.SetQUICHeaders(w.Header()); err == nil {
		// SetQUICHeaders adds Alt-Svc header advertising HTTP/3
	}
	m.inner.ServeHTTP(w, r)
}

// http3AuthHandler provides basic auth for the HTTP/3 path, matching the
// behavior of exporter-toolkit's webHandler. This is needed because
// exporter-toolkit only wraps the TCP server handler internally.
type http3AuthHandler struct {
	inner      http.Handler
	configPath string
	logger     *slog.Logger
	bcryptMtx  sync.Mutex
}

func (h *http3AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := loadWebConfig(h.configPath)
	if err != nil {
		h.logger.Error("Unable to parse configuration", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(c.Users) == 0 {
		h.inner.ServeHTTP(w, r)
		return
	}

	user, pass, auth := r.BasicAuth()
	if auth {
		hashedPassword, validUser := c.Users[user]

		if !validUser {
			// Use a fixed password hash to prevent user enumeration by timing.
			// This is a bcrypt-hashed version of "fakepassword".
			hashedPassword = "$2y$10$QOauhQNbBCuQDKes6eFzPeMqBSjb7Mr5DUmpZ/VcEd00UAV/LDeSi"
		}

		// Use a cache key for logging only; we always verify for correctness.
		_ = strings.Join([]string{
			hex.EncodeToString([]byte(user)),
			hex.EncodeToString([]byte(hashedPassword)),
			hex.EncodeToString([]byte(pass)),
		}, ":")

		h.bcryptMtx.Lock()
		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))
		h.bcryptMtx.Unlock()

		if validUser && err == nil {
			h.inner.ServeHTTP(w, r)
			return
		}
	}

	w.Header().Set("WWW-Authenticate", "Basic")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// startHTTP3Server starts a QUIC/HTTP3 server on the given addresses using
// the TLS configuration from the web config file. It returns the http3.Server
// (for use with Alt-Svc header injection), a cleanup function, and any error.
func startHTTP3Server(addresses []string, configPath string, handler http.Handler, logger *slog.Logger) (*http3.Server, func(), error) {
	c, err := loadWebConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading web config for HTTP/3: %w", err)
	}

	if !c.hasTLS() {
		return nil, nil, fmt.Errorf("HTTP/3 requires TLS configuration in %s", configPath)
	}

	// Build TLS config using the exported ConfigToTLSConfig.
	// We need to construct a web.TLSConfig to pass to it.
	webTLSConfig := &web.TLSConfig{
		TLSCertPath: c.TLSConfig.TLSCertPath,
		TLSKeyPath:  c.TLSConfig.TLSKeyPath,
		TLSCert:     c.TLSConfig.TLSCert,
		TLSKey:      c.TLSConfig.TLSKey,
		ClientCAs:   c.TLSConfig.ClientCAs,
		ClientAuth:  c.TLSConfig.ClientAuth,
		MinVersion:  c.TLSConfig.MinVersion,
		MaxVersion:  c.TLSConfig.MaxVersion,
	}

	tlsConf, err := web.ConfigToTLSConfig(webTLSConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("creating TLS config for HTTP/3: %w", err)
	}

	// QUIC requires TLS 1.3 minimum.
	if tlsConf.MinVersion < tls.VersionTLS13 {
		tlsConf.MinVersion = tls.VersionTLS13
	}

	// Set up GetConfigForClient for TLS config hot-reloading.
	tlsConf.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		newC, err := loadWebConfig(configPath)
		if err != nil {
			return nil, err
		}
		newWebTLS := &web.TLSConfig{
			TLSCertPath: newC.TLSConfig.TLSCertPath,
			TLSKeyPath:  newC.TLSConfig.TLSKeyPath,
			TLSCert:     newC.TLSConfig.TLSCert,
			TLSKey:      newC.TLSConfig.TLSKey,
			ClientCAs:   newC.TLSConfig.ClientCAs,
			ClientAuth:  newC.TLSConfig.ClientAuth,
			MinVersion:  newC.TLSConfig.MinVersion,
			MaxVersion:  newC.TLSConfig.MaxVersion,
		}
		newTLS, err := web.ConfigToTLSConfig(newWebTLS)
		if err != nil {
			return nil, err
		}
		if newTLS.MinVersion < tls.VersionTLS13 {
			newTLS.MinVersion = tls.VersionTLS13
		}
		return newTLS, nil
	}

	// Configure ALPN for HTTP/3.
	quicTLSConf := http3.ConfigureTLSConfig(tlsConf)

	h3Server := &http3.Server{
		TLSConfig: quicTLSConf,
		Handler:   handler,
		Logger:    logger,
	}

	var conns []net.PacketConn
	for _, addr := range addresses {
		if strings.HasPrefix(addr, "vsock://") {
			logger.Warn("HTTP/3 does not support vsock addresses, skipping", "address", addr)
			continue
		}

		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			// Close any already-opened connections.
			for _, c := range conns {
				c.Close()
			}
			return nil, nil, fmt.Errorf("resolving UDP address %q: %w", addr, err)
		}

		conn, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			for _, c := range conns {
				c.Close()
			}
			return nil, nil, fmt.Errorf("listening on UDP %s: %w", addr, err)
		}

		conns = append(conns, conn)
		logger.Info("HTTP/3 (QUIC) listener started", "address", addr)

		go func(c net.PacketConn) {
			if err := h3Server.Serve(c); err != nil && err != http.ErrServerClosed {
				logger.Error("HTTP/3 server error", "err", err)
			}
		}(conn)
	}

	if len(conns) == 0 {
		return nil, nil, fmt.Errorf("no valid addresses for HTTP/3 listener")
	}

	cleanup := func() {
		if err := h3Server.Shutdown(context.Background()); err != nil {
			logger.Error("HTTP/3 server shutdown error", "err", err)
		}
		for _, c := range conns {
			c.Close()
		}
	}

	return h3Server, cleanup, nil
}
