package main

import (
	"flag"
	"fmt"

	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// SkyWalking Go Agent 集成
	_ "github.com/apache/skywalking-go"
)

var configFile = flag.String("f", "etc/leaserpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		lease.RegisterLeaseServer(grpcServer, server.NewLeaseServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 将本服务注册到 Consul
	err := consul.RegisterService(c.ListenOn, c.Consul)
	if err != nil {
		logx.Errorf("consul register service %s", err)
	}

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
