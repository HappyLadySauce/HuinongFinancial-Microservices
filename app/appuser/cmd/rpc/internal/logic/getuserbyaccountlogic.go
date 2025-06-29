package logic

import (
	"context"
	"errors"
	"model"
	"rpc/appuser"
	"rpc/internal/svc"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByAccountLogic {
	return &GetUserByAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证相关
func (l *GetUserByAccountLogic) GetUserByAccount(in *appuser.GetUserByAccountReq) (*appuser.AppUserInfo, error) {
	// 参数验证
	if in.Account == "" {
		return nil, errors.New("账号不能为空")
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

	// 转换为Proto格式并返回
	return utils.ModelToProto(user), nil
}
