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

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证管理
func (l *RegisterLogic) Register(in *appuser.RegisterReq) (*appuser.RegisterResp, error) {
	l.Infof("用户注册请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		l.Infof("注册参数不完整")
		return nil, constants.ErrInvalidParams
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		l.Infof("手机号格式无效")
		return nil, constants.ErrPhoneInvalid
	}

	// 验证密码长度
	if len(in.Password) < 6 {
		l.Infof("密码长度不足")
		return nil, constants.ErrInvalidParams
	}

	// 检查用户是否已存在
	existUser, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}
	if existUser != nil {
		l.Infof("用户已存在")
		return nil, constants.ErrUserExists
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(in.Password)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 创建用户
	now := time.Now()
	newUser := &model.AppUsers{
		Phone:     in.Phone,
		Password:  hashedPassword,
		Name:      "", // 注册时姓名为空，后续完善资料
		Nickname:  "",
		Age:       0,
		Gender:    0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := l.svcCtx.AppUserModel.Insert(l.ctx, newUser)
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 获取插入的用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		userID,
		in.Phone,
		"app",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("用户注册成功, user_id: %d", userID)
	return &appuser.RegisterResp{
		Token: token,
	}, nil
}
