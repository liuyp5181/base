package cache

import (
	"fmt"
	"github.com/liuyp5181/base/config"
)

type Option func(*config.Cache)

func WithHost(Host string) Option {
	return func(o *config.Cache) {
		o.Host = Host
	}
}

func WithPort(port int) Option {
	return func(o *config.Cache) {
		o.Port = port
	}
}

func WithPass(pass string) Option {
	return func(o *config.Cache) {
		o.Pass = pass
	}
}

func WithDB(db int) Option {
	return func(o *config.Cache) {
		o.DB = db
	}
}

func connect(conf config.Cache) error {
	switch conf.Type {
	case "redis":
		return ConnectRedis(conf)
	}
	return fmt.Errorf("not found type = %v", conf.Type)
}

func Connect(name string, opts ...Option) error {
	list := config.GetConfig().Cache
	for _, v := range list {
		if v.Name == name {
			for _, o := range opts {
				o(&v)
			}
			return connect(v)
		}
	}

	return nil
}

func ConnectByConfig(conf config.Cache) error {
	return connect(conf)
}
