package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyLeaseApplicationLogic {
	return &GetMyLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyLeaseApplicationLogic) GetMyLeaseApplication() (resp *types.GetLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
