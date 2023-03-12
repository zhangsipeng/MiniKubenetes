package apiclient

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"example/Minik8s/pkg/const/urlconst"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"gopkg.in/yaml.v2"
)

func createCSR() (csrDerBytes []byte, privkeyDerBytes []byte, err error) {
	pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
	csrTemplate := x509.CertificateRequest{
		SignatureAlgorithm: x509.PureEd25519,
		EmailAddresses:     nil,
		URIs:               nil,
		ExtraExtensions:    nil,
		PublicKeyAlgorithm: x509.Ed25519,
		PublicKey:          pubkey,
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privkey)
	if err != nil {
		return nil, nil, err
	}
	privkeyDerBytes, err = x509.MarshalPKCS8PrivateKey(privkey)
	return csr, privkeyDerBytes, err
}

func getRootCAPem(apiServerAddrPort string, caHash string) ([]byte, error) {
	targetURL := "https://" + apiServerAddrPort
	request, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		return nil, err
	}
	var verifiedCA *x509.Certificate
	verifiedCA = nil
	clientWithoutCert := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					for _, rawCert := range rawCerts {
						certHash := sha256.Sum256(rawCert[:])
						if hex.EncodeToString(certHash[:]) == caHash {
							verifiedCA, err = x509.ParseCertificate(rawCert)
							if err != nil {
								return err
							}
						}
					}
					if verifiedCA == nil {
						return errors.New("unrecognized CA")
					}
					return nil
				},
			},
		},
	}
	defer clientWithoutCert.CloseIdleConnections()
	response, err := clientWithoutCert.Do(request)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: verifiedCA.Raw,
	})
	return certPEM, nil
}

func GetCertificate(token string, apiServerIP string, caHash string) (yamlConfig runtimedata.YamlConfig, err error) {
	yamlConfig = runtimedata.YamlConfig{
		APIServerIP: apiServerIP,
		Others:      map[string]string{},
	}
	apiServerAddrPort := fmt.Sprintf("%s:%d", apiServerIP, urlconst.PortNonAuth)
	rootCAPem, err := getRootCAPem(apiServerAddrPort, caHash)
	if err != nil {
		return
	}
	yamlConfig.RootCAPem = string(rootCAPem)
	err = AddAuthClientRootCA(yamlConfig)
	if err != nil {
		return
	}
	csr, privkey, err := createCSR()
	if err != nil {
		return
	}
	certRequest := runtimedata.CertRequest{
		Token:  token,
		CsrStr: csr,
	}
	responseByte, err := RequestWithAddrPort(apiServerIP, urlconst.PortNonAuth,
		"/api/v1/token/", certRequest, "POST")
	// close tcp connection to PortNonAuth of API Server
	authClient.CloseIdleConnections()
	if err != nil {
		return
	}
	var certResponse runtimedata.CertResponse
	err = json.Unmarshal(responseByte, &certResponse)
	if err != nil {
		return
	}
	cert, err := x509.ParseCertificate(certResponse.CrtStr)
	if err != nil {
		return
	}
	if len(cert.IPAddresses) != 1 {
		err = errors.New("wrong number of IPAddresses in cert")
		return
	}
	yamlConfig.ClientIP = cert.IPAddresses[0].String()
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certResponse.CrtStr,
	})
	yamlConfig.ClientCertPem = string(certPEM)
	privkeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privkey,
	})
	yamlConfig.ClientKeyPem = string(privkeyPEM)
	return
}

func AddAuthClientCertificate(yamlConfig runtimedata.YamlConfig) (err error) {
	clientTLSCert, err := tls.X509KeyPair(
		[]byte(yamlConfig.ClientCertPem), []byte(yamlConfig.ClientKeyPem))
	if err != nil {
		return
	}
	authClient.Transport.(*http.Transport).TLSClientConfig.Certificates =
		[]tls.Certificate{clientTLSCert}
	return
}

func AddAuthClientRootCA(yamlConfig runtimedata.YamlConfig) error {
	caPool := x509.NewCertPool()
	ok := caPool.AppendCertsFromPEM([]byte(yamlConfig.RootCAPem))
	if !ok {
		return errors.New("bad RootCAPem")
	}
	tlsConfig := tls.Config{
		RootCAs: caPool,
	}
	authClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
		},
	}
	return nil
}

func InitRuntimeConfigOrError(info runtimedata.InitInfo, setOther func(yamlConfig *runtimedata.YamlConfig)) (runtimeConfig runtimedata.RuntimeConfig, err error) {
	init := info.Init
	apiServerIP := info.ApiServerIP
	token := info.Token
	caHash := info.CaHash
	configFilename := info.ConfigFileName

	var yamlConfig runtimedata.YamlConfig
	runtimeConfig.YamlConfig = &yamlConfig
	if init {
		if net.ParseIP(apiServerIP) == nil {
			err = errors.New(fmt.Sprintf("%s is not a valid IP", apiServerIP))
			return
		}
		if caHash == "" {
			err = errors.New("需要指定caHash")
			return
		}
		yamlConfig, err = GetCertificate(token, apiServerIP, caHash)
		if err != nil {
			return
		}
		if setOther != nil {
			setOther(&yamlConfig)
		}
		var yamlBytes []byte
		yamlBytes, err = yaml.Marshal(yamlConfig)
		if err != nil {
			return
		}
		if err = ioutil.WriteFile(configFilename, yamlBytes, 0600); err != nil {
			return
		}
	} else {
		var yamlBytes []byte
		yamlBytes, err = ioutil.ReadFile(configFilename)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(yamlBytes, &yamlConfig)
		if err != nil {
			return
		}
		err = AddAuthClientRootCA(yamlConfig)
	}

	err = AddAuthClientCertificate(yamlConfig)
	return
}

// deprecated: use InitRuntimeConfigOrError instead
func InitRuntimeConfig(info runtimedata.InitInfo, setOther func(yamlConfig *runtimedata.YamlConfig)) runtimedata.RuntimeConfig {
	runtimeConfig, err := InitRuntimeConfigOrError(info, setOther)
	if err != nil {
		panic(err)
	}

	return runtimeConfig
}

func GetInitInfo() runtimedata.InitInfo {
	init := flag.Bool("init", false, "是否重新注册")
	apiServerIP := flag.String("api-server-ip", "127.0.0.1", "API Server IP")
	token := flag.String("token", "", "注册令牌")
	caHash := flag.String("ca-hash", "", "API Server TLS CA证书的hash")
	configFilename := flag.String("config", "default-config.yaml", "配置文件")
	secretDump := flag.String("secret-dump", "", "调试用，API Server导出的存有token和caHash的文件")
	flag.Parse()
	if *secretDump != "" {
		secretDumpContent, err := ioutil.ReadFile(*secretDump)
		if err != nil {
			log.Printf("warning: reading secret dump file %s error. %s\n", *secretDump, err.Error())
			log.Println("ignoring secret dump file.")
		} else {
			var token_, caHash_ string
			n, err := fmt.Sscanf(string(secretDumpContent),
				"token: %s\ncaHash: %s\n", &token_, &caHash_)
			if n != 2 || err != nil {
				log.Printf("warning: reading secret dump file %s error. %s\n", *secretDump, err.Error())
				log.Println("ignoring secret dump file.")
			} else {
				*token = token_
				*caHash = caHash_
			}
		}
	}
	return runtimedata.InitInfo{
		Init:           *init,
		ApiServerIP:    *apiServerIP,
		Token:          *token,
		CaHash:         *caHash,
		ConfigFileName: *configFilename,
	}
}
func GetInitInfo_v2(init bool, apiserverip string, token string, cahash string, configfilename string) runtimedata.InitInfo {
	return runtimedata.InitInfo{
		Init:           init,
		ApiServerIP:    apiserverip,
		Token:          token,
		CaHash:         cahash,
		ConfigFileName: configfilename,
	}

}
