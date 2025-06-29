package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseProductLogic {
	return &CreateLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品管理
func (l *CreateLeaseProductLogic) CreateLeaseProduct(in *leaseproduct.CreateLeaseProductReq) (*leaseproduct.CreateLeaseProductResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.CreateLeaseProductResp{}, nil
}
