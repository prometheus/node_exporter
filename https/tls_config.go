

package https 

import (
	"crypto/tls"
)

func GetTLSConfig(c func(*tls.ClientHelloInfo)(*tls.Certificate, error))(*tls.Config) {
	tlsc := &tls.Config{
		GetCertificate: c,
//		ClientAuth: tls.RequireAndVerifyClientCert,

	}
	tlsc.BuildNameToCertificate()	
 	return tlsc 
}
