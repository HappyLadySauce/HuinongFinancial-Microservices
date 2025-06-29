package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// Consul 服务注册配置
	Consul consul.Conf
	// MySQL 数据库配置
	MySQL struct {
		DataSource string
	}
	// Redis 缓存配置
	CacheConf cache.CacheConf
}
