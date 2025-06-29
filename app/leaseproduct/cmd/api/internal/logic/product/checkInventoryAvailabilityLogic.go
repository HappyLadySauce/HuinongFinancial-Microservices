package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckInventoryAvailabilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 检查库存可用性
func NewCheckInventoryAvailabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInventoryAvailabilityLogic {
	return &CheckInventoryAvailabilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckInventoryAvailabilityLogic) CheckInventoryAvailability(req *types.CheckInventoryReq) (resp *types.CheckInventoryResp, err error) {
	// todo: add your logic here and delete this line

	return
}
