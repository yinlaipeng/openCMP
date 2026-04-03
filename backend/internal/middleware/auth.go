package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取JWT密钥
		jwtSecret, exists := c.Get("jwt_secret")
		if !exists {
			logger.Error("JWT secret not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 解析Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret.(string)), nil
		})

		if err != nil || !token.Valid {
			logger.Info("Invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 获取声明
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查用户ID是否存在
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			logger.Error("User ID not found in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("user_name", claims["user_name"])
		c.Set("user_email", claims["user_email"])

		c.Next()
	}
}

// AdminOnlyMiddleware 仅管理员可访问的中间件
func AdminOnlyMiddleware(logger *zap.Logger, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 检查用户是否为管理员
		var user model.User
		if err := db.Preload("Domain").First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				logger.Error("Database error when checking admin status", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		// 检查用户是否为管理员（这里简化为检查是否为默认管理员用户，实际应用中应该检查角色）
		var adminRole model.Role
		if err := db.Where("name = ? AND type = ?", "admin", "system").First(&adminRole).Error; err != nil {
			logger.Error("Admin role not found", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// 检查用户是否拥有管理员角色
		var userRole model.UserRole
		if err := db.Where("user_id = ? AND role_id = ?", userID, adminRole.ID).First(&userRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
				c.Abort()
				return
			}
			logger.Error("Database error when checking user role", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, userName, userEmail string, jwtSecret string, expireHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   float64(userID),
		"user_name": userName,
		"user_email": userEmail,
		"exp":       time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}