package logic

import (
	"context"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(in *loanproduct.UpdateProductStatusReq) (*loanproduct.UpdateProductStatusResp, error) {
	// 参数验证
	if in.Id <= 0 {
		return &loanproduct.UpdateProductStatusResp{
			Code:    400,
			Message: "产品ID不能为空",
		}, nil
	}

	if in.Status != 1 && in.Status != 2 {
		return &loanproduct.UpdateProductStatusResp{
			Code:    400,
			Message: "状态值必须为1(上架)或2(下架)",
		}, nil
	}

	// 检查产品是否存在
	existingProduct, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &loanproduct.UpdateProductStatusResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	// 更新产品状态
	updatedProduct := &model.LoanProducts{
		Id:           existingProduct.Id,
		ProductCode:  existingProduct.ProductCode,
		Name:         existingProduct.Name,
		Type:         existingProduct.Type,
		MaxAmount:    existingProduct.MaxAmount,
		MinAmount:    existingProduct.MinAmount,
		MaxDuration:  existingProduct.MaxDuration,
		MinDuration:  existingProduct.MinDuration,
		InterestRate: existingProduct.InterestRate,
		Description:  existingProduct.Description,
		Status:       uint64(in.Status), // 只更新状态
		CreatedAt:    existingProduct.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	err = l.svcCtx.LoanProductModel.Update(l.ctx, updatedProduct)
	if err != nil {
		l.Errorf("更新产品状态失败: %v", err)
		return &loanproduct.UpdateProductStatusResp{
			Code:    500,
			Message: "更新状态失败",
		}, nil
	}

	statusText := "上架"
	if in.Status == 2 {
		statusText = "下架"
	}

	return &loanproduct.UpdateProductStatusResp{
		Code:    200,
		Message: "产品已" + statusText,
	}, nil
}
