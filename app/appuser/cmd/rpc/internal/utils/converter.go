package utils

import (
	"model"
	"rpc/appuser"
)

// ModelToProto 将数据库模型转换为Proto消息
func ModelToProto(user *model.AppUsers) *appuser.AppUserInfo {
	if user == nil {
		return nil
	}

	return &appuser.AppUserInfo{
		Id:         int64(user.Id),
		Account:    user.Phone,
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
}

// ModelsToProtos 批量转换数据库模型为Proto消息
func ModelsToProtos(users []*model.AppUsers) []*appuser.AppUserInfo {
	if len(users) == 0 {
		return nil
	}

	result := make([]*appuser.AppUserInfo, 0, len(users))
	for _, user := range users {
		if protoUser := ModelToProto(user); protoUser != nil {
			result = append(result, protoUser)
		}
	}
	return result
}
