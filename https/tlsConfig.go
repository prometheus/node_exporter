

package https 

import (
	"crypto/tls"
	
)

type TLSConfig struct {
	config *tls.Config 	

}

func GetTLSConfig(c func(*tls.ClientHelloInfo)(*tls.Certificate, error))(*tls.Config) {
	
	tlsc := &tls.Config{
	CipherSuites: []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	},
	PreferServerCipherSuites: true,
	GetCertificate: c,
	}
	
 	return tlsc 
}
