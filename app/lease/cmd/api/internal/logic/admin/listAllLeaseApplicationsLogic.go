package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLeaseApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLeaseApplicationsLogic {
	return &ListAllLeaseApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLeaseApplicationsLogic) ListAllLeaseApplications(req *types.ListLeaseApplicationsReq) (resp *types.ListLeaseApplicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
