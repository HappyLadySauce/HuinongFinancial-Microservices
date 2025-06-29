package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanApplicationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoanApplicationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanApplicationDetailLogic {
	return &GetLoanApplicationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanApplicationDetailLogic) GetLoanApplicationDetail() (resp *types.GetLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
