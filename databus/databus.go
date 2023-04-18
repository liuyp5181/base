package databus

import (
	"context"
	"encoding/json"
	"errors"
	"etcd/pkg/cache"
	"etcd/pkg/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
	"time"
)

const RECOVER_QUEUE = "databus:recover"

type DataBus interface {
	GetKey() string
	GetDB() (bool, error)
	SetDB() error
	GetCache() (bool, error)
	SetCache() error
	SetNilCache() error
}

var (
	typeList = make(map[string]reflect.Type)
	rKey     string
	wKey     string
)
var retry = errors.New("retry")

func Init(redisKey string, watchCnt int, dbs ...DataBus) {
	if len(dbs) == 0 {
		panic("databus is nil")
	}
	for _, d := range dbs {
		typeList[d.GetKey()] = reflect.TypeOf(d)
	}
	rKey = redisKey
	wKey = RECOVER_QUEUE + ":" + config.ServiceName
	for i := 0; i < watchCnt; i++ {
		go watch()
	}
}

func pushFill(d DataBus) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	r := cache.GetRedis(rKey)
	_, err = r.SAdd(context.Background(), wKey, d.GetKey()+","+string(data)).Result()
	if err != nil {
		return err
	}
	return nil
}

func popFill() (DataBus, error) {
	r := cache.GetRedis(rKey)
	ret, err := r.SPop(context.Background(), wKey).Result()
	if err != nil {
		return nil, err
	}
	idx := strings.Index(ret, ",{")
	key := ret[:idx]
	data := ret[idx+1:]

	t, ok := typeList[key]
	if !ok {
		return nil, fmt.Errorf("key=%s, not found", key)
	}
	var d = reflect.New(t).Interface()
	err = json.Unmarshal([]byte(data), &d)
	if err != nil {
		return nil, err
	}
	return d.(DataBus), nil
}

func watch() {
	exp := time.Millisecond * 10
	nilExp := time.Millisecond * 100
	var tick = time.NewTimer(exp)
	for {
		<-tick.C
		d, err := popFill()
		if err != nil {
			tick.Reset(nilExp)
			if err == redis.Nil {
				continue
			}
			// log
			continue
		}
		err = syncData(d)
		if err != nil {
			continue
		}
		tick.Reset(exp)
	}
}

func syncData(d DataBus) error {
	is, err := d.GetDB()
	if err != nil {
		return err
	}
	if !is {
		err = d.SetNilCache()
		if err != nil {
			return err
		}
	} else {
		err = d.SetCache()
		if err != nil {
			return err
		}
	}
	return nil
}

func AGet(d DataBus) (bool, error) {
	is, err := d.GetCache()
	if err != nil {
		return false, err
	}
	if !is {
		pushFill(d)
		return false, retry
	}

	return true, nil
}

func Get(d DataBus) (bool, error) {
	is, err := d.GetCache()
	if err != nil {
		return false, err
	}
	if !is {
		err = syncData(d)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func Set(d DataBus) error {
	err := d.SetDB()
	if err != nil {
		return err
	}
	err = d.SetCache()
	if err != nil {
		return err
	}
	return nil
}
