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

// PermissionHandler 权限 Handler
type PermissionHandler struct {
	service *service.PermissionService
	logger  *zap.Logger
}

// NewPermissionHandler 创建权限 Handler
func NewPermissionHandler(db *gorm.DB, logger *zap.Logger) *PermissionHandler {
	return &PermissionHandler{
		service: service.NewPermissionService(db),
		logger:  logger,
	}
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Scope       string `json:"scope" binding:"required,oneof=system domain project"` // system/domain/project
	DomainID    *uint  `json:"domain_id"`                                            // 仅当 scope=domain 时需要
	Enabled     bool   `json:"enabled"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Enabled     bool   `json:"enabled"`
}

// List 权限列表
func (h *PermissionHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// 获取筛选参数
	keyword := c.Query("keyword")
	resource := c.Query("resource")
	action := c.Query("action")
	scope := c.Query("scope")
	enabledStr := c.Query("enabled")

	var domainID *uint
	if domainIDStr := c.Query("domain_id"); domainIDStr != "" {
		if id, err := strconv.ParseUint(domainIDStr, 10, 32); err == nil {
			domainID = func() *uint { id32 := uint(id); return &id32 }()
		}
	}

	var enabled *bool
	if enabledStr != "" {
		if enabledStr == "true" {
			enabled = func() *bool { b := true; return &b }()
		} else if enabledStr == "false" {
			enabled = func() *bool { b := false; return &b }()
		}
	}

	permissions, total, err := h.service.ListPermissions(c.Request.Context(), domainID, keyword, resource, action, scope, enabled, limit, offset)
	if err != nil {
		h.logger.Error("failed to list permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": permissions,
		"total": total,
	})
}

// Get 获取权限详情
func (h *PermissionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	permission, err := h.service.GetPermission(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if permission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// Create 创建权限
func (h *PermissionHandler) Create(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission := &model.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		Scope:       req.Scope,
		DomainID:    req.DomainID,
		Enabled:     req.Enabled,
	}

	if err := h.service.CreatePermission(c.Request.Context(), permission); err != nil {
		h.logger.Error("failed to create permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// Update 更新权限
func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	permission, err := h.service.GetPermission(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if permission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permission not found"})
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission.Name = req.Name
	permission.Description = req.Description
	permission.Resource = req.Resource
	permission.Action = req.Action
	permission.Enabled = req.Enabled

	if err := h.service.UpdatePermission(c.Request.Context(), permission); err != nil {
		h.logger.Error("failed to update permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// Delete 删除权限
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeletePermission(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用权限
func (h *PermissionHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnablePermission(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用权限
func (h *PermissionHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisablePermission(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// AssignToRole 分配权限给角色
func (h *PermissionHandler) AssignToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	if err := h.service.AssignPermissionToRole(c.Request.Context(), uint(roleID), uint(permissionID)); err != nil {
		h.logger.Error("failed to assign permission to role", zap.Error(err))
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "permission already assigned to role"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission assigned to role"})
}

// RemoveFromRole 从角色移除权限
func (h *PermissionHandler) RemoveFromRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	if err := h.service.RevokePermissionFromRole(c.Request.Context(), uint(roleID), uint(permissionID)); err != nil {
		h.logger.Error("failed to revoke permission from role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission removed from role"})
}

// GetRolePermissions 获取角色的权限列表
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	permissions, err := h.service.GetRolePermissions(c.Request.Context(), uint(roleID))
	if err != nil {
		h.logger.Error("failed to get role permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": permissions,
	})
}

// GetPermissionsForResourceAction 获取指定资源和操作的权限
func (h *PermissionHandler) GetPermissionsForResourceAction(c *gin.Context) {
	resource := c.Query("resource")
	action := c.Query("action")

	if resource == "" || action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource and action are required"})
		return
	}

	permissions, err := h.service.GetPermissionsForResourceAction(c.Request.Context(), resource, action)
	if err != nil {
		h.logger.Error("failed to get permissions for resource and action", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": permissions,
	})
}
