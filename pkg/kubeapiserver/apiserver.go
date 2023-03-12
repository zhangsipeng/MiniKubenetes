package kubeapiserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"example/Minik8s/pkg/const/urlconst"
	"example/Minik8s/pkg/kubeapiserver/data"
	"example/Minik8s/pkg/kubeapiserver/etcd"
	"example/Minik8s/pkg/kubeapiserver/service"
	"example/Minik8s/pkg/kubeapiserver/ssl"
	"example/Minik8s/pkg/kubeapiserver/util"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type info struct {
	Token         string
	ServerKeyStr  []byte
	ServerCertStr []byte
	AuthKeyStr    []byte
	AuthCertStr   []byte
}

var sysPath string
var sysInfoPath string
var serverKeyPath string
var serverCertPath string
var authKeyPath string
var authCertPath string
var serverIPs []net.IP

var secretDumpFile *os.File

func StartService() {
	tlsRouter := gin.Default()
	nonTLSRouter := gin.Default()

	// create etcd client
	err := etcd.CreateClient()
	util.CheckError(err)

	ifInit := flag.Bool("init", false, "是否初始化")
	path := flag.String("sysPath", "", "系统路径")
	secretDump := flag.String("secret-dump", "", "调试用，导出token和caHash的文件")
	ifaceAddrs, err := net.InterfaceAddrs()
	var defaultServerIPs string
	if err != nil {
		util.CheckError(err)
		defaultServerIPs = ""
	} else {
		defaultServerIPList := []string{}
		for _, ifaceAddr := range ifaceAddrs {
			if ipnet, ok := ifaceAddr.(*net.IPNet); ok {
				defaultServerIPList = append(defaultServerIPList, ipnet.IP.String())
			}
		}
		defaultServerIPs = strings.Join(defaultServerIPList, ",")
	}
	serverIPStrs := flag.String("server-ip", defaultServerIPs, "服务器IP（多个IP用逗号隔开）")
	flag.Parse()

	if *path == "" {
		log.Println("请输入系统路径！")
		return
	}
	sysPath = *path

	serverKeyPath = sysPath + "server.pem"
	serverCertPath = sysPath + "server.crt"
	authKeyPath = sysPath + "auth.pem"
	authCertPath = sysPath + "auth.crt"
	sysInfoPath = sysPath + "sys_info"
	serverIPs = []net.IP{}
	for _, serverIPStr := range strings.Split(*serverIPStrs, ",") {
		ip := net.ParseIP(serverIPStr)
		if ip == nil {
			log.Printf("%s is not a valid IP\n", serverIPStr)
		} else {
			serverIPs = append(serverIPs, ip)
		}
	}
	if len(serverIPs) == 0 {
		log.Println("请提供正确的server IP！")
		return
	}

	var initInfo info
	// init system information
	if *ifInit {
		if *secretDump != "" {
			secretDumpFile, err = os.Create(*secretDump)
			if err != nil {
				panic(err)
			}
		}
		initSysInfo(&initInfo)
		if err := secretDumpFile.Close(); err != nil {
			panic(err)
		}
	} else {
		readSysInfo(&initInfo)
	}
	writeFile(&initInfo)

	// register Pod service
	service.RegisterTLSService(tlsRouter)

	// register Token service
	service.RegisterNonTLSService(nonTLSRouter)

	// start service
	config, err := generateConfig()
	util.CheckError(err)
	s := &http.Server{
		Addr:      fmt.Sprintf(":%d", urlconst.PortAuth),
		TLSConfig: config,
		Handler:   tlsRouter,
	}
	go func() {
		err = s.ListenAndServeTLS(serverCertPath, serverKeyPath)
		util.CheckError(err)
	}()
	//go func() {
	err = nonTLSRouter.RunTLS(fmt.Sprintf(":%d", urlconst.PortNonAuth),
		serverCertPath, serverKeyPath)
	util.CheckError(err)
	//}()
}

func generateConfig() (config *tls.Config, err error) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(authCertPath)
	util.CheckError(err)
	pool.AppendCertsFromPEM(caCrt)

	config = &tls.Config{
		ClientCAs:  pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	return
}

// ifInit is true
func initSysInfo(initInfo *info) {
	// clean etcd
	err := etcd.Clean()
	util.CheckError(err)

	// set token
	initToken(initInfo)

	// set serverKey pair
	initKey(initInfo)

	// save info
	initData, err := yaml.Marshal(initInfo)
	util.CheckError(err)
	err = ioutil.WriteFile(sysInfoPath, initData, 0600)
	util.CheckError(err)
}

// ifInit is false
func readSysInfo(initInfo *info) {
	// read file
	content, err := ioutil.ReadFile(sysInfoPath)
	util.CheckError(err)

	// parse file content
	err = yaml.Unmarshal(content, &initInfo)
	util.CheckError(err)
	data.SetToken(initInfo.Token)
	serverKey, err := x509.ParsePKCS1PrivateKey(initInfo.ServerKeyStr)
	util.CheckError(err)
	data.SetServerKey(serverKey)
	serverCert, err := x509.ParseCertificate(initInfo.ServerCertStr)
	util.CheckError(err)
	data.SetServerCert(serverCert)
	authKey, err := x509.ParsePKCS1PrivateKey(initInfo.AuthKeyStr)
	util.CheckError(err)
	data.SetAuthKey(authKey)
	authCert, err := x509.ParseCertificate(initInfo.AuthCertStr)
	util.CheckError(err)
	data.SetAuthCert(authCert)
}

func writeFile(initInfo *info) {
	saveKey(initInfo.ServerKeyStr, serverKeyPath)
	saveCertificate(initInfo.ServerCertStr, serverCertPath)
	saveKey(initInfo.AuthKeyStr, authKeyPath)
	saveCertificate(initInfo.AuthCertStr, authCertPath)
}

func initToken(initInfo *info) {
	tokenByte := make([]byte, 1)
	_, err := rand.Read(tokenByte)
	util.CheckError(err)
	token := base64.RawStdEncoding.EncodeToString(tokenByte)
	data.SetToken(token)
	initInfo.Token = token
	log.Println("token:", token)
	if secretDumpFile != nil {
		secretDumpFile.Write([]byte(fmt.Sprintf("token: %s\n", token)))
	}
}

func initKey(initInfo *info) {
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	util.CheckError(err)
	data.SetServerKey(serverKey)
	initInfo.ServerKeyStr = x509.MarshalPKCS1PrivateKey(serverKey)

	// TODO: hard coded DNS
	dnsNames := []string{"localhost"}
	cert := ssl.GenSignSelfCertificate(serverKey, &serverKey.PublicKey, dnsNames, serverIPs)
	serverCert, err := x509.ParseCertificate(cert)
	util.CheckError(err)
	data.SetServerCert(serverCert)
	initInfo.ServerCertStr = cert

	caHash := sha256.Sum256(serverCert.Raw)
	log.Printf("server TLS CA hash: %x\n", caHash)
	if secretDumpFile != nil {
		secretDumpFile.Write([]byte(fmt.Sprintf("caHash: %x\n", caHash)))
	}

	authKey, err := rsa.GenerateKey(rand.Reader, 2048)
	util.CheckError(err)
	data.SetAuthKey(authKey)
	initInfo.AuthKeyStr = x509.MarshalPKCS1PrivateKey(authKey)

	cert = ssl.GenSignSelfCertificate(authKey, &authKey.PublicKey, dnsNames, serverIPs)
	authCert, err := x509.ParseCertificate(cert)
	util.CheckError(err)
	data.SetAuthCert(authCert)
	initInfo.AuthCertStr = cert
}

// save kay in file
func saveKey(key []byte, file string) {
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: key}
	saveBlock(block, file)
}

// save certificate in file
func saveCertificate(cert []byte, file string) {
	block := &pem.Block{Type: "CERTIFICATE", Bytes: cert}
	saveBlock(block, file)
}

// save block in file
func saveBlock(block *pem.Block, file string) {
	certOut, _ := os.Create(file)
	err := pem.Encode(certOut, block)
	util.CheckError(err)
	err = certOut.Close()
	util.CheckError(err)
}
