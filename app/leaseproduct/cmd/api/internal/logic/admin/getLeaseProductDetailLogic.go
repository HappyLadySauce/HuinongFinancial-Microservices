package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取租赁产品详情
func NewGetLeaseProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseProductDetailLogic {
	return &GetLeaseProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaseProductDetailLogic) GetLeaseProductDetail() (resp *types.GetLeaseProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
