package logic

import (
	"context"
	"fmt"
	"time"

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
		return nil, fmt.Errorf("申请编号不能为空")
	}

	// 查询申请信息
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return nil, fmt.Errorf("申请不存在")
	}

	// 检查申请状态是否可以修改
	if application.Status != "pending" {
		return nil, fmt.Errorf("只有待审批状态的申请才可以修改")
	}

	// 更新允许修改的字段
	now := time.Now()
	if in.Purpose != "" {
		application.Purpose.String = in.Purpose
		application.Purpose.Valid = true
	}
	if in.DeliveryAddress != "" {
		application.DeliveryAddress = in.DeliveryAddress
	}
	if in.ContactPhone != "" {
		application.ContactPhone = in.ContactPhone
	}
	application.UpdatedAt = now

	// 执行更新
	err = l.svcCtx.LeaseApplicationsModel.Update(l.ctx, application)
	if err != nil {
		l.Errorf("更新申请失败: %v", err)
		return nil, fmt.Errorf("更新申请失败")
	}

	// 重新查询更新后的申请信息
	updatedApplication, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询更新后的申请失败: %v", err)
		return nil, fmt.Errorf("更新成功但查询失败")
	}

	// 转换为响应格式
	return &lease.UpdateLeaseApplicationResp{
		ApplicationInfo: &lease.LeaseApplicationInfo{
			Id:              int64(updatedApplication.Id),
			ApplicationId:   updatedApplication.ApplicationId,
			UserId:          int64(updatedApplication.UserId),
			ApplicantName:   updatedApplication.ApplicantName,
			ProductId:       int64(updatedApplication.ProductId),
			ProductCode:     updatedApplication.ProductCode,
			Name:            updatedApplication.Name,
			Type:            updatedApplication.Type,
			Machinery:       updatedApplication.Machinery,
			StartDate:       updatedApplication.StartDate.Format("2006-01-02"),
			EndDate:         updatedApplication.EndDate.Format("2006-01-02"),
			Duration:        int32(updatedApplication.Duration),
			DailyRate:       updatedApplication.DailyRate,
			TotalAmount:     updatedApplication.TotalAmount,
			Deposit:         updatedApplication.Deposit,
			DeliveryAddress: updatedApplication.DeliveryAddress,
			ContactPhone:    updatedApplication.ContactPhone,
			Purpose:         updatedApplication.Purpose.String,
			Status:          updatedApplication.Status,
			CreatedAt:       updatedApplication.CreatedAt.Unix(),
			UpdatedAt:       updatedApplication.UpdatedAt.Unix(),
		},
	}, nil
}
