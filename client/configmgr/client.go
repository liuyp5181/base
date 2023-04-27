package configmgr

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"sync"

	pb "github.com/liuyp5181/base/client/configmgr/api"
	"github.com/liuyp5181/base/log"
	"github.com/liuyp5181/base/service"
)

var configList = struct {
	sync.RWMutex
	list map[string]map[string]interface{}
}{
	list: map[string]map[string]interface{}{},
}

func Init() error {
	err := service.InitClients(pb.Greeter_ServiceDesc.ServiceName)
	if err != nil {
		return err
	}
	return nil
}

func load(val []byte, conf interface{}, confType string) error {
	vp := viper.New()
	vp.SetConfigType(confType)
	vp.AutomaticEnv()
	if err := vp.ReadConfig(bytes.NewReader(val)); err != nil {
		return fmt.Errorf("read config failed, err_msg=[%s], extend=[%s]", err.Error(), string(val))
	}

	if err := vp.Unmarshal(conf); err != nil {
		return fmt.Errorf("local config unmarshal failed, err_msg=[%s], extend=[%s]", err.Error(), string(val))
	}

	return nil
}

func Load(group, key string, confPtr interface{}, confType string) error {
	t := reflect.TypeOf(confPtr)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("conf is not ptr, kind is %v", t.Kind().String())
	}

	cc, err := service.GetClient(pb.Greeter_ServiceDesc.ServiceName)
	if err != nil {
		return err
	}
	c := pb.NewGreeterClient(cc)
	resp, err := c.Get(context.Background(), &pb.GetReq{Group: group, Key: key})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("val = ", resp.Val)

	err = load([]byte(resp.Val), confPtr, confType)
	if err != nil {
		log.Error(err)
		return err
	}

	configList.Lock()
	defer configList.Unlock()
	m, ok := configList.list[group]
	if !ok {
		m = make(map[string]interface{})
	}
	m[key] = confPtr
	configList.list[group] = m

	return nil
}

func Watch(group, key string, confPtr interface{}, confType string) error {
	t := reflect.TypeOf(confPtr)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("conf is not ptr, kind is %v", t.Kind().String())
	}
	t = t.Elem()

	list, err := service.GetClientList(pb.Greeter_ServiceDesc.ServiceName)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, cc := range list {
		c := pb.NewGreeterClient(cc)
		stream, err := c.Watch(context.Background(), &pb.WatchReq{Group: group, Key: key})
		if err != nil {
			log.Error(err)
			return err
		}
		go func(sm pb.Greeter_WatchClient) {
			for {
				res, err := sm.Recv()
				if err != nil {
					log.Error(err)
					return
				}

				switch res.Type {
				case pb.WatchType_PUT:
					cfg := reflect.New(t).Interface()
					err := load(res.Val, &cfg, confType)
					if err != nil {
						log.Error(err)
					}

					configList.Lock()
					m, ok := configList.list[res.Group]
					if !ok {
						m = make(map[string]interface{})
					}
					m[res.Key] = &cfg
					configList.list[res.Group] = m
					configList.Unlock()

				case pb.WatchType_DELETE:
					configList.Lock()
					m, ok := configList.list[res.Group]
					if ok {
						delete(m, res.Key)
						configList.list[res.Group] = m
					}
					configList.Unlock()
				}
			}
		}(stream)
	}

	return nil
}

func Get(group, key string) (interface{}, error) {
	configList.RLock()
	defer configList.RUnlock()
	m, ok := configList.list[group]
	if !ok {
		return nil, fmt.Errorf("group not found config group = %v, key = %v", group, key)
	}
	conf, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("key not found config group = %v, key = %v", group, key)
	}
	return conf, nil
}
