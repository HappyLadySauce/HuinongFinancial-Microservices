package logic

import (
	"context"
	"fmt"
	"time"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLeaseApplicationLogic {
	return &CancelLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLeaseApplicationLogic) CancelLeaseApplication(in *lease.CancelLeaseApplicationReq) (*lease.CancelLeaseApplicationResp, error) {
	// 参数验证
	if in.ApplicationId == "" {
		return nil, fmt.Errorf("申请编号不能为空")
	}

	// 查询申请信息
	application, err := l.svcCtx.LeaseApplicationsModel.FindOneByApplicationId(l.ctx, in.ApplicationId)
	if err != nil {
		l.Errorf("查询申请失败: %v", err)
		return nil, fmt.Errorf("申请不存在")
	}

	// 检查申请状态是否可以撤销
	if application.Status != "pending" {
		return nil, fmt.Errorf("只有待审批状态的申请才可以撤销")
	}

	// 更新申请状态为已撤销
	now := time.Now()
	application.Status = "cancelled"
	application.UpdatedAt = now

	err = l.svcCtx.LeaseApplicationsModel.Update(l.ctx, application)
	if err != nil {
		l.Errorf("撤销申请失败: %v", err)
		return nil, fmt.Errorf("撤销申请失败")
	}

	// 记录撤销原因 (可以考虑在future增加撤销原因字段到数据库)
	l.Infof("租赁申请已撤销 - 申请编号: %s, 撤销原因: %s", in.ApplicationId, in.Reason)

	// TODO: 可能需要调用其他服务进行后续处理
	// 比如释放库存预占、发送通知等

	return &lease.CancelLeaseApplicationResp{}, nil
}
