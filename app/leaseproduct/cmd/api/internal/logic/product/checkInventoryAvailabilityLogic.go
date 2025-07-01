package product

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproductservice"

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
	// 调用RPC服务 - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "leaseproduct-rpc", func() (*leaseproductservice.CheckInventoryAvailabilityResp, error) {
		return l.svcCtx.LeaseProductRpc.CheckInventoryAvailability(l.ctx, &leaseproductservice.CheckInventoryAvailabilityReq{
			ProductCode: req.ProductCode,
			Quantity:    req.Quantity,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换响应数据
	return &types.CheckInventoryResp{
		Available:      rpcResp.Available,
		AvailableCount: rpcResp.AvailableCount,
	}, nil
}
