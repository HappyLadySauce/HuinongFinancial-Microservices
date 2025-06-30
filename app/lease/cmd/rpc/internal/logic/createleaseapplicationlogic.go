package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"model"
	"rpc/internal/clients"
	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

type CreateLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseApplicationLogic {
	return &CreateLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租赁申请管理
func (l *CreateLeaseApplicationLogic) CreateLeaseApplication(in *lease.CreateLeaseApplicationReq) (*lease.CreateLeaseApplicationResp, error) {
	// 参数验证
	if err := l.validateCreateRequest(in); err != nil {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 1. 调用AppUser RPC验证用户信息并获取用户姓名
	userResp, err := l.svcCtx.AppUserClient.GetUserById(l.ctx, &clients.GetUserByIdReq{
		UserId: in.UserId,
	})
	if err != nil {
		l.Errorf("调用AppUser服务失败: %v", err)
		return &lease.CreateLeaseApplicationResp{
			Code:    500,
			Message: "用户信息验证失败，请稍后重试",
		}, nil
	}

	if userResp.Code != 200 {
		l.Errorf("用户信息验证失败: %s", userResp.Message)
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: userResp.Message,
		}, nil
	}

	if userResp.UserInfo == nil {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: "用户信息不存在",
		}, nil
	}

	applicantName := userResp.UserInfo.Name

	// 2. 调用LeaseProduct RPC验证产品信息和库存
	stockResp, err := l.svcCtx.LeaseProductClient.CheckInventoryAvailability(l.ctx, &clients.CheckInventoryAvailabilityReq{
		ProductCode: in.ProductCode,
		Quantity:    1,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if err != nil {
		l.Errorf("调用LeaseProduct服务失败: %v", err)
		return &lease.CreateLeaseApplicationResp{
			Code:    500,
			Message: "产品库存检查失败，请稍后重试",
		}, nil
	}

	if stockResp.Code != 200 {
		l.Errorf("产品库存检查失败: %s", stockResp.Message)
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: stockResp.Message,
		}, nil
	}

	if !stockResp.Available {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: "产品库存不足或时间段不可用",
		}, nil
	}

	// 3. 验证产品是否有效
	productResp, err := l.svcCtx.LeaseProductClient.GetLeaseProduct(l.ctx, &clients.GetLeaseProductReq{
		ProductCode: in.ProductCode,
	})
	if err != nil {
		l.Errorf("调用LeaseProduct服务获取产品信息失败: %v", err)
		return &lease.CreateLeaseApplicationResp{
			Code:    500,
			Message: "产品信息获取失败，请稍后重试",
		}, nil
	}

	if productResp.Code != 200 || productResp.Data == nil {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: "产品不存在或已下架",
		}, nil
	}

	// 解析开始和结束日期
	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: "开始日期格式错误",
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", in.EndDate)
	if err != nil {
		return &lease.CreateLeaseApplicationResp{
			Code:    400,
			Message: "结束日期格式错误",
		}, nil
	}

	// 生成申请编号
	applicationId := l.generateApplicationId()

	// 创建申请记录
	now := time.Now()
	application := &model.LeaseApplications{
		ApplicationId:   applicationId,
		UserId:          uint64(in.UserId),
		ApplicantName:   applicantName, // 从AppUser服务获取的真实用户姓名
		ProductId:       uint64(in.ProductId),
		ProductCode:     in.ProductCode,
		Name:            in.Name,
		Type:            in.Type,
		Machinery:       in.Machinery,
		StartDate:       startDate,
		EndDate:         endDate,
		Duration:        uint64(in.Duration),
		DailyRate:       in.DailyRate,
		TotalAmount:     in.TotalAmount,
		Deposit:         in.Deposit,
		DeliveryAddress: in.DeliveryAddress,
		ContactPhone:    in.ContactPhone,
		Purpose:         sql.NullString{String: in.Purpose, Valid: in.Purpose != ""},
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = l.svcCtx.LeaseApplicationsModel.Insert(l.ctx, application)
	if err != nil {
		l.Errorf("创建租赁申请失败: %v", err)
		return &lease.CreateLeaseApplicationResp{
			Code:    500,
			Message: "创建申请失败",
		}, nil
	}

	l.Infof("租赁申请创建成功 - 申请编号: %s, 用户: %s (ID: %d), 产品: %s, 时间: %s到%s",
		applicationId, applicantName, in.UserId, in.ProductCode, in.StartDate, in.EndDate)

	return &lease.CreateLeaseApplicationResp{
		Code:          200,
		Message:       "申请创建成功",
		ApplicationId: applicationId,
	}, nil
}

// validateCreateRequest 验证创建请求参数
func (l *CreateLeaseApplicationLogic) validateCreateRequest(in *lease.CreateLeaseApplicationReq) error {
	if in.UserId <= 0 {
		return fmt.Errorf("用户ID不能为空")
	}
	if in.ProductId <= 0 {
		return fmt.Errorf("产品ID不能为空")
	}
	if in.ProductCode == "" {
		return fmt.Errorf("产品编码不能为空")
	}
	if in.Name == "" {
		return fmt.Errorf("申请名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("租赁类型不能为空")
	}
	if in.StartDate == "" {
		return fmt.Errorf("开始日期不能为空")
	}
	if in.EndDate == "" {
		return fmt.Errorf("结束日期不能为空")
	}
	if in.Duration <= 0 {
		return fmt.Errorf("租期必须大于0")
	}
	if in.DailyRate <= 0 {
		return fmt.Errorf("日租金必须大于0")
	}
	if in.TotalAmount <= 0 {
		return fmt.Errorf("总金额必须大于0")
	}
	if in.ContactPhone == "" {
		return fmt.Errorf("联系电话不能为空")
	}
	if in.DeliveryAddress == "" {
		return fmt.Errorf("交付地址不能为空")
	}

	// 验证日期格式和合理性
	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return fmt.Errorf("开始日期格式错误")
	}

	endDate, err := time.Parse("2006-01-02", in.EndDate)
	if err != nil {
		return fmt.Errorf("结束日期格式错误")
	}

	if startDate.After(endDate) {
		return fmt.Errorf("开始日期不能晚于结束日期")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("开始日期不能早于今天")
	}

	return nil
}

// generateApplicationId 生成申请编号
func (l *CreateLeaseApplicationLogic) generateApplicationId() string {
	// 格式：LA + 年月日 + 6位随机数
	now := time.Now()
	dateStr := now.Format("20060102")
	randomStr := stringx.Randn(6)
	return fmt.Sprintf("LA%s%s", dateStr, randomStr)
}
