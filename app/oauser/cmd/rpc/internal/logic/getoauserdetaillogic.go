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

type GetOAUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOAUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOAUserDetailLogic {
	return &GetOAUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOAUserDetailLogic) GetOAUserDetail(in *oauser.GetOAUserByUsernameReq) (*oauser.OAUserInfo, error) {
	// 参数验证
	if in.Username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 根据用户名查询用户详情
	user, err := l.svcCtx.OaUsersModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("查询用户详情失败: %v", err)
		return nil, errors.New("查询用户详情失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(user), nil
}
