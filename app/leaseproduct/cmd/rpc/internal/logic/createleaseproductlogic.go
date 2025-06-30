package logic

import (
	"context"
	"fmt"

	"model"
	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseProductLogic {
	return &CreateLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品管理
func (l *CreateLeaseProductLogic) CreateLeaseProduct(in *leaseproduct.CreateLeaseProductReq) (*leaseproduct.CreateLeaseProductResp, error) {
	// 参数验证
	if err := l.validateCreateRequest(in); err != nil {
		return nil, err
	}

	// 检查产品编码是否已存在
	existingProduct, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err == nil && existingProduct != nil {
		return nil, fmt.Errorf("产品编码已存在")
	}

	// 创建产品记录
	product := &model.LeaseProducts{
		ProductCode:    in.ProductCode,
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
		InventoryCount: uint64(in.InventoryCount),
		AvailableCount: uint64(in.InventoryCount), // 初始可用数量等于库存数量
		Status:         1,                         // 默认上架状态
	}

	result, err := l.svcCtx.LeaseProductModel.Insert(l.ctx, product)
	if err != nil {
		l.Errorf("创建产品失败: %v", err)
		return nil, fmt.Errorf("创建产品失败")
	}

	// 获取新创建的产品ID
	productId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取新产品ID失败: %v", err)
		return nil, fmt.Errorf("创建成功但查询失败")
	}

	// 查询完整的产品信息
	createdProduct, err := l.svcCtx.LeaseProductModel.FindOne(l.ctx, uint64(productId))
	if err != nil {
		l.Errorf("查询创建的产品失败: %v", err)
		return nil, fmt.Errorf("创建成功但查询失败")
	}

	return &leaseproduct.CreateLeaseProductResp{
		Data: &leaseproduct.LeaseProductInfo{
			Id:             int64(createdProduct.Id),
			ProductCode:    createdProduct.ProductCode,
			Name:           createdProduct.Name,
			Type:           createdProduct.Type,
			Machinery:      createdProduct.Machinery,
			Brand:          createdProduct.Brand,
			Model:          createdProduct.Model,
			DailyRate:      createdProduct.DailyRate,
			Deposit:        createdProduct.Deposit,
			MaxDuration:    int32(createdProduct.MaxDuration),
			MinDuration:    int32(createdProduct.MinDuration),
			Description:    createdProduct.Description,
			InventoryCount: int32(createdProduct.InventoryCount),
			AvailableCount: int32(createdProduct.AvailableCount),
			Status:         int32(createdProduct.Status),
			CreatedAt:      createdProduct.CreatedAt.Unix(),
			UpdatedAt:      createdProduct.UpdatedAt.Unix(),
		},
	}, nil
}

// validateCreateRequest 验证创建请求参数
func (l *CreateLeaseProductLogic) validateCreateRequest(in *leaseproduct.CreateLeaseProductReq) error {
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
	if in.InventoryCount <= 0 {
		return fmt.Errorf("库存数量必须大于0")
	}
	if in.Description == "" {
		return fmt.Errorf("产品描述不能为空")
	}
	return nil
}
