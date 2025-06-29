package logic

import (
	"context"
	"errors"
	"fmt"
	"rpc/auth"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateTokenLogic) ValidateToken(in *auth.ValidateTokenReq) (*auth.ValidateTokenResp, error) {
	// 参数验证
	if in.AccessToken == "" {
		return nil, errors.New("token不能为空")
	}

	// 解析并验证JWT token
	claims, err := l.svcCtx.JWTUtil.ValidateToken(in.AccessToken)
	if err != nil {
		l.Infof("token验证失败: %v", err)
		return nil, errors.New("token无效或已过期")
	}

	// 检查Redis中的token是否存在（黑名单检查）
	tokenKey := fmt.Sprintf("token:%s:%d", claims.UserType, claims.UserID)
	storedToken, err := l.svcCtx.Redis.GetCtx(l.ctx, tokenKey)
	if err != nil {
		// Redis错误不影响验证，只记录日志
		l.Errorf("从Redis获取token失败: %v", err)
	} else if storedToken != in.AccessToken {
		// 如果Redis中的token与请求的token不匹配，可能是token已被刷新
		l.Infof("token不匹配，可能已被刷新")
		return nil, errors.New("token已失效")
	}

	return &auth.ValidateTokenResp{
		UserId:   claims.UserID,
		UserType: claims.UserType,
		Roles:    []string{}, // 暂时不实现角色
		IsValid:  true,
	}, nil
}
