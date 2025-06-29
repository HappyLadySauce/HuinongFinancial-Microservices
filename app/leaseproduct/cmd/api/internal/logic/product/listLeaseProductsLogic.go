package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取租赁产品列表
func NewListLeaseProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseProductsLogic {
	return &ListLeaseProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLeaseProductsLogic) ListLeaseProducts(req *types.ListLeaseProductsReq) (resp *types.ListLeaseProductsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
