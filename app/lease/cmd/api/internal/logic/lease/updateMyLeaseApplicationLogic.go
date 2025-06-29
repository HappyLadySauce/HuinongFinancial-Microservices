package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyLeaseApplicationLogic {
	return &UpdateMyLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyLeaseApplicationLogic) UpdateMyLeaseApplication(req *types.UpdateLeaseApplicationReq) (resp *types.UpdateLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
