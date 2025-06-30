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
		return &types.ListLeaseApprovalsResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.ListLeaseApprovalsResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	// 转换审批记录列表
	if len(rpcResp.List) > 0 {
		resp.List = make([]types.LeaseApprovalInfo, len(rpcResp.List))
		for i, item := range rpcResp.List {
			resp.List[i] = types.LeaseApprovalInfo{
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
			}
		}
	} else {
		resp.List = make([]types.LeaseApprovalInfo, 0)
	}

	return resp, nil
}
