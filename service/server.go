package service

import (
	"context"
	"etcd/pkg/config"
	"etcd/pkg/etcd"
	"etcd/pkg/log"
	"etcd/pkg/service/extend"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	HEALTHCHECK_SERVICE = "grpc.health.v1.Health"
)

var (
	version = "1.0.1"
)

type Server struct {
	name string
	sev  *grpc.Server
	lis  net.Listener
}

func (s *Server) Serve() {
	log.Info(s.name, "start")
	if err := s.sev.Serve(s.lis); err != nil {
		panic(err)
	}
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.sev.RegisterService(sd, ss)
}

func (s *Server) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.sev.GetServiceInfo()
}

func (s *Server) GetGrpcServer() *grpc.Server {
	return s.sev
}

// UnaryServerInterceptor 拦截器，相对于中间件
func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 扩展字段
	e := extend.NewContext(ctx)
	var tid = e.GetClient("trace_id")
	var uid = e.GetClient("user_id")

	var addr string
	pr, ok := peer.FromContext(ctx)
	if ok {
		addr = pr.Addr.String()
	}

	log.Infof("request  [%s] %s %s %s data: %+v", tid, uid, info.FullMethod, addr, req)

	resp, err = handler(ctx, req)
	if err != nil {
		log.Errorf("handler  [%s] %s %s err: %+v", tid, uid, info.FullMethod, err)
		return
	}

	log.Infof("response [%s] %s %s data: %+v", tid, uid, info.FullMethod, resp)
	return
}

func NewServer() *Server {
	name := config.ServiceName
	serverCfg := config.GetConfig().Server

	err := etcd.SetService(name, serverCfg.IP, serverCfg.Port, version, 100)
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverCfg.IP, serverCfg.Port))
	if err != nil {
		panic(err)
	}

	s := &Server{
		name: name,
		sev:  grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor)),
		lis:  listen,
	}

	// grpc反射
	// server端：从中获取所有的可变和不可变的服务，遍历获取所有的服务、方法、属性，添加到相应的对象中
	// client端：根据请求参数进行判断，使用不同的方法处理，并返回响应
	reflection.Register(s)

	// 心跳
	hs := health.NewServer()
	hs.SetServingStatus(HEALTHCHECK_SERVICE, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)

	return s
}
