package etcd

import (
	"etcd/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

var (
	client *clientv3.Client
	isInit bool
)

func Init(conf clientv3.Config) error {
	log.Info("etcd init", isInit)
	if isInit == true {
		return nil
	}

	// 建立一个客户端
	var err error
	client, err = clientv3.New(conf)
	if err != nil {
		return err
	}

	initService()

	isInit = true

	return nil
}

func GetClient() *clientv3.Client {
	return client
}
