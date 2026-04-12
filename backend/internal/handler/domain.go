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

// DomainHandler 域 Handler
type DomainHandler struct {
	service *service.DomainService
	logger  *zap.Logger
}

// NewDomainHandler 创建域 Handler
func NewDomainHandler(db *gorm.DB, logger *zap.Logger) *DomainHandler {
	return &DomainHandler{
		service: service.NewDomainService(db),
		logger:  logger,
	}
}

// CreateDomainRequest 创建域请求
type CreateDomainRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// List 列出域
func (h *DomainHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get filter parameters
	filters := make(map[string]interface{})
	if keyword := c.Query("keyword"); keyword != "" {
		filters["keyword"] = keyword
	}
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			filters["enabled"] = enabled
		}
	}

	domains, total, err := h.service.ListDomainsWithFilters(c.Request.Context(), filters, limit, offset)
	if err != nil {
		h.logger.Error("failed to list domains", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": domains,
		"total": total,
	})
}

// Get 获取域详情
func (h *DomainHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	domain, err := h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if domain == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// Create 创建域
func (h *DomainHandler) Create(c *gin.Context) {
	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain := &model.Domain{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}

	if err := h.service.CreateDomain(c.Request.Context(), domain); err != nil {
		h.logger.Error("failed to create domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

// Update 更新域
func (h *DomainHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	domain, err := h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if domain == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "domain not found"})
		return
	}

	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain.Name = req.Name
	domain.Description = req.Description
	domain.Enabled = req.Enabled

	if err := h.service.UpdateDomain(c.Request.Context(), domain); err != nil {
		h.logger.Error("failed to update domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain)
}

// Delete 删除域
func (h *DomainHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteDomain(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Enable 启用域
func (h *DomainHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.EnableDomain(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to enable domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

// Disable 禁用域
func (h *DomainHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DisableDomain(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to disable domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}

// GetUsers 获取域的用户列表
func (h *DomainHandler) GetUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain id"})
		return
	}

	// 检查域是否存在
	_, err = h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.service.GetDomainUsers(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get domain users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
	})
}

// GetProjects 获取域的项目列表
func (h *DomainHandler) GetProjects(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain id"})
		return
	}

	// 检查域是否存在
	_, err = h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, total, err := h.service.GetDomainProjects(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get domain projects", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": projects,
		"total": total,
	})
}

// GetRoles 获取域的角色列表
func (h *DomainHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain id"})
		return
	}

	// 检查域是否存在
	_, err = h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	roles, total, err := h.service.GetDomainRoles(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get domain roles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": roles,
		"total": total,
	})
}

// GetCloudAccounts 获取域的云账号列表
func (h *DomainHandler) GetCloudAccounts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain id"})
		return
	}

	// 检查域是否存在
	_, err = h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 当前云账号表没有直接关联到域的字段
	// 可能通过项目间接关联，或需要扩展云账号模型增加域字段
	// 暂时返回空列表，可根据实际需求调整
	var cloudAccounts []*model.CloudAccount

	c.JSON(http.StatusOK, gin.H{
		"items": cloudAccounts,
		"total": len(cloudAccounts),
	})
}

// GetOperationLogs 获取域的操作日志列表
func (h *DomainHandler) GetOperationLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain id"})
		return
	}

	// 检查域是否存在
	_, err = h.service.GetDomain(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get domain", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取与该域相关的操作日志
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	operationLogs, total, err := h.service.GetDomainOperationLogs(c.Request.Context(), uint(id), limit, offset)
	if err != nil {
		h.logger.Error("failed to get operation logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": operationLogs,
		"total": total,
	})
}
