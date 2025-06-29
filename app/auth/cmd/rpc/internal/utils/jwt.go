package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	UserType string `json:"user_type"` // "app" 或 "oa"
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTUtil JWT工具类
type JWTUtil struct {
	accessSecret string
	accessExpire int64
}

// NewJWTUtil 创建JWT工具实例
func NewJWTUtil(accessSecret string, accessExpire int64) *JWTUtil {
	return &JWTUtil{
		accessSecret: accessSecret,
		accessExpire: accessExpire,
	}
}

// GenerateToken 生成JWT token
func (j *JWTUtil) GenerateToken(userID int64, userType, username string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		UserType: userType,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "huinong-auth",
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.accessExpire))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

// ValidateToken 验证JWT token
func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
} 