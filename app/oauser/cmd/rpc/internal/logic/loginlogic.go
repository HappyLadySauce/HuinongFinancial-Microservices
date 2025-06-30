package logic

import (
	"context"
	"regexp"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"
	"rpc/oauser"

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
func (l *LoginLogic) Login(in *oauser.LoginReq) (*oauser.LoginResp, error) {
	l.Infof("后台用户登录请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		l.Infof("登录参数不完整")
		return nil, constants.ErrInvalidParams
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		l.Infof("手机号格式无效")
		return nil, constants.ErrPhoneInvalid
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
	if user.Status == constants.UserStatusDisabled {
		l.Infof("用户账号被禁用")
		return nil, constants.ErrUserDisabled
	}

	// 验证密码
	if !utils.CheckPassword(in.Password, user.PasswordHash) {
		l.Infof("密码错误")
		return nil, constants.ErrPasswordError
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		int64(user.Id),
		user.Phone,
		"oa",
		user.Role,
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("后台用户登录成功, user_id: %d", user.Id)
	return &oauser.LoginResp{
		Token: token,
	}, nil
}
