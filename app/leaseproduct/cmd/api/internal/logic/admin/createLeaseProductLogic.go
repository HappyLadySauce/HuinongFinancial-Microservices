package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建租赁产品
func NewCreateLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseProductLogic {
	return &CreateLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLeaseProductLogic) CreateLeaseProduct(req *types.CreateLeaseProductReq) (resp *types.CreateLeaseProductResp, err error) {
	// todo: add your logic here and delete this line

	return
}
