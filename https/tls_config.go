package https

import (
	"crypto/tls"
	"io/ioutil"

	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
)

type config struct {
	TLSConfig TLSStruct `yaml:"tlsConfig"`
	TestString string `yaml:"testString"`
}

type TLSStruct struct {
        RootCAs	interface{} `yaml:"rootCAs"`
        ServerName	string      `yaml:"serverName"`
        ClientAuth	interface{} `yaml:"clientAuth"`
        ClientCAs	interface{} `yaml:"clientCAs"`
        InsecureSkipVerify	bool	`yaml:"insecureSkipVerify"`
        CipherSuites	[]uint16 `yaml:"cipherSuites"`
        PreferServerCipherSuites	bool	`yaml:"preferServerCipherSuites"`
        MinVersion	uint16	`yaml:"minVersion"`
        MaxVersion	uint16	`yaml:"maxVersion"`
}

//func GetTLSConfig(c func(*tls.ClientHelloInfo)(*tls.Certificate, error))(*tls.Config) {
//	var tlsc config
//	tlsc.YamlIn("https/tls2.yml")
//	if err != nil {
//		log.Fatalf("BROKEN YAML")
//       }	
//	tlsc.TLSConfig.GetCertificate = c
//	
//	tlsc := &tls.Config{
//		GetCertificate: c,
//		ClientAuth: tls.RequireAndVerifyClientCert,
//	}
//	tlsc.TLSConfig.BuildNameToCertificate()	
// 	return  &tlsc.TLSConfig 
//}

//func NewConfig(cfg *config) (*tls.Config, error){
//	tlsConfig := &tls.Config{
//		
//	}
//}

func GetTLSConfig(cert, key string) *tls.Config {
	tlsc := &tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				return nil, err
			}
			return &cert, nil
		},
		//		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsc, err := LoadConfigFromYaml(tlsc, "https/tls2.yml")
	if err != nil {
		log.Fatalf("Config failed to load from Yaml")
	}
	tlsc.BuildNameToCertificate()
	return tlsc	
}

func LoadConfigFromYaml(cfg *tls.Config, fileName string)(*tls.Config, error){
	
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Infof("Failed to read file")
		return cfg, err
	}
	c := &config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		
		return cfg, err
	}

	log.Infoln(c.TestString)
	log.Infoln(c.TLSConfig.ServerName)
	return cfg, nil 
}
