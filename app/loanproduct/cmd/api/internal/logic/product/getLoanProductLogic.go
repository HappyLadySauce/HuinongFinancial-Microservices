package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取贷款产品详情
func NewGetLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanProductLogic {
	return &GetLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanProductLogic) GetLoanProduct() (resp *types.GetLoanProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
