package data

import "crypto/x509"

var serverCert *x509.Certificate
var authCert *x509.Certificate

func GetServerCert() *x509.Certificate {
	return serverCert
}

func SetServerCert(cert *x509.Certificate) {
	serverCert = cert
}

func GetAuthCert() *x509.Certificate {
	return authCert
}

func SetAuthCert(cert *x509.Certificate) {
	authCert = cert
}
