package apiclient

import (
	"bytes"
	"encoding/json"
	"example/Minik8s/pkg/const/urlconst"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RequestWithAddrPort(addr string, port int, relativeUrlPath string, object any, method string) (responseByte []byte, err error) {
	targetUrl := fmt.Sprintf("https://%s:%d%s", addr, port, relativeUrlPath)
	var objJSON []byte
	if object == nil {
		objJSON = []byte{}
	} else {
		objJSON, err = json.Marshal(object)
		if err != nil {
			return
		}
	}

	request, err := http.NewRequest(method, targetUrl, bytes.NewReader(objJSON))
	if err != nil {
		return
	}

	response, err := authClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	responseByte, err = ioutil.ReadAll(response.Body)

	return
}

func Request(runtimeConfig runtimedata.RuntimeConfig, relativeUrlPath string, object any, method string) []byte {
	responseByte, err := RequestWithAddrPort(runtimeConfig.YamlConfig.APIServerIP, urlconst.PortAuth,
		relativeUrlPath, object, method)
	if err != nil {
		panic(err)
	}
	return responseByte
}
