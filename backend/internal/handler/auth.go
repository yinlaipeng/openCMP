package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// AuthHandler 认证 Handler
type AuthHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

// NewAuthHandler 创建认证 Handler
func NewAuthHandler(db *gorm.DB, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: service.NewUserService(db),
		logger:      logger,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	user, err := h.userService.GetUserByName(c.Request.Context(), req.Username)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 检查用户是否启用
	if !user.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
		return
	}

	// TODO: 验证密码（演示模式跳过）
	// if !checkPassword(user.Password, req.Password) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	// 	return
	// }

	// 更新最后登录信息
	h.userService.UpdateLastLogin(c.Request.Context(), user.ID, c.ClientIP())

	// 生成 token（简化实现）
	token := generateToken(user.ID)

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: 使 token 失效
	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从 token 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 验证旧密码并更新新密码
	// if !checkPassword(user.Password, req.OldPassword) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "old password is incorrect"})
	// 	return
	// }

	err := h.userService.UpdatePassword(c.Request.Context(), userID.(uint), req.NewPassword)
	if err != nil {
		h.logger.Error("failed to change password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// generateToken 生成 token（简化实现）
func generateToken(userID uint) string {
	// 实际应该使用 JWT
	return "token-" + string(rune(userID)) + "-" + time.Now().Format("20060102")
}

// checkPassword 验证密码（简化实现）
func checkPassword(stored, input string) bool {
	// 实际应该使用 bcrypt 等加密算法
	return stored == input
}
