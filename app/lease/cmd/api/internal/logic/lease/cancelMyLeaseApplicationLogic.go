package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelMyLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMyLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMyLeaseApplicationLogic {
	return &CancelMyLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMyLeaseApplicationLogic) CancelMyLeaseApplication(req *types.CancelLeaseApplicationReq) (resp *types.CancelLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
