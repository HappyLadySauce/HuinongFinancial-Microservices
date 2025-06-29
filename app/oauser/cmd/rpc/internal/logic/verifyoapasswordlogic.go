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

type VerifyOAPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyOAPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyOAPasswordLogic {
	return &VerifyOAPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyOAPasswordLogic) VerifyOAPassword(in *oauser.VerifyOAPasswordReq) (*oauser.VerifyOAPasswordResp, error) {
	// 参数验证
	if in.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
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

	// 检查用户状态
	if user.Status == 2 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	err = l.svcCtx.PasswordUtil.VerifyPassword(user.Password, in.Password)
	if err != nil {
		l.Infof("用户[%s]密码验证失败", in.Username)
		return nil, errors.New("密码错误")
	}

	// 返回验证结果
	return &oauser.VerifyOAPasswordResp{
		UserInfo: utils.ModelToProto(user),
	}, nil
}
