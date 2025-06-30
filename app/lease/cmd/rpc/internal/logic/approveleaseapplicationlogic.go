package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLeaseApplicationLogic {
	return &ApproveLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租赁审批管理
func (l *ApproveLeaseApplicationLogic) ApproveLeaseApplication(in *lease.ApproveLeaseApplicationReq) (*lease.ApproveLeaseApplicationResp, error) {
	// 参数验证
	if err := l.validateApproveRequest(in); err != nil {
		return nil, err
	}

	// 查询申请信息
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return nil, fmt.Errorf("申请不存在")
	}

	// 检查申请状态是否可以审批
	if application.Status != "pending" {
		return nil, fmt.Errorf("申请状态不允许审批")
	}

	// 开始事务处理
	// 1. 更新申请状态
	now := time.Now()
	newStatus := "approved"
	if in.Action == "reject" {
		newStatus = "rejected"
	}

	application.Status = newStatus
	application.UpdatedAt = now

	err = l.svcCtx.LeaseApplicationsModel.Update(l.ctx, application)
	if err != nil {
		l.Errorf("更新申请状态失败: %v", err)
		return nil, fmt.Errorf("审批失败")
	}

	// 2. 创建审批记录
	approval := &model.LeaseApprovals{
		ApplicationId:    application.Id,
		AuditorId:        uint64(in.AuditorId),
		AuditorName:      in.AuditorName,
		Action:           in.Action,
		Suggestions:      sql.NullString{String: in.Suggestions, Valid: in.Suggestions != ""},
		ApprovedDuration: sql.NullInt64{Int64: int64(in.ApprovedDuration), Valid: in.Action == "approve"},
		ApprovedAmount:   sql.NullFloat64{Float64: in.ApprovedAmount, Valid: in.Action == "approve"},
		ApprovedDeposit:  sql.NullFloat64{Float64: in.ApprovedDeposit, Valid: in.Action == "approve"},
		CreatedAt:        now,
	}

	_, err = l.svcCtx.LeaseApprovalsModel.Insert(l.ctx, approval)
	if err != nil {
		l.Errorf("创建审批记录失败: %v", err)
		return nil, fmt.Errorf("审批记录创建失败")
	}

	// TODO: 如果是批准，可能需要调用其他服务执行后续操作
	// 比如更新库存、发送通知等

	return &lease.ApproveLeaseApplicationResp{}, nil
}

// validateApproveRequest 验证审批请求参数
func (l *ApproveLeaseApplicationLogic) validateApproveRequest(in *lease.ApproveLeaseApplicationReq) error {
	if in.ApplicationId == "" {
		return fmt.Errorf("申请编号不能为空")
	}
	if in.AuditorId <= 0 {
		return fmt.Errorf("审核员ID不能为空")
	}
	if in.AuditorName == "" {
		return fmt.Errorf("审核员姓名不能为空")
	}
	if in.Action != "approve" && in.Action != "reject" {
		return fmt.Errorf("审批动作必须为approve或reject")
	}
	if in.Action == "approve" {
		if in.ApprovedDuration <= 0 {
			return fmt.Errorf("批准租期必须大于0")
		}
		if in.ApprovedAmount <= 0 {
			return fmt.Errorf("批准金额必须大于0")
		}
	}
	return nil
}
