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

	// 获取筛选参数
	keyword := c.Query("keyword")
	roleType := c.Query("type")
	enabledStr := c.Query("enabled")

	var enabled *bool
	if enabledStr != "" {
		if enabledStr == "true" {
			enabled = func() *bool { b := true; return &b }()
		} else if enabledStr == "false" {
			enabled = func() *bool { b := false; return &b }()
		}
	}

	roles, total, err := h.service.ListRoles(c.Request.Context(), domainID, keyword, roleType, enabled, limit, offset)
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

// Enable 启用角色
func (h *RoleHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnableRole(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用角色
func (h *RoleHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisableRole(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// MakePublic 公开角色
func (h *RoleHandler) MakePublic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.MakeRolePublic(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to make role public", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "public"})
}

// GetRoleUsers 获取角色的用户列表
func (h *RoleHandler) GetRoleUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.service.GetRoleUsers(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get role users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// GetRoleGroups 获取角色的用户组列表
func (h *RoleHandler) GetRoleGroups(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	groups, total, err := h.service.GetRoleGroups(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get role groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"total": total,
	})
}

// GetRolePolicies 获取角色的策略列表
func (h *RoleHandler) GetRolePolicies(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	policies, err := h.service.GetRolePolicies(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get role policies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": policies,
		"total": len(policies),
	})
}

// AssignPolicyToRole 分配策略给角色
func (h *RoleHandler) AssignPolicyToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		PolicyID string `json:"policy_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignPolicyToRole(c.Request.Context(), uint(roleID), req.PolicyID); err != nil {
		h.logger.Error("failed to assign policy to role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// RevokePolicyFromRole 从角色撤销策略
func (h *RoleHandler) RevokePolicyFromRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	policyID := c.Query("policy_id")
	if policyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "policy_id is required"})
		return
	}

	if err := h.service.RevokePolicyFromRole(c.Request.Context(), uint(roleID), policyID); err != nil {
		h.logger.Error("failed to revoke policy from role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}
