package monitor

import (
	"context"
	"flag"
	"fmt"
	"github.com/liuyp5181/base/config"
	"github.com/liuyp5181/base/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	pb "github.com/liuyp5181/base/client/monitor/api"
)

const defaultMonitorPort = 6226

var monitorPort = defaultMonitorPort

func init() {
	flag.IntVar(&monitorPort, "moPort", defaultMonitorPort, "monitor port")
}

func Init() {
	//提供 /metrics HTTP 端点
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", monitorPort), nil)
		if err != nil {
			panic(err)
		}
	}()
}

func Register() error {
	cc, err := service.GetClient(pb.Greeter_ServiceDesc.ServiceName)
	if err != nil {
		return err
	}
	c := pb.NewGreeterClient(cc)
	_, err = c.Register(context.Background(), &pb.RegisterReq{
		Name: config.ServiceName,
		Port: uint32(monitorPort),
	})
	if err != nil {
		return err
	}
	return nil
}
