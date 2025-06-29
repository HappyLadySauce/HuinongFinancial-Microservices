package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建贷款产品
func NewCreateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanProductLogic {
	return &CreateLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanProductLogic) CreateLoanProduct(req *types.CreateLoanProductReq) (resp *types.CreateLoanProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
