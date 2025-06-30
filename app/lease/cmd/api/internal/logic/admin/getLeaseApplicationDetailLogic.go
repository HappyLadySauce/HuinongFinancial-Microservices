package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseApplicationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLeaseApplicationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseApplicationDetailLogic {
	return &GetLeaseApplicationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLeaseApplicationDetailLogic) GetLeaseApplicationDetail(applicationId string) (resp *types.GetLeaseApplicationResp, err error) {
	// 调用 Lease RPC 获取申请详情
	rpcResp, err := l.svcCtx.LeaseRpc.GetLeaseApplication(l.ctx, &leaseclient.GetLeaseApplicationReq{
		ApplicationId: applicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return &types.GetLeaseApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.GetLeaseApplicationResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	// 转换申请信息
	if rpcResp.ApplicationInfo != nil {
		resp.ApplicationInfo = types.LeaseApplicationInfo{
			Id:              rpcResp.ApplicationInfo.Id,
			ApplicationId:   rpcResp.ApplicationInfo.ApplicationId,
			UserId:          rpcResp.ApplicationInfo.UserId,
			ApplicantName:   rpcResp.ApplicationInfo.ApplicantName,
			ProductId:       rpcResp.ApplicationInfo.ProductId,
			ProductCode:     rpcResp.ApplicationInfo.ProductCode,
			Name:            rpcResp.ApplicationInfo.Name,
			Type:            rpcResp.ApplicationInfo.Type,
			Machinery:       rpcResp.ApplicationInfo.Machinery,
			StartDate:       rpcResp.ApplicationInfo.StartDate,
			EndDate:         rpcResp.ApplicationInfo.EndDate,
			Duration:        rpcResp.ApplicationInfo.Duration,
			DailyRate:       rpcResp.ApplicationInfo.DailyRate,
			TotalAmount:     rpcResp.ApplicationInfo.TotalAmount,
			Deposit:         rpcResp.ApplicationInfo.Deposit,
			DeliveryAddress: rpcResp.ApplicationInfo.DeliveryAddress,
			ContactPhone:    rpcResp.ApplicationInfo.ContactPhone,
			Purpose:         rpcResp.ApplicationInfo.Purpose,
			Status:          rpcResp.ApplicationInfo.Status,
			CreatedAt:       rpcResp.ApplicationInfo.CreatedAt,
			UpdatedAt:       rpcResp.ApplicationInfo.UpdatedAt,
		}
	}

	return resp, nil
}
