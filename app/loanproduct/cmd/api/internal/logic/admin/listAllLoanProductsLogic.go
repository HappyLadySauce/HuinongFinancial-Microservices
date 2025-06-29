package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLoanProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有贷款产品列表
func NewListAllLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLoanProductsLogic {
	return &ListAllLoanProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLoanProductsLogic) ListAllLoanProducts(req *types.ListLoanProductsReq) (resp *types.ListLoanProductsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
