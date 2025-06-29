package logic

import (
	"context"
	"errors"
	"fmt"
	"rpc/auth"
	"rpc/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Token管理
func (l *GenerateTokenLogic) GenerateToken(in *auth.GenerateTokenReq) (*auth.GenerateTokenResp, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}
	if in.UserType == "" {
		return nil, errors.New("用户类型不能为空")
	}

	// 验证用户类型
	if in.UserType != "app" && in.UserType != "oa" {
		return nil, errors.New("用户类型无效，必须为 app 或 oa")
	}

	// 生成用户名（从ID生成临时用户名）
	username := fmt.Sprintf("%s_user_%d", in.UserType, in.UserId)

	// 生成JWT token
	token, err := l.svcCtx.JWTUtil.GenerateToken(in.UserId, in.UserType, username)
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, errors.New("生成token失败")
	}

	// 计算过期时间
	now := time.Now().Unix()
	accessExpire := now + l.svcCtx.Config.JwtAuth.AccessExpire

	// 将token存储到Redis（可选，用于黑名单管理）
	tokenKey := fmt.Sprintf("token:%s:%d", in.UserType, in.UserId)
	err = l.svcCtx.Redis.SetexCtx(l.ctx, tokenKey, token, int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		l.Errorf("保存token到Redis失败: %v", err)
		// 这里不返回错误，因为token已经生成成功
	}

	return &auth.GenerateTokenResp{
		AccessToken:  token,
		AccessExpire: accessExpire,
		RefreshToken: "", // 暂时不实现刷新token
		RefreshExpire: 0,
	}, nil
}
