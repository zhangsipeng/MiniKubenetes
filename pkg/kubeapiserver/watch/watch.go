package watch

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type WatchEvent struct {
	Type  string
	Key   string
	Value string
}

func StartWatch(c *gin.Context, watcher *clientv3.WatchChan) {
	w := c.Writer
	// set HTTP response header
	header := w.Header()
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	// begin to monitor the channel
	for resp := range *watcher {
		err := resp.Err()
		if err != nil {
			log.Println(err)
		}
		for _, ev := range resp.Events {
			event := &WatchEvent{
				Type:  ev.Type.String(),
				Key:   string(ev.Kv.Key),
				Value: string(ev.Kv.Value),
			}
			s, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				continue
			}
			_, err = w.Write(s)
			if err != nil {
				return
			}
		}
		w.(http.Flusher).Flush()
	}
}
