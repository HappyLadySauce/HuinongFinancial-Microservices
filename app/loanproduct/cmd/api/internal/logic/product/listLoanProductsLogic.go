package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取贷款产品列表
func NewListLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanProductsLogic {
	return &ListLoanProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoanProductsLogic) ListLoanProducts(req *types.ListLoanProductsReq) (resp *types.ListLoanProductsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
