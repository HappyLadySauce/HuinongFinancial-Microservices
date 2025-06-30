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

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *oauser.UpdateUserInfoReq) (*oauser.UpdateUserInfoResp, error) {
	l.Infof("更新后台用户信息请求, user_id: %d", in.UserInfo.Id)

	// 参数验证
	if in.UserInfo == nil || in.UserInfo.Id <= 0 {
		l.Infof("用户信息参数无效")
		return nil, constants.ErrInvalidParams
	}

	// 检查用户是否存在
	existUser, err := l.svcCtx.OaUserModel.FindOne(l.ctx, uint64(in.UserInfo.Id))
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 更新用户信息
	existUser.Name = in.UserInfo.Name
	existUser.Nickname = in.UserInfo.Nickname
	existUser.Age = uint64(in.UserInfo.Age)
	existUser.Gender = uint64(in.UserInfo.Gender)
	existUser.Role = in.UserInfo.Role
	existUser.UpdatedAt = time.Now()

	err = l.svcCtx.OaUserModel.Update(l.ctx, existUser)
	if err != nil {
		l.Errorf("更新用户信息失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 构造返回的用户信息
	userInfo := &oauser.UserInfo{
		Id:        int64(existUser.Id),
		Phone:     existUser.Phone,
		Name:      existUser.Name,
		Nickname:  existUser.Nickname,
		Age:       int32(existUser.Age),
		Gender:    int32(existUser.Gender),
		Role:      existUser.Role,
		Status:    int32(existUser.Status),
		CreatedAt: existUser.CreatedAt.Unix(),
		UpdatedAt: existUser.UpdatedAt.Unix(),
	}

	l.Infof("更新后台用户信息成功, user_id: %d", existUser.Id)
	return &oauser.UpdateUserInfoResp{
		UserInfo: userInfo,
	}, nil
}
