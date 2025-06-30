package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApprovalsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoanApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApprovalsLogic {
	return &ListLoanApprovalsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoanApprovalsLogic) ListLoanApprovals(req *types.ListLoanApprovalsReq) (resp *types.ListLoanApprovalsResp, err error) {
	// 调用 Loan RPC 获取审批记录
	rpcResp, err := l.svcCtx.LoanRpc.ListLoanApprovals(l.ctx, &loanclient.ListLoanApprovalsReq{
		ApplicationId: req.ApplicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.ListLoanApprovalsResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.ListLoanApprovalsResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	// 转换审批记录列表
	if len(rpcResp.List) > 0 {
		resp.List = make([]types.LoanApprovalInfo, len(rpcResp.List))
		for i, item := range rpcResp.List {
			resp.List[i] = types.LoanApprovalInfo{
				Id:               item.Id,
				ApplicationId:    item.ApplicationId,
				AuditorId:        item.AuditorId,
				AuditorName:      item.AuditorName,
				Action:           item.Action,
				Suggestions:      item.Suggestions,
				ApprovedAmount:   item.ApprovedAmount,
				ApprovedDuration: item.ApprovedDuration,
				InterestRate:     item.InterestRate,
				CreatedAt:        item.CreatedAt,
			}
		}
	} else {
		resp.List = make([]types.LoanApprovalInfo, 0)
	}

	return resp, nil
}
