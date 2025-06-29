package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApprovalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLeaseApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApprovalsLogic {
	return &ListLeaseApprovalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLeaseApprovalsLogic) ListLeaseApprovals(req *types.ListLeaseApprovalsReq) (resp *types.ListLeaseApprovalsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
