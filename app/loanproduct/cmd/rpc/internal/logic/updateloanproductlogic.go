package logic

import (
	"context"
	"fmt"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanProductLogic {
	return &UpdateLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanProductLogic) UpdateLoanProduct(in *loanproduct.UpdateLoanProductReq) (*loanproduct.UpdateLoanProductResp, error) {
	// 参数验证
	if err := l.validateUpdateRequest(in); err != nil {
		return &loanproduct.UpdateLoanProductResp{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 检查产品是否存在
	existingProduct, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &loanproduct.UpdateLoanProductResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	// 更新产品信息
	updatedProduct := &model.LoanProducts{
		Id:           existingProduct.Id,
		ProductCode:  existingProduct.ProductCode, // 产品编码不能修改
		Name:         in.Name,
		Type:         in.Type,
		MaxAmount:    in.MaxAmount,
		MinAmount:    in.MinAmount,
		MaxDuration:  uint64(in.MaxDuration),
		MinDuration:  uint64(in.MinDuration),
		InterestRate: in.InterestRate,
		Description:  in.Description,
		Status:       existingProduct.Status, // 状态保持不变，通过专门接口修改
		CreatedAt:    existingProduct.CreatedAt, // 创建时间保持不变
		UpdatedAt:    time.Now(),
	}

	err = l.svcCtx.LoanProductModel.Update(l.ctx, updatedProduct)
	if err != nil {
		l.Errorf("更新产品失败: %v", err)
		return &loanproduct.UpdateLoanProductResp{
			Code:    500,
			Message: "更新产品失败",
		}, nil
	}

	// 查询更新后的产品信息
	product, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("查询更新后的产品失败: %v", err)
		return &loanproduct.UpdateLoanProductResp{
			Code:    500,
			Message: "更新成功但查询失败",
		}, nil
	}

	return &loanproduct.UpdateLoanProductResp{
		Code:    200,
		Message: "更新成功",
		Data: &loanproduct.LoanProductInfo{
			Id:           int64(product.Id),
			ProductCode:  product.ProductCode,
			Name:         product.Name,
			Type:         product.Type,
			MaxAmount:    product.MaxAmount,
			MinAmount:    product.MinAmount,
			MaxDuration:  int32(product.MaxDuration),
			MinDuration:  int32(product.MinDuration),
			InterestRate: product.InterestRate,
			Description:  product.Description,
			Status:       int32(product.Status),
			CreatedAt:    product.CreatedAt.Unix(),
			UpdatedAt:    product.UpdatedAt.Unix(),
		},
	}, nil
}

// validateUpdateRequest 验证更新请求参数
func (l *UpdateLoanProductLogic) validateUpdateRequest(in *loanproduct.UpdateLoanProductReq) error {
	if in.Id <= 0 {
		return fmt.Errorf("产品ID不能为空")
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
	if in.MinAmount > in.MaxAmount {
		return fmt.Errorf("最小金额不能大于最大金额")
	}
	if in.MinDuration <= 0 {
		return fmt.Errorf("最小期限必须大于0")
	}
	if in.MaxDuration <= 0 {
		return fmt.Errorf("最大期限必须大于0")
	}
	if in.MinDuration > in.MaxDuration {
		return fmt.Errorf("最小期限不能大于最大期限")
	}
	if in.InterestRate < 0 {
		return fmt.Errorf("利率不能小于0")
	}
	if in.Description == "" {
		return fmt.Errorf("产品描述不能为空")
	}
	return nil
}
