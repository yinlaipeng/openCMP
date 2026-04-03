package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	RoleIDs   []uint `json:"role_ids"`
	jwt.RegisteredClaims
}

// GenerateJWTToken 生成JWT token
func GenerateJWTToken(userID uint, userName, userEmail, secret string, expireHours int) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(expireHours))
	claims := &JWTClaims{
		UserID:    userID,
		UserName:  userName,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWTToken 解析JWT token
func ParseJWTToken(tokenString, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// HashPassword 使用bcrypt哈希密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 比较密码和哈希值
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken 生成带角色信息的 JWT token
func GenerateToken(userID uint, userName string, roleIDs []uint, secret string, expireHours int) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(expireHours))
	claims := &JWTClaims{
		UserID:   userID,
		UserName: userName,
		RoleIDs:  roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}