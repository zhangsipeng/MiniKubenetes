package service

import (
	"errors"
	etcdio "example/Minik8s/pkg/kubeapiserver/etcd/io"
	service_const "example/Minik8s/pkg/kubeapiserver/service/const"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"github.com/gin-gonic/gin"
)

func create(body []byte, prefix, name string) error {
	if name == "" {
		return errors.New(service_const.EmptyNameError)
	}

	path := prefix + name
	obj, _ := etcdio.Get(path)
	if obj != "" {
		return errors.New(service_const.AlreadyExistError)
	}

	err := etcdio.Put(path, string(body))
	return err
}

func put(body []byte, prefix, name string) error {
	if name == "" {
		return errors.New(service_const.EmptyNameError)
	}

	path := prefix + name
	obj, _ := etcdio.Get(path)
	if obj == "" {
		return errors.New(service_const.NotExistError)
	}

	err := etcdio.Put(path, string(body))
	return err
}

func get(prefix, name string) (string, error) {
	return etcdio.Get(prefix + name)
}

func list(prefix string) ([]string, error) {
	return etcdio.List(prefix)
}

func watchList(c *gin.Context, prefix string) {
	watcher := etcdio.WatchList(prefix)
	watch.StartWatch(c, watcher)
}

func watchObject(c *gin.Context, prefix, name string) {
	watcher := etcdio.WatchList(prefix + name)
	watch.StartWatch(c, watcher)
}

func delete(prefix, name string) error {
	if name == "" {
		return errors.New(service_const.EmptyNameError)
	}
	return etcdio.Delete(prefix + name)
}
