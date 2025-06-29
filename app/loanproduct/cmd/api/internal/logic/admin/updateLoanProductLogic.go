package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新贷款产品
func NewUpdateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanProductLogic {
	return &UpdateLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLoanProductLogic) UpdateLoanProduct(req *types.UpdateLoanProductReq) (resp *types.UpdateLoanProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
