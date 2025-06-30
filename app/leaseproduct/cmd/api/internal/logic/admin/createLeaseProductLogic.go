package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建租赁产品
func NewCreateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseProductLogic {
	return &CreateLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLeaseProductLogic) CreateLeaseProduct(req *types.CreateLeaseProductReq) (resp *types.CreateLeaseProductResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.LeaseProductRpc.CreateLeaseProduct(l.ctx, &leaseproduct.CreateLeaseProductReq{
		ProductCode:    req.ProductCode,
		Name:           req.Name,
		Type:           req.Type,
		Machinery:      req.Machinery,
		Brand:          req.Brand,
		Model:          req.Model,
		DailyRate:      req.DailyRate,
		Deposit:        req.Deposit,
		MaxDuration:    req.MaxDuration,
		MinDuration:    req.MinDuration,
		Description:    req.Description,
		InventoryCount: req.InventoryCount,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.CreateLeaseProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 检查RPC响应
	if rpcResp.Code != 200 {
		return &types.CreateLeaseProductResp{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 转换响应数据
	return &types.CreateLeaseProductResp{
		Code:    200,
		Message: "创建成功",
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
