package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

type AutoscalingGroupHandler struct {
	Service *service.AutoscalingGroupService
	Logger  *zap.Logger
}

func NewAutoscalingGroupHandler(service *service.AutoscalingGroupService) *AutoscalingGroupHandler {
	logger, _ := zap.NewProduction()
	return &AutoscalingGroupHandler{Service: service, Logger: logger}
}

// CreateAutoscalingGroup godoc
// @Summary 创建弹性伸缩组
// @Description 创建一个新的弹性伸缩组
// @Tags AutoscalingGroup
// @Accept json
// @Produce json
// @Param request body cloudprovider.AutoscalingGroup true "弹性伸缩组信息"
// @Success 200 {object} model.AutoscalingGroup
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /autoscaling-groups [post]
func (h *AutoscalingGroupHandler) CreateAutoscalingGroup(c *gin.Context) {
	var req cloudprovider.AutoscalingGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 设置默认状态为 Inactive
	if req.Status == "" {
		req.Status = "Inactive"
	}

	ctx := c.Request.Context()

	// 验证请求参数
	if err := h.Service.ValidateAutoscalingGroupConfig(ctx, &req); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	autoscalingGroup, err := h.Service.CreateAutoscalingGroup(ctx, &req)
	if err != nil {
		h.Logger.Error("Failed to create autoscaling group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create autoscaling group: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, autoscalingGroup)
}

// GetAutoscalingGroup godoc
// @Summary 获取弹性伸缩组详情
// @Description 根据ID获取弹性伸缩组详情
// @Tags AutoscalingGroup
// @Produce json
// @Param id path string true "弹性伸缩组ID"
// @Success 200 {object} model.AutoscalingGroup
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /autoscaling-groups/{id} [get]
func (h *AutoscalingGroupHandler) GetAutoscalingGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Autoscaling group ID is required"})
		return
	}

	ctx := c.Request.Context()

	autoscalingGroup, err := h.Service.GetAutoscalingGroupByID(ctx, id)
	if err != nil {
		h.Logger.Error("Failed to get autoscaling group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get autoscaling group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, autoscalingGroup)
}

// ListAutoscalingGroups godoc
// @Summary 获取弹性伸缩组列表
// @Description 获取弹性伸缩组列表，支持分页
// @Tags AutoscalingGroup
// @Produce json
// @Param page query int false "页码，默认为1"
// @Param page_size query int false "每页数量，默认为10"
// @Param project_id query string false "项目ID过滤"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /autoscaling-groups [get]
func (h *AutoscalingGroupHandler) ListAutoscalingGroups(c *gin.Context) {
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

	autoscalingGroups, total, err := h.Service.ListAutoscalingGroups(ctx, projectID, page, pageSize)
	if err != nil {
		h.Logger.Error("Failed to list autoscaling groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list autoscaling groups: " + err.Error()})
		return
	}

	pagination := map[string]interface{}{
		"page":       page,
		"page_size":  pageSize,
		"total":      total,
		"total_page": int((total + int64(pageSize) - 1) / int64(pageSize)),
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      autoscalingGroups,
		"pagination": pagination,
	})
}

// UpdateAutoscalingGroup godoc
// @Summary 更新弹性伸缩组
// @Description 根据ID更新弹性伸缩组信息
// @Tags AutoscalingGroup
// @Accept json
// @Produce json
// @Param id path string true "弹性伸缩组ID"
// @Param request body cloudprovider.AutoscalingGroup true "弹性伸缩组信息"
// @Success 200 {object} model.AutoscalingGroup
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /autoscaling-groups/{id} [put]
func (h *AutoscalingGroupHandler) UpdateAutoscalingGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Autoscaling group ID is required"})
		return
	}

	var req cloudprovider.AutoscalingGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证请求参数
	if err := h.Service.ValidateAutoscalingGroupConfig(ctx, &req); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	autoscalingGroup, err := h.Service.UpdateAutoscalingGroup(ctx, id, &req)
	if err != nil {
		h.Logger.Error("Failed to update autoscaling group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update autoscaling group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, autoscalingGroup)
}

// DeleteAutoscalingGroup godoc
// @Summary 删除弹性伸缩组
// @Description 根据ID删除弹性伸缩组
// @Tags AutoscalingGroup
// @Produce json
// @Param id path string true "弹性伸缩组ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /autoscaling-groups/{id} [delete]
func (h *AutoscalingGroupHandler) DeleteAutoscalingGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Autoscaling group ID is required"})
		return
	}

	ctx := c.Request.Context()

	err := h.Service.DeleteAutoscalingGroup(ctx, id)
	if err != nil {
		h.Logger.Error("Failed to delete autoscaling group", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete autoscaling group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Autoscaling group deleted successfully"})
}