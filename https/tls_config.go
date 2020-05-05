// Copyright 2019 The Prometheus Authors
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

// Package https allows the implementation of TLS.
package https

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	config_util "github.com/prometheus/common/config"
	"gopkg.in/yaml.v2"
)

var (
	errNoTLSConfig = errors.New("TLS config is not present")
)

type Config struct {
	TLSConfig  TLSStruct                     `yaml:"tls_config"`
	HTTPConfig HTTPStruct                    `yaml:"http_server_config"`
	Users      map[string]config_util.Secret `yaml:"basic_auth_users"`
}

type TLSStruct struct {
	TLSCertPath              string   `yaml:"cert_file"`
	TLSKeyPath               string   `yaml:"key_file"`
	ClientAuth               string   `yaml:"client_auth_type"`
	ClientCAs                string   `yaml:"client_ca_file"`
	CipherSuites             []cipher `yaml:"cipher_suites"`
	CurvePreferences         []curve  `yaml:"curve_preferences"`
	MinVersion               string   `yaml:"min_version"`
	MaxVersion               string   `yaml:"max_version"`
	PreferServerCipherSuites bool     `yaml:"prefer_server_cipher_suites"`
	SessionTicketsDisabled   bool     `yaml:"session_tickets_disabled"`
}

type HTTPStruct struct {
	DisableHTTP2 bool `yaml:"disable_http2"`
}

func getConfig(configPath string) (*Config, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.UnmarshalStrict(content, c)
	return c, err
}

func getTLSConfig(configPath string) (*tls.Config, error) {
	c, err := getConfig(configPath)
	if err != nil {
		return nil, err
	}
	return ConfigToTLSConfig(&c.TLSConfig)
}

// ConfigToTLSConfig generates the golang tls.Config from the TLSStruct config.
func ConfigToTLSConfig(c *TLSStruct) (*tls.Config, error) {
	if c.TLSCertPath == "" && c.TLSKeyPath == "" && c.ClientAuth == "" && c.ClientCAs == "" {
		return nil, errNoTLSConfig
	}

	if c.TLSCertPath == "" {
		return nil, errors.New("missing cert_file")
	}

	if c.TLSKeyPath == "" {
		return nil, errors.New("missing key_file")
	}

	loadCert := func() (*tls.Certificate, error) {
		cert, err := tls.LoadX509KeyPair(c.TLSCertPath, c.TLSKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load X509KeyPair")
		}
		return &cert, nil
	}

	// Confirm that certificate and key paths are valid.
	if _, err := loadCert(); err != nil {
		return nil, err
	}

	minVersion, err := pickMinVersion(c.MinVersion)
	if err != nil {
		return nil, err
	}

	maxVersion, err := pickMaxVersion(c.MaxVersion)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		MinVersion:               minVersion,
		MaxVersion:               maxVersion,
		PreferServerCipherSuites: c.PreferServerCipherSuites,
		SessionTicketsDisabled:   c.SessionTicketsDisabled,
	}

	cfg.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return loadCert()
	}

	var cf []uint16
	for _, c := range c.CipherSuites {
		cf = append(cf, (uint16)(c))
	}
	if len(cf) > 0 {
		cfg.CipherSuites = cf
	}

	var cp []tls.CurveID
	for _, c := range c.CurvePreferences {
		cp = append(cp, (tls.CurveID)(c))
	}
	if len(cp) > 0 {
		cfg.CurvePreferences = cp
	}

	if c.ClientCAs != "" {
		clientCAPool := x509.NewCertPool()
		clientCAFile, err := ioutil.ReadFile(c.ClientCAs)
		if err != nil {
			return nil, err
		}
		clientCAPool.AppendCertsFromPEM(clientCAFile)
		cfg.ClientCAs = clientCAPool
	}

	switch c.ClientAuth {
	case "RequestClientCert":
		cfg.ClientAuth = tls.RequestClientCert
	case "RequireClientCert":
		cfg.ClientAuth = tls.RequireAnyClientCert
	case "VerifyClientCertIfGiven":
		cfg.ClientAuth = tls.VerifyClientCertIfGiven
	case "RequireAndVerifyClientCert":
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	case "", "NoClientCert":
		cfg.ClientAuth = tls.NoClientCert
	default:
		return nil, errors.New("Invalid ClientAuth: " + c.ClientAuth)
	}

	if c.ClientCAs != "" && cfg.ClientAuth == tls.NoClientCert {
		return nil, errors.New("Client CA's have been configured without a Client Auth Policy")
	}

	return cfg, nil
}

// Listen starts the server on the given address. If tlsConfigPath isn't empty the server connection will be started using TLS.
func Listen(server *http.Server, tlsConfigPath string, logger log.Logger) error {
	if tlsConfigPath == "" {
		level.Info(logger).Log("msg", "TLS is disabled and it cannot be enabled on the fly.")
		return server.ListenAndServe()
	}

	if err := validateUsers(tlsConfigPath); err != nil {
		return err
	}

	// Setup basic authentication.
	var handler http.Handler = http.DefaultServeMux
	if server.Handler != nil {
		handler = server.Handler
	}
	server.Handler = &userAuthRoundtrip{
		tlsConfigPath: tlsConfigPath,
		logger:        logger,
		handler:       handler,
	}

	c, err := getConfig(tlsConfigPath)
	if err != nil {
		return err
	}
	config, err := ConfigToTLSConfig(&c.TLSConfig)
	switch err {
	case nil:
		withHTTP2 := "enabled"
		if c.HTTPConfig.DisableHTTP2 {
			withHTTP2 = "disabled"
			server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		}
		// Valid TLS config.
		level.Info(logger).Log("msg", "TLS is enabled and it cannot be disabled on the fly.", "http2", withHTTP2)
	case errNoTLSConfig:
		// No TLS config, back to plain HTTP.
		level.Info(logger).Log("msg", "TLS is disabled and it cannot be enabled on the fly.")
		return server.ListenAndServe()
	default:
		// Invalid TLS config.
		return err
	}

	server.TLSConfig = config

	// Set the GetConfigForClient method of the HTTPS server so that the config
	// and certs are reloaded on new connections.
	server.TLSConfig.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		return getTLSConfig(tlsConfigPath)
	}
	return server.ListenAndServeTLS("", "")
}

type cipher uint16

func (c *cipher) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal((*string)(&s))
	if err != nil {
		return err
	}
	for _, cs := range tls.CipherSuites() {
		if cs.Name == s {
			*c = (cipher)(cs.ID)
			return nil
		}
	}
	return errors.New("unknown cipher: " + s)
}

func (c cipher) MarshalYAML() (interface{}, error) {
	return tls.CipherSuiteName((uint16)(c)), nil
}

type curve tls.CurveID

var curves = map[string]curve{
	"CurveP256": (curve)(tls.CurveP256),
	"CurveP384": (curve)(tls.CurveP384),
	"CurveP521": (curve)(tls.CurveP521),
	"X25519":    (curve)(tls.X25519),
}

func (c *curve) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal((*string)(&s))
	if err != nil {
		return err
	}
	if curveid, ok := curves[s]; ok {
		*c = curveid
		return nil
	}
	return errors.New("unknown curve: " + s)
}

func (c *curve) MarshalYAML() (interface{}, error) {
	for s, curveid := range curves {
		if *c == curveid {
			return s, nil
		}
	}
	return fmt.Sprintf("%v", c), nil
}
