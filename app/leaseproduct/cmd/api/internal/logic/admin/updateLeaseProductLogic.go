package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproduct"

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
	// 调用RPC服务
	rpcResp, err := l.svcCtx.LeaseProductRpc.UpdateLeaseProduct(l.ctx, &leaseproduct.UpdateLeaseProductReq{
		ProductCode: req.ProductCode,
		Name:        req.Name,
		Type:        req.Type,
		Machinery:   req.Machinery,
		Brand:       req.Brand,
		Model:       req.Model,
		DailyRate:   req.DailyRate,
		Deposit:     req.Deposit,
		MaxDuration: req.MaxDuration,
		MinDuration: req.MinDuration,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.UpdateLeaseProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 检查RPC响应
	if rpcResp.Code != 200 {
		return &types.UpdateLeaseProductResp{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 转换响应数据
	return &types.UpdateLeaseProductResp{
		Code:    200,
		Message: "更新成功",
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
