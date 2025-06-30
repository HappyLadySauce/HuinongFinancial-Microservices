package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyLeaseApplicationLogic {
	return &GetMyLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyLeaseApplicationLogic) GetMyLeaseApplication(applicationId string) (resp *types.GetLeaseApplicationResp, err error) {
	// 调用 Lease RPC 获取申请详情
	rpcResp, err := l.svcCtx.LeaseRpc.GetLeaseApplication(l.ctx, &leaseclient.GetLeaseApplicationReq{
		ApplicationId: applicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换申请信息
	return &types.GetLeaseApplicationResp{
		ApplicationInfo: types.LeaseApplicationInfo{
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
		},
	}, nil
}
