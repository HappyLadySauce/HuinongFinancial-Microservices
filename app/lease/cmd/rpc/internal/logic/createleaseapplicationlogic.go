package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"appuserrpc/appuserclient"
	"leaseproductrpc/leaseproductservice"
	"model"
	"rpc/internal/breaker"
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
		return nil, err
	}

	// 1. 使用熔断器调用AppUser RPC验证用户信息并获取用户姓名
	userResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuserclient.GetUserInfoResp, error) {
		return l.svcCtx.AppUserClient.GetUserById(l.ctx, &appuserclient.GetUserByIdReq{
			UserId: in.UserId,
		})
	}, breaker.IsAcceptableError)

	if err != nil {
		l.Errorf("调用AppUser服务失败: %v", err)
		return nil, fmt.Errorf("用户信息验证失败，请稍后重试")
	}

	if userResp.UserInfo == nil {
		return nil, fmt.Errorf("用户信息不存在")
	}

	applicantName := userResp.UserInfo.Name

	// 2. 使用熔断器调用LeaseProduct RPC检查库存和验证产品信息
	stockResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "leaseproduct-rpc", func() (*leaseproductservice.CheckInventoryAvailabilityResp, error) {
		return l.svcCtx.LeaseProductClient.CheckInventoryAvailability(l.ctx, &leaseproductservice.CheckInventoryAvailabilityReq{
			ProductCode: in.ProductCode,
			Quantity:    1, // 通常租赁1台设备
			StartDate:   in.StartDate,
			EndDate:     in.EndDate,
		})
	}, breaker.IsAcceptableError)

	if err != nil {
		l.Errorf("调用LeaseProduct服务失败: %v", err)
		return nil, fmt.Errorf("产品库存检查失败，请稍后重试")
	}

	if !stockResp.Available {
		return nil, fmt.Errorf("产品库存不足或时间段不可用")
	}

	// 3. 生成申请ID
	applicationId := l.generateApplicationId()

	// 4. 创建租赁申请记录
	startDate, _ := time.Parse("2006-01-02", in.StartDate)
	endDate, _ := time.Parse("2006-01-02", in.EndDate)
	
	application := &model.LeaseApplications{
		ApplicationId:   applicationId,
		UserId:          uint64(in.UserId),
		ApplicantName:   applicantName,
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
		Status:          "pending", // 待审核
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = l.svcCtx.LeaseApplicationsModel.Insert(l.ctx, application)
	if err != nil {
		l.Errorf("创建租赁申请失败: %v", err)
		return nil, fmt.Errorf("创建申请失败，请稍后重试")
	}

	return &lease.CreateLeaseApplicationResp{
		ApplicationId: applicationId,
	}, nil
}

// 参数验证
func (l *CreateLeaseApplicationLogic) validateCreateRequest(in *lease.CreateLeaseApplicationReq) error {
	if in.UserId <= 0 {
		return fmt.Errorf("用户ID无效")
	}
	if in.ProductId <= 0 {
		return fmt.Errorf("产品ID无效")
	}
	if in.ProductCode == "" {
		return fmt.Errorf("产品编码不能为空")
	}
	if in.Name == "" {
		return fmt.Errorf("租赁名称不能为空")
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
		return fmt.Errorf("租赁期限必须大于0")
	}
	if in.TotalAmount <= 0 {
		return fmt.Errorf("租赁总额必须大于0")
	}
	if in.ContactPhone == "" {
		return fmt.Errorf("联系电话不能为空")
	}
	return nil
}

// 生成申请ID
func (l *CreateLeaseApplicationLogic) generateApplicationId() string {
	// 生成格式：LEASE + 年月日 + 6位随机数
	now := time.Now()
	dateStr := now.Format("20060102")
	randomStr := stringx.Randn(6)
	return fmt.Sprintf("LEASE%s%s", dateStr, randomStr)
}
