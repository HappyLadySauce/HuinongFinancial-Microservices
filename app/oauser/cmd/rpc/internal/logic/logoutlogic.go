package logic

import (
	"context"
	"time"

	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *oauser.LogoutReq) (*oauser.LogoutResp, error) {
	log := logger.WithContext(l.ctx)
	log.Info("后台用户注销请求")

	// 参数验证
	if in.Token == "" {
		log.Warn("Token参数为空")
		return &oauser.LogoutResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证 token 有效性
	claims, err := utils.ParseToken(in.Token, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		log.WithError(err).Warn("Token解析失败")
		return &oauser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: constants.GetMessage(constants.CodeTokenInvalid),
		}, nil
	}

	// 检查 token 是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		log.Warn("Token已过期")
		return &oauser.LogoutResp{
			Code:    constants.CodeTokenExpired,
			Message: constants.GetMessage(constants.CodeTokenExpired),
		}, nil
	}

	// 检查用户类型
	if claims.UserType != "oa" {
		log.Warn("Token用户类型不匹配")
		return &oauser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: "Token用户类型不匹配",
		}, nil
	}

	// 注销逻辑：在实际应用中，可以将 token 加入黑名单
	// 由于这是一个简单的实现，我们只记录日志
	// 在生产环境中，应该将 token 存储到 Redis 黑名单中

	log.WithField("user_id", claims.UserID).
		WithField("phone", claims.Phone).
		Info("后台用户注销成功")

	return &oauser.LogoutResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
