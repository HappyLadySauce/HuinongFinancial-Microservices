package logic

import (
	"context"
	"errors"
	"model"
	"rpc/appuser"
	"rpc/internal/svc"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *appuser.CreateUserReq) (*appuser.AppUserInfo, error) {
	// 参数验证
	if in.Account == "" {
		return nil, errors.New("账号不能为空")
	}
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if in.Name == "" {
		return nil, errors.New("姓名不能为空")
	}

	// 检查账号是否已存在
	existUser, err := l.svcCtx.AppUsersModel.FindOneByPhone(l.ctx, in.Account)
	if err != nil && err != model.ErrNotFound {
		l.Errorf("检查账号是否存在失败: %v", err)
		return nil, errors.New("系统错误")
	}
	if existUser != nil {
		return nil, errors.New("账号已存在")
	}

	// 密码加密
	hashedPassword, err := l.svcCtx.PasswordUtil.HashPassword(in.Password)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return nil, errors.New("密码加密失败")
	}

	// 创建用户数据
	newUser := &model.AppUsers{
		Phone:    in.Account,
		Password: hashedPassword,
		Name:     in.Name,
		Status:   1, // 默认状态为正常
	}

	// 插入数据库
	result, err := l.svcCtx.AppUsersModel.Insert(l.ctx, newUser)
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	// 获取插入的ID
	userId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	// 查询创建的用户完整信息
	createdUser, err := l.svcCtx.AppUsersModel.FindOne(l.ctx, uint64(userId))
	if err != nil {
		l.Errorf("查询创建的用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(createdUser), nil
}
