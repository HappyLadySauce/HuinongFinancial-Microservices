package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApprovalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLeaseApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApprovalsLogic {
	return &ListLeaseApprovalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLeaseApprovalsLogic) ListLeaseApprovals(req *types.ListLeaseApprovalsReq) (resp *types.ListLeaseApprovalsResp, err error) {
	// 调用 Lease RPC 获取审批记录
	rpcResp, err := l.svcCtx.LeaseRpc.ListLeaseApprovals(l.ctx, &leaseclient.ListLeaseApprovalsReq{
		ApplicationId: req.ApplicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换审批记录列表
	var approvals []types.LeaseApprovalInfo
	for _, item := range rpcResp.List {
		approvals = append(approvals, types.LeaseApprovalInfo{
			Id:               item.Id,
			ApplicationId:    item.ApplicationId,
			AuditorId:        item.AuditorId,
			AuditorName:      item.AuditorName,
			Action:           item.Action,
			Suggestions:      item.Suggestions,
			ApprovedDuration: item.ApprovedDuration,
			ApprovedAmount:   item.ApprovedAmount,
			ApprovedDeposit:  item.ApprovedDeposit,
			CreatedAt:        item.CreatedAt,
		})
	}

	return &types.ListLeaseApprovalsResp{
		List: approvals,
	}, nil
}
