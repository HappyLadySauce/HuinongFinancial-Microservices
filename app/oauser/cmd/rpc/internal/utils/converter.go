package utils

import (
	"model"
	"rpc/oauser"
	"strings"
)

// ModelToProto 将数据库模型转换为Proto消息
func ModelToProto(user *model.OaUsers) *oauser.OAUserInfo {
	if user == nil {
		return nil
	}

	return &oauser.OAUserInfo{
		Id:       int64(user.Id),
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Roles:    StringToRoles(user.Roles),
		Status:   int32(user.Status),
	}
}

// ModelsToProtos 将数据库模型列表转换为Proto消息列表
func ModelsToProtos(users []*model.OaUsers) []*oauser.OAUserInfo {
	result := make([]*oauser.OAUserInfo, 0, len(users))
	for _, user := range users {
		result = append(result, ModelToProto(user))
	}
	return result
}

// RolesToString 将角色数组转换为数据库存储的字符串
func RolesToString(roles []string) string {
	if len(roles) == 0 {
		return ""
	}
	return strings.Join(roles, ",")
}

// StringToRoles 将数据库存储的字符串转换为角色数组
func StringToRoles(roleStr string) []string {
	if roleStr == "" {
		return []string{}
	}
	return strings.Split(roleStr, ",")
} 