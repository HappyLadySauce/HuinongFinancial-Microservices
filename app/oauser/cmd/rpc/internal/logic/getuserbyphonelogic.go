package logic

import (
	"context"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByPhoneLogic {
	return &GetUserByPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *GetUserByPhoneLogic) GetUserByPhone(in *oauser.GetUserInfoReq) (*oauser.GetUserInfoResp, error) {
	l.Infof("获取后台用户信息请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" {
		l.Infof("手机号参数为空")
		return nil, constants.ErrInvalidParams
	}

	// 查询用户信息
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 构造返回的用户信息
	userInfo := &oauser.UserInfo{
		Id:        int64(user.Id),
		Phone:     user.Phone,
		Name:      user.Name,
		Nickname:  user.Nickname,
		Age:       int32(user.Age),
		Gender:    int32(user.Gender),
		Role:      user.Role,
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}

	l.Infof("获取后台用户信息成功, user_id: %d", user.Id)

	return &oauser.GetUserInfoResp{
		UserInfo: userInfo,
	}, nil
}
