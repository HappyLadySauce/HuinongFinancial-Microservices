package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLeaseProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有租赁产品列表
func NewListAllLeaseProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLeaseProductsLogic {
	return &ListAllLeaseProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLeaseProductsLogic) ListAllLeaseProducts(req *types.ListLeaseProductsReq) (resp *types.ListLeaseProductsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
