package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLoanApplicationLogic {
	return &ApproveLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 贷款审批管理
func (l *ApproveLoanApplicationLogic) ApproveLoanApplication(in *loan.ApproveLoanApplicationReq) (*loan.ApproveLoanApplicationResp, error) {
	// 参数验证
	if err := l.validateApproveRequest(in); err != nil {
		return &loan.ApproveLoanApplicationResp{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 查询申请信息
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &loan.ApproveLoanApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 检查申请状态是否可以审批
	if application.Status != "pending" {
		return &loan.ApproveLoanApplicationResp{
			Code:    400,
			Message: "申请状态不允许审批",
		}, nil
	}

	// 开始事务处理
	// 1. 更新申请状态
	now := time.Now()
	newStatus := "approved"
	if in.Action == "reject" {
		newStatus = "rejected"
	}

	updatedApplication := &model.LoanApplications{
		Id:            application.Id,
		ApplicationId: application.ApplicationId,
		UserId:        application.UserId,
		ApplicantName: application.ApplicantName,
		ProductId:     application.ProductId,
		Name:          application.Name,
		Type:          application.Type,
		Amount:        application.Amount,
		Duration:      application.Duration,
		Purpose:       application.Purpose,
		Status:        newStatus,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     now,
	}

	err = l.svcCtx.LoanApplicationsModel.Update(l.ctx, updatedApplication)
	if err != nil {
		l.Errorf("更新申请状态失败: %v", err)
		return &loan.ApproveLoanApplicationResp{
			Code:    500,
			Message: "审批失败",
		}, nil
	}

	// 2. 创建审批记录
	approval := &model.LoanApprovals{
		ApplicationId:    application.Id,
		AuditorId:        uint64(in.AuditorId),
		AuditorName:      in.AuditorName,
		Action:           in.Action,
		Suggestions:      sql.NullString{String: in.Suggestions, Valid: in.Suggestions != ""},
		ApprovedAmount:   sql.NullFloat64{Float64: in.ApprovedAmount, Valid: in.ApprovedAmount > 0},
		ApprovedDuration: sql.NullInt64{Int64: int64(in.ApprovedDuration), Valid: in.ApprovedDuration > 0},
		InterestRate:     sql.NullFloat64{Float64: in.InterestRate, Valid: in.InterestRate > 0},
		CreatedAt:        now,
	}

	_, err = l.svcCtx.LoanApprovalsModel.Insert(l.ctx, approval)
	if err != nil {
		l.Errorf("创建审批记录失败: %v", err)
		return &loan.ApproveLoanApplicationResp{
			Code:    500,
			Message: "审批记录创建失败",
		}, nil
	}

	// TODO: 如果是批准，可能需要调用其他服务执行后续操作
	// 比如发送通知等

	message := "审批成功"
	if in.Action == "approve" {
		message = "申请已批准"
	} else {
		message = "申请已拒绝"
	}

	return &loan.ApproveLoanApplicationResp{
		Code:    200,
		Message: message,
	}, nil
}

// validateApproveRequest 验证审批请求参数
func (l *ApproveLoanApplicationLogic) validateApproveRequest(in *loan.ApproveLoanApplicationReq) error {
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
		if in.ApprovedAmount <= 0 {
			return fmt.Errorf("批准金额必须大于0")
		}
		if in.ApprovedDuration <= 0 {
			return fmt.Errorf("批准期限必须大于0")
		}
		if in.InterestRate < 0 {
			return fmt.Errorf("利率不能小于0")
		}
	}
	return nil
}
