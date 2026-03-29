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

// RoleHandler 角色 Handler
type RoleHandler struct {
	service *service.RoleService
	logger  *zap.Logger
}

// NewRoleHandler 创建角色 Handler
func NewRoleHandler(db *gorm.DB, logger *zap.Logger) *RoleHandler {
	return &RoleHandler{
		service: service.NewRoleService(db),
		logger:  logger,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	DomainID    uint   `json:"domain_id"`
	Type        string `json:"type"`
}

// List 列出角色
func (h *RoleHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var domainID *uint
	if d := c.Query("domain_id"); d != "" {
		id, _ := strconv.ParseUint(d, 10, 32)
		uid := uint(id)
		domainID = &uid
	}

	roles, total, err := h.service.ListRoles(c.Request.Context(), domainID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": roles,
		"total": total,
	})
}

// Get 获取角色详情
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Create 创建角色
func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := &model.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		DomainID:    req.DomainID,
		Type:        req.Type,
		Enabled:     true,
	}

	if err := h.service.CreateRole(c.Request.Context(), role); err != nil {
		h.logger.Error("failed to create role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// Update 更新角色
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role.Name = req.Name
	role.DisplayName = req.DisplayName
	role.Description = req.Description
	role.DomainID = req.DomainID
	role.Type = req.Type

	if err := h.service.UpdateRole(c.Request.Context(), role); err != nil {
		h.logger.Error("failed to update role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Delete 删除角色
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ListPermissions 列出权限（支持筛选、搜索）
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resource := c.Query("resource")
	action := c.Query("action")
	permissionType := c.Query("type")
	keyword := c.Query("keyword")

	permissions, total, err := h.service.ListPermissions(c.Request.Context(), resource, action, permissionType, keyword, limit, offset)
	if err != nil {
		h.logger.Error("failed to list permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": permissions,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetPermission 获取权限详情
func (h *RoleHandler) GetPermission(c *gin.Context) {
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

// CreatePermission 创建权限
func (h *RoleHandler) CreatePermission(c *gin.Context) {
	var req model.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if req.Resource == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource is required"})
		return
	}
	if req.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action is required"})
		return
	}

	if err := h.service.CreatePermission(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to create permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// UpdatePermission 更新权限
func (h *RoleHandler) UpdatePermission(c *gin.Context) {
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

	// 系统权限不可更新
	if permission.Type == "system" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不可更新系统权限"})
		return
	}

	var req model.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission.DisplayName = req.DisplayName
	permission.Description = req.Description

	if err := h.service.UpdatePermission(c.Request.Context(), permission); err != nil {
		h.logger.Error("failed to update permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// DeletePermission 删除权限
func (h *RoleHandler) DeletePermission(c *gin.Context) {
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

	// 系统权限不可删除
	if permission.Type == "system" {
		c.JSON(http.StatusForbidden, gin.H{"error": "不可删除系统权限"})
		return
	}

	if err := h.service.DeletePermission(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete permission", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "权限已被使用，无法删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// AssignPermission 分配权限给角色
func (h *RoleHandler) AssignPermission(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignPermissionToRole(c.Request.Context(), uint(roleID), req.PermissionID); err != nil {
		h.logger.Error("failed to assign permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// GetPermissions 获取角色权限
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	permissions, err := h.service.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get role permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": permissions,
		"total": len(permissions),
	})
}

// RevokePermission 撤销角色权限
func (h *RoleHandler) RevokePermission(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	permissionID := c.Query("permission_id")
	if permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "permission_id is required"})
		return
	}

	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission_id"})
		return
	}

	if err := h.service.RevokePermissionFromRole(c.Request.Context(), uint(roleID), uint(pid)); err != nil {
		h.logger.Error("failed to revoke permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}

// ListResources 列出资源类型
func (h *RoleHandler) ListResources(c *gin.Context) {
	resources, err := h.service.GetPermissionResources(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list resources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resources,
		"total": len(resources),
	})
}

// ListActions 列出操作类型
func (h *RoleHandler) ListActions(c *gin.Context) {
	actions, err := h.service.GetPermissionActions(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list actions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": actions,
		"total": len(actions),
	})
}
