package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLoanApplicationLogic {
	return &ApproveLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLoanApplicationLogic) ApproveLoanApplication(req *types.ApproveLoanApplicationReq) (resp *types.ApproveLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
