package service

//
//import (
//	"bytes"
//	"context"
//	configmgrpb "etcd/api/configmgr"
//	"etcd/pkg/log"
//	"fmt"
//	"github.com/spf13/viper"
//	"reflect"
//	"sync"
//)
//
//var configList sync.Map
//
//func loadConfig(val []byte, conf interface{}, confType string) error {
//	vp := viper.New()
//	vp.SetConfigType(confType)
//	vp.AutomaticEnv()
//	if err := vp.ReadConfig(bytes.NewReader(val)); err != nil {
//		return fmt.Errorf("read config failed, err_msg=[%s], extend=[%s]", err.Error(), string(val))
//	}
//
//	if err := vp.Unmarshal(conf); err != nil {
//		return fmt.Errorf("local config unmarshal failed, err_msg=[%s], extend=[%s]", err.Error(), string(val))
//	}
//
//	return nil
//}
//
//func LoadConfig(key string, conf interface{}, confType string) error {
//	t := reflect.TypeOf(conf)
//	if t.Kind() != reflect.Ptr {
//		return fmt.Errorf("conf is not ptr, kind is %v", t.Kind().String())
//	}
//
//	cc := getClient(configmgrpb.Configmgr_ServiceDesc.ServiceName)
//	if cc == nil {
//		return fmt.Errorf("not found service, name = %s", configmgrpb.Configmgr_ServiceDesc.ServiceName)
//	}
//	c := configmgrpb.NewConfigmgrClient(cc)
//	resp, err := c.GetConfig(context.Background(), &configmgrpb.GetConfigReq{Key: key})
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	log.Info("val = ", resp.Val)
//
//	err = loadConfig([]byte(resp.Val), conf, confType)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	configList.Store(key, conf)
//
//	return nil
//}
//
//func WatchConfig(key string, conf interface{}, confType string) error {
//	t := reflect.TypeOf(conf)
//	if t.Kind() == reflect.Ptr {
//		t = t.Elem()
//	}
//	cfg := reflect.New(t).Interface()
//
//	list := GetClientList(configmgrpb.Configmgr_ServiceDesc.ServiceName)
//	if len(list) == 0 {
//		return fmt.Errorf("not found service, name = %s", configmgrpb.Configmgr_ServiceDesc.ServiceName)
//	}
//
//	for _, cc := range list {
//		c := configmgrpb.NewConfigmgrClient(cc)
//		stream, err := c.Watch(context.Background(), &configmgrpb.WatchReq{Key: key})
//		if err != nil {
//			log.Error(err)
//			return err
//		}
//		go func(sm configmgrpb.Configmgr_WatchClient) {
//			for {
//				res, err := sm.Recv()
//				if err != nil {
//					log.Error(err)
//					return
//				}
//				switch res.Type {
//				case configmgrpb.WatchType_PUT:
//					err := loadConfig(res.Val, &cfg, confType)
//					if err != nil {
//						log.Error(err)
//					}
//					configList.Store(res.Key, cfg)
//				case configmgrpb.WatchType_DELETE:
//					configList.Delete(res.Key)
//				}
//			}
//		}(stream)
//	}
//
//	return nil
//}
//
//func GetConfig(name string) interface{} {
//	val, ok := configList.Load(name)
//	if ok {
//		return val
//	}
//	return nil
//}
