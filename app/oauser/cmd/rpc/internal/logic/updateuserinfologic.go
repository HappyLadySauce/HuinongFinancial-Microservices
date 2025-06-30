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

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *oauser.UpdateUserInfoReq) (*oauser.UpdateUserInfoResp, error) {
	l.Infof("更新后台用户信息请求, phone: %s", in.UserInfo.Phone)

	// 参数验证
	if in.UserInfo == nil || in.UserInfo.Phone == "" {
		l.Infof("用户信息参数无效")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	userInfo := in.UserInfo

	// 查找现有用户
	existingUser, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, userInfo.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return &oauser.UpdateUserInfoResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if existingUser.Status == constants.UserStatusDisabled {
		l.Infof("用户账号被禁用")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 更新用户信息
	existingUser.Name = userInfo.Name
	existingUser.Nickname = userInfo.Nickname
	existingUser.Age = uint64(userInfo.Age)
	existingUser.Gender = uint64(userInfo.Gender)
	existingUser.UpdatedAt = time.Now()

	// 如果提供了角色信息，也更新角色（需要权限控制）
	if userInfo.Role != "" {
		if userInfo.Role == constants.RoleAdmin || userInfo.Role == constants.RoleOperator {
			existingUser.Role = userInfo.Role
		} else {
			l.Infof("无效的用户角色")
			return &oauser.UpdateUserInfoResp{
				Code:    constants.CodeInvalidParams,
				Message: "无效的用户角色，只支持 admin 或 operator",
			}, nil
		}
	}

	// 保存到数据库
	err = l.svcCtx.OaUserModel.Update(l.ctx, existingUser)
	if err != nil {
		l.Errorf("更新用户信息失败: %v", err)
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 构造返回的用户信息
	updatedUserInfo := &oauser.UserInfo{
		Id:        int64(existingUser.Id),
		Phone:     existingUser.Phone,
		Name:      existingUser.Name,
		Nickname:  existingUser.Nickname,
		Age:       int32(existingUser.Age),
		Gender:    int32(existingUser.Gender),
		Role:      existingUser.Role,
		Status:    int32(existingUser.Status),
		CreatedAt: existingUser.CreatedAt.Unix(),
		UpdatedAt: existingUser.UpdatedAt.Unix(),
	}

	l.Infof("更新后台用户信息成功, user_id: %d", existingUser.Id)

	return &oauser.UpdateUserInfoResp{
		Code:     constants.CodeSuccess,
		Message:  constants.GetMessage(constants.CodeSuccess),
		UserInfo: updatedUserInfo,
	}, nil
}
