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
	db      *gorm.DB
	logger  *zap.Logger
}

// NewDomainHandler 创建域 Handler
func NewDomainHandler(db *gorm.DB, logger *zap.Logger) *DomainHandler {
	return &DomainHandler{
		service: service.NewDomainService(db),
		db:      db,
		logger:  logger,
	}
}

// DomainResponse 域响应结构（带统计字段）
type DomainResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	ParentID    *uint  `json:"parent_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// 统计字段
	UserCount    int  `json:"user_count"`
	GroupCount   int  `json:"group_count"`
	ProjectCount int  `json:"project_count"`
	RoleCount    int  `json:"role_count"`
	PolicyCount  int  `json:"policy_count"`
	IdpCount     int  `json:"idp_count"`
	IsSSO        bool `json:"is_sso"`
	CanDelete    bool `json:"can_delete"`
	CanUpdate    bool `json:"can_update"`
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
	details := c.Query("details") == "true"

	// Get filter parameters
	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if description := c.Query("description"); description != "" {
		filters["description"] = description
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

	// 如果需要详细信息，添加统计字段
	if details {
		var responses []DomainResponse
		for _, d := range domains {
			resp := h.buildDomainResponse(d)
			responses = append(responses, resp)
		}
		c.JSON(http.StatusOK, gin.H{
			"items": responses,
			"total": total,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": domains,
		"total": total,
	})
}

// buildDomainResponse 构建域响应（带统计字段）
func (h *DomainHandler) buildDomainResponse(d *service.DomainWithAuthSourceCount) DomainResponse {
	resp := DomainResponse{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Enabled:     d.Enabled,
		ParentID:    d.ParentID,
		CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		IsSSO:       false, // 暂不支持 SSO
		CanDelete:   !h.isDefaultDomainByID(d.ID),
		CanUpdate:   !h.isDefaultDomainByID(d.ID),
	}

	// 计算统计字段
	resp.UserCount = h.countUsers(d.ID)
	resp.GroupCount = h.countGroups(d.ID)
	resp.ProjectCount = h.countProjects(d.ID)
	resp.RoleCount = h.countRoles(d.ID)
	resp.PolicyCount = h.countPolicies(d.ID)
	resp.IdpCount = d.AuthSourceCount // 使用已有的认证源数量

	return resp
}

// isDefaultDomainByID 判断是否为默认域（根据ID）
func (h *DomainHandler) isDefaultDomainByID(id uint) bool {
	return id == 1
}

// isDefaultDomain 判断是否为默认域（根据名称）
func (h *DomainHandler) isDefaultDomain(d *model.Domain) bool {
	return d.Name == "default" || d.Name == "Default" || d.ID == 1
}

// countUsers 计算域用户数量
func (h *DomainHandler) countUsers(domainID uint) int {
	var count int64
	h.db.Model(&model.User{}).Where("domain_id = ?", domainID).Count(&count)
	return int(count)
}

// countGroups 计算域组数量
func (h *DomainHandler) countGroups(domainID uint) int {
	var count int64
	h.db.Model(&model.Group{}).Where("domain_id = ?", domainID).Count(&count)
	return int(count)
}

// countProjects 计算域项目数量
func (h *DomainHandler) countProjects(domainID uint) int {
	var count int64
	h.db.Model(&model.Project{}).Where("domain_id = ?", domainID).Count(&count)
	return int(count)
}

// countRoles 计算域角色数量
func (h *DomainHandler) countRoles(domainID uint) int {
	var count int64
	h.db.Model(&model.Role{}).Where("domain_id = ?", domainID).Count(&count)
	return int(count)
}

// countPolicies 计算域策略数量
func (h *DomainHandler) countPolicies(domainID uint) int {
	// 系统策略对所有域可见，这里计算系统策略 + 域策略
	var systemCount, domainCount int64
	h.db.Model(&model.Policy{}).Where("scope = ?", "system").Count(&systemCount)
	h.db.Model(&model.Policy{}).Where("domain_id = ? OR scope = ?", domainID, "domain").Count(&domainCount)
	return int(systemCount + domainCount)
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
