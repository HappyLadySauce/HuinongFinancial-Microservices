package logic

import (
	"context"
	"fmt"

	"model"
	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanProductLogic {
	return &CreateLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品管理
func (l *CreateLoanProductLogic) CreateLoanProduct(in *loanproduct.CreateLoanProductReq) (*loanproduct.CreateLoanProductResp, error) {
	// 参数验证
	if err := l.validateCreateRequest(in); err != nil {
		return nil, err
	}

	// 检查产品编码是否已存在
	existingProduct, err := l.svcCtx.LoanProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err == nil && existingProduct != nil {
		return nil, fmt.Errorf("产品编码已存在")
	}

	// 创建产品记录
	product := &model.LoanProducts{
		ProductCode:  in.ProductCode,
		Name:         in.Name,
		Type:         in.Type,
		MaxAmount:    in.MaxAmount,
		MinAmount:    in.MinAmount,
		MaxDuration:  uint64(in.MaxDuration),
		MinDuration:  uint64(in.MinDuration),
		InterestRate: in.InterestRate,
		Description:  in.Description,
		Status:       1, // 默认上架状态
	}

	result, err := l.svcCtx.LoanProductModel.Insert(l.ctx, product)
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
	createdProduct, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(productId))
	if err != nil {
		l.Errorf("查询创建的产品失败: %v", err)
		return nil, fmt.Errorf("创建成功但查询失败")
	}

	return &loanproduct.CreateLoanProductResp{
		Data: &loanproduct.LoanProductInfo{
			Id:           int64(createdProduct.Id),
			ProductCode:  createdProduct.ProductCode,
			Name:         createdProduct.Name,
			Type:         createdProduct.Type,
			MaxAmount:    createdProduct.MaxAmount,
			MinAmount:    createdProduct.MinAmount,
			MaxDuration:  int32(createdProduct.MaxDuration),
			MinDuration:  int32(createdProduct.MinDuration),
			InterestRate: createdProduct.InterestRate,
			Description:  createdProduct.Description,
			Status:       int32(createdProduct.Status),
			CreatedAt:    createdProduct.CreatedAt.Unix(),
			UpdatedAt:    createdProduct.UpdatedAt.Unix(),
		},
	}, nil
}

// validateCreateRequest 验证创建请求参数
func (l *CreateLoanProductLogic) validateCreateRequest(in *loanproduct.CreateLoanProductReq) error {
	if in.ProductCode == "" {
		return fmt.Errorf("产品编码不能为空")
	}
	if in.Name == "" {
		return fmt.Errorf("产品名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("产品类型不能为空")
	}
	if in.MaxAmount <= 0 {
		return fmt.Errorf("最大金额必须大于0")
	}
	if in.MinAmount <= 0 {
		return fmt.Errorf("最小金额必须大于0")
	}
	if in.MinAmount >= in.MaxAmount {
		return fmt.Errorf("最小金额不能大于等于最大金额")
	}
	if in.MaxDuration <= 0 {
		return fmt.Errorf("最大期限必须大于0")
	}
	if in.MinDuration <= 0 {
		return fmt.Errorf("最小期限必须大于0")
	}
	if in.MinDuration >= in.MaxDuration {
		return fmt.Errorf("最小期限不能大于等于最大期限")
	}
	if in.InterestRate <= 0 {
		return fmt.Errorf("利率必须大于0")
	}
	if in.Description == "" {
		return fmt.Errorf("产品描述不能为空")
	}
	return nil
}
