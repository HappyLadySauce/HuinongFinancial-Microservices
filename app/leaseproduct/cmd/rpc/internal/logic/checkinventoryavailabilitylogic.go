package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckInventoryAvailabilityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckInventoryAvailabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInventoryAvailabilityLogic {
	return &CheckInventoryAvailabilityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 库存检查
func (l *CheckInventoryAvailabilityLogic) CheckInventoryAvailability(in *leaseproduct.CheckInventoryAvailabilityReq) (*leaseproduct.CheckInventoryAvailabilityResp, error) {
	// todo: add your logic here and delete this line

	return &leaseproduct.CheckInventoryAvailabilityResp{}, nil
}
