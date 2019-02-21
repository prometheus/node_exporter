package https

import (
	"crypto/tls"
)

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
	tlsc.BuildNameToCertificate()
	return tlsc
}
