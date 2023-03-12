package apiclient

import (
	"encoding/json"
	"example/Minik8s/pkg/const/urlconst"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func WatchAPIWithRelativePath(runtimeConfig runtimedata.RuntimeConfig, relativeUrlPath string, outEvent chan watch.WatchEvent) {
	targetUrl := fmt.Sprintf("https://%s:%d%s",
		runtimeConfig.YamlConfig.APIServerIP, urlconst.PortAuth, relativeUrlPath)
	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	params := make(url.Values)
	params.Add("watch", "true")

	request.URL.RawQuery = params.Encode()

	response, err := authClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	for {
		var watchEvent watch.WatchEvent
		err := decoder.Decode(&watchEvent)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
			}
		} else {
			outEvent <- watchEvent
		}
	}
}
