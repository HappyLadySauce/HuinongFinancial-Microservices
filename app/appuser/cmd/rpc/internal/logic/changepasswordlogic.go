package logic

import (
	"context"
	"regexp"
	"time"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"

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
func (l *ChangePasswordLogic) ChangePassword(in *appuser.ChangePasswordReq) (*appuser.ChangePasswordResp, error) {
	l.Infof("修改密码请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" || in.OldPassword == "" || in.NewPassword == "" {
		l.Infof("修改密码参数不完整")
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		l.Infof("手机号格式无效")
		return &appuser.ChangePasswordResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 验证新密码长度
	if len(in.NewPassword) < 6 {
		l.Infof("新密码长度不足")
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeInvalidParams,
			Message: "新密码长度不能少于6位",
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return &appuser.ChangePasswordResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		l.Infof("用户状态异常")
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 验证旧密码
	if !utils.CheckPassword(in.OldPassword, user.Password) {
		l.Infof("旧密码错误")
		return &appuser.ChangePasswordResp{
			Code:    constants.CodePasswordError,
			Message: "旧密码错误",
		}, nil
	}

	// 加密新密码
	hashedNewPassword, err := utils.HashPassword(in.NewPassword)
	if err != nil {
		l.Errorf("新密码加密失败: %v", err)
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 更新密码
	user.Password = hashedNewPassword
	user.UpdatedAt = time.Now()
	err = l.svcCtx.AppUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("更新密码失败: %v", err)
		return &appuser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	l.Infof("修改密码成功, user_id: %d", user.Id)
	return &appuser.ChangePasswordResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
