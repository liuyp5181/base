package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	serviceKey = "services"
)

type Service struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Version string `json:"version"`
	Power   int    `json:"power"`
}

var (
	serLease   clientv3.Lease
	serLeaseID clientv3.LeaseID
)

func initService() {
	//设置租约过期时间为20秒
	var ctx, cancle = context.WithCancel(context.Background())
	serLease = clientv3.NewLease(client)
	leaseRes, err := serLease.Grant(ctx, 20)
	if err != nil {
		panic(err)
	}
	serLeaseID = leaseRes.ID
	//续租时间约为自动租约的三分之一时间，extend.TODO官方定义为是你不知道要传什么
	keepaliveRes, err := serLease.KeepAlive(ctx, serLeaseID) // context的时候就用这个
	if err != nil {
		panic(err)
	}

	go func() {
		defer cancle()
		for {
			select {
			case ret := <-keepaliveRes:
				if ret == nil {
					fmt.Println("服务发现续租失败", time.Now().Format("2006-01-02 15:04:05"), ret)
				}
			}
		}
	}()
}

// WatcherService 负责将监听到的put、delete请求存放到指定list
func WatcherService(name string) clientv3.WatchChan {
	key := fmt.Sprintf("%s/%s", serviceKey, name)
	return client.Watch(context.Background(), key, clientv3.WithPrefix())
}

func GetService(name string) ([]Service, error) {
	key := fmt.Sprintf("%s/%s", serviceKey, name)
	resp, err := client.Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}
	var list = make([]Service, 0, len(resp.Kvs))
	for _, v := range resp.Kvs {
		if v != nil {
			var s Service
			err = json.Unmarshal(v.Value, &s)
			if err != nil {
				return nil, err
			}
			list = append(list, s)
		}
	}
	return list, nil
}

func SetService(name string, ip string, port int, version string, power int) error {
	key := fmt.Sprintf("%s/%s/%s:%d", serviceKey, name, ip, port)

	val, err := json.Marshal(Service{
		Key:     key,
		Name:    name,
		IP:      ip,
		Port:    port,
		Version: version,
		Power:   power,
	})

	kv := clientv3.NewKV(client)
	ctx := context.Background()

	_, err = kv.Put(ctx, key, string(val), clientv3.WithLease(serLeaseID)) //把服务的key绑定到租约下面
	if err != nil {
		return err
	}

	return nil
}

func PrintService() {
	resp, err := client.Get(context.Background(), serviceKey, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, v := range resp.Kvs {
		if v != nil {
			fmt.Println("print", string(v.Key))
		}
	}
}
