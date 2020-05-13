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

// +build go1.14

package https

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sync"
	"testing"
	"time"
)

var (
	port       = getPort()
	testlogger = &testLogger{}

	ErrorMap = map[string]*regexp.Regexp{
		"HTTP Response to HTTPS":       regexp.MustCompile(`server gave HTTP response to HTTPS client`),
		"No such file":                 regexp.MustCompile(`no such file`),
		"Invalid argument":             regexp.MustCompile(`invalid argument`),
		"YAML error":                   regexp.MustCompile(`yaml`),
		"Invalid ClientAuth":           regexp.MustCompile(`invalid ClientAuth`),
		"TLS handshake":                regexp.MustCompile(`tls`),
		"HTTP Request to HTTPS server": regexp.MustCompile(`HTTP`),
		"Invalid CertPath":             regexp.MustCompile(`missing cert_file`),
		"Invalid KeyPath":              regexp.MustCompile(`missing key_file`),
		"ClientCA set without policy":  regexp.MustCompile(`Client CA's have been configured without a Client Auth Policy`),
		"Bad password":                 regexp.MustCompile(`hashedSecret too short to be a bcrypted password`),
		"Unauthorized":                 regexp.MustCompile(`Unauthorized`),
		"Forbidden":                    regexp.MustCompile(`Forbidden`),
		"Handshake failure":            regexp.MustCompile(`handshake failure`),
		"Unknown cipher":               regexp.MustCompile(`unknown cipher`),
		"Unknown curve":                regexp.MustCompile(`unknown curve`),
		"Unknown TLS version":          regexp.MustCompile(`unknown TLS version`),
		"No HTTP2 cipher":              regexp.MustCompile(`TLSConfig.CipherSuites is missing an HTTP/2-required`),
		"Incompatible TLS version":     regexp.MustCompile(`protocol version not supported`),
	}
)

type testLogger struct{}

func (t *testLogger) Log(keyvals ...interface{}) error {
	return nil
}

func getPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	p := listener.Addr().(*net.TCPAddr).Port
	return fmt.Sprintf(":%v", p)
}

type TestInputs struct {
	Name                string
	Server              func() *http.Server
	UseNilServer        bool
	YAMLConfigPath      string
	ExpectedError       *regexp.Regexp
	UseTLSClient        bool
	ClientMaxTLSVersion uint16
	CipherSuites        []uint16
	ActualCipher        uint16
	CurvePreferences    []tls.CurveID
	Username            string
	Password            string
}

func TestYAMLFiles(t *testing.T) {
	testTables := []*TestInputs{
		{
			Name:           `path to config yml invalid`,
			YAMLConfigPath: "somefile",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `empty config yml`,
			YAMLConfigPath: "testdata/tls_config_empty.yml",
			ExpectedError:  nil,
		},
		{
			Name:           `invalid config yml (invalid structure)`,
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			ExpectedError:  ErrorMap["YAML error"],
		},
		{
			Name:           `invalid config yml (invalid key)`,
			YAMLConfigPath: "testdata/tls_config_junk_key.yml",
			ExpectedError:  ErrorMap["YAML error"],
		},
		{
			Name:           `invalid config yml (cert path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_empty.bad.yml",
			ExpectedError:  ErrorMap["Invalid CertPath"],
		},
		{
			Name:           `invalid config yml (key path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_empty.bad.yml",
			ExpectedError:  ErrorMap["Invalid KeyPath"],
		},
		{
			Name:           `invalid config yml (cert path and key path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			ExpectedError:  ErrorMap["Invalid CertPath"],
		},
		{
			Name:           `invalid config yml (cert path invalid)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_invalid.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (key path invalid)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_invalid.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (cert path and key path invalid)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_invalid.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (invalid ClientAuth)`,
			YAMLConfigPath: "testdata/tls_config_noAuth.bad.yml",
			ExpectedError:  ErrorMap["ClientCA set without policy"],
		},
		{
			Name:           `invalid config yml (invalid ClientCAs filepath)`,
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (invalid user list)`,
			YAMLConfigPath: "testdata/tls_config_auth_user_list_invalid.bad.yml",
			ExpectedError:  ErrorMap["Bad password"],
		},
		{
			Name:           `invalid config yml (bad cipher)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_inventedCiphers.bad.yml",
			ExpectedError:  ErrorMap["Unknown cipher"],
		},
		{
			Name:           `invalid config yml (bad curves)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_inventedCurves.bad.yml",
			ExpectedError:  ErrorMap["Unknown curve"],
		},
		{
			Name:           `invalid config yml (bad TLS version)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_wrongTLSVersion.bad.yml",
			ExpectedError:  ErrorMap["Unknown TLS version"],
		},
	}
	for _, testInputs := range testTables {
		t.Run(testInputs.Name, testInputs.Test)
	}
}

func TestServerBehaviour(t *testing.T) {
	testTables := []*TestInputs{
		{
			Name:           `empty string YAMLConfigPath and default client`,
			YAMLConfigPath: "",
			ExpectedError:  nil,
		},
		{
			Name:           `empty string YAMLConfigPath and TLS client`,
			YAMLConfigPath: "",
			UseTLSClient:   true,
			ExpectedError:  ErrorMap["HTTP Response to HTTPS"],
		},
		{
			Name:           `valid tls config yml and default client`,
			YAMLConfigPath: "testdata/tls_config_noAuth.good.yml",
			ExpectedError:  ErrorMap["HTTP Request to HTTPS server"],
		},
		{
			Name:           `valid tls config yml and tls client`,
			YAMLConfigPath: "testdata/tls_config_noAuth.good.yml",
			UseTLSClient:   true,
			ExpectedError:  nil,
		},
		{
			Name:                `valid tls config yml with TLS 1.1 client`,
			YAMLConfigPath:      "testdata/tls_config_noAuth.good.yml",
			UseTLSClient:        true,
			ClientMaxTLSVersion: tls.VersionTLS11,
			ExpectedError:       ErrorMap["Incompatible TLS version"],
		},
		{
			Name:           `valid tls config yml with all ciphers`,
			YAMLConfigPath: "testdata/tls_config_noAuth_allCiphers.good.yml",
			UseTLSClient:   true,
			ExpectedError:  nil,
		},
		{
			Name:           `valid tls config yml with some ciphers`,
			YAMLConfigPath: "testdata/tls_config_noAuth_someCiphers.good.yml",
			UseTLSClient:   true,
			CipherSuites:   []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
			ExpectedError:  nil,
		},
		{
			Name:           `valid tls config yml with no common cipher`,
			YAMLConfigPath: "testdata/tls_config_noAuth_someCiphers.good.yml",
			UseTLSClient:   true,
			CipherSuites:   []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
			ExpectedError:  ErrorMap["Handshake failure"],
		},
		{
			Name:           `valid tls config yml with multiple client ciphers`,
			YAMLConfigPath: "testdata/tls_config_noAuth_someCiphers.good.yml",
			UseTLSClient:   true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			ActualCipher:  tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			ExpectedError: nil,
		},
		{
			Name:           `valid tls config yml with multiple client ciphers, client chooses cipher`,
			YAMLConfigPath: "testdata/tls_config_noAuth_someCiphers_noOrder.good.yml",
			UseTLSClient:   true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			ActualCipher:  tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			ExpectedError: nil,
		},
		{
			Name:           `valid tls config yml with all curves`,
			YAMLConfigPath: "testdata/tls_config_noAuth_allCurves.good.yml",
			UseTLSClient:   true,
			ExpectedError:  nil,
		},
		{
			Name:             `valid tls config yml with some curves`,
			YAMLConfigPath:   "testdata/tls_config_noAuth_someCurves.good.yml",
			UseTLSClient:     true,
			CurvePreferences: []tls.CurveID{tls.CurveP521},
			ExpectedError:    nil,
		},
		{
			Name:             `valid tls config yml with no common curves`,
			YAMLConfigPath:   "testdata/tls_config_noAuth_someCurves.good.yml",
			UseTLSClient:     true,
			CurvePreferences: []tls.CurveID{tls.CurveP384},
			ExpectedError:    ErrorMap["Handshake failure"],
		},
		{
			Name:           `valid tls config yml with non-http2 ciphers`,
			YAMLConfigPath: "testdata/tls_config_noAuth_noHTTP2.good.yml",
			UseTLSClient:   true,
			ExpectedError:  nil,
		},
		{
			Name:           `valid tls config yml with non-http2 ciphers but http2 enabled`,
			YAMLConfigPath: "testdata/tls_config_noAuth_noHTTP2Cipher.bad.yml",
			UseTLSClient:   true,
			ExpectedError:  ErrorMap["No HTTP2 cipher"],
		},
	}
	for _, testInputs := range testTables {
		t.Run(testInputs.Name, testInputs.Test)
	}
}

func TestConfigReloading(t *testing.T) {
	errorChannel := make(chan error, 1)
	var once sync.Once
	recordConnectionError := func(err error) {
		once.Do(func() {
			errorChannel <- err
		})
	}
	defer func() {
		if recover() != nil {
			recordConnectionError(errors.New("Panic in test function"))
		}
	}()

	goodYAMLPath := "testdata/tls_config_noAuth.good.yml"
	badYAMLPath := "testdata/tls_config_noAuth.good.blocking.yml"

	server := &http.Server{
		Addr: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		}),
	}
	defer func() {
		server.Close()
	}()

	go func() {
		defer func() {
			if recover() != nil {
				recordConnectionError(errors.New("Panic starting server"))
			}
		}()
		err := Listen(server, badYAMLPath, testlogger)
		recordConnectionError(err)
	}()

	client := getTLSClient()

	TestClientConnection := func() error {
		time.Sleep(250 * time.Millisecond)
		r, err := client.Get("https://localhost" + port)
		if err != nil {
			return (err)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return (err)
		}
		if string(body) != "Hello World!" {
			return (errors.New(string(body)))
		}
		return (nil)
	}

	err := TestClientConnection()
	if err == nil {
		recordConnectionError(errors.New("connection accepted but should have failed"))
	} else {
		swapFileContents(goodYAMLPath, badYAMLPath)
		defer swapFileContents(goodYAMLPath, badYAMLPath)
		err = TestClientConnection()
		if err != nil {
			recordConnectionError(errors.New("connection failed but should have been accepted"))
		} else {

			recordConnectionError(nil)
		}
	}

	err = <-errorChannel
	if err != nil {
		t.Errorf(" *** Failed test: %s *** Returned error: %v", "TestConfigReloading", err)
	}
}

func (test *TestInputs) Test(t *testing.T) {
	errorChannel := make(chan error, 1)
	var once sync.Once
	recordConnectionError := func(err error) {
		once.Do(func() {
			errorChannel <- err
		})
	}
	defer func() {
		if recover() != nil {
			recordConnectionError(errors.New("Panic in test function"))
		}
	}()

	var server *http.Server
	if test.UseNilServer {
		server = nil
	} else {
		server = &http.Server{
			Addr: port,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World!"))
			}),
		}
		defer func() {
			server.Close()
		}()
	}
	go func() {
		defer func() {
			if recover() != nil {
				recordConnectionError(errors.New("Panic starting server"))
			}
		}()
		err := Listen(server, test.YAMLConfigPath, testlogger)
		recordConnectionError(err)
	}()

	ClientConnection := func() (*http.Response, error) {
		var client *http.Client
		var proto string
		if test.UseTLSClient {
			client = getTLSClient()
			t := client.Transport.(*http.Transport)
			t.TLSClientConfig.MaxVersion = test.ClientMaxTLSVersion
			if len(test.CipherSuites) > 0 {
				t.TLSClientConfig.CipherSuites = test.CipherSuites
			}
			if len(test.CurvePreferences) > 0 {
				t.TLSClientConfig.CurvePreferences = test.CurvePreferences
			}
			proto = "https"
		} else {
			client = http.DefaultClient
			proto = "http"
		}
		req, err := http.NewRequest("GET", proto+"://localhost"+port, nil)
		if err != nil {
			t.Error(err)
		}
		if test.Username != "" {
			req.SetBasicAuth(test.Username, test.Password)
		}
		return client.Do(req)
	}
	go func() {
		time.Sleep(250 * time.Millisecond)
		r, err := ClientConnection()
		if err != nil {
			recordConnectionError(err)
			return
		}

		if test.ActualCipher != 0 {
			if r.TLS.CipherSuite != test.ActualCipher {
				recordConnectionError(
					fmt.Errorf("bad cipher suite selected. Expected: %s, got: %s",
						tls.CipherSuiteName(r.TLS.CipherSuite),
						tls.CipherSuiteName(test.ActualCipher),
					),
				)
			}
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			recordConnectionError(err)
			return
		}
		if string(body) != "Hello World!" {
			recordConnectionError(errors.New(string(body)))
			return
		}
		recordConnectionError(nil)
	}()
	err := <-errorChannel
	if test.isCorrectError(err) == false {
		if test.ExpectedError == nil {
			t.Logf("Expected no error, got error: %v", err)
		} else {
			t.Logf("Expected error matching regular expression: %v", test.ExpectedError)
			t.Logf("Got: %v", err)
		}
		t.Fail()
	}
}

func (test *TestInputs) isCorrectError(returnedError error) bool {
	switch {
	case returnedError == nil && test.ExpectedError == nil:
	case returnedError != nil && test.ExpectedError != nil && test.ExpectedError.MatchString(returnedError.Error()):
	default:
		return false
	}
	return true
}

func getTLSClient() *http.Client {
	cert, err := ioutil.ReadFile("testdata/tls-ca-chain.pem")
	if err != nil {
		panic("Unable to start TLS client. Check cert path")
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: func() *x509.CertPool {
					caCertPool := x509.NewCertPool()
					caCertPool.AppendCertsFromPEM(cert)
					return caCertPool
				}(),
			},
		},
	}
	return client
}

func swapFileContents(file1, file2 string) error {
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		return err
	}
	content2, err := ioutil.ReadFile(file2)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file1, content2, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file2, content1, 0644)
	if err != nil {
		return err
	}
	return nil
}

func TestUsers(t *testing.T) {
	testTables := []*TestInputs{
		{
			Name:           `without basic auth`,
			YAMLConfigPath: "testdata/tls_config_users_noTLS.good.yml",
			ExpectedError:  ErrorMap["Unauthorized"],
		},
		{
			Name:           `with correct basic auth`,
			YAMLConfigPath: "testdata/tls_config_users_noTLS.good.yml",
			Username:       "dave",
			Password:       "dave123",
			ExpectedError:  nil,
		},
		{
			Name:           `without basic auth and TLS`,
			YAMLConfigPath: "testdata/tls_config_users.good.yml",
			UseTLSClient:   true,
			ExpectedError:  ErrorMap["Unauthorized"],
		},
		{
			Name:           `with correct basic auth and TLS`,
			YAMLConfigPath: "testdata/tls_config_users.good.yml",
			UseTLSClient:   true,
			Username:       "dave",
			Password:       "dave123",
			ExpectedError:  nil,
		},
		{
			Name:           `with another correct basic auth and TLS`,
			YAMLConfigPath: "testdata/tls_config_users.good.yml",
			UseTLSClient:   true,
			Username:       "carol",
			Password:       "carol123",
			ExpectedError:  nil,
		},
		{
			Name:           `with bad password and TLS`,
			YAMLConfigPath: "testdata/tls_config_users.good.yml",
			UseTLSClient:   true,
			Username:       "dave",
			Password:       "bad",
			ExpectedError:  ErrorMap["Forbidden"],
		},
		{
			Name:           `with bad username and TLS`,
			YAMLConfigPath: "testdata/tls_config_users.good.yml",
			UseTLSClient:   true,
			Username:       "nonexistent",
			Password:       "nonexistent",
			ExpectedError:  ErrorMap["Forbidden"],
		},
	}
	for _, testInputs := range testTables {
		t.Run(testInputs.Name, testInputs.Test)
	}
}
