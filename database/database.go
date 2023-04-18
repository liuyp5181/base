package database

import (
	"fmt"
	"github.com/liuyp5181/base/config"
)

type Option func(*config.Database)

func WithHost(Host string) Option {
	return func(o *config.Database) {
		o.Host = Host
	}
}

func WithPort(port int) Option {
	return func(o *config.Database) {
		o.Port = port
	}
}

func WithUser(user string) Option {
	return func(o *config.Database) {
		o.User = user
	}
}

func WithPass(pass string) Option {
	return func(o *config.Database) {
		o.Pass = pass
	}
}

func WithDB(db string) Option {
	return func(o *config.Database) {
		o.DB = db
	}
}

func connect(conf config.Database) error {
	switch conf.Type {
	case "mysql":
		return NewMysql(conf)
	}
	return fmt.Errorf("not found type = %v", conf.Type)
}

func Connect(name string, opts ...Option) error {
	list := config.GetConfig().Database
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

func ConnectByConfig(conf config.Database) error {
	return connect(conf)
}
