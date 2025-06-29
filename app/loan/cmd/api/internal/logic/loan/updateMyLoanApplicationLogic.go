package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyLoanApplicationLogic {
	return &UpdateMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyLoanApplicationLogic) UpdateMyLoanApplication(req *types.UpdateLoanApplicationReq) (resp *types.UpdateLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
