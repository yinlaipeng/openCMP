package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	// jwt/v5 used for ErrTokenExpired check
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/pkg/utils"
)

// AuthMiddleware 认证中间件（使用 typed JWTClaims，支持 domain_id 注入）
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required", "code": "TOKEN_MISSING"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format", "code": "TOKEN_INVALID"})
			c.Abort()
			return
		}

		// 使用 typed claims 解析，同时验证签名和过期
		claims, err := utils.ParseJWTToken(tokenString, jwtSecret.(string))
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.Info("Token expired", zap.Error(err))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired", "code": "TOKEN_EXPIRED"})
			} else {
				logger.Info("Invalid token", zap.Error(err))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "code": "TOKEN_INVALID"})
			}
			c.Abort()
			return
		}

		// 注入用户基础信息
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)
		c.Set("user_email", claims.UserEmail)

		// 注入域信息（LDAP/域级认证源登录时存在）
		if claims.DomainID > 0 {
			c.Set("domain_id", claims.DomainID)
		}
		if claims.AuthSourceID > 0 {
			c.Set("auth_source_id", claims.AuthSourceID)
		}

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

