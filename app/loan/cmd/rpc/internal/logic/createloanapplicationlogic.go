package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"appuserrpc/appuserclient"
	"loanproductrpc/loanproductservice"
	"model"
	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

type CreateLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanApplicationLogic {
	return &CreateLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 贷款申请管理
func (l *CreateLoanApplicationLogic) CreateLoanApplication(in *loan.CreateLoanApplicationReq) (*loan.CreateLoanApplicationResp, error) {
	// 参数验证
	if err := l.validateCreateRequest(in); err != nil {
		return nil, err
	}

	// 1. 调用AppUser RPC验证用户信息并获取用户姓名
	userResp, err := l.svcCtx.AppUserClient.GetUserById(l.ctx, &appuserclient.GetUserByIdReq{
		UserId: in.UserId,
	})
	if err != nil {
		l.Errorf("调用AppUser服务失败: %v", err)
		return nil, fmt.Errorf("用户信息验证失败，请稍后重试")
	}

	if userResp.UserInfo == nil {
		return nil, fmt.Errorf("用户信息不存在")
	}

	applicantName := userResp.UserInfo.Name

	// 2. 调用LoanProduct RPC验证产品信息
	productResp, err := l.svcCtx.LoanProductClient.GetLoanProduct(l.ctx, &loanproductservice.GetLoanProductReq{
		Id: in.ProductId,
	})
	if err != nil {
		l.Errorf("调用LoanProduct服务失败: %v", err)
		return nil, fmt.Errorf("产品信息验证失败，请稍后重试")
	}

	if productResp.Data == nil {
		return nil, fmt.Errorf("产品不存在")
	}

	product := productResp.Data

	// 3. 验证申请金额是否在产品限额内
	if in.Amount < product.MinAmount || in.Amount > product.MaxAmount {
		return nil, fmt.Errorf("申请金额应在%.2f到%.2f之间", product.MinAmount, product.MaxAmount)
	}

	// 4. 验证申请期限是否在产品范围内
	if int32(in.Duration) < product.MinDuration || int32(in.Duration) > product.MaxDuration {
		return nil, fmt.Errorf("申请期限应在%d到%d个月之间", product.MinDuration, product.MaxDuration)
	}

	// 5. 验证产品状态
	if product.Status != 1 {
		return nil, fmt.Errorf("产品已下架，无法申请")
	}

	// 生成申请编号
	applicationId := l.generateApplicationId()

	// 创建申请记录
	now := time.Now()
	application := &model.LoanApplications{
		ApplicationId: applicationId,
		UserId:        uint64(in.UserId),
		ApplicantName: applicantName, // 从AppUser服务获取的真实用户姓名
		ProductId:     uint64(in.ProductId),
		Name:          in.Name,
		Type:          in.Type,
		Amount:        in.Amount,
		Duration:      uint64(in.Duration),
		Purpose:       sql.NullString{String: in.Purpose, Valid: in.Purpose != ""},
		Status:        "pending",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	_, err = l.svcCtx.LoanApplicationsModel.Insert(l.ctx, application)
	if err != nil {
		l.Errorf("创建贷款申请失败: %v", err)
		return nil, fmt.Errorf("创建申请失败")
	}

	l.Infof("贷款申请创建成功 - 申请编号: %s, 用户: %s (ID: %d), 产品: %s (ID: %d), 金额: %.2f, 期限: %d个月",
		applicationId, applicantName, in.UserId, product.Name, in.ProductId, in.Amount, in.Duration)

	return &loan.CreateLoanApplicationResp{
		ApplicationId: applicationId,
	}, nil
}

// validateCreateRequest 验证创建请求参数
func (l *CreateLoanApplicationLogic) validateCreateRequest(in *loan.CreateLoanApplicationReq) error {
	if in.UserId <= 0 {
		return fmt.Errorf("用户ID不能为空")
	}
	if in.ProductId <= 0 {
		return fmt.Errorf("产品ID不能为空")
	}
	if in.Name == "" {
		return fmt.Errorf("申请名称不能为空")
	}
	if in.Type == "" {
		return fmt.Errorf("贷款类型不能为空")
	}
	if in.Amount <= 0 {
		return fmt.Errorf("申请金额必须大于0")
	}
	if in.Duration <= 0 {
		return fmt.Errorf("贷款期限必须大于0")
	}
	if in.Purpose == "" {
		return fmt.Errorf("贷款用途不能为空")
	}
	return nil
}

// generateApplicationId 生成申请编号
func (l *CreateLoanApplicationLogic) generateApplicationId() string {
	// 格式：LN + 年月日 + 6位随机数
	now := time.Now()
	dateStr := now.Format("20060102")
	randomStr := stringx.Randn(6)
	return fmt.Sprintf("LN%s%s", dateStr, randomStr)
}
