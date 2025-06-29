package utils

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// go-zero JWT中用户ID的标准键名
	JwtUserIdKey = "user_id"
)

// GetUserIdFromCtx 从context中获取用户ID
func GetUserIdFromCtx(ctx context.Context) (int64, error) {
	// go-zero框架会将JWT中的数据存放到context中
	// 优先尝试标准的键名
	possibleKeys := []string{"user_id", "sub", "userId", "uid"}

	var value interface{}
	for _, key := range possibleKeys {
		value = ctx.Value(key)
		if value != nil {
			logx.Infof("Found admin info with key: %s, value: %v, type: %T", key, value, value)
			break
		}
	}

	if value == nil {
		logx.Error("管理员未登录：无法从context中获取用户信息")
		return 0, errors.New("管理员未登录")
	}

	// 根据不同类型处理值
	switch v := value.(type) {
	case json.Number:
		// go-zero框架中，JWT的数值会被解析为json.Number类型
		userId, err := v.Int64()
		if err != nil {
			logx.Errorf("JSON Number类型管理员ID解析失败: %s, error: %v", v, err)
			return 0, errors.New("管理员ID格式错误")
		}
		logx.Infof("解析JSON Number管理员ID成功: %d", userId)
		return userId, nil
	case string:
		userId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logx.Errorf("管理员ID字符串解析失败: %s, error: %v", v, err)
			return 0, errors.New("管理员ID格式错误")
		}
		logx.Infof("解析字符串管理员ID成功: %d", userId)
		return userId, nil
	case int64:
		logx.Infof("获取int64管理员ID成功: %d", v)
		return v, nil
	case int:
		userId := int64(v)
		logx.Infof("转换int管理员ID成功: %d", userId)
		return userId, nil
	case float64:
		userId := int64(v)
		logx.Infof("转换float64管理员ID成功: %d", userId)
		return userId, nil
	default:
		logx.Errorf("不支持的管理员ID类型: %T, value: %v", v, v)
		return 0, errors.New("管理员ID类型错误")
	}
}
