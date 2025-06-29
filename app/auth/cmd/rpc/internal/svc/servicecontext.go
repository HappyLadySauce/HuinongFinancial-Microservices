package svc

import (
	"rpc/internal/config"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	Redis   *redis.Redis
	JWTUtil *utils.JWTUtil

	// 其他服务 RPC 客户端
	AppUserRpc zrpc.Client
	OaUserRpc  zrpc.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.CacheConf[0].Host,
		Pass: c.CacheConf[0].Pass,
		Type: c.CacheConf[0].Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:     c,
		Redis:      rds,
		JWTUtil:    utils.NewJWTUtil(c.JwtAuth.AccessSecret, c.JwtAuth.AccessExpire),
		AppUserRpc: zrpc.MustNewClient(c.AppUserRpc),
		OaUserRpc:  zrpc.MustNewClient(c.OaUserRpc),
	}
}
