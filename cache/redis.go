package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/liuyp5181/base/config"
	"reflect"
)

var redisList = make(map[string]*redis.Client)

func ConnectRedis(cfg config.Cache) error {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // redis地址
		Password: cfg.Pass,                                 // redis密码，没有则留空
		DB:       cfg.DB,                                   // 默认数据库，默认是0
	})
	err := c.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	redisList[cfg.Name] = c
	return nil
}

func GetRedis(name string) *redis.Client {
	return redisList[name]
}

func HMGet(name, key string, vPtr interface{}, fields ...string) error {
	r := GetRedis(name)

	t := reflect.TypeOf(vPtr)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("vPtr is not ptr, is %v", t.Kind())
	}
	if len(fields) == 0 {
		t = t.Elem()
		for i := 0; i < t.NumField(); i++ {
			// 获取每个成员的结构体字段类型
			f := t.Field(i)
			fields = append(fields, f.Name)
		}
	}

	vals, err := r.HMGet(context.Background(), key, fields...).Result()
	if err != nil {
		return err
	}
	v := reflect.ValueOf(vPtr).Elem()
	for i, vl := range vals {
		if vl == nil {
			continue
		}
		v.FieldByName(fields[i]).Set(reflect.ValueOf(vl))
	}

	return nil
}

func HMSet(name, key string, vPtr interface{}, fields ...string) error {
	r := GetRedis(name)

	t := reflect.TypeOf(vPtr)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("vPtr is not ptr, is %v", t.Kind())
	}
	if len(fields) == 0 {
		t = t.Elem()
		for i := 0; i < t.NumField(); i++ {
			// 获取每个成员的结构体字段类型
			f := t.Field(i)
			fields = append(fields, f.Name)
		}
	}

	var vals []interface{}
	v := reflect.ValueOf(vPtr).Elem()
	for _, f := range fields {
		vals = append(vals, f, v.FieldByName(f).Interface())
	}

	_, err := r.HMSet(context.Background(), key, vals).Result()
	return err
}
