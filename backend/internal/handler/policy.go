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

// PolicyHandler 策略 Handler
type PolicyHandler struct {
	service *service.PolicyService
	logger  *zap.Logger
}

// NewPolicyHandler 创建策略 Handler
func NewPolicyHandler(db *gorm.DB, logger *zap.Logger) *PolicyHandler {
	return &PolicyHandler{
		service: service.NewPolicyService(db),
		logger:  logger,
	}
}

// ListPoliciesRequest 列出策略请求
type ListPoliciesRequest struct {
	Scope          string `form:"scope"`
	DomainID       string `form:"domain_id"`
	ShowFailReason bool   `form:"show_fail_reason"`
	Details        bool   `form:"details"`
	SummaryStats   bool   `form:"summary_stats"`
	Limit          int    `form:"limit"`
	Offset         int    `form:"offset"`
}

// List 列出策略
func (h *PolicyHandler) List(c *gin.Context) {
	var req ListPoliciesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	var domainID *string
	if req.DomainID != "" {
		domainID = &req.DomainID
	}

	policies, total, err := h.service.ListPolicies(c.Request.Context(), req.Scope, domainID, req.Limit, req.Offset)
	if err != nil {
		h.logger.Error("failed to list policies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   policies,
		"limit":  req.Limit,
		"offset": req.Offset,
		"total":  total,
	})
}

// GetPolicy 获取策略详情
func (h *PolicyHandler) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	policy, err := h.service.GetPolicy(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if policy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// CreatePolicy 创建策略
func (h *PolicyHandler) CreatePolicy(c *gin.Context) {
	var req model.Policy
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if req.Scope == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scope is required"})
		return
	}
	if req.Policy == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "policy content is required"})
		return
	}

	// 检查策略名称是否已存在
	exists, err := h.service.CheckPolicyExists(c.Request.Context(), req.Name)
	if err != nil {
		h.logger.Error("failed to check policy exists", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "policy name already exists"})
		return
	}

	if err := h.service.CreatePolicy(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to create policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算权限（是否可删除/更新）
	req.CanDelete = !req.IsSystem
	req.CanUpdate = !req.IsSystem
	c.JSON(http.StatusCreated, req)
}

// UpdatePolicy 更新策略
func (h *PolicyHandler) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePolicy(c.Request.Context(), id, req); err != nil {
		h.logger.Error("failed to update policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DeletePolicy 删除策略
func (h *PolicyHandler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	if err := h.service.DeletePolicy(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete policy", zap.Error(err))
		// 根据错误类型返回不同的状态码
		if err.Error() == "不可删除系统策略" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// AssignPolicyToRole 分配策略给角色
func (h *PolicyHandler) AssignPolicyToRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
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

// GetRolePolicies 获取角色的策略列表
func (h *PolicyHandler) GetRolePolicies(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	policies, err := h.service.GetRolePolicies(c.Request.Context(), uint(roleID))
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

// RevokePolicyFromRole 从角色撤销策略
func (h *PolicyHandler) RevokePolicyFromRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
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
