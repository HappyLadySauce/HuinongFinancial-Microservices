package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanApplicationLogic {
	return &GetLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoanApplicationLogic) GetLoanApplication(in *loan.GetLoanApplicationReq) (*loan.GetLoanApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return &loan.GetLoanApplicationResp{
			Code:    400,
			Message: "申请编号不能为空",
		}, nil
	}

	// 根据申请编号查询申请信息
	application, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return &loan.GetLoanApplicationResp{
			Code:    404,
			Message: "申请不存在",
		}, nil
	}

	// 转换为响应格式
	applicationInfo := &loan.LoanApplicationInfo{
		Id:            int64(application.Id),
		ApplicationId: application.ApplicationId,
		UserId:        int64(application.UserId),
		ApplicantName: application.ApplicantName,
		ProductId:     int64(application.ProductId),
		Name:          application.Name,
		Type:          application.Type,
		Amount:        application.Amount,
		Duration:      int32(application.Duration),
		Purpose:       application.Purpose.String,
		Status:        application.Status,
		CreatedAt:     application.CreatedAt.Unix(),
		UpdatedAt:     application.UpdatedAt.Unix(),
	}

	return &loan.GetLoanApplicationResp{
		Code:            200,
		Message:         "查询成功",
		ApplicationInfo: applicationInfo,
	}, nil
}
