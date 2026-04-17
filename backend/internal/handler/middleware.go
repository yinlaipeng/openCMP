package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opencmp/opencmp/internal/service"
	"github.com/opencmp/opencmp/pkg/cloudprovider"
)

// MiddlewareHandler 中间件资源 Handler
type MiddlewareHandler struct {
	service *service.MiddlewareService
	logger  *zap.Logger
}

// NewMiddlewareHandler 创建中间件资源 Handler
func NewMiddlewareHandler(db *gorm.DB, logger *zap.Logger) *MiddlewareHandler {
	return &MiddlewareHandler{
		service: service.NewMiddlewareService(db),
		logger:  logger,
	}
}

// CreateKafkaRequest 创建 Kafka 请求
type CreateKafkaRequest struct {
	AccountID     uint              `json:"account_id" binding:"required"`
	Name          string            `json:"name" binding:"required"`
	Version       string            `json:"version" binding:"required"`
	Storage       int               `json:"storage"`    // GB
	Bandwidth     int               `json:"bandwidth"`  // MB/s
	Retention     int               `json:"retention"`  // hours
	VPCID         string            `json:"vpc_id" binding:"required"`
	SubnetID      string            `json:"subnet_id" binding:"required"`
	ZoneID        string            `json:"zone_id"`
	Tags          map[string]string `json:"tags"`
}

// CreateKafka 创建 Kafka 实例
func (h *MiddlewareHandler) CreateKafka(c *gin.Context) {
	var req CreateKafkaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.KafkaConfig{
		Name:      req.Name,
		Version:   req.Version,
		Storage:   req.Storage,
		Bandwidth: req.Bandwidth,
		Retention: req.Retention,
		VPCID:     req.VPCID,
		SubnetID:  req.SubnetID,
		ZoneID:    req.ZoneID,
		Tags:      req.Tags,
	}

	instance, err := h.service.CreateKafka(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create kafka instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// ListKafka 列出 Kafka 实例
func (h *MiddlewareHandler) ListKafka(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	filter := cloudprovider.KafkaFilter{
		InstanceID: c.Query("instance_id"),
		Status:     c.Query("status"),
		Version:    c.Query("version"),
		RegionID:   c.Query("region_id"),
	}

	instances, err := h.service.ListKafka(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list kafka instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instances)
}

// GetKafka 获取 Kafka 实例详情
func (h *MiddlewareHandler) GetKafka(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	instance, err := h.service.GetKafka(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to get kafka instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// DeleteKafka 删除 Kafka 实例
func (h *MiddlewareHandler) DeleteKafka(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.DeleteKafka(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to delete kafka instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kafka instance deleted successfully"})
}

// KafkaActionRequest Kafka 操作请求
type KafkaActionRequest struct {
	Action string `json:"action" binding:"required"` // start/stop
}

// KafkaAction Kafka 实例操作
func (h *MiddlewareHandler) KafkaAction(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req KafkaActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.KafkaAction(c.Request.Context(), uint(accountID), instanceID, req.Action)
	if err != nil {
		h.logger.Error("failed to execute kafka action", zap.Error(err), zap.String("action", req.Action))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kafka action executed successfully"})
}

// ResizeKafkaRequest Kafka 调整规格请求
type ResizeKafkaRequest struct {
	Storage   int `json:"storage"`    // GB
	Bandwidth int `json:"bandwidth"`  // MB/s
}

// ResizeKafka 调整 Kafka 规格
func (h *MiddlewareHandler) ResizeKafka(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req ResizeKafkaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	spec := cloudprovider.KafkaSpec{
		Storage:   req.Storage,
		Bandwidth: req.Bandwidth,
	}

	err = h.service.ResizeKafka(c.Request.Context(), uint(accountID), instanceID, spec)
	if err != nil {
		h.logger.Error("failed to resize kafka instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "kafka instance resized successfully"})
}

// CreateElasticsearchRequest 创建 Elasticsearch 请求
type CreateElasticsearchRequest struct {
	AccountID     uint              `json:"account_id" binding:"required"`
	Name          string            `json:"name" binding:"required"`
	Version       string            `json:"version" binding:"required"`
	InstanceType  string            `json:"instance_type" binding:"required"`
	NodeCount     int               `json:"node_count"`
	Storage       int               `json:"storage"` // GB
	VPCID         string            `json:"vpc_id" binding:"required"`
	SubnetID      string            `json:"subnet_id" binding:"required"`
	ZoneID        string            `json:"zone_id"`
	Tags          map[string]string `json:"tags"`
}

// CreateElasticsearch 创建 Elasticsearch 实例
func (h *MiddlewareHandler) CreateElasticsearch(c *gin.Context) {
	var req CreateElasticsearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := cloudprovider.ElasticsearchConfig{
		Name:         req.Name,
		Version:      req.Version,
		InstanceType: req.InstanceType,
		NodeCount:    req.NodeCount,
		Storage:      req.Storage,
		VPCID:        req.VPCID,
		SubnetID:     req.SubnetID,
		ZoneID:       req.ZoneID,
		Tags:         req.Tags,
	}

	instance, err := h.service.CreateElasticsearch(c.Request.Context(), req.AccountID, config)
	if err != nil {
		h.logger.Error("failed to create elasticsearch instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

// ListElasticsearch 列出 Elasticsearch 实例
func (h *MiddlewareHandler) ListElasticsearch(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	filter := cloudprovider.ElasticsearchFilter{
		InstanceID:   c.Query("instance_id"),
		Status:       c.Query("status"),
		Version:      c.Query("version"),
		InstanceType: c.Query("instance_type"),
		RegionID:     c.Query("region_id"),
	}

	instances, err := h.service.ListElasticsearch(c.Request.Context(), uint(accountID), filter)
	if err != nil {
		h.logger.Error("failed to list elasticsearch instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instances)
}

// GetElasticsearch 获取 Elasticsearch 实例详情
func (h *MiddlewareHandler) GetElasticsearch(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	instance, err := h.service.GetElasticsearch(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to get elasticsearch instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instance)
}

// DeleteElasticsearch 删除 Elasticsearch 实例
func (h *MiddlewareHandler) DeleteElasticsearch(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.DeleteElasticsearch(c.Request.Context(), uint(accountID), instanceID)
	if err != nil {
		h.logger.Error("failed to delete elasticsearch instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "elasticsearch instance deleted successfully"})
}

// ElasticsearchActionRequest Elasticsearch 操作请求
type ElasticsearchActionRequest struct {
	Action string `json:"action" binding:"required"` // start/stop
}

// ElasticsearchAction Elasticsearch 实例操作
func (h *MiddlewareHandler) ElasticsearchAction(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req ElasticsearchActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	err = h.service.ElasticsearchAction(c.Request.Context(), uint(accountID), instanceID, req.Action)
	if err != nil {
		h.logger.Error("failed to execute elasticsearch action", zap.Error(err), zap.String("action", req.Action))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "elasticsearch action executed successfully"})
}

// ResizeElasticsearchRequest Elasticsearch 调整规格请求
type ResizeElasticsearchRequest struct {
	NodeCount int `json:"node_count"`
	Storage   int `json:"storage"` // GB
}

// ResizeElasticsearch 调整 Elasticsearch 规格
func (h *MiddlewareHandler) ResizeElasticsearch(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	instanceID := c.Param("id")

	var req ResizeElasticsearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	spec := cloudprovider.ElasticsearchSpec{
		NodeCount: req.NodeCount,
		Storage:   req.Storage,
	}

	err = h.service.ResizeElasticsearch(c.Request.Context(), uint(accountID), instanceID, spec)
	if err != nil {
		h.logger.Error("failed to resize elasticsearch instance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "elasticsearch instance resized successfully"})
}