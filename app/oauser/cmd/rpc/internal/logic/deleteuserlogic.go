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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *DeleteUserLogic) DeleteUser(in *oauser.DeleteUserReq) (*oauser.DeleteUserResp, error) {
	l.Infof("删除后台用户请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" {
		l.Infof("手机号不能为空")
		return nil, constants.ErrInvalidParams
	}

	// 验证调用者权限（简化实现，实际应该验证 JWT token）
	// 注：现在从 JWT 认证上下文获取调用者信息，此处省略 token 验证

	// 查找要删除的用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 禁用用户（软删除）
	user.Status = constants.UserStatusDisabled
	user.UpdatedAt = time.Now()

	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("禁用用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("后台用户删除成功, user_id: %d", user.Id)
	return &oauser.DeleteUserResp{}, nil
}
