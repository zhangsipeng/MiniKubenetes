package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var config = clientv3.Config{
	Endpoints:   []string{"localhost:2379"},
	DialTimeout: 5 * time.Second,
}

var Client *clientv3.Client

// CreateClient /* Initialize etcd client */
func CreateClient() error {
	cli, err := clientv3.New(config)

	if err != nil {
		log.Println(err)
		return err
	}

	Client = cli
	return nil
}

func CloseClient() error {
	err := Client.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetWatcher /* Get the watch chan of the specified resource */
func GetWatcher(path string, usePrefix bool) *clientv3.WatchChan {
	var watcher clientv3.WatchChan
	if usePrefix {
		watcher = Client.Watch(context.Background(), path, clientv3.WithPrefix())
	} else {
		watcher = Client.Watch(context.Background(), path)
	}
	return &watcher
}

func Clean() error {
	_, err := Client.Delete(context.TODO(), "", clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
