package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul    consul.Conf
	CacheConf cache.CacheConf
	JwtAuth   struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshExpire int64
	}

	// 其他服务 RPC 客户端配置
	AppUserRpc zrpc.RpcClientConf
	OaUserRpc  zrpc.RpcClientConf
}
