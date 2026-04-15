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
	userService       *service.UserService
	authSourceService *service.AuthSourceService
	logger            *zap.Logger
	jwtSecret         string
	jwtExpire         int
}

// NewAuthHandler 创建认证 Handler
func NewAuthHandler(db *gorm.DB, logger *zap.Logger, jwtSecret string, jwtExpire int) *AuthHandler {
	return &AuthHandler{
		userService:       service.NewUserService(db),
		authSourceService: service.NewAuthSourceService(db),
		logger:            logger,
		jwtSecret:         jwtSecret,
		jwtExpire:         jwtExpire,
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

	ctx := c.Request.Context()

	// 统一认证入口：先本地，后 LDAP
	user, authSource, err := h.authSourceService.AuthenticateUser(ctx, req.Username, req.Password, nil)
	if err != nil {
		// 区分用户禁用和密码错误
		if err.Error() == "user is disabled" {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 获取用户的角色
	roleIDs, err := h.userService.GetUserRoleIDs(ctx, user.ID)
	if err != nil {
		h.logger.Error("failed to get user roles", zap.Error(err))
		// 如果无法获取角色，继续登录但不包含角色信息
		roleIDs = []uint{}
	}

	// 更新最后登录信息
	h.userService.UpdateLastLogin(ctx, user.ID, c.ClientIP())

	// 生成 JWT token
	// 如果是通过认证源（LDAP 等）登录，在 token 中携带 auth_source_id 和 domain_id
	var token string
	if authSource != nil {
		token, err = utils.GenerateTokenWithExtra(user.ID, user.Name, roleIDs, user.DomainID, authSource.ID, h.jwtSecret, h.jwtExpire)
	} else {
		token, err = utils.GenerateToken(user.ID, user.Name, roleIDs, h.jwtSecret, h.jwtExpire)
	}
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

// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req struct {
		DisplayName  string `json:"display_name"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Remark       string `json:"remark"`
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

	// 更新用户信息
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Remark != "" {
		user.Remark = req.Remark
	}

	err = h.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		h.logger.Error("failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated", "user": user})
}
