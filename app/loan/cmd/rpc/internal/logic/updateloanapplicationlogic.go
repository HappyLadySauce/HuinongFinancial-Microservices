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

type UpdateLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanApplicationLogic {
	return &UpdateLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanApplicationLogic) UpdateLoanApplication(in *loan.UpdateLoanApplicationReq) (*loan.UpdateLoanApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &loan.UpdateLoanApplicationResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 查询申请信息
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &loan.UpdateLoanApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 检查申请状态是否可以修改
	if application.Status != "pending" {
		return &loan.UpdateLoanApplicationResp{
			Code:    400,
			Message: "只有待审批状态的申请才可以修改",
		}, nil
	}

	// 只允许更新部分字段
	updatedApplication := &model.LoanApplications{
		Id:            application.Id,
		ApplicationId: application.ApplicationId,
		UserId:        application.UserId,
		ApplicantName: application.ApplicantName,
		ProductId:     application.ProductId,
		Name:          application.Name,
		Type:          application.Type,
		Amount:        in.Amount,   // 可更新
		Duration:      uint64(in.Duration), // 可更新
		Purpose:       sql.NullString{String: in.Purpose, Valid: in.Purpose != ""}, // 可更新
		Status:        application.Status,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     time.Now(),
	}

	// 如果字段为空或为0，保持原值
	if in.Amount <= 0 {
		updatedApplication.Amount = application.Amount
	}
	if in.Duration <= 0 {
		updatedApplication.Duration = application.Duration
	}
	if in.Purpose == "" {
		updatedApplication.Purpose = application.Purpose
	}

	err = l.svcCtx.LoanApplicationsModel.Update(l.ctx, updatedApplication)
	if err != nil {
		l.Errorf("更新申请失败: %v", err)
		return &loan.UpdateLoanApplicationResp{
			Code:    500,
			Message: "更新申请失败",
		}, nil
	}

	// 查询更新后的申请信息
	updatedApp, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询更新后的申请失败: %v", err)
		return &loan.UpdateLoanApplicationResp{
			Code:    500,
			Message: "更新成功但查询失败",
		}, nil
	}

	// 转换为响应格式
	applicationInfo := &loan.LoanApplicationInfo{
		Id:            int64(updatedApp.Id),
		ApplicationId: updatedApp.ApplicationId,
		UserId:        int64(updatedApp.UserId),
		ApplicantName: updatedApp.ApplicantName,
		ProductId:     int64(updatedApp.ProductId),
		Name:          updatedApp.Name,
		Type:          updatedApp.Type,
		Amount:        updatedApp.Amount,
		Duration:      int32(updatedApp.Duration),
		Purpose:       updatedApp.Purpose.String,
		Status:        updatedApp.Status,
		CreatedAt:     updatedApp.CreatedAt.Unix(),
		UpdatedAt:     updatedApp.UpdatedAt.Unix(),
	}

	return &loan.UpdateLoanApplicationResp{
		Code:            200,
		Message:         "更新成功",
		ApplicationInfo: applicationInfo,
	}, nil
}
