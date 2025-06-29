package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMyLoanApplicationLogic {
	return &CancelMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMyLoanApplicationLogic) CancelMyLoanApplication(req *types.CancelLoanApplicationReq) (resp *types.CancelLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
