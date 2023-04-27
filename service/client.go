package service

import (
	"context"
	"fmt"
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/etcd"
	"github.com/liuyp5181/base/log"
	"github.com/liuyp5181/base/service/extend"
	"github.com/liuyp5181/base/service/proxy"
	"github.com/liuyp5181/base/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	Server     *etcd.Service
	Conn       *grpc.ClientConn
	proxy      *proxy.Proxy
	cancelFunc context.CancelFunc
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

func (c *Client) close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
}
