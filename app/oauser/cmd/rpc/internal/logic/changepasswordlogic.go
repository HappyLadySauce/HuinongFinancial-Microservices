package logic

import (
	"context"
	"regexp"
	"time"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证管理
func (l *ChangePasswordLogic) ChangePassword(in *oauser.ChangePasswordReq) (*oauser.ChangePasswordResp, error) {
	l.Infof("修改后台用户密码请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" || in.OldPassword == "" || in.NewPassword == "" {
		l.Infof("修改密码参数不完整")
		return nil, constants.ErrInvalidParams
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		l.Infof("手机号格式无效")
		return nil, constants.ErrPhoneInvalid
	}

	// 验证新密码长度
	if len(in.NewPassword) < 6 {
		l.Infof("新密码长度不足")
		return nil, constants.ErrInvalidParams
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 检查用户状态
	if user.Status != constants.UserStatusNormal {
		l.Infof("用户状态异常")
		return nil, constants.ErrUserDisabled
	}

	// 验证旧密码
	if !utils.CheckPassword(in.OldPassword, user.PasswordHash) {
		l.Infof("旧密码错误")
		return nil, constants.ErrPasswordError
	}

	// 加密新密码
	hashedNewPassword, err := utils.HashPassword(in.NewPassword)
	if err != nil {
		l.Errorf("新密码加密失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 更新密码
	user.PasswordHash = hashedNewPassword
	user.UpdatedAt = time.Now()
	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("更新密码失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("修改密码成功, user_id: %d", user.Id)
	return &oauser.ChangePasswordResp{}, nil
}
