package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproduct"

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
	// 调用RPC服务
	rpcResp, err := l.svcCtx.LeaseProductRpc.CheckInventoryAvailability(l.ctx, &leaseproduct.CheckInventoryAvailabilityReq{
		ProductCode: req.ProductCode,
		Quantity:    req.Quantity,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.CheckInventoryResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 返回响应
	return &types.CheckInventoryResp{
		Code:           rpcResp.Code,
		Message:        rpcResp.Message,
		Available:      rpcResp.Available,
		AvailableCount: rpcResp.AvailableCount,
	}, nil
}
