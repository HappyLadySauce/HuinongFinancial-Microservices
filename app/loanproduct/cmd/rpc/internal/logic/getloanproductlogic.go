package logic

import (
	"context"
	"fmt"

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
	// 参数验证：ID和产品编码不能都为空
	if in.Id <= 0 && in.ProductCode == "" {
		return nil, fmt.Errorf("产品ID或产品编码不能都为空")
	}

	var product *model.LoanProducts
	var err error

	// 根据参数查询产品
	if in.Id > 0 {
		// 通过ID查询
		product, err = l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	} else {
		// 通过产品编码查询
		product, err = l.svcCtx.LoanProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	}

	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return nil, fmt.Errorf("产品不存在")
	}

	return &loanproduct.GetLoanProductResp{
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
