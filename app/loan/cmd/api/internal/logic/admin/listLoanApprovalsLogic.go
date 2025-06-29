package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApprovalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoanApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApprovalsLogic {
	return &ListLoanApprovalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoanApprovalsLogic) ListLoanApprovals(req *types.ListLoanApprovalsReq) (resp *types.ListLoanApprovalsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
