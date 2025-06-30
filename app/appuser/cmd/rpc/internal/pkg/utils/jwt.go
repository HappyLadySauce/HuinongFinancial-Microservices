package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken 生成 JWT token - 兼容 go-zero JWT 中间件
// 根据 go-zero 官方文档，使用 jwt.MapClaims
func GenerateToken(userID int64, phone, userType, secret string, expireSeconds int64) (string, error) {
	now := time.Now()
	iat := now.Unix()
	exp := now.Add(time.Duration(expireSeconds) * time.Second).Unix()

	claims := make(jwt.MapClaims)
	claims["exp"] = exp
	claims["iat"] = iat
	claims["user_id"] = userID
	claims["phone"] = phone
	claims["user_type"] = userType

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT token - 兼容 go-zero JWT 中间件
func ParseToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
