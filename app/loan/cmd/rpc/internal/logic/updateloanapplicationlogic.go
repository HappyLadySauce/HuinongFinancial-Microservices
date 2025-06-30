package logic

import (
	"context"
	"fmt"
	"time"

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
		return nil, fmt.Errorf("申请编号不能为空")
	}

	// 查询申请是否存在
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return nil, fmt.Errorf("申请不存在")
	}

	// 检查申请状态是否可以修改
	if application.Status != "pending" {
		return nil, fmt.Errorf("只有待审核状态的申请才能修改")
	}

	// 验证更新参数
	if in.Amount <= 0 {
		return nil, fmt.Errorf("申请金额必须大于0")
	}
	if in.Duration <= 0 {
		return nil, fmt.Errorf("贷款期限必须大于0")
	}

	// 更新申请信息
	application.Amount = in.Amount
	application.Duration = uint64(in.Duration)
	application.Purpose.String = in.Purpose
	application.Purpose.Valid = in.Purpose != ""
	application.UpdatedAt = time.Now()

	err = l.svcCtx.LoanApplicationsModel.Update(l.ctx, application)
	if err != nil {
		l.Errorf("更新申请失败: %v", err)
		return nil, fmt.Errorf("更新申请失败")
	}

	// 查询更新后的申请信息
	updatedApplication, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询更新后的申请失败: %v", err)
		return nil, fmt.Errorf("更新成功但查询失败")
	}

	// 构造响应
	return &loan.UpdateLoanApplicationResp{
		ApplicationInfo: &loan.LoanApplicationInfo{
			Id:            int64(updatedApplication.Id),
			ApplicationId: updatedApplication.ApplicationId,
			UserId:        int64(updatedApplication.UserId),
			ApplicantName: updatedApplication.ApplicantName,
			ProductId:     int64(updatedApplication.ProductId),
			Name:          updatedApplication.Name,
			Type:          updatedApplication.Type,
			Amount:        updatedApplication.Amount,
			Duration:      int32(updatedApplication.Duration),
			Purpose:       updatedApplication.Purpose.String,
			Status:        updatedApplication.Status,
			CreatedAt:     updatedApplication.CreatedAt.Unix(),
			UpdatedAt:     updatedApplication.UpdatedAt.Unix(),
		},
	}, nil
}
