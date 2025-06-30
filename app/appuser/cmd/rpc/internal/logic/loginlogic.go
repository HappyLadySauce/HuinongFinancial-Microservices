package logic

import (
	"context"
	"regexp"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证管理
func (l *LoginLogic) Login(in *appuser.LoginReq) (*appuser.LoginResp, error) {
	l.Infof("用户登录请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		l.Infof("登录参数不完整")
		return &appuser.LoginResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		l.Infof("手机号格式无效")
		return &appuser.LoginResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return &appuser.LoginResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &appuser.LoginResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if user.Status == 2 {
		l.Infof("用户账号被冻结")
		return &appuser.LoginResp{
			Code:    constants.CodeUserFrozen,
			Message: constants.GetMessage(constants.CodeUserFrozen),
		}, nil
	}
	if user.Status == 3 {
		l.Infof("用户账号被禁用")
		return &appuser.LoginResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 验证密码
	if !utils.CheckPassword(in.Password, user.Password) {
		l.Infof("密码错误")
		return &appuser.LoginResp{
			Code:    constants.CodePasswordError,
			Message: constants.GetMessage(constants.CodePasswordError),
		}, nil
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		int64(user.Id),
		user.Phone,
		"app",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return &appuser.LoginResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	l.Infof("用户登录成功, user_id: %d", user.Id)
	return &appuser.LoginResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
		Token:   token,
	}, nil
}
