package logic

import (
	"context"
	"errors"
	"model"
	"rpc/internal/svc"
	"rpc/internal/utils"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOAUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOAUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOAUserLogic {
	return &DeleteOAUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteOAUserLogic) DeleteOAUser(in *oauser.GetOAUserByUsernameReq) (*oauser.OAUserInfo, error) {
	// 参数验证
	if in.Username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 根据用户名查询用户
	user, err := l.svcCtx.OaUsersModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 软删除：设置状态为禁用
	err = l.svcCtx.OaUsersModel.UpdateStatusById(l.ctx, int64(user.Id), 2)
	if err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, errors.New("删除用户失败")
	}

	// 记录删除日志
	l.Infof("用户[%s]已被删除（设置为禁用状态）", user.Username)

	// 查询更新后的用户信息
	updatedUser, err := l.svcCtx.OaUsersModel.FindOne(l.ctx, user.Id)
	if err != nil {
		l.Errorf("查询删除后的用户失败: %v", err)
		return nil, errors.New("删除用户失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(updatedUser), nil
}
