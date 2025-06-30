package loan

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanApplicationLogic {
	return &CreateLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanApplicationLogic) CreateLoanApplication(req *types.CreateLoanApplicationReq) (resp *types.CreateLoanApplicationResp, err error) {
	// 获取当前用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析用户ID失败: %v", err)
		return &types.CreateLoanApplicationResp{
			Code:    400,
			Message: "用户ID无效",
		}, nil
	}

	// 调用 Loan RPC 创建申请
	rpcResp, err := l.svcCtx.LoanRpc.CreateLoanApplication(l.ctx, &loanclient.CreateLoanApplicationReq{
		UserId:    userId,
		ProductId: req.ProductId,
		Name:      req.Name,
		Type:      req.Type,
		Amount:    req.Amount,
		Duration:  req.Duration,
		Purpose:   req.Purpose,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.CreateLoanApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.CreateLoanApplicationResp{
		Code:          rpcResp.Code,
		Message:       rpcResp.Message,
		ApplicationId: rpcResp.ApplicationId,
	}, nil
}
