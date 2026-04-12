package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/model"
	"github.com/opencmp/opencmp/internal/service"
)

// UserHandler 用户 Handler
type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

// NewUserHandler 创建用户 Handler
func NewUserHandler(db *gorm.DB, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service.NewUserService(db),
		logger:  logger,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name         string `json:"name" binding:"required"`
	DisplayName  string `json:"display_name"`
	Remark       string `json:"remark"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password" binding:"required,min=8"`
	DomainID     uint   `json:"domain_id" binding:"required"`
	ConsoleLogin *bool  `json:"console_login,omitempty"`
	MFAEnabled   *bool  `json:"mfa_enabled,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	DisplayName  string `json:"display_name"`
	Remark       string `json:"remark"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	ConsoleLogin *bool  `json:"console_login,omitempty"`
	MFAEnabled   *bool  `json:"mfa_enabled,omitempty"`
}

// List 列出用户
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var domainID *uint
	if d := c.Query("domain_id"); d != "" {
		id, _ := strconv.ParseUint(d, 10, 32)
		uid := uint(id)
		domainID = &uid
	}

	keyword := c.Query("keyword")
	email := c.Query("email")
	enabledStr := c.Query("enabled")

	var enabled *bool
	if enabledStr != "" {
		if enabledStr == "true" {
			enabled = func() *bool { b := true; return &b }()
		} else if enabledStr == "false" {
			enabled = func() *bool { b := false; return &b }()
		}
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), domainID, keyword, email, enabled, limit, offset)
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// Get 获取用户详情
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Remark:       req.Remark,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password, // TODO: 加密
		DomainID:     req.DomainID,
		Enabled:      true,
		ConsoleLogin: true, // Default to allowing console login
	}

	// Set console login if provided in request
	if req.ConsoleLogin != nil {
		user.ConsoleLogin = *req.ConsoleLogin
	}

	// Set MFA enabled if provided in request
	if req.MFAEnabled != nil {
		user.MFAEnabled = *req.MFAEnabled
	}

	if err := h.service.CreateUser(c.Request.Context(), user); err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用用户
func (h *UserHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnableUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用用户
func (h *UserHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisableUser(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// Update 更新用户信息
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.DisplayName = req.DisplayName
	user.Remark = req.Remark
	user.Email = req.Email
	user.Phone = req.Phone

	// Update console login if provided
	if req.ConsoleLogin != nil {
		user.ConsoleLogin = *req.ConsoleLogin
	}

	// Update MFA enabled if provided
	if req.MFAEnabled != nil {
		user.MFAEnabled = *req.MFAEnabled
	}

	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		h.logger.Error("failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ResetPassword 重置用户密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetUserPassword(c.Request.Context(), uint(id), req.Password); err != nil {
		h.logger.Error("failed to reset user password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

// GetRoles 获取用户角色
func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	domainID := c.Query("domain_id")
	var did uint = 1 // 默认域
	if domainID != "" {
		did64, _ := strconv.ParseUint(domainID, 10, 32)
		did = uint(did64)
	}

	roles, err := h.service.GetUserRoles(c.Request.Context(), uint(id), did)
	if err != nil {
		h.logger.Error("failed to get user roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": roles,
		"total": len(roles),
	})
}

// AssignRole 分配角色给用户
func (h *UserHandler) AssignRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		RoleID   uint `json:"role_id" binding:"required"`
		DomainID uint `json:"domain_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignUserRole(c.Request.Context(), uint(id), req.RoleID, req.DomainID); err != nil {
		h.logger.Error("failed to assign role to user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// RevokeRole 撤销用户角色
func (h *UserHandler) RevokeRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	roleID := c.Query("role_id")
	domainID := c.Query("domain_id")
	if roleID == "" || domainID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id and domain_id are required"})
		return
	}

	rid64, _ := strconv.ParseUint(roleID, 10, 32)
	did64, _ := strconv.ParseUint(domainID, 10, 32)

	if err := h.service.RevokeUserRole(c.Request.Context(), uint(id), uint(rid64), uint(did64)); err != nil {
		h.logger.Error("failed to revoke role from user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}

// GetGroups 获取用户组
func (h *UserHandler) GetGroups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	groups, err := h.service.GetUserGroups(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get user groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"total": len(groups),
	})
}

// JoinGroup 加入用户组
func (h *UserHandler) JoinGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		GroupID uint `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddUserToGroup(c.Request.Context(), uint(id), req.GroupID); err != nil {
		h.logger.Error("failed to add user to group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined"})
}

// LeaveGroup 离开用户组
func (h *UserHandler) LeaveGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	groupID := c.Query("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id is required"})
		return
	}

	gid, _ := strconv.ParseUint(groupID, 10, 32)

	if err := h.service.RemoveUserFromGroup(c.Request.Context(), uint(id), uint(gid)); err != nil {
		h.logger.Error("failed to remove user from group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left"})
}

// AssignUserToProjectRequest 将用户分配到项目请求
type AssignUserToProjectRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
	RoleID    uint `json:"role_id" binding:"required"`
}

// AssignUserToProject 将用户分配到项目
func (h *UserHandler) AssignUserToProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req AssignUserToProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignUserToProject(c.Request.Context(), uint(id), req.ProjectID, req.RoleID); err != nil {
		h.logger.Error("failed to assign user to project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned to project"})
}

// RemoveUserFromProject 从项目移除用户
func (h *UserHandler) RemoveUserFromProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	pid, _ := strconv.ParseUint(projectID, 10, 32)

	if err := h.service.RemoveUserFromProject(c.Request.Context(), uint(id), uint(pid)); err != nil {
		h.logger.Error("failed to remove user from project", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed from project"})
}

// GetUserProjects 获取用户所属的项目
func (h *UserHandler) GetUserProjects(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	projects, err := h.service.GetUserProjects(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get user projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": projects,
		"total": len(projects),
	})
}
