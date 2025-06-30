package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApprovalsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLeaseApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApprovalsLogic {
	return &ListLeaseApprovalsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLeaseApprovalsLogic) ListLeaseApprovals(in *lease.ListLeaseApprovalsReq) (*lease.ListLeaseApprovalsResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &lease.ListLeaseApprovalsResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 先查询申请是否存在
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &lease.ListLeaseApprovalsResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 查询审批记录
	approvals, err := l.svcCtx.LeaseApprovalsModel.FindByApplicationId(l.ctx, int64(application.Id))
	if err != nil {
		l.Errorf("查询审批记录失败: %v", err)
		return &lease.ListLeaseApprovalsResp{
			Code:    500,
			Message: "查询审批记录失败",
		}, nil
	}

	// 转换为响应格式
	var approvalList []*lease.LeaseApprovalInfo
	for _, approval := range approvals {
		approvalInfo := &lease.LeaseApprovalInfo{
			Id:               int64(approval.Id),
			ApplicationId:    int64(approval.ApplicationId),
			AuditorId:        int64(approval.AuditorId),
			AuditorName:      approval.AuditorName,
			Action:           approval.Action,
			Suggestions:      approval.Suggestions.String,
			ApprovedDuration: int32(approval.ApprovedDuration.Int64),
			ApprovedAmount:   approval.ApprovedAmount.Float64,
			ApprovedDeposit:  approval.ApprovedDeposit.Float64,
			CreatedAt:        approval.CreatedAt.Unix(),
		}
		approvalList = append(approvalList, approvalInfo)
	}

	// 如果没有数据，返回空列表
	if approvalList == nil {
		approvalList = make([]*lease.LeaseApprovalInfo, 0)
	}

	return &lease.ListLeaseApprovalsResp{
		Code:    200,
		Message: "查询成功",
		List:    approvalList,
	}, nil
}
