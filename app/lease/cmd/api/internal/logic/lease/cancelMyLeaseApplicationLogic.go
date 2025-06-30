package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

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
	// 调用 Lease RPC 撤销申请
	_, err = l.svcCtx.LeaseRpc.CancelLeaseApplication(l.ctx, &leaseclient.CancelLeaseApplicationReq{
		ApplicationId: req.ApplicationId,
		Reason:        req.Reason,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.CancelLeaseApplicationResp{}, nil
}
