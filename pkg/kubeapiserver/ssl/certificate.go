package ssl

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"example/Minik8s/pkg/kubeapiserver/util"
	"math/big"
	"net"
	"time"
)

// get certificate template
func getCertTemplate(dnsNames []string, ipAddresses []net.IP) *x509.Certificate {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"SJTU"},
			Country:      []string{"CN"},
			Province:     []string{"ShangHai"},
			Locality:     []string{"ShangHai"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
	}

	return template
}

// GenSignSelfCertificate generate sign self certificate
func GenSignSelfCertificate(priKey, pubKey interface{}, dnsNames []string, ipAddresses []net.IP) []byte {
	template := getCertTemplate(dnsNames, ipAddresses)

	cert, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, priKey)
	util.CheckError(err)
	return cert
}

// GenCertificate generate certificate
func GenCertificate(pubKey, priKey interface{}, parentCert *x509.Certificate, dnsNames []string, ipAddresses []net.IP) []byte {
	template := getCertTemplate(dnsNames, ipAddresses)
	cert, err := x509.CreateCertificate(rand.Reader, template, parentCert, pubKey, priKey)
	util.CheckError(err)
	return cert
}
