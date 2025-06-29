package logic

import (
	"context"
	"errors"
	"fmt"
	"rpc/auth"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InvalidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInvalidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvalidateTokenLogic {
	return &InvalidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InvalidateTokenLogic) InvalidateToken(in *auth.ValidateTokenReq) (*auth.InvalidateTokenResp, error) {
	// 参数验证
	if in.AccessToken == "" {
		return nil, errors.New("token不能为空")
	}

	// 解析JWT token获取用户信息
	claims, err := l.svcCtx.JWTUtil.ValidateToken(in.AccessToken)
	if err != nil {
		l.Infof("token解析失败: %v", err)
		// 即使token无效，我们也认为失效操作成功
		return &auth.InvalidateTokenResp{
			Success: true,
		}, nil
	}

	// 从Redis中删除token
	tokenKey := fmt.Sprintf("token:%s:%d", claims.UserType, claims.UserID)
	_, err = l.svcCtx.Redis.DelCtx(l.ctx, tokenKey)
	if err != nil {
		l.Errorf("从Redis删除token失败: %v", err)
		return nil, errors.New("token失效操作失败")
	}

	// 记录日志
	l.Infof("用户[%s:%d]的token已被失效", claims.UserType, claims.UserID)

	return &auth.InvalidateTokenResp{
		Success: true,
	}, nil
}
