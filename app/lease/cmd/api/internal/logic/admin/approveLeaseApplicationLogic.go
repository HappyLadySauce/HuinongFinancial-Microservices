package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLeaseApplicationLogic {
	return &ApproveLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLeaseApplicationLogic) ApproveLeaseApplication(req *types.ApproveLeaseApplicationReq) (resp *types.ApproveLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
