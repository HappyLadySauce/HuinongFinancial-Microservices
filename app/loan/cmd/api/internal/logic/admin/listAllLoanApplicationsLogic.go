package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLoanApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLoanApplicationsLogic {
	return &ListAllLoanApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLoanApplicationsLogic) ListAllLoanApplications(req *types.ListLoanApplicationsReq) (resp *types.ListLoanApplicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
