package logic

import (
	"context"
	"fmt"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApprovalsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApprovalsLogic {
	return &ListLoanApprovalsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanApprovalsLogic) ListLoanApprovals(in *loan.ListLoanApprovalsReq) (*loan.ListLoanApprovalsResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return nil, fmt.Errorf("申请编号不能为空")
	}

	// 先查询申请是否存在
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return nil, fmt.Errorf("申请不存在")
	}

	// 查询审批记录
	approvals, err := l.svcCtx.LoanApprovalsModel.FindByApplicationId(l.ctx, int64(application.Id))
	if err != nil {
		l.Errorf("查询审批记录失败: %v", err)
		return nil, fmt.Errorf("查询审批记录失败")
	}

	// 转换为响应格式
	var approvalList []*loan.LoanApprovalInfo
	for _, approval := range approvals {
		approvalInfo := &loan.LoanApprovalInfo{
			Id:               int64(approval.Id),
			ApplicationId:    int64(approval.ApplicationId),
			AuditorId:        int64(approval.AuditorId),
			AuditorName:      approval.AuditorName,
			Action:           approval.Action,
			Suggestions:      approval.Suggestions.String,
			ApprovedAmount:   approval.ApprovedAmount.Float64,
			ApprovedDuration: int32(approval.ApprovedDuration.Int64),
			InterestRate:     approval.InterestRate.Float64,
			CreatedAt:        approval.CreatedAt.Unix(),
		}
		approvalList = append(approvalList, approvalInfo)
	}

	// 如果没有数据，返回空列表
	if approvalList == nil {
		approvalList = make([]*loan.LoanApprovalInfo, 0)
	}

	return &loan.ListLoanApprovalsResp{
		List: approvalList,
	}, nil
}
