package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproductservice"

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

func (l *GetLeaseProductDetailLogic) GetLeaseProductDetail(productCode string) (resp *types.GetLeaseProductResp, err error) {
	// 调用RPC服务 - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "leaseproduct-rpc", func() (*leaseproductservice.GetLeaseProductResp, error) {
		return l.svcCtx.LeaseProductRpc.GetLeaseProduct(l.ctx, &leaseproductservice.GetLeaseProductReq{
			ProductCode: productCode,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换响应数据
	return &types.GetLeaseProductResp{
		Data: types.LeaseProductInfo{
			Id:             rpcResp.Data.Id,
			ProductCode:    rpcResp.Data.ProductCode,
			Name:           rpcResp.Data.Name,
			Type:           rpcResp.Data.Type,
			Machinery:      rpcResp.Data.Machinery,
			Brand:          rpcResp.Data.Brand,
			Model:          rpcResp.Data.Model,
			DailyRate:      rpcResp.Data.DailyRate,
			Deposit:        rpcResp.Data.Deposit,
			MaxDuration:    rpcResp.Data.MaxDuration,
			MinDuration:    rpcResp.Data.MinDuration,
			Description:    rpcResp.Data.Description,
			InventoryCount: rpcResp.Data.InventoryCount,
			AvailableCount: rpcResp.Data.AvailableCount,
			Status:         rpcResp.Data.Status,
			CreatedAt:      rpcResp.Data.CreatedAt,
			UpdatedAt:      rpcResp.Data.UpdatedAt,
		},
	}, nil
}
