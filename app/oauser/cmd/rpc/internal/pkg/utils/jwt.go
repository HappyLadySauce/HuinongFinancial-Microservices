package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT 载荷
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Phone    string `json:"phone"`
	UserType string `json:"user_type"` // app, oa
	Role     string `json:"role"`      // admin, operator (仅对 oa 用户有效)
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
// role 参数：对于 B 端用户传入 admin/operator，对于 C 端用户可以传入空字符串
func GenerateToken(userID int64, phone, userType, role, secret string, expireSeconds int64) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Phone:    phone,
		UserType: userType,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "huinong-financial",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT token
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// IsAdmin 判断用户是否为管理员
func (c *JWTClaims) IsAdmin() bool {
	return c.UserType == "oa" && c.Role == "admin"
}

// IsOperator 判断用户是否为操作员
func (c *JWTClaims) IsOperator() bool {
	return c.UserType == "oa" && c.Role == "operator"
}

// IsOAUser 判断是否为后台用户
func (c *JWTClaims) IsOAUser() bool {
	return c.UserType == "oa"
}

// HasRole 检查用户是否具有指定角色
func (c *JWTClaims) HasRole(role string) bool {
	return c.UserType == "oa" && c.Role == role
}

// HasAnyRole 检查用户是否具有任意一个指定角色
func (c *JWTClaims) HasAnyRole(roles []string) bool {
	if c.UserType != "oa" {
		return false
	}
	for _, role := range roles {
		if c.Role == role {
			return true
		}
	}
	return false
}

// ValidateAndGetClaims 验证 token 并返回 claims（组合验证方法）
func ValidateAndGetClaims(tokenString, secret string) (*JWTClaims, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// JWTUtils JWT 工具结构体
type JWTUtils struct {
	secret string
}

// NewJWTUtils 创建 JWT 工具实例
func NewJWTUtils(secret string) *JWTUtils {
	return &JWTUtils{
		secret: secret,
	}
}

// ValidateAndGetClaims 验证 token 并返回 claims
func (j *JWTUtils) ValidateAndGetClaims(tokenString string) (*JWTClaims, error) {
	return ValidateAndGetClaims(tokenString, j.secret)
}

// GenerateToken 生成 token
func (j *JWTUtils) GenerateToken(userID int64, phone, userType, role string, expireSeconds int64) (string, error) {
	return GenerateToken(userID, phone, userType, role, j.secret, expireSeconds)
}
