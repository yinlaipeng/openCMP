package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/utils"
)

// AuthHandler 认证 Handler
type AuthHandler struct {
	userService *service.UserService
	logger      *zap.Logger
	jwtSecret   string
	jwtExpire   int
}

// NewAuthHandler 创建认证 Handler
func NewAuthHandler(db *gorm.DB, logger *zap.Logger, jwtSecret string, jwtExpire int) *AuthHandler {
	return &AuthHandler{
		userService: service.NewUserService(db),
		logger:      logger,
		jwtSecret:   jwtSecret,
		jwtExpire:   jwtExpire,
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

	// 验证密码
	if !checkPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 获取用户的角色
	roleIDs, err := h.userService.GetUserRoleIDs(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to get user roles", zap.Error(err))
		// 如果无法获取角色，继续登录但不包含角色信息
		roleIDs = []uint{}
	}

	// 更新最后登录信息
	h.userService.UpdateLastLogin(c.Request.Context(), user.ID, c.ClientIP())

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Name, roleIDs, h.jwtSecret, h.jwtExpire)
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: 实现 token 黑名单机制
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

	// 获取当前用户信息
	user, err := h.userService.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// 验证旧密码
	if !checkPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "old password is incorrect"})
		return
	}

	err = h.userService.UpdatePassword(c.Request.Context(), userID.(uint), req.NewPassword)
	if err != nil {
		h.logger.Error("failed to change password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// checkPassword 验证密码
func checkPassword(stored, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(input))
	return err == nil
}
