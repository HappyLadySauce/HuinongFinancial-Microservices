package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取租赁产品详情
func NewGetLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseProductLogic {
	return &GetLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaseProductLogic) GetLeaseProduct() (resp *types.GetLeaseProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
