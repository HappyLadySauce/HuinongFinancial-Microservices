package lease

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyLeaseApplicationLogic {
	return &UpdateMyLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyLeaseApplicationLogic) UpdateMyLeaseApplication(req *types.UpdateLeaseApplicationReq) (resp *types.UpdateLeaseApplicationResp, err error) {
	// 调用 Lease RPC 更新申请
	rpcResp, err := l.svcCtx.LeaseRpc.UpdateLeaseApplication(l.ctx, &leaseclient.UpdateLeaseApplicationReq{
		ApplicationId:   req.ApplicationId,
		Purpose:         req.Purpose,
		DeliveryAddress: req.DeliveryAddress,
		ContactPhone:    req.ContactPhone,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换申请信息
	return &types.UpdateLeaseApplicationResp{
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
