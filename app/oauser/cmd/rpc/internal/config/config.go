package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
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

	// JWT 配置 (重命名避免与 zrpc.RpcServerConf.Auth 冲突)
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	// Logrus 日志配置
	Logger struct {
		ServiceName string
		Mode        string // file, console
		Path        string
		Level       string
		KeepDays    int
		Compress    bool
		MaxSize     int
		MaxBackups  int
	}
}
