package svc

import (
	"appuserrpc/appuserclient"
	"leaseproductrpc/leaseproductservice"
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	LeaseApplicationsModel model.LeaseApplicationsModel
	LeaseApprovalsModel    model.LeaseApprovalsModel

	// RPC 客户端 - 通过consul服务发现调用其他服务
	LeaseProductClient leaseproductservice.LeaseProductService
	AppUserClient      appuserclient.AppUser
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                 c,
		LeaseApplicationsModel: model.NewLeaseApplicationsModel(conn, c.CacheConf),
		LeaseApprovalsModel:    model.NewLeaseApprovalsModel(conn, c.CacheConf),

		// 通过consul服务发现初始化RPC客户端
		LeaseProductClient: leaseproductservice.NewLeaseProductService(zrpc.MustNewClient(c.LeaseProductRpc)),
		AppUserClient:      appuserclient.NewAppUser(zrpc.MustNewClient(c.AppUserRpc)),
	}
}
