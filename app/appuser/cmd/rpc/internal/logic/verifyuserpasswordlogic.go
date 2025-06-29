package logic

import (
	"context"
	"errors"
	"model"
	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyUserPasswordLogic {
	return &VerifyUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyUserPasswordLogic) VerifyUserPassword(in *appuser.VerifyUserPasswordReq) (*appuser.VerifyUserPasswordResp, error) {
	// 参数验证
	if in.Account == "" {
		return nil, errors.New("账号不能为空")
	}
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 根据账号查询用户
	user, err := l.svcCtx.AppUsersModel.FindOneByPhone(l.ctx, in.Account)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 验证密码
	err = l.svcCtx.PasswordUtil.VerifyPassword(user.Password, in.Password)
	if err != nil {
		l.Infof("用户[%s]密码验证失败", in.Account)
		return nil, errors.New("密码错误")
	}

	// 返回验证结果
	return &appuser.VerifyUserPasswordResp{
		UserId: int64(user.Id),
		Status: int32(user.Status),
	}, nil
}
