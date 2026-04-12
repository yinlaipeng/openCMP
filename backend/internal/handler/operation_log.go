package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/utils/pagination"
)

type OperationLogHandler struct {
	Service *service.OperationLogService
}

// NewOperationLogHandler creates a new instance of OperationLogHandler
func NewOperationLogHandler(svc *service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{
		Service: svc,
	}
}

// GetOperationLogs retrieves operation logs with pagination and filtering
// @Summary Get operation logs
// @Description Retrieve operation logs with optional filtering and pagination
// @Tags OperationLogs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default 1
// @Param limit query int false "Page size" default 20
// @Param resource_name query string false "Resource name filter"
// @Param resource_type query string false "Resource type filter"
// @Param operation_type query string false "Operation type filter"
// @Param service_type query string false "Service type filter"
// @Param risk_level query string false "Risk level filter"
// @Param result query string false "Result filter"
// @Param operator query string false "Operator filter"
// @Success 200 {object} pagination.Pagination
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-logs [get]
func (h *OperationLogHandler) GetOperationLogs(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Prepare filter
	filter := make(map[string]interface{})
	if resource_name := c.Query("resource_name"); resource_name != "" {
		filter["resource_name"] = resource_name
	}
	if resource_type := c.Query("resource_type"); resource_type != "" {
		filter["resource_type"] = resource_type
	}
	if operation_type := c.Query("operation_type"); operation_type != "" {
		filter["operation_type"] = operation_type
	}
	if service_type := c.Query("service_type"); service_type != "" {
		filter["service_type"] = service_type
	}
	if risk_level := c.Query("risk_level"); risk_level != "" {
		filter["risk_level"] = risk_level
	}
	if result := c.Query("result"); result != "" {
		filter["result"] = result
	}
	if operator := c.Query("operator"); operator != "" {
		filter["operator"] = operator
	}
	if project_id_str := c.Query("project_id"); project_id_str != "" {
		if project_id, err := strconv.ParseUint(project_id_str, 10, 32); err == nil {
			filter["project_id"] = uint(project_id)
		}
	}
	if domain_id_str := c.Query("domain_id"); domain_id_str != "" {
		if domain_id, err := strconv.ParseUint(domain_id_str, 10, 32); err == nil {
			filter["domain_id"] = uint(domain_id)
		}
	}
	if user_id_str := c.Query("user_id"); user_id_str != "" {
		if user_id, err := strconv.ParseUint(user_id_str, 10, 32); err == nil {
			filter["user_id"] = uint(user_id)
		}
	}

	pg := &pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	result, err := h.Service.GetOperationLogs(filter, pg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetResourceOperationLogs retrieves operation logs for a specific resource
// @Summary Get operation logs for a specific resource
// @Description Retrieve operation logs for a specific resource type and ID
// @Tags OperationLogs
// @Accept json
// @Produce json
// @Param resource_type path string true "Resource type (e.g., auth_source, user, project, etc.)"
// @Param resource_id path int true "Resource ID"
// @Param page query int false "Page number" default 1
// @Param limit query int false "Page size" default 20
// @Success 200 {object} pagination.Pagination
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operation-logs/{resource_type}/{resource_id} [get]
func (h *OperationLogHandler) GetResourceOperationLogs(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceIDStr := c.Param("resource_id")

	resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	pg := &pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	result, err := h.Service.GetResourceOperationLogs(resourceType, uint(resourceID), pg)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resource operation logs not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
