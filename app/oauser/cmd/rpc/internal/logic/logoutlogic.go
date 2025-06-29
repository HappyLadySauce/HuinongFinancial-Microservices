package logic

import (
	"context"

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

	// 验证 token 有效性和获取用户信息
	claims, err := utils.ValidateAndGetClaims(in.Token, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		log.WithError(err).Warn("Token验证失败")
		// 根据错误类型返回不同的错误码
		if err.Error() == "token expired" {
			return &oauser.LogoutResp{
				Code:    constants.CodeTokenExpired,
				Message: constants.GetMessage(constants.CodeTokenExpired),
			}, nil
		}
		return &oauser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: constants.GetMessage(constants.CodeTokenInvalid),
		}, nil
	}

	// 检查用户类型
	if !claims.IsOAUser() {
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
		WithField("role", claims.Role).
		Info("后台用户注销成功")

	return &oauser.LogoutResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
