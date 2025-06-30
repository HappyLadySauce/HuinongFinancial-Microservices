package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseProductLogic {
	return &GetLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品查询
func (l *GetLeaseProductLogic) GetLeaseProduct(in *leaseproduct.GetLeaseProductReq) (*leaseproduct.GetLeaseProductResp, error) {
	// 参数验证
	if in.ProductCode == "" {
		return &leaseproduct.GetLeaseProductResp{
			Code:    400,
			Message: "产品编码不能为空",
		}, nil
	}

	// 根据产品编码查询产品信息
	product, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &leaseproduct.GetLeaseProductResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	return &leaseproduct.GetLeaseProductResp{
		Code:    200,
		Message: "查询成功",
		Data: &leaseproduct.LeaseProductInfo{
			Id:             int64(product.Id),
			ProductCode:    product.ProductCode,
			Name:           product.Name,
			Type:           product.Type,
			Machinery:      product.Machinery,
			Brand:          product.Brand,
			Model:          product.Model,
			DailyRate:      product.DailyRate,
			Deposit:        product.Deposit,
			MaxDuration:    int32(product.MaxDuration),
			MinDuration:    int32(product.MinDuration),
			Description:    product.Description,
			InventoryCount: int32(product.InventoryCount),
			AvailableCount: int32(product.AvailableCount),
			Status:         int32(product.Status),
			CreatedAt:      product.CreatedAt.Unix(),
			UpdatedAt:      product.UpdatedAt.Unix(),
		},
	}, nil
}
