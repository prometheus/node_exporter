package https

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"
)

type TestInputs struct {
	Name           string
	Server         func() *http.Server
	httpServer     *http.Server
	YAMLConfigPath string
	httpClient     *http.Client
	Client         func() *http.Client
	ConnectionURL  string
	ExpectedResult bool
	ExpectedError  *regexp.Regexp
}

func TestListen(t *testing.T) {
	logging := testing.Verbose()

	port := ":9100"

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

	BaseServer := func() *http.Server {
		return &http.Server{
			Addr: port,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World!"))
			}),
		}
	}

	ErrorMap := map[string]*regexp.Regexp{
		"HTTP Response to HTTPS Client": regexp.MustCompile(`server gave HTTP response to HTTPS client`),
		"Server Panic":                  regexp.MustCompile(`Panic starting server`),
		"No such file":                  regexp.MustCompile(`no such file`),
		"YAML error":                    regexp.MustCompile(`yaml`),
		"Invalid ClientAuth":            regexp.MustCompile(`ClientAuth`),
		"TLS handshake":                 regexp.MustCompile(`tls`),
		"Malformed response":            regexp.MustCompile(`malformed HTTP`),
	}

	testTables := []*TestInputs{
		{
			Name:           `empty string YAMLConfigPath and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: true,
		},
		{
			Name:           `empty string YAMLConfigPath and TLS client`,
			Server:         BaseServer,
			YAMLConfigPath: "",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["HTTP Response to HTTPS Client"],
		},
		{
			Name:           `nil Server and default client`,
			Server:         nil,
			YAMLConfigPath: "",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["Server Panic"],
		},
		{
			Name:           `nil Server and tls client`,
			Server:         nil,
			YAMLConfigPath: "",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["Server Panic"],
		},
		{
			Name:           `empty config.yml and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_empty.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `empty config.yml and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_empty.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `valid tls config yml and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth.good.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: true,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (cert path invalid) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (key path invalid) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (cert path and key path invalid) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (cert path empty) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (key path empty) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_keyPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (cert path and key path empty) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (cert path and key path empty) and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth_certPath_keyPath_empty.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `invalid tls config yml (invalid structure) and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["YAML error"],
		},
		{
			Name:           `invalid tls config yml (invalid structure) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_junk.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["YAML error"],
		},
		{
			Name:           `bad config yml path and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "somefile",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `bad config yml path and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "somefile",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `bad config yml (invalid ClientAuth) and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["Invalid ClientAuth"],
		},
		{
			Name:           `bad config yml (invalid ClientAuth) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_noAuth.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["Invalid ClientAuth"],
		},
		{
			Name:           `bad config yml (invalid ClientCAs filepath) and default client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			Client:         DefaultClient,
			ConnectionURL:  httpPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
		{
			Name:           `bad config yml (invalid ClientCAs filepath) and tls client`,
			Server:         BaseServer,
			YAMLConfigPath: "testdata/tls_config_auth_clientCAs_invalid.bad.yml",
			Client:         TLSClient,
			ConnectionURL:  httpsPath,
			ExpectedResult: false,
			ExpectedError:  ErrorMap["No such file"],
		},
	}

	numberOfTests, failedTests := len(testTables), 0

	if logging {
		log.Printf("Running %v tests:", numberOfTests)
	}

	logsDisabled := func(disable bool) {
		if logging {
			return
		}
		if disable {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
		} else {
			log.SetFlags(1)
			log.SetOutput(os.Stderr)
		}
	}
	for _, test := range testTables {
		test.httpClient = test.Client()
		if test.Server != nil {
			test.httpServer = test.Server()
		}
		if logging {
			log.Printf(" *** Running test: %s", test.Name)
		}

		logsDisabled(true)
		actualResult, err := runTest(test)
		logsDisabled(false)

		switch {
		case actualResult && test.ExpectedResult:
		case test.ExpectedResult == false && test.ExpectedError.MatchString(err.Error()):
		default:
			t.Fail()
			failedTests++
			if logging {
				log.Printf(" *** Failed test: %s", test.Name)
				log.Printf(" *** Returned error: %s", err.Error())
				if test.ExpectedError != nil {
					log.Printf(" *** Expected error: %s", test.ExpectedError.String())
				}
			}
		}
	}
	log.Printf("Passed %v of %v tests", numberOfTests-failedTests, numberOfTests)
}

func runTest(test *TestInputs) (bool, error) {
	connectionEstablished := make(chan bool, 1)
	var errorMessage error

	var once sync.Once
	recordConnectionResult := func(status bool, err error) {
		once.Do(func() {
			errorMessage = err
			connectionEstablished <- status
		})
	}

	defer func() {
		if recover() != nil {
			recordConnectionResult(false, errors.New("Panic in test function"))
		}
	}()

	go func() {
		defer func() {
			if recover() != nil {
				log.Printf("recovering")
				recordConnectionResult(false, errors.New("Panic starting server"))
			}
		}()
		err := Listen(test.httpServer, test.YAMLConfigPath)
		recordConnectionResult(false, err)
	}()

	go func() {
		time.Sleep(800 * time.Millisecond)
		r, err := test.httpClient.Get(test.ConnectionURL)
		if err != nil {
			recordConnectionResult(false, err)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			recordConnectionResult(false, err)
			return
		}
		if string(body) != "Hello World!" {
			recordConnectionResult(false, err)
			return
		}
		recordConnectionResult(true, nil)
	}()

	defer func() {
		if test.httpServer != nil {
			test.httpServer.Close()
		}
	}()

	return <-connectionEstablished, errorMessage
}
