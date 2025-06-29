package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaseProductLogic {
	return &UpdateLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLeaseProductLogic) UpdateLeaseProduct(in *leaseproduct.UpdateLeaseProductReq) (*leaseproduct.UpdateLeaseProductResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.UpdateLeaseProductResp{}, nil
}
