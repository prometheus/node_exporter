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

type TLSStruct {
	Certificates                interface{} `yaml:"certificates"`
	NameToCertificate           interface{} `yaml:"nameToCertificate"`
	GetCertificate              interface{} `yaml:"getCertificate"`
	GetClientCertificate        interface{} `yaml:"getClientCertificate"`
        GetConfigForClient          interface{} `yaml:"getConfigForClient"`
        VerifyPeerCertificate       interface{} `yaml:"verifyPeerCertificate"`
        RootCAs                     interface{} `yaml:"rootCAs"`
        NextProtos                  interface{} `yaml:"nextProtos"`
        ServerName                  string      `yaml:"serverName"`
        ClientAuth                  interface{} `yaml:"clientAuth"`
        ClientCAs                   interface{} `yaml:"clientCAs"`
        InsecureSkipVerify          bool        `yaml:"insecureSkipVerify"`
        CipherSuites                interface{} `yaml:"cipherSuites"`
        PreferServerCipherSuites    interface{} `yaml:"preferServerCipherSuites"`
        SessionTicketsDisabled      interface{} `yaml:"sessionTicketsDisabled"`
        SessionTicketKey            interface{} `yaml:"sessionTicketKey"`
        ClientSessionCache          interface{} `yaml:"clientSessionCache"`
        MinVersion                  interface{} `yaml:"minVersion"`
        MaxVersion                  interface{} `yaml:"maxVersion"`
        CurvePreferences            interface{} `yaml:"curvePreferences"`
        DynamicRecordSizingDisabled interface{} `yaml:"dynamicRecordSizingDisabled"`
        Renegotiation               interface{} `yaml:"renegotiation"`
}

func GetTLSConfig(c func(*tls.ClientHelloInfo)(*tls.Certificate, error))(*tls.Config) {
	var tlsc config
	tlsc.YamlIn("https/tls2.yml")
//	if err != nil {
//		log.Fatalf("BROKEN YAML")
//       }	
	tlsc.TLSConfig.GetCertificate = c
	
//	tlsc := &tls.Config{
//		GetCertificate: c,
//		ClientAuth: tls.RequireAndVerifyClientCert,
//	}
	tlsc.TLSConfig.BuildNameToCertificate()	
 	return  &tlsc.TLSConfig 
}
func NewConfig(cfg *config) (*tls.Config, error){
	
}

func (cfg *config) YamlIn(fileName string)(*config){
	var defaultC config

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Infof("Failed to read file")
		return &defaultC
	}
	 	
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return &defaultC
	}
	log.Infoln(cfg.TestString)
	return cfg
}
