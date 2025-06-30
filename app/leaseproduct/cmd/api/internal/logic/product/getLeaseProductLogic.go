package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproduct"

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

func (l *GetLeaseProductLogic) GetLeaseProduct(req *types.GetLeaseProductReq) (resp *types.GetLeaseProductResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.LeaseProductRpc.GetLeaseProduct(l.ctx, &leaseproduct.GetLeaseProductReq{
		ProductCode: req.ProductCode,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.GetLeaseProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 检查RPC响应
	if rpcResp.Code != 200 {
		return &types.GetLeaseProductResp{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 转换响应数据
	return &types.GetLeaseProductResp{
		Code:    200,
		Message: "查询成功",
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
