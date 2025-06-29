package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除租赁产品
func NewDeleteLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLeaseProductLogic {
	return &DeleteLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLeaseProductLogic) DeleteLeaseProduct() (resp *types.DeleteLeaseProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
