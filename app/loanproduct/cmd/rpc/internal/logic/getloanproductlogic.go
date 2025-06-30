package logic

import (
	"context"

	"model"
	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanProductLogic {
	return &GetLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品查询
func (l *GetLoanProductLogic) GetLoanProduct(in *loanproduct.GetLoanProductReq) (*loanproduct.GetLoanProductResp, error) {
	// 参数验证 - 支持ID或产品编码查询
	if in.Id <= 0 && in.ProductCode == "" {
		return &loanproduct.GetLoanProductResp{
			Code:    400,
			Message: "产品ID或产品编码不能都为空",
		}, nil
	}

	var product *model.LoanProducts
	var err error

	// 根据查询条件选择查询方式
	if in.ProductCode != "" {
		// 优先使用产品编码查询
		product, err = l.svcCtx.LoanProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	} else {
		// 使用产品ID查询
		product, err = l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	}

	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &loanproduct.GetLoanProductResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	return &loanproduct.GetLoanProductResp{
		Code:    200,
		Message: "查询成功",
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
