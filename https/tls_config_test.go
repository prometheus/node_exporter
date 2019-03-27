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
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var (
	logging bool
	port    = ":5000"
)

type TestInputs struct {
	Name           string
	Server         *http.Server
	YAMLConfigPath string
	httpClient     *http.Client
	Client         func() *http.Client
	ConnectionURL  string
	ExpectedResult bool
}

func TestListen(t *testing.T) {
	logging = testing.Verbose()

	httpPath := "http://localhost" + port
	httpsPath := "https://localhost" + port

	DefaultClient := func() *http.Client {
		return http.DefaultClient
	}

	TLSClient := func() *http.Client {
		cert, err := ioutil.ReadFile("testdata/tls-ca-chain.pem")
		if err != nil {
			log.Fatal("Unable to start TLS client. Check cert path")
		}
		return &http.Client{
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
	}

	testTables := []*TestInputs{
		{
			Name:           `empty string YAMLConfigPath and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: true,
		},
		{
			Name:           `empty string YAMLConfigPath and TLS client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `nil Server and default client`,
			Server:         nil,
			YAMLConfigPath: "",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `nil Server and tls client`,
			Server:         nil,
			YAMLConfigPath: "",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `empty config.yml and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_empty.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `empty config.yml and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_empty.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `valid tls config yml and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth.good.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: true,
		},
		{
			Name:           `invalid tls config yml (cert path invalid) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (key path invalid) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (cert path and key path invalid) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (cert path empty) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (key path empty) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (cert path and key path empty) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (cert path and key path empty) and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (invalid structure) and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `invalid tls config yml (invalid structure) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml path and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/file_does_not_exist",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml path and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/file_does_not_exist",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (invalid ClientAuth) and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (invalid ClientAuth) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_noAuth.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (invalid ClientCAs filepath) and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (invalid ClientCAs filepath) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (ClientCAs not provided) and default client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_missing.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
		},
		{
			Name:           `bad config yml (ClientCAs not provided) and tls client`,
			Server:         &http.Server{Addr: port},
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_missing.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
		},
	}

	if logging {
		log.Printf("Running %v tests:", len(testTables))
	}
	for _, test := range testTables {
		if logging {
			log.Printf("Running test: %s: ", test.Name)
		}
		if test.Server != nil {
			test.Server.Handler = &handler{test.Name}
		}
		test.httpClient = test.Client()
		if runTest(test) != test.ExpectedResult {
			t.Fail()
			if logging {
				log.Println("Failed test: " + test.Name)
			}
		}
	}
}

func runTest(test *TestInputs) (passed bool) {
	connectedOK := make(chan bool, 1)

	var once sync.Once
	connected := func(b bool, msg string) {
		once.Do(func() {
			if logging && msg != "" {
				log.Println(msg)
			}
			connectedOK <- b
		})
	}

	defer func() {
		if recover() != nil {
			connected(false, "recovered in runTest")
		}
	}()

	defer func() {
		test.Server.Close()
	}()

	// If goroutines block, return false
	go func() {
		time.Sleep(5 * time.Second)
		connected(false, "test timed out")
	}()

	// Start Server with provided YAMLConfigPath
	go func() {
		defer func() {
			recover()
		}()
		Listen(test.Server, test.YAMLConfigPath)
	}()

	// Try to connect with provided nonTLSClient and URL
	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			recover()
		}()
		r, err := test.httpClient.Get(test.ConnectionURL)
		if err != nil {
			connected(false, err.Error())
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			connected(false, err.Error())
			return
		}
		if string(body) != string(test.Name) {
			connected(false, "server response not as expected")
			return
		}
		connected(true, "")
	}()

	// wait for first response
	passed = <-connectedOK
	close(connectedOK)
	return passed
}

type handler struct {
	ResponseText string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.ResponseText))
}
