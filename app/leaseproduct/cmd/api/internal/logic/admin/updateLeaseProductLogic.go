package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新租赁产品
func NewUpdateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaseProductLogic {
	return &UpdateLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLeaseProductLogic) UpdateLeaseProduct(req *types.UpdateLeaseProductReq) (resp *types.UpdateLeaseProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
