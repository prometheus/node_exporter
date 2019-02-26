package https

import (
	"crypto/tls"
	"io/ioutil"

	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
)

type config struct {
	TLSCertPath string `yaml:"tlsCertPath"`
	TLSKeyPath string `yaml:"tlsKeyPath"`
	TLSConfig TLSStruct `yaml:"tlsConfig"`
}

type TLSStruct struct {
//	RootCAs	interface{} `yaml:"rootCAs"`
	ServerName	string	`yaml:"serverName"`
	ClientAuth	string	`yaml:"clientAuth"`
//	ClientCAs	interface{} `yaml:"clientCAs"`
	InsecureSkipVerify	bool	`yaml:"insecureSkipVerify"`
	CipherSuites	[]uint16 `yaml:"cipherSuites"`
	PreferServerCipherSuites	bool	`yaml:"preferServerCipherSuites"`
	MinVersion	uint16	`yaml:"minVersion"`
	MaxVersion	uint16	`yaml:"maxVersion"`
}

func GetTLSConfig( configPath string) *tls.Config {
	tlsc := &tls.Config{}
	tlsc, err := LoadConfigFromYaml(tlsc, configPath)
	if err != nil {
		log.Fatalf("Config failed to load from Yaml", err)
	}
	tlsc.BuildNameToCertificate()
	return tlsc	
}

func LoadConfigFromYaml(cfg *tls.Config, fileName string)(*tls.Config, error){
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}
	c := &config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return cfg, err
	}
	if len(c.TLSCertPath) > 0 {
		cfg.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(c.TLSCertPath, c.TLSKeyPath)
			if err != nil {
				return nil, err
			}
			return &cert, nil
		}
	}
	if len(c.TLSConfig.ServerName) > 0 {
		cfg.ServerName = c.TLSConfig.ServerName
	}
	if (c.TLSConfig.InsecureSkipVerify)==true{
		cfg.InsecureSkipVerify = true
	}		
	if len(c.TLSConfig.CipherSuites) > 0 {
		cfg.CipherSuites = c.TLSConfig.CipherSuites
	}
	if (c.TLSConfig.PreferServerCipherSuites)==true{
		cfg.PreferServerCipherSuites = c.TLSConfig.PreferServerCipherSuites
	}
	if (c.TLSConfig.MinVersion) != 0 {
		cfg.MinVersion = c.TLSConfig.MinVersion
	}
	if (c.TLSConfig.MaxVersion) != 0 {
                cfg.MaxVersion = c.TLSConfig.MaxVersion
        }
//	if len(c.TLSConfig.RootCAs) > 0{
//	//Append root ca's to certpool
//	}
	if len(c.TLSConfig.ClientAuth) > 0 {
		switch s := (c.TLSConfig.ClientAuth); s{
		case "RequestClientCert":
			cfg.ClientAuth = tls.RequestClientCert
		case "RequireClientCert":
			cfg.ClientAuth = tls.RequireAnyClientCert
		case "VerifyClientCertIfGiven":
			cfg.ClientAuth = tls.VerifyClientCertIfGiven
		case "RequireAndVerifyClientCert":
			cfg.ClientAuth = tls.RequireAndVerifyClientCert	
		default: 
			cfg.ClientAuth = tls.NoClientCert
		}
	}
	return cfg, nil 
}
