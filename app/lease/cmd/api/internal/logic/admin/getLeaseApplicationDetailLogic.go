package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseApplicationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLeaseApplicationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseApplicationDetailLogic {
	return &GetLeaseApplicationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaseApplicationDetailLogic) GetLeaseApplicationDetail() (resp *types.GetLeaseApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
