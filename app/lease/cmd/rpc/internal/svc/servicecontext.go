package svc

import (
	"model"
	"rpc/internal/clients"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	LeaseApplicationsModel model.LeaseApplicationsModel
	LeaseApprovalsModel    model.LeaseApprovalsModel

	// RPC 客户端 - 通过consul服务发现调用其他服务
	LeaseProductClient clients.LeaseProductClient
	AppUserClient      clients.AppUserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                 c,
		LeaseApplicationsModel: model.NewLeaseApplicationsModel(conn, c.CacheConf),
		LeaseApprovalsModel:    model.NewLeaseApprovalsModel(conn, c.CacheConf),

		// 通过consul服务发现初始化RPC客户端
		LeaseProductClient: clients.NewLeaseProductClient(zrpc.MustNewClient(c.LeaseProductRpc)),
		AppUserClient:      clients.NewAppUserClient(zrpc.MustNewClient(c.AppUserRpc)),
	}
}
