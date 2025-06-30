package logic

import (
	"context"
	"fmt"
	"time"

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
	if in.Id <= 0 {
		return nil, fmt.Errorf("产品ID不能为空")
	}

	// 查询产品是否存在
	product, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
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
	product.MinAmount = in.MinAmount
	product.MaxAmount = in.MaxAmount
	product.MinDuration = uint64(in.MinDuration)
	product.MaxDuration = uint64(in.MaxDuration)
	product.InterestRate = in.InterestRate
	product.Description = in.Description
	product.UpdatedAt = time.Now()

	err = l.svcCtx.LoanProductModel.Update(l.ctx, product)
	if err != nil {
		l.Errorf("更新产品失败: %v", err)
		return nil, fmt.Errorf("更新产品失败")
	}

	// 查询更新后的产品信息
	updatedProduct, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("查询更新后的产品失败: %v", err)
		return nil, fmt.Errorf("更新成功但查询失败")
	}

	return &loanproduct.UpdateLoanProductResp{
		Data: &loanproduct.LoanProductInfo{
			Id:           int64(updatedProduct.Id),
			Name:         updatedProduct.Name,
			Type:         updatedProduct.Type,
			MinAmount:    updatedProduct.MinAmount,
			MaxAmount:    updatedProduct.MaxAmount,
			MinDuration:  int32(updatedProduct.MinDuration),
			MaxDuration:  int32(updatedProduct.MaxDuration),
			InterestRate: updatedProduct.InterestRate,
			Description:  updatedProduct.Description,
			Status:       int32(updatedProduct.Status),
			CreatedAt:    updatedProduct.CreatedAt.Unix(),
			UpdatedAt:    updatedProduct.UpdatedAt.Unix(),
		},
	}, nil
}

// validateUpdateRequest 验证更新请求参数
func (l *UpdateLoanProductLogic) validateUpdateRequest(in *loanproduct.UpdateLoanProductReq) error {
	if in.Name == "" {
		return fmt.Errorf("产品名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("产品类型不能为空")
	}
	if in.MinAmount <= 0 {
		return fmt.Errorf("最小金额必须大于0")
	}
	if in.MaxAmount <= 0 {
		return fmt.Errorf("最大金额必须大于0")
	}
	if in.MinAmount >= in.MaxAmount {
		return fmt.Errorf("最小金额不能大于等于最大金额")
	}
	if in.MinDuration <= 0 {
		return fmt.Errorf("最小期限必须大于0")
	}
	if in.MaxDuration <= 0 {
		return fmt.Errorf("最大期限必须大于0")
	}
	if in.MinDuration >= in.MaxDuration {
		return fmt.Errorf("最小期限不能大于等于最大期限")
	}
	if in.InterestRate <= 0 {
		return fmt.Errorf("利率必须大于0")
	}
	return nil
}
