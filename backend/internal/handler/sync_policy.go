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

// SyncPolicyHandler 同步策略 Handler
type SyncPolicyHandler struct {
	service *service.SyncPolicyService
	logger  *zap.Logger
}

// NewSyncPolicyHandler 创建同步策略 Handler
func NewSyncPolicyHandler(db *gorm.DB, logger *zap.Logger) *SyncPolicyHandler {
	return &SyncPolicyHandler{
		service: service.NewSyncPolicyService(db),
		logger:  logger,
	}
}

// CreateSyncPolicyRequest 创建同步策略请求
type CreateSyncPolicyRequest struct {
	Name     string      `json:"name" binding:"required"`
	Remarks  string      `json:"remarks"`
	Status   string      `json:"status"`
	Enabled  bool        `json:"enabled"`
	Rules    []RuleData  `json:"rules" binding:"required"`
	Scope    string      `json:"scope" binding:"required"`
	DomainID uint        `json:"domain_id" binding:"required"`
}

// UpdateSyncPolicyRequest 更新同步策略请求
type UpdateSyncPolicyRequest struct {
	Name     string      `json:"name"`
	Remarks  string      `json:"remarks"`
	Status   string      `json:"status"`
	Enabled  bool        `json:"enabled"`
	Rules    []RuleData  `json:"rules"`
	Scope    string      `json:"scope"`
	DomainID uint        `json:"domain_id"`
}

// RuleData 规则数据结构
type RuleData struct {
	ID                uint         `json:"id,omitempty"`
	ConditionType     string       `json:"condition_type" binding:"required"`
	ResourceMapping   string       `json:"resource_mapping" binding:"required"`
	TargetProjectID   *uint        `json:"target_project_id"`
	TargetProjectName string       `json:"target_project_name"`
	Tags              []TagData    `json:"tags"`
}

// TagData 标签数据结构
type TagData struct {
	ID       uint   `json:"id,omitempty"`
	TagKey   string `json:"tag_key" binding:"required"`
	TagValue string `json:"tag_value" binding:"required"`
}

// Create 创建同步策略
func (h *SyncPolicyHandler) Create(c *gin.Context) {
	var req CreateSyncPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := &model.SyncPolicy{
		Name:     req.Name,
		Remarks:  req.Remarks,
		Status:   req.Status,
		Enabled:  req.Enabled,
		Scope:    req.Scope,
		DomainID: req.DomainID,
	}

	// 转换规则和标签数据
	var rules []model.Rule
	var ruleTags []model.RuleTag
	for _, ruleData := range req.Rules {
		rule := model.Rule{
			ID:                ruleData.ID,
			ConditionType:     ruleData.ConditionType,
			ResourceMapping:   ruleData.ResourceMapping,
			TargetProjectID:   ruleData.TargetProjectID,
			TargetProjectName: ruleData.TargetProjectName,
		}
		rules = append(rules, rule)

		// 转换标签数据
		for _, tagData := range ruleData.Tags {
			tag := model.RuleTag{
				ID:       tagData.ID,
				TagKey:   tagData.TagKey,
				TagValue: tagData.TagValue,
			}
			ruleTags = append(ruleTags, tag)
		}
	}

	if err := h.service.CreateSyncPolicy(c.Request.Context(), policy, rules, ruleTags); err != nil {
		h.logger.Error("failed to create sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, policy)
}

// List 列出同步策略
func (h *SyncPolicyHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	policies, total, err := h.service.ListSyncPolicies(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list sync policies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": policies,
		"total": total,
	})
}

// Get 获取同步策略详情
func (h *SyncPolicyHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	policy, err := h.service.GetSyncPolicy(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if policy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// Update 更新同步策略
func (h *SyncPolicyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	policy, err := h.service.GetSyncPolicy(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if policy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "policy not found"})
		return
	}

	var req UpdateSyncPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		policy.Name = req.Name
	}
	if req.Remarks != "" {
		policy.Remarks = req.Remarks
	}
	if req.Status != "" {
		policy.Status = req.Status
	}
	policy.Enabled = req.Enabled
	if req.Scope != "" {
		policy.Scope = req.Scope
	}
	if req.DomainID != 0 {
		policy.DomainID = req.DomainID
	}

	// 转换规则和标签数据
	var rules []model.Rule
	var ruleTags []model.RuleTag
	for _, ruleData := range req.Rules {
		rule := model.Rule{
			ID:                ruleData.ID,
			ConditionType:     ruleData.ConditionType,
			ResourceMapping:   ruleData.ResourceMapping,
			TargetProjectID:   ruleData.TargetProjectID,
			TargetProjectName: ruleData.TargetProjectName,
		}
		rules = append(rules, rule)

		// 转换标签数据
		for _, tagData := range ruleData.Tags {
			tag := model.RuleTag{
				ID:       tagData.ID,
				TagKey:   tagData.TagKey,
				TagValue: tagData.TagValue,
			}
			ruleTags = append(ruleTags, tag)
		}
	}

	if err := h.service.UpdateSyncPolicy(c.Request.Context(), policy, rules, ruleTags); err != nil {
		h.logger.Error("failed to update sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// Delete 删除同步策略
func (h *SyncPolicyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteSyncPolicy(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete sync policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// UpdateStatus 更新同步策略状态
func (h *SyncPolicyHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ToggleSyncPolicyStatus(c.Request.Context(), uint(id), req.Enabled); err != nil {
		h.logger.Error("failed to update sync policy status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}