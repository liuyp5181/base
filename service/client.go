package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/etcd"
	"github.com/liuyp5181/base/log"
	"github.com/liuyp5181/base/service/extend"
	"github.com/liuyp5181/base/service/proxy"
	"github.com/liuyp5181/base/util"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Server *etcd.Service
	Conn   *grpc.ClientConn
	proxy  *proxy.Proxy
}

type Clients struct {
	sync.RWMutex
	list map[string][]*Client
	m    map[string]*Client
}

var clients = &Clients{
	list: map[string][]*Client{},
	m:    map[string]*Client{},
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

func setClient(key string, s etcd.Service) error {
	if s.Name == config.ServiceName && s.IP == config.GetConfig().Server.IP && s.Port == config.GetConfig().Server.Port {
		return nil
	}

	f := func() bool {
		clients.Lock()
		defer clients.Unlock()
		c, ok := clients.m[key]
		if ok {
			c.Server = &s
		}
		return ok
	}
	if f() {
		return nil
	}

	c, err := newClient(&s)
	if err != nil {
		return err
	}
	fmt.Println("set", key, c)
	clients.Lock()
	defer clients.Unlock()
	clients.m[key] = c
	clients.list[c.Server.Name] = append(clients.list[c.Server.Name], c)

	return nil
}

func delClient(key string) {
	fmt.Println("del", key)
	clients.Lock()
	defer clients.Unlock()
	c := clients.m[key]
	if c == nil {
		return
	}
	s := c.Server
	delete(clients.m, key)
	list := clients.list[s.Name]
	for i, v := range list {
		if s.IP == v.Server.IP && s.Port == v.Server.Port {
			clients.list[s.Name] = append(list[:i], list[i+1:]...)
			return
		}
	}
}

func getClient(name string) *Client {
	rand.Seed(time.Now().UnixNano())
	clients.RLock()
	defer clients.RUnlock()
	list := clients.list[name]
	fmt.Println("getClient", list, len(list))
	if len(list) == 0 {
		return nil
	}

	var max int
	for _, v := range list {
		max += v.Server.Power
	}
	r := rand.Intn(max)
	var index int
	for _, v := range list {
		if v.Server.Power == 0 {
			continue
		}
		index += v.Server.Power
		if r < index {
			return v
		}
	}

	return nil
}

// unaryClientInterceptor 拦截器，相对于中间件
func unaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	e := extend.NewContext(ctx)
	var tid string
	var uid string
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		ss := md.Get("trace_id")
		if len(ss) > 0 {
			tid = ss[0]
		}
		ss = md.Get("user_id")
		if len(ss) > 0 {
			uid = ss[0]
		}
	}
	if tid == "" {
		tid = util.GenerateId("trace_id", req)
		e.SetClient("trace_id", tid)
	}
	if uid == "" {
		uid = config.ServiceName
		e.SetClient("user_id", uid)
	}

	log.Infof("request  [%s] %s %s data: %+v", tid, uid, method, req)

	if err := invoker(e.Ctx, method, req, reply, cc, opts...); err != nil {
		log.Errorf("invoker  [%s] %s %s err: %v", tid, uid, method, err)
		return err
	}

	log.Infof("response [%s] %s %s data: %+v", tid, uid, method, reply)
	return nil
}

func newClient(s *etcd.Service) (*Client, error) {
	log.Info("newClient", s)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.IP, s.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"HealthCheckConfig": {"ServiceName": "%s"}}`, HEALTHCHECK_SERVICE)),
		grpc.WithUnaryInterceptor(unaryClientInterceptor))
	if err != nil {
		return nil, fmt.Errorf("dial err = %v", err)
	}
	p := proxy.NewClient(context.Background(), conn)
	c := &Client{Conn: conn, proxy: p, Server: s}
	return c, nil
}

func initClient(name string) error {
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

	go watchClient(name)

	return nil
}

func watchClient(name string) {
	rch := etcd.WatcherService(name)
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
				delClient(string(ev.Kv.Key))
			}
		}
	}
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

func GetClient(name string) (*Client, error) {
	c := getClient(name)
	fmt.Println("GetClient", c)
	if c == nil {
		return nil, fmt.Errorf("not found, name = %v", name)
	}

	// todo ping
	return c, nil
}

func GetClientList(name string) []*Client {
	clients.RLock()
	defer clients.RUnlock()
	list := clients.list[name]
	fmt.Println("GetClientList", list, len(list))
	return list
}

func (c *Client) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return c.Conn.Invoke(ctx, method, args, reply, opts...)
}

func (c *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return c.Conn.NewStream(ctx, desc, method, opts...)
}

func (c *Client) Proxy(ctx context.Context, methodName string, message []byte, opts ...grpc.CallOption) ([]byte, error) {
	rsp, err := c.proxy.Call(ctx, c.Server.Name, methodName, message, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
