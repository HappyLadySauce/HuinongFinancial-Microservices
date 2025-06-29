package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordUtil 密码工具
type PasswordUtil struct{}

// NewPasswordUtil 创建密码工具实例
func NewPasswordUtil() *PasswordUtil {
	return &PasswordUtil{}
}

// HashPassword 密码加密
func (p *PasswordUtil) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hash), nil
}

// VerifyPassword 验证密码
func (p *PasswordUtil) VerifyPassword(hashedPassword, password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if hashedPassword == "" {
		return errors.New("hashed password cannot be empty")
	}
	
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
} 