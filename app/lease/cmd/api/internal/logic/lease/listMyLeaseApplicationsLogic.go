package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyLeaseApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLeaseApplicationsLogic {
	return &ListMyLeaseApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyLeaseApplicationsLogic) ListMyLeaseApplications(req *types.ListLeaseApplicationsReq) (resp *types.ListLeaseApplicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
