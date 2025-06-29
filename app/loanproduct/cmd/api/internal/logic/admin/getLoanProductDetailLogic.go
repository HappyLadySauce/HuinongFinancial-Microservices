package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取贷款产品详情
func NewGetLoanProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanProductDetailLogic {
	return &GetLoanProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanProductDetailLogic) GetLoanProductDetail() (resp *types.GetLoanProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
