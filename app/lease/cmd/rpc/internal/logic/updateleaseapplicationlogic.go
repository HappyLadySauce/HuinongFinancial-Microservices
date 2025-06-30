package logic

import (
	"context"
	"database/sql"
	"time"

	"model"
	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaseApplicationLogic {
	return &UpdateLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLeaseApplicationLogic) UpdateLeaseApplication(in *lease.UpdateLeaseApplicationReq) (*lease.UpdateLeaseApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &lease.UpdateLeaseApplicationResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 查询申请信息
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &lease.UpdateLeaseApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 检查申请状态是否可以修改
	if application.Status != "pending" {
		return &lease.UpdateLeaseApplicationResp{
			Code:    400,
			Message: "只有待审批状态的申请才可以修改",
		}, nil
	}

	// 只允许更新部分字段
	updatedApplication := &model.LeaseApplications{
		Id:              application.Id,
		ApplicationId:   application.ApplicationId,
		UserId:          application.UserId,
		ApplicantName:   application.ApplicantName,
		ProductId:       application.ProductId,
		ProductCode:     application.ProductCode,
		Name:            application.Name,
		Type:            application.Type,
		Machinery:       application.Machinery,
		StartDate:       application.StartDate,
		EndDate:         application.EndDate,
		Duration:        application.Duration,
		DailyRate:       application.DailyRate,
		TotalAmount:     application.TotalAmount,
		Deposit:         application.Deposit,
		DeliveryAddress: in.DeliveryAddress,                                          // 可更新
		ContactPhone:    in.ContactPhone,                                             // 可更新
		Purpose:         sql.NullString{String: in.Purpose, Valid: in.Purpose != ""}, // 可更新
		Status:          application.Status,
		CreatedAt:       application.CreatedAt,
		UpdatedAt:       time.Now(),
	}

	// 如果字段为空，保持原值
	if in.DeliveryAddress == "" {
		updatedApplication.DeliveryAddress = application.DeliveryAddress
	}
	if in.ContactPhone == "" {
		updatedApplication.ContactPhone = application.ContactPhone
	}
	if in.Purpose == "" {
		updatedApplication.Purpose = application.Purpose
	}

	err = l.svcCtx.LeaseApplicationsModel.Update(l.ctx, updatedApplication)
	if err != nil {
		l.Errorf("更新申请失败: %v", err)
		return &lease.UpdateLeaseApplicationResp{
			Code:    500,
			Message: "更新申请失败",
		}, nil
	}

	// 查询更新后的申请信息
	updatedApp, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询更新后的申请失败: %v", err)
		return &lease.UpdateLeaseApplicationResp{
			Code:    500,
			Message: "更新成功但查询失败",
		}, nil
	}

	// 转换为响应格式
	applicationInfo := &lease.LeaseApplicationInfo{
		Id:              int64(updatedApp.Id),
		ApplicationId:   updatedApp.ApplicationId,
		UserId:          int64(updatedApp.UserId),
		ApplicantName:   updatedApp.ApplicantName,
		ProductId:       int64(updatedApp.ProductId),
		ProductCode:     updatedApp.ProductCode,
		Name:            updatedApp.Name,
		Type:            updatedApp.Type,
		Machinery:       updatedApp.Machinery,
		StartDate:       updatedApp.StartDate.Format("2006-01-02"),
		EndDate:         updatedApp.EndDate.Format("2006-01-02"),
		Duration:        int32(updatedApp.Duration),
		DailyRate:       updatedApp.DailyRate,
		TotalAmount:     updatedApp.TotalAmount,
		Deposit:         updatedApp.Deposit,
		DeliveryAddress: updatedApp.DeliveryAddress,
		ContactPhone:    updatedApp.ContactPhone,
		Purpose:         updatedApp.Purpose.String,
		Status:          updatedApp.Status,
		CreatedAt:       updatedApp.CreatedAt.Unix(),
		UpdatedAt:       updatedApp.UpdatedAt.Unix(),
	}

	return &lease.UpdateLeaseApplicationResp{
		Code:            200,
		Message:         "更新成功",
		ApplicationInfo: applicationInfo,
	}, nil
}
