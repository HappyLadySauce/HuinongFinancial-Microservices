package logic

import (
	"context"
	"fmt"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaseProductLogic {
	return &UpdateLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLeaseProductLogic) UpdateLeaseProduct(in *leaseproduct.UpdateLeaseProductReq) (*leaseproduct.UpdateLeaseProductResp, error) {
	// 参数验证
	if err := l.validateUpdateRequest(in); err != nil {
		return &leaseproduct.UpdateLeaseProductResp{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 检查产品是否存在
	existingProduct, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &leaseproduct.UpdateLeaseProductResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	// 更新产品信息
	updatedProduct := &model.LeaseProducts{
		Id:             existingProduct.Id,
		ProductCode:    existingProduct.ProductCode, // 产品编码不能修改
		Name:           in.Name,
		Type:           in.Type,
		Machinery:      in.Machinery,
		Brand:          in.Brand,
		Model:          in.Model,
		DailyRate:      in.DailyRate,
		Deposit:        in.Deposit,
		MaxDuration:    uint64(in.MaxDuration),
		MinDuration:    uint64(in.MinDuration),
		Description:    in.Description,
		InventoryCount: existingProduct.InventoryCount, // 库存数量保持不变
		AvailableCount: existingProduct.AvailableCount, // 可用数量保持不变
		Status:         uint64(in.Status),
		CreatedAt:      existingProduct.CreatedAt, // 创建时间保持不变
		UpdatedAt:      time.Now(),
	}

	err = l.svcCtx.LeaseProductModel.Update(l.ctx, updatedProduct)
	if err != nil {
		l.Errorf("更新产品失败: %v", err)
		return &leaseproduct.UpdateLeaseProductResp{
			Code:    500,
			Message: "更新产品失败",
		}, nil
	}

	// 查询更新后的产品信息
	product, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询更新后的产品失败: %v", err)
		return &leaseproduct.UpdateLeaseProductResp{
			Code:    500,
			Message: "更新成功但查询失败",
		}, nil
	}

	return &leaseproduct.UpdateLeaseProductResp{
		Code:    200,
		Message: "更新成功",
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

// validateUpdateRequest 验证更新请求参数
func (l *UpdateLeaseProductLogic) validateUpdateRequest(in *leaseproduct.UpdateLeaseProductReq) error {
	if in.ProductCode == "" {
		return fmt.Errorf("产品编码不能为空")
	}
	if in.Name == "" {
		return fmt.Errorf("产品名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("产品类型不能为空")
	}
	if in.DailyRate <= 0 {
		return fmt.Errorf("日租金必须大于0")
	}
	if in.MinDuration <= 0 {
		return fmt.Errorf("最小租期必须大于0")
	}
	if in.MaxDuration <= 0 {
		return fmt.Errorf("最大租期必须大于0")
	}
	if in.MinDuration > in.MaxDuration {
		return fmt.Errorf("最小租期不能大于最大租期")
	}
	if in.Status != 1 && in.Status != 2 {
		return fmt.Errorf("状态值必须为1(上架)或2(下架)")
	}
	if in.Description == "" {
		return fmt.Errorf("产品描述不能为空")
	}
	return nil
}
