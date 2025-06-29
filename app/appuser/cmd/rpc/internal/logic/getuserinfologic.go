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

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *GetUserInfoLogic) GetUserInfo(in *appuser.GetUserInfoReq) (*appuser.AppUserInfo, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	// 根据ID查询用户
	user, err := l.svcCtx.AppUsersModel.FindOne(l.ctx, uint64(in.UserId))
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
