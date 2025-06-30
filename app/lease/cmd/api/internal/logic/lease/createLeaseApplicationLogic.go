package lease

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseApplicationLogic {
	return &CreateLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLeaseApplicationLogic) CreateLeaseApplication(req *types.CreateLeaseApplicationReq) (resp *types.CreateLeaseApplicationResp, err error) {
	// 获取当前用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析用户ID失败: %v", err)
		return nil, err
	}

	// 调用 Lease RPC 创建申请
	rpcResp, err := l.svcCtx.LeaseRpc.CreateLeaseApplication(l.ctx, &leaseclient.CreateLeaseApplicationReq{
		UserId:          userId,
		ProductId:       req.ProductId,
		ProductCode:     req.ProductCode,
		Name:            req.Name,
		Type:            req.Type,
		Machinery:       req.Machinery,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Duration:        req.Duration,
		DailyRate:       req.DailyRate,
		TotalAmount:     req.TotalAmount,
		Deposit:         req.Deposit,
		DeliveryAddress: req.DeliveryAddress,
		ContactPhone:    req.ContactPhone,
		Purpose:         req.Purpose,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.CreateLeaseApplicationResp{
		ApplicationId: rpcResp.ApplicationId,
	}, nil
}
