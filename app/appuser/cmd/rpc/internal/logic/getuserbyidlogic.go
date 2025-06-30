package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *appuser.GetUserByIdReq) (*appuser.GetUserInfoResp, error) {
	// 参数验证
	if in.UserId <= 0 {
		return &appuser.GetUserInfoResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	// 根据用户ID查询用户信息
	user, err := l.svcCtx.AppUserModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &appuser.GetUserInfoResp{
			Code:    404,
			Message: "用户不存在",
		}, nil
	}

	// 检查用户状态
	if user.Status != 1 {
		return &appuser.GetUserInfoResp{
			Code:    400,
			Message: "用户状态异常",
		}, nil
	}

	// 转换为响应格式
	userInfo := &appuser.UserInfo{
		Id:         int64(user.Id),
		Phone:      user.Phone,
		Name:       user.Name,
		Nickname:   user.Nickname,
		Age:        int32(user.Age),
		Gender:     int32(user.Gender),
		Occupation: user.Occupation,
		Address:    user.Address,
		Income:     user.Income,
		Status:     int32(user.Status),
		CreatedAt:  user.CreatedAt.Unix(),
		UpdatedAt:  user.UpdatedAt.Unix(),
	}

	return &appuser.GetUserInfoResp{
		Code:     200,
		Message:  "查询成功",
		UserInfo: userInfo,
	}, nil
}
