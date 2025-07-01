package logic

import (
	"context"

	"leaseproductrpc/internal/svc"
	"leaseproductrpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseProductLogic {
	return &GetLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品查询
func (l *GetLeaseProductLogic) GetLeaseProduct(in *leaseproduct.GetLeaseProductReq) (*leaseproduct.GetLeaseProductResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.GetLeaseProductResp{}, nil
}
