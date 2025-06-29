package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanApplicationLogic {
	return &CreateLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanApplicationLogic) CreateLoanApplication(req *types.CreateLoanApplicationReq) (resp *types.CreateLoanApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
