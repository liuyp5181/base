package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/etcd"
	"github.com/liuyp5181/base/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"math/rand"
	"sync"
	"time"
)

type Clients struct {
	sync.RWMutex
	list map[string][]*Client
	m    map[string]*Client
}

var clients = &Clients{
	list: map[string][]*Client{},
	m:    map[string]*Client{},
}

func (cs *Clients) isExist(name string) bool {
	cs.RLock()
	defer cs.RUnlock()
	if len(cs.list[name]) > 0 {
		return true
	}
	return false
}

func (cs *Clients) setCancel(name string, cancelFunc context.CancelFunc) {
	cs.Lock()
	defer cs.Unlock()
	cs.m[name].cancelFunc = cancelFunc
}

func (cs *Clients) setClient(key string, c *Client) {
	cs.Lock()
	defer cs.Unlock()
	cs.m[key] = c
	cs.list[c.Server.Name] = append(cs.list[c.Server.Name], c)
}

func (cs *Clients) delClient(key string) {
	fmt.Println("del", key)
	cs.Lock()
	defer cs.Unlock()
	c := cs.m[key]
	if c == nil {
		return
	}
	s := c.Server
	delete(cs.m, key)
	list := cs.list[s.Name]
	for i, v := range list {
		if s.IP == v.Server.IP && s.Port == v.Server.Port {
			cs.list[s.Name] = append(list[:i], list[i+1:]...)
			return
		}
	}
	c.close()
}

func (cs *Clients) getClientList(name string) []*Client {
	clients.RLock()
	defer clients.RUnlock()
	return clients.list[name]
}

func (cs *Clients) closeClients(name string) {
	clients.Lock()
	defer clients.Unlock()
	list := clients.list[name]
	for _, c := range list {
		s := c.Server
		delete(clients.m, s.Key)
		c.close()
	}
	clients.list[name] = nil
}

func initClient(name string) error {
	if clients.isExist(name) {
		return nil
	}

	list, err := etcd.GetService(name)
	if err != nil {
		return err
	}
	fmt.Println("initClient", len(list))

	for _, s := range list {
		fmt.Println("initClient", s)
		err = setClient(s.Key, s)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	go watchClient(ctx, name)

	clients.setCancel(name, cancel)

	return nil
}

func watchClient(cancelCtx context.Context, name string) {
	rch := etcd.WatcherService(cancelCtx, name)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Println("watch", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))

			switch ev.Type {
			case mvccpb.PUT:
				var s etcd.Service
				err := json.Unmarshal(ev.Kv.Value, &s)
				if err != nil {
					continue
				}
				err = setClient(string(ev.Kv.Key), s)
				if err != nil {
					continue
				}
			case mvccpb.DELETE:
				clients.delClient(string(ev.Kv.Key))
			}
		}
	}
	fmt.Println("watchClient is close", name)
}

func PrintClient() {
	clients.RLock()
	defer clients.RUnlock()
	for k, c := range clients.m {
		log.Info("PrintClient", k, c, c.Server)
	}
	for _, l := range clients.list {
		for i, v := range l {
			log.Info("PrintClient", i, v, v.Server)
		}
	}
}

func setClient(key string, s etcd.Service) error {
	if s.Name == config.ServiceName && s.IP == config.GetConfig().Server.IP && s.Port == config.GetConfig().Server.Port {
		return nil
	}

	c, err := newClient(&s)
	if err != nil {
		return err
	}
	fmt.Println("set", key, c)

	clients.setClient(key, c)

	return nil
}

func InitClients(serviceName ...string) error {
	if len(serviceName) == 0 {
		return initClient("")
	}
	for _, n := range serviceName {
		if err := initClient(n); err != nil {
			return err
		}
	}
	return nil
}

func CloseClients(name string) {
	clients.closeClients(name)
}

func GetClient(name string) (*Client, error) {
	list := clients.getClientList(name)
	fmt.Println("getClientList", list, len(list))
	if len(list) == 0 {
		return nil, fmt.Errorf("service is nil, name = %v", name)
	}

	var max int
	for _, v := range list {
		max += v.Server.Power
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(max)
	var index int
	for _, v := range list {
		if v.Server.Power == 0 {
			continue
		}
		index += v.Server.Power
		if r < index {
			return v, nil
		}
	}

	return nil, fmt.Errorf("not found service, name = %v", name)
}

func GetClientList(name string) ([]*Client, error) {
	list := clients.getClientList(name)
	fmt.Println("GetClientList", list, len(list))
	if len(list) == 0 {
		return nil, fmt.Errorf("not found service, name = %s", name)
	}
	return list, nil
}
