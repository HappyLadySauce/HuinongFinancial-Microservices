package main

import (
	"flag"
	"fmt"

	"rpc/appuser"
	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/appuserrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		appuser.RegisterAppUserServer(grpcServer, server.NewAppUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	err := consul.RegisterService(c.ListenOn, c.Consul)
	if err != nil {
		logx.Errorf("consul register service %s", err)
	}

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
