package logic

import (
	"context"

	"rpc/auth"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppLoginLogic {
	return &AppLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *AppLoginLogic) AppLogin(in *auth.AppLoginReq) (*auth.LoginResp, error) {
	// 1. 参数验证
	if in.Account == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "账号和密码不能为空")
	}

	// 2. 调用AppUser RPC验证用户凭据
	l.Infof("验证用户登录: account=%s", in.Account)

	// 通过服务发现调用 AppUser RPC
	verifyResp, err := l.verifyAppUserPassword(in.Account, in.Password)
	if err != nil {
		l.Errorf("验证用户密码失败: %v", err)
		return nil, status.Error(codes.Unauthenticated, "账号或密码错误")
	}

	// 检查用户状态
	if verifyResp.Status != 1 {
		l.Errorf("用户状态异常: userId=%d, status=%d", verifyResp.UserId, verifyResp.Status)
		return nil, status.Error(codes.Unauthenticated, "账号状态异常")
	}

	userId := verifyResp.UserId
	l.Infof("用户验证成功: account=%s, userId=%d", in.Account, userId)

	// 3. 生成JWT Token
	tokenResp, err := l.generateToken(userId, "app", []string{"user"})
	if err != nil {
		l.Errorf("生成Token失败: %v", err)
		return nil, status.Error(codes.Internal, "登录失败")
	}

	l.Infof("C端用户登录成功, userId: %d, account: %s", userId, in.Account)

	return &auth.LoginResp{
		AccessToken:   tokenResp.AccessToken,
		AccessExpire:  tokenResp.AccessExpire,
		RefreshToken:  tokenResp.RefreshToken,
		RefreshExpire: tokenResp.RefreshExpire,
		UserId:        userId,
		UserType:      "app",
	}, nil
}

// verifyAppUserPassword 通过服务发现调用 AppUser 服务验证密码
func (l *AppLoginLogic) verifyAppUserPassword(account, password string) (*VerifyUserPasswordResp, error) {
	// 暂时返回模拟数据，后续会通过服务发现实现真实调用
	l.Infof("调用 AppUser 服务验证密码: account=%s", account)

	// TODO: 通过服务发现调用真实的 AppUser RPC 服务
	// 这里暂时使用模拟逻辑
	if account == "13800138000" && password == "123456" {
		return &VerifyUserPasswordResp{
			UserId: 1001,
			Status: 1,
		}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "账号或密码错误")
}

// 定义调用 AppUser 服务所需的消息类型
type VerifyUserPasswordReq struct {
	Account  string
	Password string
}

type VerifyUserPasswordResp struct {
	UserId int64
	Status int32
}

// generateToken 生成JWT Token
func (l *AppLoginLogic) generateToken(userId int64, userType string, roles []string) (*auth.GenerateTokenResp, error) {
	generateLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	return generateLogic.GenerateToken(&auth.GenerateTokenReq{
		UserId:   userId,
		UserType: userType,
		Roles:    roles,
	})
}
