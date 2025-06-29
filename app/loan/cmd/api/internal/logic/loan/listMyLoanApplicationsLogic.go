package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyLoanApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLoanApplicationsLogic {
	return &ListMyLoanApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyLoanApplicationsLogic) ListMyLoanApplications(req *types.ListLoanApplicationsReq) (resp *types.ListLoanApplicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
