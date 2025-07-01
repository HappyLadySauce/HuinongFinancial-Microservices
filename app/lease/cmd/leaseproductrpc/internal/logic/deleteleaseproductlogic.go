package logic

import (
	"context"

	"leaseproductrpc/internal/svc"
	"leaseproductrpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLeaseProductLogic {
	return &DeleteLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLeaseProductLogic) DeleteLeaseProduct(in *leaseproduct.DeleteLeaseProductReq) (*leaseproduct.DeleteLeaseProductResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.DeleteLeaseProductResp{}, nil
}
