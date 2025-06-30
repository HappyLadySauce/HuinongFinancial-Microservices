package admin

import (
	"context"

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

func (l *GetLeaseProductDetailLogic) GetLeaseProductDetail(req *types.GetLeaseProductDetailReq) (resp *types.GetLeaseProductResp, err error) {
	// 调用 LeaseProduct RPC 获取产品详情
	rpcResp, err := l.svcCtx.LeaseProductRpc.GetLeaseProduct(l.ctx, &leaseproductservice.GetLeaseProductReq{
		ProductCode: req.ProductCode,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用LeaseProduct RPC失败: %v", err)
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
