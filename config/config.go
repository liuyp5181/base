package config

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/liuyp5181/base/etcd"
	"github.com/liuyp5181/base/log"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	defaultConfigPath = "./config.yaml"
)

type Global struct {
	Namespace string `mapstructure:"namespace"`
}

type Server struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

type Database struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	DB   string `yaml:"db"`
}

type Cache struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Pass string `yaml:"pass"`
	DB   int    `yaml:"db"`
}

type Conf struct {
	Etcd     []etcd.Config `mapstructure:"Etcd"`
	Log      *log.Config   `mapstructure:"Log"`
	Global   *Global       `mapstructure:"global"`
	Server   Server        `mapstructure:"server"`
	Database []Database    `mapstructure:"database"`
	Cache    []Cache       `mapstructure:"cache"`
}

var (
	confPath    = defaultConfigPath
	localData   []byte
	vp          *viper.Viper
	cfg         Conf
	ServiceName = "Service"
)

func init() {
	flag.StringVar(&confPath, "conf", defaultConfigPath, "config file path")
}

func Init() {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, config-file=[%s], err_msg=[%s]", confPath, err.Error()))
	}
	localData = data

	err = readConfig()
	if err != nil {
		panic(fmt.Sprintf("readConfig failed, err_msg=[%s], content=[\n%s]", err.Error(), string(localData)))
	}

	err = LoadLocal(&cfg)
	if err != nil {
		panic(fmt.Sprintf("LoadLocal failed, err_msg=[%s] content=[\n%s]", err.Error(), string(localData)))
	}

	fmt.Println("config = ", string(data))
	fmt.Println("config = ", cfg)

	if cfg.Global.Namespace == "" {
		panic("namespace is nil")
	}

	if cfg.Log != nil {
		err = log.Init(cfg.Log)
		if err != nil {
			panic(fmt.Sprintf("init Log failed, config=[%+v], err_msg=[%s]", cfg.Log, err.Error()))
		}
	}

	var points []string
	for _, v := range cfg.Etcd {
		points = append(points, fmt.Sprintf("%s:%d", v.IP, v.Port))
	}

	var ec = clientv3.Config{
		Endpoints:   points,
		DialTimeout: 5 * time.Second,
	}
	err = etcd.Init(ec)
	if err != nil {
		panic(fmt.Sprintf("init Etcd failed, config=[%+v], err_msg=[%s]", ec, err.Error()))
	}
}

func readConfig() error {
	vp = viper.New()
	vp.SetConfigType("yaml")
	vp.AutomaticEnv()
	if err := vp.ReadConfig(bytes.NewReader(localData)); err != nil {
		return err
	}

	for _, k := range vp.AllKeys() {
		value := vp.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			tp := strings.TrimPrefix(value, "${")
			ts := strings.TrimSuffix(tp, "}")
			e := os.Getenv(ts)
			if len(e) == 0 {
				return fmt.Errorf("get env fail, key = %v", value)
			}
			vp.Set(k, e)
		}
	}
	return nil
}

// LoadLocal load config by local config file
func LoadLocal(conf interface{}) error {
	if err := vp.Unmarshal(conf); err != nil {
		return fmt.Errorf("local config unmarshal failed, err_msg=[%s], content=[\n%s]", err.Error(), string(localData))
	}

	return nil
}

func GetConfig() Conf {
	return cfg
}
