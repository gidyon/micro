package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

// GetCert returns tls cert
func GetCert(crt, key string) (*tls.Certificate, *x509.CertPool, error) {
	ca, err := ioutil.ReadFile(crt)
	if err != nil {
		return nil, nil, err
	}
	pair, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, nil, err
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(ca)
	if !ok {
		return nil, nil, errors.New("failed to add cert to pool")
	}
	return &pair, certPool, nil
}
