package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除贷款产品
func NewDeleteLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanProductLogic {
	return &DeleteLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLoanProductLogic) DeleteLoanProduct() (resp *types.DeleteLoanProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
