package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

type HostTemplateHandler struct {
	Service *service.HostTemplateService
	Logger  *zap.Logger
}

func NewHostTemplateHandler(service *service.HostTemplateService) *HostTemplateHandler {
	logger, _ := zap.NewProduction()
	return &HostTemplateHandler{Service: service, Logger: logger}
}

// CreateHostTemplate godoc
// @Summary 创建主机模版
// @Description 创建一个新的主机模版
// @Tags HostTemplate
// @Accept json
// @Produce json
// @Param request body cloudprovider.HostTemplate true "主机模版信息"
// @Success 200 {object} model.HostTemplate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /host-templates [post]
func (h *HostTemplateHandler) CreateHostTemplate(c *gin.Context) {
	var req cloudprovider.HostTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 设置默认状态为 Draft
	if req.Status == "" {
		req.Status = "Draft"
	}

	ctx := c.Request.Context()

	// 验证请求参数
	if err := h.Service.ValidateHostTemplateConfig(ctx, &req); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	hostTemplate, err := h.Service.CreateHostTemplate(ctx, &req)
	if err != nil {
		h.Logger.Error("Failed to create host template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create host template: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hostTemplate)
}

// GetHostTemplate godoc
// @Summary 获取主机模版详情
// @Description 根据ID获取主机模版详情
// @Tags HostTemplate
// @Produce json
// @Param id path string true "主机模版ID"
// @Success 200 {object} model.HostTemplate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /host-templates/{id} [get]
func (h *HostTemplateHandler) GetHostTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host template ID is required"})
		return
	}

	ctx := c.Request.Context()

	hostTemplate, err := h.Service.GetHostTemplateByID(ctx, id)
	if err != nil {
		h.Logger.Error("Failed to get host template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get host template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, hostTemplate)
}

// ListHostTemplates godoc
// @Summary 获取主机模版列表
// @Description 获取主机模版列表，支持分页
// @Tags HostTemplate
// @Produce json
// @Param page query int false "页码，默认为1"
// @Param page_size query int false "每页数量，默认为10"
// @Param project_id query string false "项目ID过滤"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /host-templates [get]
func (h *HostTemplateHandler) ListHostTemplates(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	projectID := c.Query("project_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	ctx := c.Request.Context()

	hostTemplates, total, err := h.Service.ListHostTemplates(ctx, projectID, page, pageSize)
	if err != nil {
		h.Logger.Error("Failed to list host templates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list host templates: " + err.Error()})
		return
	}

	pagination := map[string]interface{}{
		"page":       page,
		"page_size":  pageSize,
		"total":      total,
		"total_page": int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      hostTemplates,
		"pagination": pagination,
	})
}

// UpdateHostTemplate godoc
// @Summary 更新主机模版
// @Description 根据ID更新主机模版信息
// @Tags HostTemplate
// @Accept json
// @Produce json
// @Param id path string true "主机模版ID"
// @Param request body cloudprovider.HostTemplate true "主机模版信息"
// @Success 200 {object} model.HostTemplate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /host-templates/{id} [put]
func (h *HostTemplateHandler) UpdateHostTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host template ID is required"})
		return
	}

	var req cloudprovider.HostTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证请求参数
	if err := h.Service.ValidateHostTemplateConfig(ctx, &req); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	hostTemplate, err := h.Service.UpdateHostTemplate(ctx, id, &req)
	if err != nil {
		h.Logger.Error("Failed to update host template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update host template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, hostTemplate)
}

// DeleteHostTemplate godoc
// @Summary 删除主机模版
// @Description 根据ID删除主机模版
// @Tags HostTemplate
// @Produce json
// @Param id path string true "主机模版ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /host-templates/{id} [delete]
func (h *HostTemplateHandler) DeleteHostTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host template ID is required"})
		return
	}

	ctx := c.Request.Context()

	err := h.Service.DeleteHostTemplate(ctx, id)
	if err != nil {
		h.Logger.Error("Failed to delete host template", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete host template: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Host template deleted successfully"})
}