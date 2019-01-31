package https

import (
	"crypto/tls"
	"github.com/prometheus/common/log"
)

type WrappedCertificate struct {
	certificate *tls.Certificate
	CertPath	string
	KeyPath		string
}

func (c *WrappedCertificate) GetCertificate(clientHello *tls.ClientHelloInfo ) (*tls.Certificate, error){
	log.Infoln("Client Hello Received")
	if len(c.KeyPath) <= 0 {
		c.KeyPath = c.CertPath
	}
	c.LoadCertificates(c.CertPath, c.KeyPath)
	
	return c.certificate, nil
}

func (c *WrappedCertificate) LoadCertificates(certPath, keyPath string) error{
	certAndKey, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}
	log.Infoln("Loading Certs using https lib")
	c.certificate = &certAndKey
	return nil
}
