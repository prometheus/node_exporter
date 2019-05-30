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
	port = getPort()

	ErrorMap = map[string]*regexp.Regexp{
		"HTTP Response to HTTPS":       regexp.MustCompile(`server gave HTTP response to HTTPS client`),
		"No such file":                 regexp.MustCompile(`no such file`),
		"Invalid argument":             regexp.MustCompile(`invalid argument`),
		"YAML error":                   regexp.MustCompile(`yaml`),
		"Invalid ClientAuth":           regexp.MustCompile(`ClientAuth`),
		"TLS handshake":                regexp.MustCompile(`tls`),
		"HTTP Request to HTTPS server": regexp.MustCompile(`HTTP`),
	}
)

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
	Name           string
	Server         func() *http.Server
	UseNilServer   bool
	YAMLConfigPath string
	ExpectedError  *regexp.Regexp
	UseTLSClient   bool
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
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (invalid structure)`,
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			ExpectedError:  ErrorMap["YAML error"],
		},
		{
			Name:           `invalid config yml (cert path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_empty.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (key path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_empty.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid config yml (cert path and key path empty)`,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
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
			ExpectedError:  ErrorMap["Invalid ClientAuth"],
		},
		{
			Name:           `invalid config yml (invalid ClientCAs filepath)`,
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			ExpectedError:  ErrorMap["No such file"],
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
		err := Listen(server, badYAMLPath)
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
		recordConnectionError(errors.New("Connection accepted but should have failed."))
	} else {
		swapFileContents(goodYAMLPath, badYAMLPath)
		defer swapFileContents(goodYAMLPath, badYAMLPath)
		err = TestClientConnection()
		if err != nil {
			recordConnectionError(errors.New("Connection failed but should have been accepted."))
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
		err := Listen(server, test.YAMLConfigPath)
		recordConnectionError(err)
	}()

	var ClientConnection func() (*http.Response, error)
	if test.UseTLSClient {
		ClientConnection = func() (*http.Response, error) {
			client := getTLSClient()
			return client.Get("https://localhost" + port)
		}
	} else {
		ClientConnection = func() (*http.Response, error) {
			client := http.DefaultClient
			return client.Get("http://localhost" + port)
		}
	}
	go func() {
		time.Sleep(250 * time.Millisecond)
		r, err := ClientConnection()
		if err != nil {
			recordConnectionError(err)
			return
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
