package logic

import (
	"context"
	"time"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(in *oauser.UpdateUserStatusReq) (*oauser.UpdateUserStatusResp, error) {
	l.Infof("更新用户状态请求, phone: %s, status: %d", in.Phone, in.Status)

	// 参数验证
	if in.Phone == "" {
		l.Infof("手机号不能为空")
		return nil, constants.ErrInvalidParams
	}

	// 验证状态值有效性
	if in.Status != constants.UserStatusNormal && in.Status != constants.UserStatusDisabled {
		l.Infof("状态值无效")
		return nil, constants.ErrInvalidParams
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 更新用户状态
	user.Status = uint64(in.Status)
	user.UpdatedAt = time.Now()

	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("更新用户状态失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("用户状态更新成功, user_id: %d, new_status: %d", user.Id, in.Status)
	return &oauser.UpdateUserStatusResp{
		Status: in.Status,
	}, nil
}
