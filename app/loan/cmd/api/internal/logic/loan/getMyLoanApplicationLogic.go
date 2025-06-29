package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyLoanApplicationLogic {
	return &GetMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyLoanApplicationLogic) GetMyLoanApplication() (resp *types.GetLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
