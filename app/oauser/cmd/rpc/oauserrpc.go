package main

import (
	"flag"
	"fmt"

	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/oauserrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		oauser.RegisterOaUserServer(grpcServer, server.NewOaUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 将本服务注册到 Consul
	err := consul.RegisterService(c.ListenOn, c.Consul)
	if err != nil {
		logx.Errorf("consul register service error: %s", err)
	} else {
		logx.Infof("consul register service success: %s", c.Consul.Key)
	}

	defer s.Stop()

	fmt.Printf("Starting OAUser RPC server at %s...\n", c.ListenOn)
	s.Start()
}
