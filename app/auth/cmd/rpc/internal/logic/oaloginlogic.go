package logic

import (
	"context"

	"rpc/auth"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OALoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOALoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OALoginLogic {
	return &OALoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// B端管理员登录
func (l *OALoginLogic) OALogin(in *auth.OALoginReq) (*auth.LoginResp, error) {
	// 1. 参数验证
	if in.Username == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "用户名和密码不能为空")
	}

	// 2. 调用OAUser RPC验证管理员凭据
	l.Infof("验证管理员登录: username=%s", in.Username)

	// 通过服务发现调用 OAUser RPC
	verifyResp, err := l.verifyOAUserPassword(in.Username, in.Password)
	if err != nil {
		l.Errorf("验证管理员密码失败: %v", err)
		return nil, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}

	// 检查用户状态
	if verifyResp.UserInfo.Status != 1 {
		l.Errorf("管理员状态异常: userId=%d, status=%d", verifyResp.UserInfo.Id, verifyResp.UserInfo.Status)
		return nil, status.Error(codes.Unauthenticated, "账号状态异常")
	}

	userId := verifyResp.UserInfo.Id
	l.Infof("管理员验证成功: username=%s, userId=%d", in.Username, userId)

	// 3. 生成JWT Token (管理员拥有admin角色)
	tokenResp, err := l.generateToken(userId, "oa", []string{"admin"})
	if err != nil {
		l.Errorf("生成Token失败: %v", err)
		return nil, status.Error(codes.Internal, "登录失败")
	}

	l.Infof("B端管理员登录成功, userId: %d, username: %s", userId, in.Username)

	return &auth.LoginResp{
		AccessToken:   tokenResp.AccessToken,
		AccessExpire:  tokenResp.AccessExpire,
		RefreshToken:  tokenResp.RefreshToken,
		RefreshExpire: tokenResp.RefreshExpire,
		UserId:        userId,
		UserType:      "oa",
	}, nil
}

// verifyOAUserPassword 通过服务发现调用 OAUser 服务验证密码
func (l *OALoginLogic) verifyOAUserPassword(username, password string) (*VerifyOAPasswordResp, error) {
	// 暂时返回模拟数据，后续会通过服务发现实现真实调用
	l.Infof("调用 OAUser 服务验证密码: username=%s", username)

	// TODO: 通过服务发现调用真实的 OAUser RPC 服务
	// 这里暂时使用模拟逻辑
	if username == "admin" && password == "admin123" {
		return &VerifyOAPasswordResp{
			UserInfo: &OAUserInfo{
				Id:       2001,
				Username: "admin",
				Name:     "系统管理员",
				Status:   1,
			},
		}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "用户名或密码错误")
}

// 定义调用 OAUser 服务所需的消息类型
type OAUserInfo struct {
	Id       int64
	Username string
	Name     string
	Status   int32
}

type VerifyOAPasswordResp struct {
	UserInfo *OAUserInfo
}

// generateToken 生成JWT Token
func (l *OALoginLogic) generateToken(userId int64, userType string, roles []string) (*auth.GenerateTokenResp, error) {
	generateLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	return generateLogic.GenerateToken(&auth.GenerateTokenReq{
		UserId:   userId,
		UserType: userType,
		Roles:    roles,
	})
}
