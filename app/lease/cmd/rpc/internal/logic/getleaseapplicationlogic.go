package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseApplicationLogic {
	return &GetLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLeaseApplicationLogic) GetLeaseApplication(in *lease.GetLeaseApplicationReq) (*lease.GetLeaseApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &lease.GetLeaseApplicationResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 根据申请编号查询申请信息
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &lease.GetLeaseApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 转换为响应格式
	applicationInfo := &lease.LeaseApplicationInfo{
		Id:              int64(application.Id),
		ApplicationId:   application.ApplicationId,
		UserId:          int64(application.UserId),
		ApplicantName:   application.ApplicantName,
		ProductId:       int64(application.ProductId),
		ProductCode:     application.ProductCode,
		Name:            application.Name,
		Type:            application.Type,
		Machinery:       application.Machinery,
		StartDate:       application.StartDate.Format("2006-01-02"),
		EndDate:         application.EndDate.Format("2006-01-02"),
		Duration:        int32(application.Duration),
		DailyRate:       application.DailyRate,
		TotalAmount:     application.TotalAmount,
		Deposit:         application.Deposit,
		DeliveryAddress: application.DeliveryAddress,
		ContactPhone:    application.ContactPhone,
		Purpose:         application.Purpose.String,
		Status:          application.Status,
		CreatedAt:       application.CreatedAt.Unix(),
		UpdatedAt:       application.UpdatedAt.Unix(),
	}

	return &lease.GetLeaseApplicationResp{
		Code:            200,
		Message:         "查询成功",
		ApplicationInfo: applicationInfo,
	}, nil
}
