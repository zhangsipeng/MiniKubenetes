package etcdio

import (
	"context"
	"example/Minik8s/pkg/kubeapiserver/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

func Put(key, val string) error {
	_, err := etcd.Client.Put(context.TODO(), key, val)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func Get(key string) (string, error) {
	resp, err := etcd.Client.Get(context.TODO(), key)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// there's no such object
	if resp.Count == 0 {
		return "", nil
	}

	object := string(resp.Kvs[0].Value)

	return object, nil
}

func List(prefix string) ([]string, error) {
	resp, err := etcd.Client.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	list := make([]string, 0)

	// put the object into list
	for _, ev := range resp.Kvs {
		list = append(list, string(ev.Value))
	}

	return list, nil
}

func Delete(key string) error {
	_, err := etcd.Client.Delete(context.TODO(), key)
	return err
}

func Watch(key string) *clientv3.WatchChan {
	return etcd.GetWatcher(key, false)
}

func WatchList(prefix string) *clientv3.WatchChan {
	return etcd.GetWatcher(prefix, true)
}
