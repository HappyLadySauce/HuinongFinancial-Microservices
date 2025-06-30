package logic

import (
	"context"
	"database/sql"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLoanApplicationLogic {
	return &CancelLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLoanApplicationLogic) CancelLoanApplication(in *loan.CancelLoanApplicationReq) (*loan.CancelLoanApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &loan.CancelLoanApplicationResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 查询申请信息
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &loan.CancelLoanApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 检查申请状态是否可以撤销
	if application.Status != "pending" {
		return &loan.CancelLoanApplicationResp{
			Code:    400,
			Message: "只有待审批状态的申请才可以撤销",
		}, nil
	}

	// 更新申请状态为已撤销
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
		Status:        "cancelled",
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     time.Now(),
	}

	err = l.svcCtx.LoanApplicationsModel.Update(l.ctx, updatedApplication)
	if err != nil {
		l.Errorf("撤销申请失败: %v", err)
		return &loan.CancelLoanApplicationResp{
			Code:    500,
			Message: "撤销申请失败",
		}, nil
	}

	// 记录撤销原因（可以创建一个审批记录）
	if in.Reason != "" {
		approval := &model.LoanApprovals{
			ApplicationId:    application.Id,
			AuditorId:        uint64(application.UserId), // 用户自己撤销
			AuditorName:      application.ApplicantName,
			Action:           "cancel",
			Suggestions:      sql.NullString{String: "用户撤销: " + in.Reason, Valid: true},
			ApprovedAmount:   sql.NullFloat64{Valid: false},
			ApprovedDuration: sql.NullInt64{Valid: false},
			InterestRate:     sql.NullFloat64{Valid: false},
			CreatedAt:        time.Now(),
		}

		_, err = l.svcCtx.LoanApprovalsModel.Insert(l.ctx, approval)
		if err != nil {
			l.Errorf("创建撤销记录失败: %v", err)
			// 不影响主流程，只记录日志
		}
	}

	return &loan.CancelLoanApplicationResp{
		Code:    200,
		Message: "申请已撤销",
	}, nil
}
