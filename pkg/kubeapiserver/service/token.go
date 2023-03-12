package service

import (
	"crypto/x509"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/data"
	"example/Minik8s/pkg/kubeapiserver/ssl"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkToken(c *gin.Context) {
	var body runtimedata.CertRequest
	err := c.BindJSON(&body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Unauthorized")
		return
	}

	// check token
	if body.Token != data.GetToken() {
		log.Println("token err!")
		c.JSON(http.StatusBadRequest, "Unauthorized")
		return
	}

	csr, err := x509.ParseCertificateRequest(body.CsrStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Unauthorized")
		return
	}

	err = csr.CheckSignature()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Unauthorized")
		return
	}

	pubKey := csr.PublicKey

	remoteIP := net.ParseIP(c.ClientIP())
	var resBody runtimedata.CertResponse
	resBody.CrtStr = ssl.GenCertificate(pubKey, data.GetAuthKey(), data.GetAuthCert(),
		nil, []net.IP{remoteIP})

	c.JSON(http.StatusOK, resBody)
}
