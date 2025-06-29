package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLeaseProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseProductsLogic {
	return &ListLeaseProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLeaseProductsLogic) ListLeaseProducts(in *leaseproduct.ListLeaseProductsReq) (*leaseproduct.ListLeaseProductsResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.ListLeaseProductsResp{}, nil
}
