package logic

import (
	"context"
	"fmt"
	"time"

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
	if in.ProductCode == "" {
		return nil, fmt.Errorf("产品编码不能为空")
	}

	// 查询产品是否存在
	product, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return nil, fmt.Errorf("产品不存在")
	}

	// 验证更新参数
	if err := l.validateUpdateRequest(in); err != nil {
		return nil, err
	}

	// 更新产品信息
	product.Name = in.Name
	product.Type = in.Type
	product.Machinery = in.Machinery
	product.Brand = in.Brand
	product.Model = in.Model
	product.DailyRate = in.DailyRate
	product.Deposit = in.Deposit
	product.MaxDuration = uint64(in.MaxDuration)
	product.MinDuration = uint64(in.MinDuration)
	product.Description = in.Description
	product.Status = uint64(in.Status)
	product.UpdatedAt = time.Now()

	err = l.svcCtx.LeaseProductModel.Update(l.ctx, product)
	if err != nil {
		l.Errorf("更新产品失败: %v", err)
		return nil, fmt.Errorf("更新产品失败")
	}

	// 查询更新后的产品信息
	updatedProduct, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询更新后的产品失败: %v", err)
		return nil, fmt.Errorf("更新成功但查询失败")
	}

	return &leaseproduct.UpdateLeaseProductResp{
		Data: &leaseproduct.LeaseProductInfo{
			Id:             int64(updatedProduct.Id),
			ProductCode:    updatedProduct.ProductCode,
			Name:           updatedProduct.Name,
			Type:           updatedProduct.Type,
			Machinery:      updatedProduct.Machinery,
			Brand:          updatedProduct.Brand,
			Model:          updatedProduct.Model,
			DailyRate:      updatedProduct.DailyRate,
			Deposit:        updatedProduct.Deposit,
			MaxDuration:    int32(updatedProduct.MaxDuration),
			MinDuration:    int32(updatedProduct.MinDuration),
			Description:    updatedProduct.Description,
			InventoryCount: int32(updatedProduct.InventoryCount),
			AvailableCount: int32(updatedProduct.AvailableCount),
			Status:         int32(updatedProduct.Status),
			CreatedAt:      updatedProduct.CreatedAt.Unix(),
			UpdatedAt:      updatedProduct.UpdatedAt.Unix(),
		},
	}, nil
}

// validateUpdateRequest 验证更新请求参数
func (l *UpdateLeaseProductLogic) validateUpdateRequest(in *leaseproduct.UpdateLeaseProductReq) error {
	if in.Name == "" {
		return fmt.Errorf("产品名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("产品类型不能为空")
	}
	if in.DailyRate <= 0 {
		return fmt.Errorf("日租金必须大于0")
	}
	if in.Deposit < 0 {
		return fmt.Errorf("押金不能小于0")
	}
	if in.MaxDuration <= 0 {
		return fmt.Errorf("最大租期必须大于0")
	}
	if in.MinDuration <= 0 {
		return fmt.Errorf("最小租期必须大于0")
	}
	if in.MinDuration >= in.MaxDuration {
		return fmt.Errorf("最小租期不能大于等于最大租期")
	}
	if in.Status != 1 && in.Status != 2 {
		return fmt.Errorf("状态值必须为1(上架)或2(下架)")
	}
	return nil
}
