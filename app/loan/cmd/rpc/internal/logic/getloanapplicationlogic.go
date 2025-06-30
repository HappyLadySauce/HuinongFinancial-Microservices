package logic

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("申请ID不能为空")
	}

	// 根据申请ID查询贷款申请
	loanApplication, err := l.svcCtx.LoanApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询贷款申请失败: %v", err)
		return nil, fmt.Errorf("贷款申请不存在")
	}

	// 构造响应
	return &loan.GetLoanApplicationResp{
		ApplicationInfo: &loan.LoanApplicationInfo{
			Id:            int64(loanApplication.Id),
			ApplicationId: loanApplication.ApplicationId,
			UserId:        int64(loanApplication.UserId),
			ApplicantName: loanApplication.ApplicantName,
			ProductId:     int64(loanApplication.ProductId),
			Name:          loanApplication.Name,
			Type:          loanApplication.Type,
			Amount:        loanApplication.Amount,
			Duration:      int32(loanApplication.Duration),
			Purpose:       loanApplication.Purpose.String,
			Status:        loanApplication.Status,
			CreatedAt:     loanApplication.CreatedAt.Unix(),
			UpdatedAt:     loanApplication.UpdatedAt.Unix(),
		},
	}, nil
}
