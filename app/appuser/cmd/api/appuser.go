package main

import (
	"flag"
	"fmt"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"

	// SkyWalking Go Agent 集成
	_ "github.com/apache/skywalking-go"
)

var configFile = flag.String("f", "etc/appuser.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting AppUser API server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("SkyWalking Agent: Enabled (Auto-instrumentation)\n")
	server.Start()
}
