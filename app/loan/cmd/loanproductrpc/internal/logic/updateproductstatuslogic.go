package logic

import (
	"context"

	"loanproductrpc/internal/svc"
	"loanproductrpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(in *loanproduct.UpdateProductStatusReq) (*loanproduct.UpdateProductStatusResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.UpdateProductStatusResp{}, nil
}
