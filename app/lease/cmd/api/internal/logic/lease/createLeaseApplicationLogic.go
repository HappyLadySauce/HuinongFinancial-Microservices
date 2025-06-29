package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseApplicationLogic {
	return &CreateLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLeaseApplicationLogic) CreateLeaseApplication(req *types.CreateLeaseApplicationReq) (resp *types.CreateLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
